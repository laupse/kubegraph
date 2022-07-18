# Kubegraph

This repository contains an "exporter" that will expose Kubernetes ressources (pod, replicaset, deployment, statefulset, daemonset) as [Grafana Node Graph](https://grafana.com/docs/grafana/latest/visualizations/node-graph/) data source using this [Grafana Plugin](https://grafana.com/grafana/plugins/hamedkarbasi93-nodegraphapi-datasource/)

This utilitary can help visualizing how this ressources interact and link with  other in a dynamic way. 

Great for demonstrating Kubernetes capabilities like rollout.

# Usage

## Deployment
You need a Kubernetes cluster (it could be something like kind or minkube ofc) with Grafana deploy on it. 

Grafana also needs this [Grafana Plugin](https://grafana.com/grafana/plugins/hamedkarbasi93-nodegraphapi-datasource/) to be installed

Then you can deploy this using :

``` bash
kubectl apply -f manifests/deploy.yml
```

Do not forget to add the Node Graph Api using the service dns created by the kubectl command above

## Querying

You can query this datasource directly inside a node panel. It will then show kubernetes ressources in a graph mode. 

You can reduce the scope by adding this query parameter : 

* `ns` for querying specific namespace 
* `selector` for querying specific resources with this label

Example: `ns=proxy&selector=app=nginx` will make the data source return resources in `proxy` namesapce with `app=nginx` labels 
