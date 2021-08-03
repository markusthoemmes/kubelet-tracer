# Kubelet Tracer

As the name suggests, `kubelet-tracer` is a little tool that analyzes the logs of the kubelet for a specific pod to be able to reason about what exactly is going on in the kubelet for that respective pod. The tool has a focus on the timing of the kubelet's actions to be able to potentially improve the kubelet's pod/container bringup time and to reason about performance gaps.

[![asciicast](https://asciinema.org/a/I0hZ3sPCpJPW2oPYqmAZx8L1l.svg)](https://asciinema.org/a/I0hZ3sPCpJPW2oPYqmAZx8L1l)

## `kind` "integration"

Since my workflow is currently centered around `kind`, there's a few helpers here that make iterating on Kubernetes changes a breeze.

### `rekind.sh`

This script builds a node image from a fresh Kubernetes source (located at `~/go/src/k8s.io/kubernetes`) and deploys a `kind` cluster with it. The `kind` configuration in `kind.yaml` makes sure, that the newly built image is used and that the kubelet logging is configured accordingly (JSON logging and log level raised to the max). The script is idempotent (hence the "re" prefix).

### `analyze_kind.sh`

This script takes the pod you want to analyze as an argument (in the form of `$ns/$name`) and then takes care of exporting the respective logs from `kind`, piping them into `kubelet-tracer` and cleaning up after itself.