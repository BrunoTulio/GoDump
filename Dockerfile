FROM golang:1.21.0
ENV POSTGRES_VERSION=15

WORKDIR /go/src
ENV PATH="${PATH}:/go/bin"

RUN apt-get update     \
  && apt-get install --no-install-recommends -y \
  lsb-release \
  ca-certificates \
  curl \
  &&  sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list' \
  && wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
  && apt-get update \
  && apt-get install --no-install-recommends -y  \
  postgresql-client-15  \
  postgresql-client-14  \
  postgresql-client-13  \
  postgresql-client-12  \
  postgresql-client-11  \
  postgresql-client-10  \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*




CMD [ "tail", "-f", "/dev/null" ]