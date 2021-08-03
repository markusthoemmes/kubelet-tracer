# Kubelet Tracer

As the name suggests, `kubelet-tracer` is a little tool that analyzes the logs of the kubelet for a specific pod to be able to reason about what exactly is going on in the kubelet for that respective pod. The tool has a focus on the timing of the kubelet's actions to be able to potentially improve the kubelet's pod/container bringup time and to reason about performance gaps.

[![asciicast](https://asciinema.org/a/I0hZ3sPCpJPW2oPYqmAZx8L1l.svg)](https://asciinema.org/a/I0hZ3sPCpJPW2oPYqmAZx8L1l)