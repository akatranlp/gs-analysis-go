#############################################
# Development Image with air
#############################################
FROM cosmtrek/air

WORKDIR /backend
COPY ./go.mod ./go.sum ./
RUN go mod download && \
    go install github.com/mattn/go-sqlite3 \
