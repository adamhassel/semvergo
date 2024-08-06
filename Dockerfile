FROM golang:1.22-alpine AS build
WORKDIR /build
COPY . .
RUN go build -mod=vendor -o semvergo cmd/semvergo.go

FROM busybox
WORKDIR /semvergo
COPY --from=build /build/semvergo /bin/semvergo

ENTRYPOINT ["/bin/semvergo"]