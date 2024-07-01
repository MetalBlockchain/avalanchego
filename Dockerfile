# Changes to the minimum golang version must also be replicated in
# scripts/build_avalanche.sh
# Dockerfile (here)
# README.md
# go.mod
# ============= Compilation Stage ================
FROM golang:1.21.9-bullseye AS builder

WORKDIR /build
# Copy and download avalanche dependencies using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build metalgo
ARG RACE_FLAG=""
RUN ./scripts/build.sh ${RACE_FLAG}

# ============= Cleanup Stage ================
FROM debian:11-slim AS execution

# Maintain compatibility with previous images
RUN mkdir -p /metalgo/build
WORKDIR /metalgo/build

# Copy the executables into the container
COPY --from=builder /build/build/ .

CMD [ "./metalgo" ]
