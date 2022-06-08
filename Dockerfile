FROM golang:latest
LABEL maintainer="Raphael Ebner"
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./
RUN go build


FROM debian:latest
LABEL maintainer="Raphael Ebner"
WORKDIR /app

ENV CRON_TIME="* * * * *"
ENV CERT_OWNER_ID="0"
ENV CERT_GROUP_ID="0"
ENV ON_START=1

COPY --from=0 /app/traefik-ssl-certificate-exporter ./
RUN apt-get update && apt-get install -y cron
COPY entrypoint.sh ./
RUN chmod +x ./entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]
