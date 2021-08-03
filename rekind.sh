#!/bin/bash

kind delete cluster
kind build node-image --image docker.io/markusthoemmes/kind:head --kube-root ~/go/src/k8s.io/kubernetes
kind create cluster --wait 30s --config kind.yaml
notify-send "Kind cluster done!"
