# ktop

CLI to inspect the current CPU usages of the pods within a certain Kubernetes namespace.
Use to inspect if the current load fits your requested limits.

# How to use?
```
go run cmd/ktop.go -n <namespace>
```

# Example Output
```
+--------+-----------+------------------+---------------+-------------+-------------+------------------+----------------+----------------+
| Node   | Namespace | Pod              | CPU (Request) | CPU (Limit) | CPU (Usage) | Memory (Request) | Memory (Limit) | Memory (Usage) |
+--------+-----------+------------------+---------------+-------------+-------------+------------------+----------------+----------------+
| node-1 | todo      | app-1            | 200m          | 2500m       | 26m         | 1152m            | 1536m          | 750055424000m  |
| node-2 | todo      | app-2            | 200m          | 2500m       | 23m         | 1152m            | 1536m          | 753393664000m  |
| node-3 | todo      | event-consumer-1 | 200m          | 2500m       | 39m         | 378m             | 1012m          | 204410880000m  |
| node-4 | todo      | event-consumer-2 | 200m          | 2500m       | 26m         | 378m             | 1012m          | 153833472000m  |
+--------+-----------+------------------+---------------+-------------+-------------+------------------+----------------+----------------+
```
