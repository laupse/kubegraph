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
	idField := Field{
		FieldName: "id",
		Type:      "string",
	}
	titleField := Field{
		FieldName: "title",
		Type:      "string",
	}
	arcReadyField := Field{
		FieldName: "arc__ready",
		Type:      "number",
		Color:     "green",
	}
	arcNotReadyField := Field{
		FieldName: "arc__not_ready",
		Type:      "number",
		Color:     "orange",
	}
	arcMissingField := Field{
		FieldName: "arc__missing",
		Type:      "number",
		Color:     "red",
	}
	sourceField := Field{
		FieldName: "source",
		Type:      "string",
	}
	targetField := Field{
		FieldName: "target",
		Type:      "string",
	}

	fields.EdgeFields = append(fields.EdgeFields, idField)
	fields.EdgeFields = append(fields.EdgeFields, sourceField)
	fields.EdgeFields = append(fields.EdgeFields, targetField)

	fields.NodeFields = append(fields.NodeFields, idField)
	fields.NodeFields = append(fields.NodeFields, titleField)
	fields.NodeFields = append(fields.NodeFields, arcReadyField)
	fields.NodeFields = append(fields.NodeFields, arcNotReadyField)
	fields.NodeFields = append(fields.NodeFields, arcMissingField)

	return fields
}
