# Use a specific version!
FROM google/cloud-sdk:latest

ENV CLOUDSDK_CORE_PROJECT=antifreeze-dev

RUN apt update -qq && \
  apt install -y default-jre \
  google-cloud-sdk-datastore-emulator

EXPOSE 8080

CMD ["gcloud", "beta", "emulators", "datastore", "start", "--host-port=0.0.0.0:8080"]
