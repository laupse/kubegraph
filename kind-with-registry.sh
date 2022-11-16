#!/bin/sh
set -o errexit

if [ "$(docker inspect -f '{{.State.Running}}' "kubegraph-control-plane" 2>/dev/null || true)" != 'true' ]; then
  # create a cluster with the local registry enabled in containerd
  cat <<EOF |  kind create cluster --name kubegraph --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    # port forward 80 on the host to 80 on this node
    extraPortMappings:
      - containerPort: 30000
        hostPort: 30000
containerdConfigPatches:
  - |-
    [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
      endpoint = ["http://${reg_name}:5000"]
EOF
fi

# create registry container unless it already exists
reg_name='kind-registry'
reg_port='5001'
if [ "$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
  docker run -d --restart=always \
    -p "127.0.0.1:${reg_port}:5000" \
    --name "${reg_name}" \
    --network kind \
    registry:2
fi

proxy_name='proxy-registry'
if [ "$(docker inspect -f '{{.State.Running}}' "${proxy_name}" 2>/dev/null || true)" != 'true' ]; then
  docker run -d --restart=always \
    -p "127.0.0.1:5002:5002" \
    --name "${proxy_name}" \
    --network kind \
    --volume ${PWD}/Caddyfile:/etc/caddy/Caddyfile \
    caddy:2 
fi

docker cp "${proxy_name}":/data/caddy//pki/authorities/local/root.crt .

dagger_name='dagger-buildkitd'
if [ "$(docker inspect -f '{{.State.Running}}' "${dagger_name}" 2>/dev/null || true)" != 'true' ]; then
  docker run -d --restart=always \
    -v $PWD/root.crt:/etc/ssl/certs/root.crt \
    --name "${dagger_name}" \
    --network kind \
    --privileged moby/buildkit:v0.10.5
fi

kind export kubeconfig --name kubegraph --kubeconfig kind-ci.yaml

# Document the local registry
# https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF



