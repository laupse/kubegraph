package graph

type Field struct {
	FieldName string `json:"field_name"`
	Type      string `json:"type"`
	Color     string `json:"color,omitempty"`
}

type Fields struct {
	EdgeFields []Field `json:"edges_fields"`
	NodeFields []Field `json:"nodes_fields"`
}

func GetFields() Fields {
	var fields Fields
	fields.EdgeFields = []Field{
		{FieldName: "id", Type: "string"},
		{FieldName: "source", Type: "string"},
		{FieldName: "target", Type: "string"},
	}
	fields.NodeFields = []Field{
		{FieldName: "id", Type: "string"},
		{FieldName: "title", Type: "string"},
		{FieldName: "mainstat", Type: "string"},
		{FieldName: "arc__ready", Type: "string", Color: "green"},
		{FieldName: "arc__not_ready", Type: "string", Color: "orange"},
		{FieldName: "arc__missing", Type: "string", Color: "red"},
		{FieldName: "mainstat", Type: "string"},
	}

	return fields
}
