FROM golang:1.19 as builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go generate
RUN GOOS=linux go build -o shipply-client-backend -trimpath -ldflags "-w -s"

FROM alpine:3

WORKDIR /app

RUN apk add --no-cache ca-certificates gcompat aws-cli

COPY --from=builder /build/deploy/init.sh ./
COPY --from=builder /build/shipply-client-backend ./

EXPOSE 8080

ENTRYPOINT [ "./init.sh" ]

CMD [ "./shipply-client-backend" ]