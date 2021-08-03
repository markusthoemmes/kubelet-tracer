#!/bin/bash

folder=$(kind export logs 2> /dev/null)
go run main.go --stop-after-deletion --pod "$1" < "$folder/kind-control-plane/kubelet.log"
rm -rf folder