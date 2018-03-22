FROM golang:1.9.3 AS build

ENV REPO=${GOPATH}/src/github.com/NilsG-S/antifreeze-back-end/

RUN cd / && \
curl -L -o dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
chmod +x dep && \
mv dep /usr/local/go/bin/dep

COPY Gopkg.lock Gopkg.toml ${REPO}
WORKDIR ${REPO}
RUN dep ensure --vendor-only

COPY . ${REPO}
RUN make && mv bin/antifreeze-back-end /bin/antifreeze

# Use Alpine?
FROM ubuntu:16.04

ENV DATASTORE_PROJECT_ID=antifreezedev

COPY --from=build /bin/antifreeze /bin/antifreeze

EXPOSE 8081

CMD ["antifreeze"]
