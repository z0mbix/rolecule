FROM docker.io/amazonlinux:2

RUN yum install -y curl procps-ng systemd python3 python3-pip util-linux && \
  python3 -m pip install ansible ansible-core

# # Install goss (https://github.com/goss-org/goss)
RUN curl -sSL https://github.com/goss-org/goss/releases/latest/download/goss-linux-amd64 -o /usr/local/bin/goss && \
  chmod +rx /usr/local/bin/goss

# VOLUME [ "/sys/fs/cgroup" ]
CMD ["/sbin/init"]
