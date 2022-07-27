package services

import (
	"reflect"

	"github.com/laupse/kubegraph/application/entity"
	"github.com/laupse/kubegraph/ports"
)

type Field struct {
	FieldName string `json:"field_name"`
	Type      string `json:"type"`
	Color     string `json:"color,omitempty"`
}

type Fields struct {
	EdgeFields []Field `json:"edges_fields"`
	NodeFields []Field `json:"nodes_fields"`
}

type GraphService struct {
	repository ports.GraphRepository
}

func NewGraphService(dataRepository ports.GraphRepository) *GraphService {
	return &GraphService{
		repository: dataRepository,
	}

}

func (GraphService *GraphService) GetData(ns, selector string) (*entity.GraphData, error) {
	var err error
	var edges []entity.Edge
	var nodes []entity.Node
	graphData := &entity.GraphData{
		Nodes: []entity.Node{
			{Id: "CLUSTER", Title: "CLUSTER", ArcBlue: 1},
		},
	}
	graphData.Nodes = append(graphData.Nodes, nodes...)

	nodes, edges, err = GraphService.repository.GetPods(ns, selector)
	if err != nil {
		return nil, err
	}
	graphData.Nodes = append(graphData.Nodes, nodes...)
	graphData.Edges = append(graphData.Edges, edges...)

	nodes, edges, err = GraphService.repository.GetReplicasets(ns, selector)
	if err != nil {
		return nil, err
	}
	graphData.Nodes = append(graphData.Nodes, nodes...)
	graphData.Edges = append(graphData.Edges, edges...)

	nodes, edges, err = GraphService.repository.GetDeployments(ns, selector)
	if err != nil {
		return nil, err
	}
	graphData.Nodes = append(graphData.Nodes, nodes...)
	graphData.Edges = append(graphData.Edges, edges...)

	return graphData, err
}

func (GraphService *GraphService) GetFields() *Fields {
	fields := &Fields{}

	node := entity.Node{}
	n := reflect.TypeOf(node)
	for i := 0; i < n.NumField(); i++ {
		field := Field{
			FieldName: n.Field(i).Tag.Get("json"),
			Type:      "number",
		}
		color, ok := n.Field(i).Tag.Lookup("color")
		if ok {
			field.Color = color
		}
		if n.Field(i).Type.Name() == "string" {
			field.Type = "string"
		}
		fields.NodeFields = append(fields.NodeFields, field)
	}

	edge := entity.Edge{}
	e := reflect.TypeOf(edge)
	for i := 0; i < e.NumField(); i++ {
		field := Field{
			FieldName: e.Field(i).Tag.Get("json"),
			Type:      "string",
		}
		fields.EdgeFields = append(fields.EdgeFields, field)
	}

	return fields
}
