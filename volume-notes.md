# Notes on the volume subsystem:

## Leading questions:
- Why 17 -> 77ms. What happens? --> Internal retry of 100ms
- Why multiple waits for the PollImmediate call and not just one?

## Example logs:

```
17	volumemanager/volume_manager.go:394	Waiting for volumes to attach and mount for pod
77	populator/desired_state_of_world_populator.go:349	Added volume to desired state
177	reconciler/reconciler.go:224	operationExecutor.VerifyControllerAttachedVolume started for volume "kube-api-access-4vhx2" (UniqueName: "kubernetes.io/projected/443ca847-0495-4d5d-ab92-06dd3016df3f-kube-api-access-4vhx2") pod "basic-1c906d43-6d4d-47a2-bad6-ec44f63a00d1" (UID: "443ca847-0495-4d5d-ab92-06dd3016df3f") 
278	reconciler/reconciler.go:254	Starting operationExecutor.MountVolume for volume "kube-api-access-4vhx2" (UniqueName: "kubernetes.io/projected/443ca847-0495-4d5d-ab92-06dd3016df3f-kube-api-access-4vhx2") pod "basic-1c906d43-6d4d-47a2-bad6-ec44f63a00d1" (UID: "443ca847-0495-4d5d-ab92-06dd3016df3f") 
278	reconciler/reconciler.go:269	operationExecutor.MountVolume started for volume "kube-api-access-4vhx2" (UniqueName: "kubernetes.io/projected/443ca847-0495-4d5d-ab92-06dd3016df3f-kube-api-access-4vhx2") pod "basic-1c906d43-6d4d-47a2-bad6-ec44f63a00d1" (UID: "443ca847-0495-4d5d-ab92-06dd3016df3f") 
318	volumemanager/volume_manager.go:425	All volumes are attached and mounted for pod
```

## Findings:

Internal retry: 100ms (desired state of the world populator)
External check: 300ms (WaitForAttachAndMount)

Mount operations etc. are executed in somewhat like a "thread pool", which is polled on a 100ms timer.