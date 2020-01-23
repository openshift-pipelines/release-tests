FROM ubuntu


#install golang 

RUN apt-get update
RUN apt-get install -y curl && apt-get install -q -y \
    gnupg2 \
    git \
    vim \
    openssh-client \
    curl \
    wget
RUN rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.13.3

RUN curl -sSL https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz \
		| tar -v -C /usr/local -xz

ENV PATH /usr/local/go/bin:$PATH
RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

# install oc
RUN wget https://github.com/openshift/origin/releases/download/v3.11.0/openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit.tar.gz && \
    tar xvzf openshift*.tar.gz && \
    cd openshift-origin-client-tools*/ && \
    mv  oc kubectl  /usr/local/bin/ &&  \
    oc version

# Install gauge

RUN apt-key adv --keyserver hkp://pool.sks-keyservers.net --recv-keys 023EDB0B || \
    apt-key adv --keyserver hkp://pool.sks-keyservers.net --recv-keys 023EDB0B || \
    apt-key adv --keyserver hkp://pool.sks-keyservers.net --recv-keys 023EDB0B && \
    echo deb https://dl.bintray.com/gauge/gauge-deb stable main | tee -a /etc/apt/sources.list

RUN apt-get update && apt-get install gauge && apt-get install -y build-essential

# Install go gauge plugins
RUN gauge install go && \
    gauge install html-report && \
    gauge install screenshot && \
    gauge config check_updates false && \
    gauge telemetry off

ENV GO111MODULE on
ENV CGO_ENABLED 0


WORKDIR /go/src/github.com/openshift-pipelines/release-tests


COPY . .
