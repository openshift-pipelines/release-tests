# Using base image which has requirements (go, gauge) installed
FROM quay.io/openshift-pipeline/ci

# Set WORKDIR to /root/release-tests and copy the tests
RUN mkdir /root/release-tests
WORKDIR /root/release-tests
COPY . .

CMD ["/bin/bash"]