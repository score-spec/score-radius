#!/bin/bash
set -o errexit

## Kind cluster
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF

## Radius
rad install kubernetes
kubectl wait deployments/controller -n radius-system --for condition=Available --timeout=90s