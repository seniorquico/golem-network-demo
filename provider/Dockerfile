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
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.config && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.config/ya-runtime-salad && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local/bin && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local/lib && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local/lib/yagna && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local/lib/yagna/plugins && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local/share && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local/share/ya-provider && \
    install -o ubuntu -g ubuntu -d /staging-rootfs/home/ubuntu/.local/share/yagna && \
    install -o ubuntu -g ubuntu -d /staging-resources

ADD "https://github.com/golemfactory/yagna/releases/download/${YAGNA_VERSION}/golem-provider-linux-${YAGNA_VERSION}.tar.gz" \
    golem-provider.tar.gz
ADD "https://github.com/golemfactory/ya-installer-resources/releases/latest/download/resources.tar.gz" \
    golem-resources.tar.gz

RUN tar -xvf golem-provider.tar.gz --strip-components=1 -C /staging-rootfs/home/ubuntu/.local/bin \
        golem-provider-linux-${YAGNA_VERSION}/golemsp \
        golem-provider-linux-${YAGNA_VERSION}/ya-provider \
        golem-provider-linux-${YAGNA_VERSION}/yagna && \
    chown -R ubuntu:ubuntu /staging-rootfs/home/ubuntu/.local/bin/ && \
    tar -xvf golem-provider.tar.gz --strip-components=2 -C /staging-rootfs/home/ubuntu/.local/lib/yagna/plugins \
        golem-provider-linux-${YAGNA_VERSION}/plugins/exe-unit && \
    chown -R ubuntu:ubuntu /staging-rootfs/home/ubuntu/.local/lib/yagna/plugins/ && \
    echo "Stats allow" > /staging-rootfs/home/ubuntu/.local/share/yagna/CONSENT && \
    chown -R ubuntu:ubuntu /staging-rootfs/home/ubuntu/.local/share/yagna/ && \
    tar -xvf golem-resources.tar.gz -C /staging-resources && \
    chown -R ubuntu:ubuntu /staging-resources/

COPY --chown=ubuntu:ubuntu ./ya-runtime-salad/conf/ya-runtime-salad.json /staging-rootfs/home/ubuntu/.local/lib/yagna/plugins/ya-runtime-salad.json
COPY --chown=ubuntu:ubuntu ./ya-runtime-salad/ya-runtime-salad /staging-rootfs/home/ubuntu/.local/lib/yagna/plugins/ya-runtime-salad

FROM ubuntu:$UBUNTU_RELEASE AS resource_builder
SHELL ["/bin/bash", "-oeux", "pipefail", "-c"]

COPY --from=chisel_builder /staging-rootfs /
COPY --from=chisel_builder /staging-resources /home/ubuntu/resources

USER ubuntu
WORKDIR /home/ubuntu
ENV PATH="/home/ubuntu/.local/bin:${PATH}"

RUN ya-provider pre-install && \
    golemsp manifest-bundle add /home/ubuntu/resources && \
    ya-provider rule set outbound everyone --mode whitelist

FROM scratch

COPY --from=chisel_builder /staging-rootfs /
COPY --from=resource_builder /home/ubuntu/.local/share/ya-provider /home/ubuntu/.local/share/ya-provider

USER ubuntu
WORKDIR /home/ubuntu
ENV PATH="/home/ubuntu/.local/bin:${PATH}"
ENTRYPOINT ["/home/ubuntu/.local/bin/golemsp"]
