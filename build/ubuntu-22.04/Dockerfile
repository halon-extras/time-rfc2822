FROM --platform=linux/amd64 ubuntu:22.04
LABEL org.opencontainers.image.authors="support@halon.io"

ARG HALON_REPO_USER
ARG HALON_REPO_PASS

RUN apt-get update

RUN apt-get install -y wget gnupg
RUN apt-get install -y apt-transport-https

RUN wget -qO - https://raw.githubusercontent.com/halon/changelog/master/pgp-keys/7F0A73B5.asc | apt-key add -
RUN echo "deb https://repo.halon.io/ jammy stable" >> /etc/apt/sources.list.d/halon.list
RUN echo "machine repo.halon.io login ${HALON_REPO_USER} password ${HALON_REPO_PASS}" >> /etc/apt/auth.conf
RUN apt-get update && apt-get install -y halon=5.10.0

RUN apt-get install -y build-essential cmake git curl file

RUN echo -n "UBU2204" > /OSRELEASE.txt

RUN mkdir /tmp/halon
RUN cd /tmp/halon \
    && curl -O -L https://go.dev/dl/go1.25.5.linux-amd64.tar.gz \
    && echo "9e9b755d63b36acf30c12a9a3fc379243714c1c6d3dd72861da637f336ebb35b go1.25.5.linux-amd64.tar.gz" | sha256sum -c \
    && tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz
RUN /usr/local/go/bin/go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@latest
ENV PATH="$PATH:/usr/local/go/bin:/root/go/bin"

RUN mkdir /licenses
RUN cp /usr/local/go/LICENSE /licenses/LICENSE.go

RUN rm -rf /tmp/halon

COPY build.sh /build.sh
CMD ["/build.sh"]
