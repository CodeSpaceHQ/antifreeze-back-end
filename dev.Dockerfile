FROM ubuntu:16.04 AS build

# Combine these to reduce layers?
ENV PATH=${PATH}:/usr/local/go/bin
ENV GOPATH=${HOME}/go
ENV GOBIN=${HOME}/go/bin

RUN cd ~ && \
  apt update -qq && \
  apt install -y curl make git && \
  curl -o go.tar.gz https://dl.google.com/go/go1.9.3.linux-amd64.tar.gz && \
  tar -C /usr/local/ -xzf go.tar.gz && \
  curl -L -o dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
  chmod +x dep && \
  mv dep /usr/local/go/bin/dep

COPY Gopkg.lock Gopkg.toml ${HOME}/go/src/github.com/NilsG-S/antifreeze-back-end/
WORKDIR ${HOME}/go/src/github.com/NilsG-S/antifreeze-back-end/
RUN dep ensure --vendor-only

COPY . ${HOME}/go/src/github.com/NilsG-S/antifreeze-back-end/
RUN make

FROM google/cloud-sdk:latest

ENV CLOUDSDK_CORE_PROJECT=antifreeze-dev

RUN apt update -qq && \
  apt install -y default-jre \
  google-cloud-sdk-datastore-emulator

COPY --from=build ${HOME}/go/src/github.com/NilsG-S/antifreeze-back-end/bin/antifreeze-back-end /bin/antifreeze

CMD ["nohup", "gcloud", "beta", "emulators", "datastore", "start", "&", "&&", "antifreeze"]
