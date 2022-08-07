# syntax=docker/dockerfile:1.3-labs

FROM golang:1.19

ENV GO111MODULE="on"
#ENV GOOS="linux"
#ENV CGO_ENABLED=0

RUN <<EOT
go install github.com/cweill/gotests/gotests@latest
go install github.com/fatih/gomodifytags@latest
go install github.com/josharian/impl@latest
go install github.com/haya14busa/goplay/cmd/goplay@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/tools/gopls@latest
go install github.com/ramya-rao-a/go-outline@latest
go install github.com/x-motemen/gore/cmd/gore@latest
EOT

# WORKDIR /usr/local/src

# COPY . .
