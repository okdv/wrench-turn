# Contributing
This will cover everything you need to know to get started on contributing.

## Finding an issue
As with most projects, it's best to start with simple, non-urgent issues that provide an easy introduction into the project. This is no different. See [our good first issues](https://github.com/okdv/wrench-turn/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22) for a nice place to start. There is something for everyone, rather frontend (Svelte, Typescript), backend (Go, SQLite), or infrastructure (Docker, Github Actions, NGINX).

Once you have some experience, and feel comfortable taking lead on features or more complex issues. Check out our Beta roadmap [here](https://github.com/users/okdv/projects/7/views/1), these are a fantastic place to start, especially if you'd like to become a maintainer.

## Run the app in development mode
Eventually will containerize dev environment, but haven't yet. So for now, contributors should run this bare metal or within their own containers. 

### Backend
1) [Install go](https://go.dev/dl/)
**Note:** check the image version used in [backend.Dockerfile](https://github.com/okdv/wrench-turn/blob/develop/backend.Dockerfile) if unsure which version to use. Usually assume latest stable version. 
2) cd into `wrench-turn` directory, backend is located at services root level directory
3) Run `go build`, may need to adjust commands for your particular OS
4) Run generated build, `./wrench-turn`, `./wrench-turn.exe`, etc.
5) Should startup messages logged, something like "WrenchTurn server listening on port 8080", may also notice a ./data/sqlite-dev.db file 
6) Backend API should be available at `http://localhost:8080`

### Frontend 
1) [Install node](https://nodejs.org/en/download) or [nvm](https://github.com/nvm-sh/nvm)
**Note:** check the image version used in [frontend.Dockerfile](https://github.com/okdv/wrench-turn/blob/develop/frontend.Dockerfile) if unsure which version to use. Usually assume latest stable version. 
2) cd into `wrench-turn/frontend`
3) Run `npm run dev` 
4) Frontend should be available at `http://localhost:5173`

## Commit messages 
Please follow this git commit message format: `<type: chore, feat, fix>: <description> <rel issue number>`

e.g. `chore: cleaning up comments`, `feat: adding calendar feature #13`, `fix: home button deletes all data #16`

## How to contribute
Only maintainers can contribute directly to the develop branch. Everyone else should fork the repo and contribute there, like so:

1) Fork repository 
2) Clone your fork
    - `git clone https://github.com/your-username/wrench-turn.git && cd wrench-turn`
3) Create a branch off of the develop branch for what you're working on at that time
    - `git checkout -b super-cool-feature develop`
4) Commit and push
    - `git add . && git commit -m "feat: super awesome feature #420" && git push origin super-cool-feature`
    - **Note:** git pre-commit hook will attempt to run `go build`, `go test` successfully before committing 
5) Open a PR requesting to merge your feature branch with our develop branch, and wait for review
