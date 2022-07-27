package k8s

import (
	"context"
	"fmt"

	"github.com/laupse/kubegraph/application/entity"
	"github.com/laupse/kubegraph/utils"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"k8s.io/client-go/tools/clientcmd"
)

type K8sRepository struct {
	clientset *dynamic.Interface
}

func NewK8sRepository() (*K8sRepository, error) {
	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig"))
	if err != nil {
		return nil, err
	}

	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8sRepository{clientset: &clientset}, nil
}

func getRessources(dynamic dynamic.Interface, schema schema.GroupVersionResource, ns, selector string) (*unstructured.UnstructuredList, error) {
	options := metav1.ListOptions{
		LabelSelector: selector,
	}

	resources, err := dynamic.Resource(schema).Namespace(ns).List(context.TODO(), options)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func extractFieldsFromRessource(resource *unstructured.Unstructured, ressourceType string, addClusterAsOwner bool, arcs entity.Arcs) (node entity.Node, edges []entity.Edge) {
	metadata, _, _ := unstructured.NestedFieldNoCopy(resource.Object, "metadata")
	name := metadata.(map[string]interface{})["name"].(string)
	uuid := metadata.(map[string]interface{})["uid"].(string)
	node = entity.NewNode(uuid, name, ressourceType, arcs)
	if addClusterAsOwner {
		edges = append(edges, entity.Edge{
			Id:     fmt.Sprintf("%s:%s", uuid, "CLUSTER"),
			Target: uuid,
			Source: "CLUSTER",
		})
	}
	ownerReference, ok, _ := unstructured.NestedFieldNoCopy(resource.Object, "metadata", "ownerReferences")
	if ok {
		for _, owner := range ownerReference.([]interface{}) {
			ownerUuid := owner.(map[string]interface{})["uid"].(string)
			edges = append(edges, entity.Edge{
				Id:     fmt.Sprintf("%s:%s", uuid, ownerUuid),
				Target: uuid,
				Source: ownerUuid,
			})
		}
	}
	return
}

func computePodArc(state string, isready bool) (arcs entity.Arcs) {
	arcs.Blue = 0
	if state != "Pending" && state != "Running" {
		arcs.Red = 1
		return
	}
	if isready {
		arcs.Green = 1
		return
	}
	arcs.Orange = 1
	return

}

func (k8s *K8sRepository) GetPods(ns, selector string) (nodes []entity.Node, edges []entity.Edge, err error) {
	var podsGroupVersionResource = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	pods, err := getRessources(*k8s.clientset, podsGroupVersionResource, ns, selector)
	if err != nil {
		return nil, nil, err
	}
	for _, pod := range pods.Items {
		conditions, _, _ := unstructured.NestedFieldNoCopy(pod.Object, "status", "conditions")
		phase, _, _ := unstructured.NestedFieldNoCopy(pod.Object, "status", "phase")
		podIsReady := true
		for _, condition := range conditions.([]interface{}) {
			if condition.(map[string]interface{})["status"] != "True" {
				podIsReady = false
				break
			}
		}

		arcs := computePodArc(phase.(string), podIsReady)
		node, currentEdges := extractFieldsFromRessource(&pod, "POD", false, arcs)
		nodes = append(nodes, node)
		edges = append(edges, currentEdges...)
	}
	return nodes, edges, nil
}

func computeReplicasetArc(replicas, readyReplicas, currentReplicas int64) (arcs entity.Arcs) {
	arcs.Blue = 0
	if replicas == 0 {
		arcs.Blue = 1
		return
	}
	ready := utils.RoundFloat(float64(readyReplicas)/float64(replicas), 2)
	notready := utils.RoundFloat(float64((currentReplicas-readyReplicas))/float64(replicas), 2)

	arcs.Green = ready
	arcs.Orange = notready
	arcs.Red = utils.RoundFloat(1-(ready+notready), 2)
	return
}

func (k8s *K8sRepository) GetReplicasets(ns, selector string) (nodes []entity.Node, edges []entity.Edge, err error) {
	var rssGroupVersionResource = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"}
	rss, err := getRessources(*k8s.clientset, rssGroupVersionResource, ns, selector)
	if err != nil {
		return nil, nil, err
	}

	for _, rs := range rss.Items {
		replicas, _, _ := unstructured.NestedFieldNoCopy(rs.Object, "spec", "replicas")
		availableReplicas, _, _ := unstructured.NestedFieldNoCopy(rs.Object, "status", "availableReplicas")
		readyReplicas, _, _ := unstructured.NestedFieldNoCopy(rs.Object, "status", "readyReplicas")
		arcs := computeReplicasetArc(replicas.(int64), availableReplicas.(int64), readyReplicas.(int64))
		node, currentEdges := extractFieldsFromRessource(&rs, "RS", false, arcs)
		nodes = append(nodes, node)
		edges = append(edges, currentEdges...)
	}

	return nodes, edges, nil
}

func (k8s *K8sRepository) GetDeployments(ns, selector string) (nodes []entity.Node, edges []entity.Edge, err error) {
	var deploysGroupVersionResource = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	rss, err := getRessources(*k8s.clientset, deploysGroupVersionResource, ns, selector)
	if err != nil {
		return nil, nil, err
	}

	for _, rs := range rss.Items {
		node, currentEdges := extractFieldsFromRessource(&rs, "DEPLOY", true, entity.Arcs{Blue: 1})
		nodes = append(nodes, node)
		edges = append(edges, currentEdges...)
	}

	return nodes, edges, nil
}

func (k8s *K8sRepository) GetDaemonSets(ns, selector string) (nodes []entity.Node, edges []entity.Edge, err error) {
	var daemonsetsGroupVersionResource = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}
	rss, err := getRessources(*k8s.clientset, daemonsetsGroupVersionResource, ns, selector)
	if err != nil {
		return nil, nil, err
	}

	for _, rs := range rss.Items {
		node, currentEdges := extractFieldsFromRessource(&rs, "DS", true, entity.Arcs{Blue: 1})
		nodes = append(nodes, node)
		edges = append(edges, currentEdges...)
	}

	return nodes, edges, nil
}

func (k8s *K8sRepository) GetStatefulSets(ns, selector string) (nodes []entity.Node, edges []entity.Edge, err error) {
	var deploysGroupVersionResource = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}
	rss, err := getRessources(*k8s.clientset, deploysGroupVersionResource, ns, selector)
	if err != nil {
		return nil, nil, err
	}

	for _, rs := range rss.Items {
		node, currentEdges := extractFieldsFromRessource(&rs, "STS", true, entity.Arcs{Blue: 1})
		nodes = append(nodes, node)
		edges = append(edges, currentEdges...)
	}

	return nodes, edges, nil
}
