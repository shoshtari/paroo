FROM golang:1.23.2 AS build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /opt/paroo main.go


# FROM scratch
# COPY --from=build /opt/paroo /opt/paroo
# COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/opt/paroo"]
