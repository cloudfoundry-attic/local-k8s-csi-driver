FROM alpine
LABEL maintainers="CF Persistence Team Authors"
LABEL description="Local Kubernetes CSI Driver"

COPY ./_output/local-k8s-csi-driver /local-k8s-csi-driver
COPY ./_output/csc /csc
ENTRYPOINT ["/local-k8s-csi-driver"]
