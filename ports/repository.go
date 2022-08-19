package ports

import "github.com/laupse/kubegraph/application/entity"

type GraphRepository interface {
	GetPods(ns string, selector string) ([]entity.Node, []entity.Edge, error)
	GetReplicasets(ns string, selector string) ([]entity.Node, []entity.Edge, error)
	GetDeployments(ns string, selector string) ([]entity.Node, []entity.Edge, error)
	GetDaemonSets(ns string, selector string) ([]entity.Node, []entity.Edge, error)
	GetStatefulSets(ns string, selector string) ([]entity.Node, []entity.Edge, error)
	GetJobs(ns string, selector string) ([]entity.Node, []entity.Edge, error)
}
