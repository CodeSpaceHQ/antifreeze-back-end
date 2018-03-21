# Replace this with golang to simplify?
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

FROM ubuntu:16.04

# `db` references the database emulator
ENV DATASTORE_EMULATOR_HOST=db:8082 DATASTORE_PROJECT_ID=antifreezedev

# Note that this method doesn't carry over any static files
COPY --from=build /bin/antifreeze /bin/antifreeze

EXPOSE 8081

CMD ["antifreeze"]
