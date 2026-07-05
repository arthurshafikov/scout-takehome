FROM golang:1.25.5 AS build

ENV CODE_DIR /go/src/
ENV BIN_FILE /go/src/app

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

ADD ./cmd ${CODE_DIR}cmd
ADD ./internal ${CODE_DIR}internal
ADD ./configs ${CODE_DIR}configs

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/app/main.go

FROM alpine:3.9 AS base

LABEL MAINTAINERS="arthurshafikov"

RUN apk add --no-cache tzdata

ENV BIN_FILE /go/src/app
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIGS_FOLDER /go/src/configs
COPY --from=build ${CONFIGS_FOLDER} ${CONFIGS_FOLDER}

FROM base AS dev
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait /wait
RUN chmod +x /wait
CMD /wait && \
    ${BIN_FILE} -cfgFolder ${CONFIGS_FOLDER} -env /go/src/

FROM base AS prod

EXPOSE 8080

CMD ${BIN_FILE} -cfgFolder ${CONFIGS_FOLDER} -env /go/src/
