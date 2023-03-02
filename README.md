# Stardog Userrole Operator

[![Build](https://img.shields.io/github/workflow/status/vshn/stardog-userrole-operator/Build)][build]
![Go version](https://img.shields.io/github/go-mod/go-version/vshn/stardog-userrole-operator)
![Kubernetes version](https://img.shields.io/badge/k8s-v1.18-blue)
![Stardog version](https://img.shields.io/badge/Stardog-v7-blue)
[![Version](https://img.shields.io/github/v/release/vshn/stardog-userrole-operator)][releases]
[![GitHub downloads](https://img.shields.io/github/downloads/vshn/stardog-userrole-operator/total)][releases]
[![Docker image](https://img.shields.io/docker/pulls/vshn/stardog-userrole-operator)][dockerhub]
[![License](https://img.shields.io/github/license/vshn/stardog-userrole-operator)][license]

A Kubernetes operator to manage Stardog users and roles.

## Generating the REST client

The package stardogrest is a REST client generated by [autorest](http://azure.github.io/autorest/) based on the [stardogrest/stardog_swagger.yaml](stardogrest/stardog_swagger.yaml) file. If the stardog REST API changes, the [stardogrest/stardog_swagger.yaml](stardogrest/stardog_swagger.yaml) should be updated to reflect the changes, and then autorest should be run again with the following command:

```
make autorest
```

[build]: https://github.com/vshn/stardog-userrole-operator/actions?query=workflow%3ABuild
[releases]: https://github.com/vshn/stardog-userrole-operator/releases
[license]: https://github.com/vshn/stardog-userrole-operator/blob/master/LICENSE
[dockerhub]: https://hub.docker.com/r/vshn/stardog-userrole-operator

## Local Development

### Prerequisites

- minikube
- helm
- kubectl
- Stardog image registry credentials
- Stardog license

### Init local dev environment

```
minikube start

# Create values.yaml (fill in actual image registry credentials)
cat <<EOF > values.yaml
image:
  username: "foo"
  password: "supersecret"

stardog:
  adminPassword: "admin"

replicaCount: 1

zookeeper:
  enabled: false
EOF

# Provide Stardog license key
kubectl create secret generic stardog-license --from-file stardog-license-key.bin=stardog-license-key.bin

# Install Stardog
helm repo add appuio https://charts.appuio.ch
helm install stardog appuio/stardog -f values.yaml
```
