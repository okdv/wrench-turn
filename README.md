# Wrench Turn
Self-hosted, flexible maintenance tracker for your cars, bikes, and everything in between

**WARNING**: This project is under active development in an alpha state, effectively an MVP. Would not recommend using yet unless you plan to contribute in some way, otherwise a beta release will be coming soon. This will be where everyday self hosters would likely benefit from the service. 

## Running (Production)

Clone repo, open in terminal 
`git clone https://github.com/okdv/wrench-turn.git`
`cd wrench-turn`

### Docker

Using Docker Compose 

`sudo docker-compose up` 

### Bare metal

Must have Golang and Node installed 

#### Backend

Build with go

`go build`

Start backend

`./wrench-turn`

Open frontend directory

`cd frontend`

Install

`npm install`

#### Frontend (static build)

Build frontend (creates static app with svelte-static-adapter)

`npm run build`

Run the generated build dir (/wrench-turn/frontend/build) with the desired web server, such as NGINX

#### Frontend (svelte preview)

If you simply want to use Node to render the frontend for development, run

`npm run preview`

## Contributing

1) Fork repository 
2) Clone your fork
    - `git clone https://github.com/your-username/wrench-turn.git && cd wrench-turn`
3) Create a branch for what you're working on at that time
    - `git checkout -b super-cool-feature`
4) Commit and push
    - `git add . && git commit -m "feat: super awesome feature #420" && git push origin super-cool-feature`
    - **Note:** git pre-commit hook will attempt to run `go build`, `go test` successfully before committing 
5) Open a PR requesting to merge your feature branch with our develop branch, and wait for review

Please follow this git commit message format

`<type: chore, feat, fix>: <description> <rel issue number>`

e.g. `chore: cleaning up comments`, `feat: adding calendar feature #13`, `fix: home button deletes all data #16`