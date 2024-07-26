# Start by building the application.
FROM golang:1.22 as build

WORKDIR /go/src/app
COPY go.mod *.go ./
COPY static static

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/app /
EXPOSE 80
ENTRYPOINT ["/app", "-addr", ":80"]
