# Prometheus tips

Query number of pods of type StatefulSet in namespace prometheus:

```
count(kube_pod_info{namespace="prometheus", created_by_kind="StatefulSet"})
```
