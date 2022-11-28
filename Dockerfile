FROM golang:1.17.3-alpine3.14 AS dataeng
COPY . /w
WORKDIR /w

RUN go build github.com/Mirantis/dataops-dataeng/dataengctl/bin

FROM alpine:3.14.2
COPY --from=dataengctl /w/dataeng dataengctl/bin/
ENTRYPOINT [ "dataengctl/bin" ]
