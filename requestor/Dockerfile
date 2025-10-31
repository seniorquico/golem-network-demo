ARG UBUNTU_RELEASE=24.04
ARG CHISEL_VERSION=v1.3.0
ARG YAGNA_VERSION=v0.17.6

FROM ubuntu:$UBUNTU_RELEASE AS chisel_builder
ARG TARGETARCH UBUNTU_RELEASE CHISEL_VERSION YAGNA_VERSION
SHELL ["/bin/bash", "-oeux", "pipefail", "-c"]

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    apt-get clean -y && \
    rm -rf /var/lib/apt/lists/*

ADD "https://github.com/canonical/chisel/releases/download/${CHISEL_VERSION}/chisel_${CHISEL_VERSION}_linux_${TARGETARCH}.tar.gz" \
    chisel.tar.gz

RUN tar -xvf chisel.tar.gz -C /usr/bin/

RUN mkdir /staging-rootfs && \
    chisel cut --release "ubuntu-$UBUNTU_RELEASE" --root /staging-rootfs \
        base-files_base \
        base-files_release-info \
        base-files_chisel \
        busybox_bins \
        ca-certificates_data

RUN cp /etc/passwd /etc/group /staging-rootfs/etc && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local/bin && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local/share && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local/share/yagna && \
    install -o ubuntu -g ubuntu -d /staging-resources

ADD "https://github.com/golemfactory/yagna/releases/download/${YAGNA_VERSION}/golem-requestor-linux-${YAGNA_VERSION}.tar.gz" \
    golem-requestor.tar.gz

RUN tar -xvf golem-requestor.tar.gz --strip-components=1 -C /staging-rootfs/home/ubuntu/.local/bin \
        golem-requestor-linux-${YAGNA_VERSION}/gftp \
        golem-requestor-linux-${YAGNA_VERSION}/yagna && \
    chown -R ubuntu:ubuntu /staging-rootfs/home/ubuntu/.local/bin/ && \
    echo "Stats allow" > /staging-rootfs/home/ubuntu/.local/share/yagna/CONSENT && \
    chown -R ubuntu:ubuntu /staging-rootfs/home/ubuntu/.local/share/yagna/

FROM scratch

COPY --from=chisel_builder /staging-rootfs /

USER ubuntu
WORKDIR /home/ubuntu
ENV PATH="/home/ubuntu/.local/bin:${PATH}"
ENTRYPOINT ["/home/ubuntu/.local/bin/yagna"]
