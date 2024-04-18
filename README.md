# syncwatch-be

Prerequisites:
Have docker installed to execute tasks. For I'm using https://taskfile.dev/installation

Create DB like so:
`docker run -d --name db -p 5432:5432 -e POSTGRES_PASSWORD=pswd123 -e POSTGRES_USER=postgres -e POSTGRES_DB=syncwatch postgres`

// fill environment variables