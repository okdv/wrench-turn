# Wrench Turn
Self-hosted, flexible maintenance tracker for your cars, bikes, and everything in between

**WARNING**: This project is under active development in an alpha state, effectively an MVP. Would not recommend using yet unless you plan to contribute in some way, otherwise a beta release will be coming soon. This will be where everyday self hosters would likely benefit from the service. 

## Branches 
`develop`: The latest, unstable, branch. Where all commits are made to, and where all contributing branches are branched from. Though tested first, likelihood of broken things is fairly good. 

`main`: stable channel/branch. should pretty much always match the latest release branch 

`vX.X.X[-stage]`: particular version/release, e.g. v1.0.0-alpha, v.4.2.0

## Running (Production)
1) Clone repo, open in terminal: `git clone https://github.com/okdv/wrench-turn.git && cd wrench-turn`
2) Create `.env.production` file from `.env.development`: `cp .env.development .env.production`
3) Edit `.env.production` accordingly, need to change the below, but other vars may need editing depending on implementation:
    - `NODE_ENV=production`
    - `JWT_KEY=YOUR_CUSTOM_SECRET_KEY_DONT_COMMIT_OR_LEAVE_DEFAULT`
### Docker (recommended)
4) Edit `compose.yaml` accordingly, mainly port mapping. `.env.production` should match port mapping in `compose.yaml`.
5) Build and run `sudo docker-compose up` 

### Bare metal (not recommended)
#### Backend
4) [Install go](https://go.dev/dl/)
**Note:** check the image version used in [backend.Dockerfile](https://github.com/okdv/wrench-turn/blob/develop/backend.Dockerfile) if unsure which version to use. Usually assume latest stable version. 
5) Run `go build`, may need to adjust commands for your particular OS
6) Run generated build, `./wrench-turn`, `./wrench-turn.exe`, etc.

#### Frontend
7) [Install node](https://nodejs.org/en/download) or [nvm](https://github.com/nvm-sh/nvm)
**Note:** check the image version used in [frontend.Dockerfile](https://github.com/okdv/wrench-turn/blob/develop/frontend.Dockerfile) if unsure which version to use. Usually assume latest stable version. 
8) open frontend `cd frontend`
9) Install node modules `npm i`
10) Run node build `npm run build`

##### Static (recommended)
11) Run the generated build dir (/wrench-turn/frontend/build) with the desired web server, such as NGINX

##### Using Node (not recommended) 
11) Run `npm run preview`