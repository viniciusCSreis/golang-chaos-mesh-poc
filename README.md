# Go Micro-services

Test project with skaffold,kind,chaos mech and nats

# Prerequisites

- Docker
- Kubectl
- skaffold [install skaffold](https://skaffold.dev/docs/install/)
- Kind [Install Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

# Config Cluster

```bash
    ./scripts/init-local-cluster.sh
```

# Run micro services

wait nats cluster to be running. 
(if the nats is not running the pods will restart until the nats is running)

```bash
    ./scripts/run-local.sh
```

# Run experiments

```bash
    kubectl apply -f experimets/delay-nats.yaml
```

List authors is fast: `curl localhost:3000/authors`

Create authors is slow and 
return timeout error: `curl -i -X POST localhost:3000/authors`

# List running experiments

To list experiments you can access the chaos mesh ui:
 <http://localhost:2333>

# Delete experiments

```bash
    kubectl delete -f experimets/delay-nats.yaml
```

now create is fast `curl -i -X POST localhost:3000/authors`


# Delete cluster

```bash
    kind delete cluster
```