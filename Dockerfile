FROM golang:1.20 as common-build-stage

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \ 
    go build -v -o ./server ./main.go

# Development build stage
FROM common-build-stage as development-build-stage

WORKDIR /usr/src/app

CMD ["./server"]


# Production build stage
FROM golang:1.20 as production-build-stage

WORKDIR /usr/src/app

COPY --from=common-build-stage --chown=go:go /usr/src/app/server ./

CMD ["./server"]