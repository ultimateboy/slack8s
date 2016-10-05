FROM scratch
ADD slack8s /slack8s

ENTRYPOINT ["/slack8s"]
