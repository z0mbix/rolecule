FROM ubuntu:22.04

ENV container docker
ENV DEBIAN_FRONTEND noninteractive

RUN sed -i 's/# deb/deb/g' /etc/apt/sources.list

# hadolint ignore=DL3008
RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates software-properties-common systemd curl cron gpg-agent less sudo bash iproute2 net-tools vim \
  && apt-add-repository -y ppa:ansible/ansible 1>/dev/null \
  && apt-get install -y --no-install-recommends ansible ansible-lint \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN dpkg-reconfigure ca-certificates

WORKDIR /lib/systemd/system/sysinit.target.wants/
# hadolint ignore=SC2010,SC2086
RUN ls | grep -v systemd-tmpfiles-setup | xargs rm -f $1

RUN rm -f /lib/systemd/system/multi-user.target.wants/* \
  /etc/systemd/system/*.wants/* \
  /lib/systemd/system/local-fs.target.wants/* \
  /lib/systemd/system/sockets.target.wants/*udev* \
  /lib/systemd/system/sockets.target.wants/*initctl* \
  /lib/systemd/system/basic.target.wants/* \
  /lib/systemd/system/anaconda.target.wants/* \
  /lib/systemd/system/plymouth* \
  /lib/systemd/system/systemd-update-utmp*

# Install goss (https://github.com/goss-org/goss)
RUN curl -sSL https://github.com/goss-org/goss/releases/latest/download/goss-linux-amd64 -o /usr/local/bin/goss && \
  chmod +rx /usr/local/bin/goss

WORKDIR /
RUN systemctl set-default multi-user.target
ENV init /lib/systemd/systemd
VOLUME [ "/sys/fs/cgroup" ]

ENTRYPOINT ["/lib/systemd/systemd"]
