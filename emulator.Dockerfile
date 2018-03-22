FROM google/cloud-sdk:193.0.0

ENV CLOUDSDK_CORE_PROJECT=antifreezedev

RUN apt update -qq && \
  apt install -y default-jre \
  google-cloud-sdk-datastore-emulator

# Exposes this port to other containers **only**
EXPOSE 8082

CMD ["gcloud", "beta", "emulators", "datastore", "start", "--host-port=0.0.0.0:8082"]
