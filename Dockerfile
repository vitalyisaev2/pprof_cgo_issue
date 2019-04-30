FROM fedora:29

RUN dnf update -y
RUN dnf install -y openssl-devel gcc git wget tar gdb

RUN wget https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.12.4.linux-amd64.tar.gz

ENV PATH="${PATH}:/usr/local/go/bin"
ENV GO111MODULE="on"

RUN git clone https://github.com/vitalyisaev2/pprof_cgo_issue && \
    cd pprof_cgo_issue && \
    go build

WORKDIR ./pprof_cgo_issue
