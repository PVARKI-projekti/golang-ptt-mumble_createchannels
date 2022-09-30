FROM golang:alpine AS BUILD

LABEL maintainer="eero.afheurlin@iki.fi"

ENV GO111MODULE=on

WORKDIR /app
COPY ./ /app

# This removes debug information from the binary
# Assumes go 1.10+
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -ldflags="-w -s" -o app "github.com/rambo/mumble_createchannels"


FROM gcr.io/distroless/static
COPY --from=BUILD /app/app /app
ENV PORT 3000
EXPOSE 3000
ENTRYPOINT [ "/app" ]
