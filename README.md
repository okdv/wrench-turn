# Wrench Turn

## Running

If necessary

`cd wrench-turn`

### Docker (Production)

Using Docker Compose 

`sudo docker-compose up` 

### Local

Must have Golang and Node installed 

Build backend

`go build`

Start backend

`./wrench-turn`

Open frontend directory

`cd frontend`

Install

`npm install`

If you simply want to use Node to render the frontend for development, run

`npm run dev`

#### Building frontend locally (Production)

Build frontend (creates static app with svelte-static-adapter)

`npm run build`

Run the generated build dir (/wrench-turn/frontend/build) with the desired web server, such as NGINX

## Contributing

Pull requests should be made on the `develop` branch 

git pre-commit hook will attempt to run `go build`, `go test` successfully before committing 

Please follow this git commit message format

`<type: chore, feat, fix>: <description> <rel issue number>`

e.g. `chore: cleaning up comments`, `feat: adding calendar feature #13`, `fix: home button deletes all data #16`