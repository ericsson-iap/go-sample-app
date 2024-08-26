FROM armdocker.rnd.ericsson.se/dockerhub-ericsson-remote/rockylinux:8.6

ARG USER_ID=60577
ARG USER_NAME="eric-sdk"

COPY target/eric-oss-hello-world-go-app /

# Enable GRPC logging
ENV GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info GOCACHE=/tmp/.gocache

ARG APP_VERSION
LABEL \
    adp.app.version=$APP_VERSION

RUN echo "$USER_ID:x:$USER_ID:0:An Identity for $USER_NAME:/nonexistent:/bin/false" >>/etc/passwd
RUN echo "$USER_ID:!::0:::::" >>/etc/shadow

USER $USER_ID

CMD ["/eric-oss-hello-world-go-app"]
