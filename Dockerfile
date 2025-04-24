FROM golang:1.23

WORKDIR /src

# cosmtrk/air is a project auto-build tool
RUN go install github.com/air-verse/air@latest

# pressly/goose is a database migrations tool
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# sqlc is for SQL-to-Go code generation
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# set up to run as a normal user
RUN useradd user && mkdir /home/user && chown user:user /home/user && chown user:user /src
USER user
ENV GOPATH /home/user/go

# Copy the Go Modules manifests
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
COPY --chown=user ./go.mod go.mod
COPY --chown=user ./go.sum go.sum
RUN go mod download

COPY --chown=user . .

CMD ["air"]
