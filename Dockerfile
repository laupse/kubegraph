FROM golang:1.19.2-alpine3.16 as build

ENV GO111MODULE=on

WORKDIR /go/src/github.com/laupse/kubegraph
COPY . .

WORKDIR /go/src/github.com/laupse/kubegraph/cmd
RUN go build -o /go/bin/kubegraph

FROM alpine:3.16

COPY --from=build /go/bin/kubegraph /go/bin/

CMD [ "/go/bin/kubegraph" ]