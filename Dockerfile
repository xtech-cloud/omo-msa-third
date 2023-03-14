FROM alpine:3.11
ADD omo.msa.third /usr/bin/omo.msa.third
ENV MSA_REGISTRY_PLUGIN
ENV MSA_REGISTRY_ADDRESS
ENTRYPOINT [ "omo.msa.third" ]
