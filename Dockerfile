#############################################
# Builder web
#############################################
FROM node:20.11.1-alpine3.19 as builder-web

WORKDIR /app/build
COPY ./web/package.json ./web/yarn.lock ./
RUN yarn --pure-lockfile

COPY ./web ./
RUN yarn build

#############################################
# Builder go
#############################################
FROM golang:1.22-alpine3.19 as builder-go

RUN apk add --no-cache gcc musl-dev
RUN go install github.com/vektra/mockery/v2@v2.40.1

WORKDIR /app/build

COPY ./go.mod ./go.sum ./

# Install go-sqlite3 so the build does not take too long
RUN go mod download && \
    go install github.com/mattn/go-sqlite3

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./api ./api

RUN CGO_ENABLED=1 GOOS=linux go build -o ./gs-analysis cmd/gs-analysis/main.go

#############################################
# Runtime image
#############################################
FROM alpine:3.19 as release

ENV FRONTEND_PATH=/public
ENV PORT=3000
EXPOSE 3000

RUN adduser -D gorunner
USER gorunner

WORKDIR /app
RUN mkdir db

COPY --chown=gorunner:gorunner --from=builder-go /app/build/gs-analysis ./gs-analysis
COPY --chown=gorunner:gorunner --from=builder-web /app/build/dist ./public

ENTRYPOINT /app/gs-analysis
