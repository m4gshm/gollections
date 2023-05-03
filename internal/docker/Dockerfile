FROM buildpack-deps:bullseye-scm

RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
	; \
	rm -rf /var/lib/apt/lists/*

ENV OUT_DIR /usr/local/go
ENV PATH /usr/local/go/bin:$PATH

ENV GOLANG_VERSION 1.18rc1

COPY go/bin/linux_amd64 $OUT_DIR/bin
COPY go/src $OUT_DIR/src
COPY go/pkg/linux_amd64 $OUT_DIR/pkg/linux_amd64
COPY go/pkg/tool/linux_amd64 $OUT_DIR/pkg/tool/linux_amd64
COPY go/pkg/include $OUT_DIR/pkg/include
COPY go/pkg/obj $OUT_DIR/pkg/obj
COPY go/lib $OUT_DIR/lib


ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH