# This is based on the debezium server docker file from: https://github.com/debezium/docker-images/tree/main/server
#
# NOTE: If you want to run some of the connectors such as DB2 or Oracle, then you
# _must_ spin up your own image to include the drivers

FROM registry.access.redhat.com/ubi8/openjdk-11

LABEL maintainer="Debezium Community"

#
# Set the version, home directory, and MD5 hash.
#
ENV DEBEZIUM_VERSION=2.0.0-SNAPSHOT \
    SERVER_HOME=/debezium \
    MAVEN_REPO_CENTRAL="https://repo1.maven.org/maven2" \
    LD_LIBRARY_PATH="/instantclient_21_3"

#
# Create a directory for Debezium Server
#
USER root
RUN microdnf -y install gzip && \
    microdnf -y install unzip && \
    microdnf -y install curl && \
    microdnf -y install libaio && \
    microdnf clean all && \
    mkdir $SERVER_HOME && \
    chmod 755 $SERVER_HOME

#
# Change ownership and switch user
#
RUN chown -R jboss $SERVER_HOME && \
    chgrp -R jboss $SERVER_HOME
USER jboss

RUN mkdir $SERVER_HOME/conf && \
    mkdir $SERVER_HOME/data

#
# Download and install Debezium Server.
# This uses a local build, but you can pull from io/debezium/debezium-server-dist/$DEBEZIUM_VERSION/debezium-server-dist-$DEBEZIUM_VERSION.tar.gz
#
COPY debezium-server-dist-2.0.0-SNAPSHOT.tar.gz /tmp/debezium.tar.gz

#
# Verify the contents and then install ...
#
RUN tar xzf /tmp/debezium.tar.gz -C $SERVER_HOME --strip-components 1

USER root
RUN rm -f /tmp/debezium.tar.gz

# Include the Oracle specific components. I'm using an alternate version here
# which requires a custom build
COPY debezium-connector-oracle-2.0.0-SNAPSHOT.jar /debezium/lib/debezium-connector-oracle-2.0.0-SNAPSHOT.jar
RUN curl -o /tmp/instantclient-basic-linux.x64-21.3.0.0.0.zip -O https://download.oracle.com/otn_software/linux/instantclient/213000/instantclient-basic-linux.x64-21.3.0.0.0.zip

RUN unzip /tmp/instantclient-basic-linux.x64-21.3.0.0.0.zip -d / && \
    cp /instantclient_21_3/xstreams.jar /debezium/lib/xstreams.jar && \
    cp /instantclient_21_3/ojdbc8.jar /debezium/lib/ojdbc8.jar && \
    rm /tmp/instantclient-basic-linux.x64-21.3.0.0.0.zip

#
# Allow random UID to use Debezium Server
#
RUN chmod -R g+w,o+w $SERVER_HOME

# Set the working directory to the Debezium Server home directory
WORKDIR $SERVER_HOME

#
# Expose the ports and set up volumes for the data, transaction log, and configuration
#
EXPOSE 8080
VOLUME ["/debezium/conf","/debezium/data"]

CMD ["/debezium/run.sh"]
