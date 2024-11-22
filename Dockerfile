# =========================================================
FROM golang:1.23-alpine AS build

WORKDIR /src
COPY . .

RUN go mod tidy && CGO_ENABLED=0 go build -o /src/jobd .

# =========================================================
FROM alpine:3.20

COPY --from=build /src/jobd /jobd

ENTRYPOINT ["/jobd"]
# =========================================================

