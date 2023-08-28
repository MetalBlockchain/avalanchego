# syntax=docker/dockerfile:experimental

# This Dockerfile is meant to be used with the build_local_dep_image.sh script
# in order to build an image using the local version of coreth

# Changes to the minimum golang version must also be replicated in
# scripts/build_avalanche.sh
# scripts/local.Dockerfile (here)
# Dockerfile
# README.md
# go.mod
FROM golang:1.19.12-bullseye

RUN mkdir -p /go/src/github.com/ava-labs

WORKDIR $GOPATH/src/github.com/ava-labs
COPY metalgo metalgo

WORKDIR $GOPATH/src/github.com/ava-labs/metalgo
RUN ./scripts/build_avalanche.sh

RUN ln -sv $GOPATH/src/github.com/ava-labs/avalanche-byzantine/ /metalgo
