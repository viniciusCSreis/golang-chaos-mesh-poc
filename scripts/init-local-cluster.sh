#!/bin/bash

#Create cluster

kind delete cluster
kind create cluster

CURRENT_CONTEXT=$(kubectl config current-context)

if [ "$CURRENT_CONTEXT" != 'kind-kind' ]; then
  echo "CURRENT_CONTEXT: $CURRENT_CONTEXT != kind-kind"
  exit 1
fi

#Instal chaos-mech

curl -sSL https://mirrors.chaos-mesh.org/v0.9.1/install.sh | bash -s -- --local kind

# Install nats
kubectl apply -f manifests/nats/00-prereqs.yaml
kubectl apply -f manifests/nats/10-deployment.yaml

NATS_OPERATOR_STATUS="wait"
while [  "$NATS_OPERATOR_STATUS" == 'wait'  ]; do
  echo "waiting NATS_OPERATOR_STATUS"
  kubectl get NatsCluster 2> /dev/null && NATS_OPERATOR_STATUS="done"
  sleep 1
done

kubectl apply -f manifests/nats/20-nats-cluster.yaml
kubectl apply -f manifests/nats/30-nats-service-role.yaml