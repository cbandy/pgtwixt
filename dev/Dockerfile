ARG DOCKER_RELEASE

FROM docker.io/library/golang:${DOCKER_RELEASE} AS golang
FROM docker.io/library/python:${DOCKER_RELEASE} AS python

FROM docker.io/library/buildpack-deps:${DOCKER_RELEASE}

COPY --from=golang /usr/local /usr/local
COPY --from=python /usr/local /usr/local

# https://wiki.postgresql.org/wiki/Apt
# https://www.postgresql.org/support/versioning
ARG DOCKER_RELEASE
RUN curl -sS https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
 && echo > /etc/apt/sources.list.d/pgdg.list "deb http://apt.postgresql.org/pub/repos/apt/ ${DOCKER_RELEASE}-pgdg main" \
 && apt-get update \
 && apt-get install --no-install-recommends --yes postgresql-common \
 && sed -i 's/^#\{,1\}\(create_main_cluster\).*/\1 = false/' /etc/postgresql-common/createcluster.conf \
 && apt-get install --no-install-recommends --yes \
    bash-completion less sudo \
    $(for v in 9.4 9.5 9.6 10 11; do echo "postgresql-${v}"; done)

ARG DOCKER_GROUP
ARG DOCKER_USER
RUN groupadd --gid ${DOCKER_GROUP} dev \
 && useradd  --gid ${DOCKER_GROUP} --uid ${DOCKER_USER} --create-home --groups postgres,sudo --no-log-init dev \
 && echo > /etc/sudoers.d/dev "dev ALL=(ALL) NOPASSWD:ALL"

COPY ./entry /opt/entry
ENTRYPOINT ["/opt/entry"]
CMD ["bash"]

USER dev
WORKDIR /home/dev/code

COPY --chown=dev ./pip.conf ../.config/pip/pip.conf
