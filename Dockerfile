FROM golang:1.12.6-stretch

RUN apt-get update -y &&\
    apt-get install -y \
        git

ENV PATH="${GOPATH}/bin:${PATH}"

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR ${GOPATH}/src/github.com/dev-schueppchen/Kirby

ADD . .

RUN mkdir -p /etc/config

RUN dep ensure -v

RUN go build -v -o /usr/bin/kirby -ldflags "\
	-X github.com/dev-schueppchen/Kirby/internal/static.AppVersion=$(git describe --tags) \
	-X github.com/dev-schueppchen/Kirby/internal/static.AppCommit=$(git rev-parse HEAD) \
        -X github.com/dev-schueppchen/Kirby/internal/static.Release=TRUE" \
            ./cmd/kirby/*.go

CMD kirby \
        -c /etc/config/config.yml
