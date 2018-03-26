FROM google/cloud-sdk:193.0.0

ENV CLOUDSDK_CORE_PROJECT=antifreezedev

# Exposes this port to other containers **only**
EXPOSE 8082

# For some reason, data seems to presist between rebuilds...publish port for gcloud reset?
# Or, use a volume so the data file can be deleted
CMD ["gcloud", "beta", "emulators", "datastore", "start", "--host-port=0.0.0.0:8082"]
