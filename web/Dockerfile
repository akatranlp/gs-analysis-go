#############################################
# Development Image with nodejs
#############################################
FROM node:20.11.1-alpine as builder

RUN apk --no-cache --update-cache --available upgrade \
    && apk add git bash yarn

USER node
WORKDIR /app

CMD yarn && yarn dev
