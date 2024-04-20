FROM node:20-alpine AS build 

ARG NODE_ENV=production
ENV NODE_ENV=production

WORKDIR /app 

COPY ./frontend /app/

RUN npm ci 
RUN npm run build 

FROM nginx:1.25-alpine 

WORKDIR /etc/nginx
COPY --from=build /app/nginx.conf .

WORKDIR /usr/share/nginx/html 
RUN rm -rf ./* 
COPY --from=build /app/build .

ENTRYPOINT ["nginx", "-g", "daemon off;"]