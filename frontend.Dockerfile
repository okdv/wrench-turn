FROM node:20-alpine AS build 

WORKDIR /app 
COPY ./frontend /app/
RUN npm ci 
RUN npm run build 

FROM nginx:1.25-alpine 

WORKDIR /usr/share/nginx/html 
RUN rm -rf ./* 
COPY --from=build /app/build .
ENTRYPOINT ["nginx", "-g", "daemon off;"]