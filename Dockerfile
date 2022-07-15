FROM golang:1.18.3-alpine3.16 as build

ENV GO111MODULE=on

WORKDIR /go/src/github.com/laupse/kubegraph
COPY . .

RUN go build -o /go/bin/kubegraph

FROM alpine:3.16

COPY --from=build /go/bin/kubegraph /go/bin/

CMD [ "/go/bin/kubegraph" ]