FROM ubuntu:16.04 AS build

ENV PATH=${PATH}:/usr/local/go/bin
ENV GOPATH=${HOME}/go
ENV GOBIN=${HOME}/go/bin

RUN cd ~ && \
  apt update -qq && \
  apt install -y curl make && \
  curl -o go.tar.gz https://dl.google.com/go/go1.9.3.linux-amd64.tar.gz && \
  tar -C /usr/local/ -xzf go.tar.gz && \
  curl -L -o dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
  chmod +x dep && \
  mv dep /usr/local/go/bin/dep

COPY Gopkg.lock Gopkg.toml ${HOME}/go/src/github.com/NilsG-S/antifreeze-back-end/
WORKDIR ${HOME}/go/src/github.com/NilsG-S/antifreeze-back-end/
RUN dep ensure

COPY . ${HOME}/go/src/github.com/NilsG-S/antifreeze-back-end/
RUN make

FROM ubuntu:16:04

COPY --from=build ${HOME}/go/src/github.com/NilsG-S/antifreeze-back-end/bin/antifreeze-back-end antifreeze

# RUN export CLOUD_SDK_REPO="cloud-sdk-$(lsb_release -c -s)" && \
#   echo "deb http://packages.cloud.google.com/apt $CLOUD_SDK_REPO main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
#   curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - && \
#   apt-get update && sudo apt-get install google-cloud-sdk

CMD ["antifreeze"]
