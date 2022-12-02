ARG VARIANT=1.19-bullseye
FROM golang:${VARIANT}

# RUN go get -x <your-dependency-or-tool>
RUN go install github.com/cweill/gotests/gotests@latest
RUN go install github.com/fatih/gomodifytags@latest
RUN go install github.com/josharian/impl@latest
RUN go install github.com/haya14busa/goplay/cmd/goplay@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/1pkg/goutline@latest
RUN go install github.com/1pkg/gopium@latest
RUN go install github.com/mgechev/revive@latest
RUN go install github.com/go-critic/go-critic/cmd/gocritic@latest
RUN go install mvdan.cc/gofumpt@latest

RUN apt-get update && apt-get -y install libx11-dev