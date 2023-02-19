
FROM alpine:3.15

MAINTAINER Arshad Zackeriya

RUN apk add --no-cache \
		ca-certificates

RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

ENV GOLANG_VERSION 1.13.3

RUN set -eux; \
	apk add --no-cache --virtual .build-deps \
		bash \
		gcc \
		musl-dev \
		openssl \
		go \
	; \
	export \
		GOROOT_BOOTSTRAP="$(go env GOROOT)" \
		GOOS="$(go env GOOS)" \
		GOARCH="$(go env GOARCH)" \
		GOHOSTOS="$(go env GOHOSTOS)" \
		GOHOSTARCH="$(go env GOHOSTARCH)" \
	; \

	apkArch="$(apk --print-arch)"; \
	case "$apkArch" in \
		armhf) export GOARM='6' ;; \
		x86) export GO386='387' ;; \
	esac; \
	\
	wget -O go.tgz "https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz"; \
	echo '4f7123044375d5c404280737fbd2d0b17064b66182a65919ffe20ffe8620e3df *go.tgz' | sha256sum -c -; \
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz; \
	\
	cd /usr/local/go/src; \
	./make.bash; \
	\
	rm -rf \
		/usr/local/go/pkg/bootstrap \
		/usr/local/go/pkg/obj \
	; \
	apk del .build-deps; \
	\
	export PATH="/usr/local/go/bin:$PATH"; \
	go version


ENV GOPATH /codepool
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p ${GOPATH}/src/aws-secret-manager-test ${GOPATH}/bin


COPY . ${GOPATH}/src/aws-secret-manager-test/

RUN chmod -R 775 ${GOPATH}/src/aws-secret-manager-test/

WORKDIR ${GOPATH}/src/aws-secret-manager-test/
RUN apk add glide git
RUN echo N | glide init
RUN glide install
RUN glide up
RUN go build main.go