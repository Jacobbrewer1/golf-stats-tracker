FROM golang:1.21.1

LABEL org.opencontainers.image.source='https://github.com/Jacobbrewer1/golf-stats-tracker'
LABEL org.opencontainers.image.description="A simple golf stats tracker that allows you to track your golf scores and stats."
LABEL org.opencontainers.image.licenses='GNU General Public License v3.0'

WORKDIR /cmd

# Copy the binary from the build
COPY ./bin/app /cmd/app

RUN ["chmod", "+x", "./app"]

ENTRYPOINT ["/cmd/app"]
