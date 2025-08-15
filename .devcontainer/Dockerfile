FROM --platform=linux/amd64 ubuntu:20.04
LABEL org.opencontainers.image.authors="support@halon.io"

ARG HALON_REPO_USER
ARG HALON_REPO_PASS

RUN apt-get update

RUN apt-get install -y wget gnupg
RUN apt-get install -y apt-transport-https

RUN wget -qO - https://raw.githubusercontent.com/halon/changelog/master/pgp-keys/7F0A73B5.asc | apt-key add -
RUN echo "deb https://repo.halon.io/ focal stable" >> /etc/apt/sources.list.d/halon.list
RUN echo "machine repo.halon.io login ${HALON_REPO_USER} password ${HALON_REPO_PASS}" >> /etc/apt/auth.conf
RUN apt-get update && apt-get install -y halon=5.10.0

RUN /usr/bin/install -d /var/run/halon
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH

RUN apt-get install -y build-essential cmake git curl

RUN echo -n "UBU2004" > /OSRELEASE.txt

RUN mkdir /tmp/halon
RUN cd /tmp/halon \
    && curl -O -L https://go.dev/dl/go1.17.6.linux-amd64.tar.gz \
    && echo "231654bbf2dab3d86c1619ce799e77b03d96f9b50770297c8f4dff8836fc8ca2 go1.17.6.linux-amd64.tar.gz" | sha256sum -c \
    && tar -C /usr/local -xzf go1.17.6.linux-amd64.tar.gz

RUN mkdir /licenses
RUN cp /usr/local/go/LICENSE /licenses/LICENSE.go

RUN rm -rf /tmp/halon

RUN apt-get install -y supervisor
COPY .devcontainer/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
CMD ["/usr/bin/supervisord"]