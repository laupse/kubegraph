from requests import get


def get_data():   
    return get(url="http://localhost:30000/api/graph/data").json()

def get_fields():   
    return get(url="http://localhost:30000/api/graph/fields").json()

class TestClassInstance:
    data = get_data()
    fields = get_fields()

    def test_length(self):
        assert len(self.data["edges"]) + 1 == len(self.data["nodes"])

    def test_target_and_source_are_node(self):
        nodes_id = [x["id"] for x in self.data["nodes"]]
        for edge in self.data["edges"]: 
            assert edge["target"] in nodes_id
            assert edge["source"] in nodes_id

    def test_nodes_has_all_fields(self):
        nodes_fields_name = sorted([x["field_name"] for x in self.fields["nodes_fields"]])
        print(self.fields["nodes_fields"])
        print(nodes_fields_name)
        for node in self.data["nodes"]:
            keys = sorted(list(node.keys()))
            print(keys)
            for i,field in enumerate(keys):
                assert field == nodes_fields_name[i]
    
    def test_edges_has_all_fields(self):
        edges_fields_name = sorted([x["field_name"] for x in self.fields["edges_fields"]])
        for edge in self.data["edges"]:
            keys = sorted(list(edge.keys()))
            for i,field in enumerate(keys):
                assert field == edges_fields_name[i]


