FROM google/cloud-sdk:latest

ENV CLOUDSDK_CORE_PROJECT=antifreeze-dev

RUN apt update -qq && \
  apt install -y default-jre \
  google-cloud-sdk-datastore-emulator

CMD ["gcloud", "beta", "emulators", "datastore", "start"]
