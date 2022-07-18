package graph

import (
	"context"

	"github.com/laupse/kubegraph/utils"
	"github.com/spf13/viper"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Edge struct {
	Id     string `json:"id"`
	Target string `json:"target"`
	Source string `json:"source"`
}

type Node struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	MainStat    string  `json:"mainstat"`
	ArcReady    float64 `json:"arc__ready"`
	ArcNotReady float64 `json:"arc__not_ready"`
	ArcMissing  float64 `json:"arc__missing"`
}

type GraphData struct {
	Edges []Edge `json:"edges"`
	Nodes []Node `json:"nodes"`
}

func newSimpleNode(id, name, mainStat string) Node {
	return Node{
		Id:          id,
		Title:       name,
		MainStat:    mainStat,
		ArcReady:    1,
		ArcNotReady: 0,
		ArcMissing:  0,
	}

}

func podsToNodesAndEdges(graphData *GraphData, pods v1.PodList) {
	for _, pod := range pods.Items {
		podIsReady := true
		for _, condition := range pod.Status.Conditions {
			if condition.Status != "True" {
				podIsReady = false
				break
			}
		}
		ready, notready, missing := utils.ComputePodArc(string(pod.Status.Phase), podIsReady)
		graphData.Nodes = append(graphData.Nodes, Node{
			Id:          string(pod.UID),
			Title:       pod.Name,
			MainStat:    "POD",
			ArcReady:    ready,
			ArcNotReady: notready,
			ArcMissing:  missing,
		})
		for _, owner := range pod.OwnerReferences {
			graphData.Edges = append(graphData.Edges, Edge{
				Id:     string(uuid.NewUUID()),
				Target: string(pod.UID),
				Source: string(owner.UID),
			})
		}
	}
}

func GetDataFromNamespace(ns string, selector string) (*GraphData, error) {
	graphData := &GraphData{}
	options := metav1.ListOptions{
		LabelSelector: selector,
	}

	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig"))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), options)
	if err != nil {
		return nil, err
	}

	podsToNodesAndEdges(graphData, *pods)

	rss, err := clientset.AppsV1().ReplicaSets(ns).List(context.TODO(), options)
	if err != nil {
		return nil, err
	}

	for _, rs := range rss.Items {
		readyReplica, notReadyReplica, missingReplica := utils.ComputeReplicasetArc(*rs.Spec.Replicas, rs.Status.AvailableReplicas, rs.Status.Replicas)
		graphData.Nodes = append(graphData.Nodes, Node{
			Id:          string(rs.UID),
			Title:       rs.Name,
			MainStat:    "RS",
			ArcReady:    readyReplica,
			ArcNotReady: notReadyReplica,
			ArcMissing:  missingReplica,
		})
		for _, owner := range rs.OwnerReferences {
			graphData.Edges = append(graphData.Edges, Edge{
				Id:     string(uuid.NewUUID()),
				Target: string(rs.UID),
				Source: string(owner.UID),
			})
		}
	}

	deploys, err := clientset.AppsV1().Deployments(ns).List(context.TODO(), options)
	if err != nil {
		return nil, err
	}

	for _, deploy := range deploys.Items {
		graphData.Nodes = append(graphData.Nodes, newSimpleNode(string(deploy.UID), deploy.Name, "DEPLOY"))
	}

	dss, err := clientset.AppsV1().DaemonSets(ns).List(context.TODO(), options)
	if err != nil {
		return nil, err
	}

	for _, ds := range dss.Items {
		graphData.Nodes = append(graphData.Nodes, newSimpleNode(string(ds.UID), ds.Name, "DS"))
	}

	stss, err := clientset.AppsV1().StatefulSets(ns).List(context.TODO(), options)
	if err != nil {
		return nil, err
	}

	for _, sts := range stss.Items {
		graphData.Nodes = append(graphData.Nodes, newSimpleNode(string(sts.UID), sts.Name, "STS"))
	}

	return graphData, nil
}
