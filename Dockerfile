FROM debian:stretch-slim
ARG SOURCE_LOCATION=./build
EXPOSE 8888

# install dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    apt-transport-https \
    curl \
    ca-certificates \
    && apt-get clean \
    && apt-get autoremove \ 
    && rm -rf /var/lib/apt/lists/*

# create new user
RUN useradd --create-home drophere

# create new directory
RUN mkdir -p /home/drophere/drophere-service

# specify directory
WORKDIR /home/drophere/drophere-service
COPY ${SOURCE_LOCATION} .

# change owner to user "drophere"
RUN chown -R drophere:drophere .

USER drophere
RUN chmod +x drophere-service

CMD ["./drophere-service"]
