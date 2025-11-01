FROM --platform=$BUILDPLATFORM golang:alpine AS build

ARG TARGETARCH
WORKDIR /src

RUN apk --no-cache add git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go run github.com/swaggo/swag/cmd/swag@latest init --parseDependency

ENV CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH}
RUN go build -ldflags="-s -w" -o /src/mozhi ./   # adjust path if needed

FROM scratch AS runtime
WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /src/mozhi /app/mozhi

EXPOSE 3000
ENTRYPOINT ["/app/mozhi", "serve"]
