FROM docker.io/rockylinux/rockylinux:9.1

RUN dnf install -y systemd systemd-libs util-linux ansible-core procps-ng && \
  rm -f /lib/systemd/system/multi-user.target.wants/*;\
  rm -f /etc/systemd/system/*.wants/*;\
  rm -f /lib/systemd/system/local-fs.target.wants/*; \
  rm -f /lib/systemd/system/sockets.target.wants/*udev*; \
  rm -f /lib/systemd/system/sockets.target.wants/*initctl*; \
  rm -f /lib/systemd/system/basic.target.wants/*;\
  rm -f /lib/systemd/system/anaconda.target.wants/*;

# Install goss (https://github.com/goss-org/goss)
RUN curl -sSL https://github.com/goss-org/goss/releases/latest/download/goss-linux-amd64 -o /usr/local/bin/goss && \
  chmod +rx /usr/local/bin/goss

VOLUME [ "/sys/fs/cgroup" ]
CMD ["/usr/sbin/init"]
