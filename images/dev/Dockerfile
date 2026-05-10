FROM golang:1.26.3

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
    rm -rf /var/lib/apt/lists/*

RUN curl -o /tmp/python.tar.xz https://www.python.org/ftp/python/3.14.4/Python-3.14.4.tar.xz && \
    cd /tmp && tar -xvf python.tar.xz && cd Python-3.14.4 && \
    ./configure --enable-optimizations --enable-shared

RUN cd /tmp/Python-3.14.4 && make install

RUN ldconfig && \
    rm -rf /var/lib/apt/lists/*
