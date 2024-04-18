package json

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:revive,cognitive-complexity
func TestReadJSON(t *testing.T) {
	type Payload struct {
		Name string `json:"name"`
	}

	cases := []struct {
		name           string
		contentType    string
		body           string
		expectedError  *MalformedRequest
		expectedOutput Payload
	}{
		{
			name:        "Well formatted payload",
			contentType: "application/json",
			body:        `{"name":"example"}`,
			expectedOutput: Payload{
				Name: "example",
			},
			expectedError: nil,
		},
		{
			name:        "Bad content type",
			contentType: "text/plain",
			body:        `{"name":"example"}`,
			expectedError: &MalformedRequest{
				Status: http.StatusUnsupportedMediaType,
				Msg:    "Content-Type header is not application/json",
			},
		},
		{
			name:        "Field mismatch",
			contentType: "application/json",
			body:        `{"badname":"example"}`,
			expectedError: &MalformedRequest{
				Status: http.StatusBadRequest,
				Msg:    "Request body contains unknown field \"badname\"",
			},
			expectedOutput: Payload{},
		},
		{
			name:        "Badly formatted JSON ending",
			contentType: "application/json",
			body:        `{"name":"example"`,
			expectedError: &MalformedRequest{
				Status: http.StatusBadRequest,
				Msg:    "Request body contains badly-formed JSON",
			},
		},
		{
			name:        "Badly formatted JSON, missing quotes",
			contentType: "application/json",
			body:        `{name":"example"}`,
			expectedError: &MalformedRequest{
				Status: http.StatusBadRequest,
				Msg:    "Request body contains badly-formed JSON (at position 2)",
			},
		},
		{
			name:        "Too large payload",
			contentType: "application/json",
			body:        `{"name":"` + strings.Repeat("a", maxPayloadSize) + `"}`,
			expectedError: &MalformedRequest{
				Status: http.StatusRequestEntityTooLarge,
				Msg:    "Request body must not be larger than 1MB",
			},
		},
		{
			name:        "Payload can not be empty",
			contentType: "application/json",
			body:        ``,
			expectedError: &MalformedRequest{
				Status: http.StatusBadRequest,
				Msg:    "Request body must not be empty",
			},
		},
		{
			name:        "Payload containing invalid value",
			contentType: "application/json",
			body:        `{"name":1}`,
			expectedError: &MalformedRequest{
				Status: http.StatusBadRequest,
				Msg:    "Request body contains an invalid value for the \"name\" field (at position 9)",
			},
		},
		{
			name:        "Payload containing more than one json",
			contentType: "application/json",
			body:        `{"name":"john"}{"name":"doe"}`,
			expectedError: &MalformedRequest{
				Status: http.StatusBadRequest,
				Msg:    "Request body must only contain a single JSON object",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tc.body)))
			req.Header.Set("Content-Type", tc.contentType)
			w := httptest.NewRecorder()

			var dst Payload
			var err *MalformedRequest
			errors.As(ReadJSON(w, req, &dst), &err)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError.Status, err.Status)
				if tc.expectedError.Msg != "" {
					assert.Equal(t, tc.expectedError.Msg, err.Msg)
				}
			}
		})
	}
}

//nolint:revive,cognitive-complexity
func TestWriteJSON(t *testing.T) {
	data := map[string]string{"key": "value"}
	headers := http.Header{"Test-Header": []string{"Test-Value"}}

	tests := []struct {
		name      string
		status    int
		data      any
		headers   http.Header
		err       error
		expHeader http.Header
		expBody   string
	}{
		{
			name:      "ValidData",
			status:    200,
			data:      data,
			headers:   headers,
			err:       nil,
			expHeader: headers,
			expBody:   `{"key":"value"}`,
		},
		{
			name:      "ErrorData",
			status:    500,
			data:      make(chan int), // Not serializable
			headers:   headers,
			err:       errors.New("json: unsupported type: chan int"),
			expHeader: headers,
			expBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			err := WriteJSON(recorder, tt.status, tt.data, tt.headers)

			if tt.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				result := recorder.Result() // nolint
				assert.Equal(t, tt.status, result.StatusCode)

				// Headers
				assert.Equal(t, "application/json", result.Header["Content-Type"][0])
				for name, values := range tt.expHeader {
					assert.Equal(t, values, result.Header[name])
				}

				// Body
				body, _ := io.ReadAll(result.Body)
				assert.JSONEq(t, tt.expBody, string(body))
			}
		})
	}
}
