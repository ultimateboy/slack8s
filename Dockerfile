FROM scratch
ADD slack8s /slack8s
COPY ./ld-linux-x86-64.so.2 /lib64/ld-linux-x86-64.so.2
COPY ./libc.so.6 /lib64/libc.so.6
COPY ./libpthread.so.0 /lib64/libpthread.so.0

ENTRYPOINT ["/slack8s"]
