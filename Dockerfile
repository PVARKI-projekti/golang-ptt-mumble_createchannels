FROM golang:alpine AS BUILD
ENV GO111MODULE=on

WORKDIR /app
COPY ./ /app

# This removes debug information from the binary
# Assumes go 1.10+
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -ldflags="-w -s" -o app "github.com/PVARKI-projekti/golang-ptt-mumble_createchannels"


FROM gcr.io/distroless/static as production
COPY --from=BUILD /app/app /app
ENTRYPOINT [ "/app" ]
