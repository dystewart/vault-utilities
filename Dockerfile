# BUilds on and updates the opf-toolbox container and adds go functionality

FROM quay.io/operate-first/opf-toolbox:v0.8.0

# Install Go
RUN curl -O https://dl.google.com/go/go1.17.1.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.17.1.linux-amd64.tar.gz \
    && rm go1.17.1.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:$PATH"

# Install any necessary build tools or dependencies
RUN yum update -y
