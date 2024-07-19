FROM --platform=$BUILDPLATFORM golang:alpine AS build

ARG TARGETARCH

WORKDIR /src
RUN apk --no-cache add git
COPY . .

ENV GOPRIVATE=codeberg.org/aryak/libmozhi
#RUN go mod download
RUN go run github.com/swaggo/swag/cmd/swag@latest init --parseDependency
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /src/mozhi

FROM alpine:3.16 AS bin

WORKDIR /app
COPY --from=build /src/mozhi .

EXPOSE 3000

CMD ["/app/mozhi", "serve"]
