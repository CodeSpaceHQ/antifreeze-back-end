# Use a specific version!
FROM google/cloud-sdk:latest

ENV CLOUDSDK_CORE_PROJECT=antifreezedev

RUN apt update -qq && \
  apt install -y default-jre \
  google-cloud-sdk-datastore-emulator

# Mount a directory?
EXPOSE 8082

CMD ["gcloud", "beta", "emulators", "datastore", "start", "--host-port=0.0.0.0:8082"]
