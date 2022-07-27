package entity

type Edge struct {
	Id     string `json:"id"`
	Target string `json:"target"`
	Source string `json:"source"`
}

type Node struct {
	Id        string  `json:"id"`
	Title     string  `json:"title"`
	MainStat  string  `json:"mainstat"`
	ArcGreen  float64 `json:"arc__green" color:"green"`
	ArcRed    float64 `json:"arc__red" color:"red"`
	ArcOrange float64 `json:"arc__orange" color:"orange"`
	ArcBlue   float64 `json:"arc__blue" color:"blue"`
}

type Arcs struct {
	Green  float64
	Red    float64
	Orange float64
	Blue   float64
}

type GraphData struct {
	Edges []Edge `json:"edges"`
	Nodes []Node `json:"nodes"`
}

func NewNode(id, name, mainStat string, arcs Arcs) Node {
	return Node{
		Id:        id,
		Title:     name,
		MainStat:  mainStat,
		ArcGreen:  arcs.Green,
		ArcBlue:   arcs.Blue,
		ArcRed:    arcs.Red,
		ArcOrange: arcs.Orange,
	}

}
