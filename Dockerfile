# Using base image which has requirements (go, gauge) installed
FROM quay.io/openshift-pipeline/ci

# Set this var to install gauge plugins at custom path
ENV GAUGE_HOME=/tmp

# Add timeout to ignore runner connection error
RUN gauge config runner_connection_timeout 600000 && \
    gauge config runner_request_timeout 300000

# Copy the tests into /tmp/release-tests
RUN mkdir /tmp/release-tests
WORKDIR /tmp/release-tests
COPY . .

# Set required permissions for OpenShift usage
RUN chgrp -R 0 /tmp && \
    chmod -R g=u /tmp

CMD ["/bin/bash"]