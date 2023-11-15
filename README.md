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