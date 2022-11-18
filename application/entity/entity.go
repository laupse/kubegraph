package entity

import (
	"fmt"
	"math/big"

	"encoding/json"
)

type Edge struct {
	Id     string `json:"id"`
	Target string `json:"target"`
	Source string `json:"source"`
}

type Node struct {
	Id        string   `json:"id"`
	Title     string   `json:"title"`
	MainStat  string   `json:"mainstat"`
	ArcGreen  Rational `json:"arc__green" color:"green"`
	ArcRed    Rational `json:"arc__red" color:"red"`
	ArcOrange Rational `json:"arc__orange" color:"orange"`
	ArcBlue   Rational `json:"arc__blue" color:"blue"`
}

type Arcs struct {
	Green  Rational `json:"green"`
	Red    Rational `json:"red"`
	Orange Rational `json:"orange"`
	Blue   Rational `json:"blue"`
}

type GraphData struct {
	Edges []Edge `json:"edges"`
	Nodes []Node `json:"nodes"`
}

type Rational struct {
	big.Rat
}

func NewRational(a, b int64) Rational {

	return Rational{
		*big.NewRat(a, b),
	}
}

func (a *Arcs) MarshalJSON() ([]byte, error) {
	fmt.Println("lol")
	return json.Marshal(&struct {
		Green  string `json:"green"`
		Red    string `json:"red"`
		Orange string `json:"orange"`
		Blue   string `json:"blue"`
	}{
		Green:  a.Green.FloatString(2),
		Red:    a.Red.FloatString(2),
		Orange: a.Orange.FloatString(2),
		Blue:   a.Blue.FloatString(2),
	})
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
