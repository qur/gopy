FROM golang:1.20

RUN apt-get update && \
    apt-get install -y \
      build-essential \
      zlib1g-dev \
      libncurses5-dev \
      libgdbm-dev \
      libnss3-dev \
      libssl-dev \
      libreadline-dev \
      libffi-dev && \
    curl -o /tmp/python.tar.xz https://www.python.org/ftp/python/3.11.2/Python-3.11.2.tar.xz && \
    cd /tmp && tar -xvf python.tar.xz && cd Python-3.11.2 && \
    ./configure --enable-optimizations && \
    make install && \
    rm -rf /var/lib/apt/lists/*
