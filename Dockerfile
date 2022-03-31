# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
FROM golang:1.17.1-alpine3.13 as builder
WORKDIR /go/src/gitlab.apulis.com.cn/hjl/blank-web-app

ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add git pkgconfig build-base && \
    go get -u github.com/kevinburke/go-bindata/go-bindata

# Cache go modules
COPY go.mod .
COPY go.sum .
ADD . .
RUN go mod download
RUN make build

FROM alpine:3.11
COPY --from=builder /go/src/gitlab.apulis.com.cn/hjl/blank-web-app/bin/blankWebApp /root/blankWebApp
COPY --from=builder /go/src/gitlab.apulis.com.cn/hjl/blank-web-app/configs/config.yaml /root
WORKDIR /root
ENTRYPOINT ["./blankWebApp", "run"]
