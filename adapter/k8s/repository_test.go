package k8s

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestComputeReplicasetArcAllready(t *testing.T) {
	arcs := computeReplicasetArc(2, 2, 2)
	if arcs.Blue.Num().Int64() != 0 {
		t.Errorf("Arc Blue got %d; attendu 0", arcs.Blue.Num().Int64())
	}
	if arcs.Green.Num().Int64() != 1 {
		t.Errorf("Arc Green got %d; attendu 1", arcs.Green.Num().Int64())
	}
	if arcs.Red.Num().Int64() != 0 {
		t.Errorf("Arc Red got %d; attendu 0", arcs.Red.Num().Int64())
	}
	if arcs.Orange.Num().Int64() != 0 {
		t.Errorf("Arc Orange got %d; attendu 0", arcs.Orange.Num().Int64())
	}
	content, err := json.Marshal(arcs)
	fmt.Println(arcs)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(arcs.Blue.FloatString(2))
	fmt.Println(string(content))
}

func TestComputeReplicasetArcNotAllready(t *testing.T) {
	arcs := computeReplicasetArc(5, 3, 5)
	if arcs.Blue.Num().Int64() != 0 {
		t.Errorf("Arc Blue got %d; attendu 0", arcs.Blue.Num().Int64())
	}
	if arcs.Green.Num().Int64() != 3 {
		t.Errorf("Arc Green got %d; attendu 3", arcs.Green.Num().Int64())
	}
	if arcs.Red.Num().Int64() != 0 {
		t.Errorf("Arc Red got %d; attendu 0", arcs.Red.Num().Int64())
	}
	if arcs.Orange.Num().Int64() != 2 {
		t.Errorf("Arc Orange got %d; attendu 2", arcs.Orange.Num().Int64())
	}
}

func TestComputeReplicasetArcNotAllreadyWithMissing(t *testing.T) {
	arcs := computeReplicasetArc(5, 3, 4)
	if arcs.Blue.Num().Int64() != 0 {
		t.Errorf("Arc Blue got %d; attendu 0", arcs.Blue.Num().Int64())
	}
	if arcs.Green.Num().Int64() != 3 {
		t.Errorf("Arc Green got %d; attendu 3", arcs.Green.Num().Int64())
	}
	if arcs.Red.Num().Int64() != 1 {
		t.Errorf("Arc Red got %d; attendu 1", arcs.Red.Num().Int64())
	}
	if arcs.Orange.Num().Int64() != 1 {
		t.Errorf("Arc Orange got %d; attendu 1", arcs.Orange.Num().Int64())
	}
	content, err := json.Marshal(arcs)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(content))

}

func TestComputeReplicasetArcZeroReplica(t *testing.T) {
	arcs := computeReplicasetArc(0, 0, 0)
	if arcs.Blue.Num().Int64() != 1 {
		t.Errorf("Arc Blue got %d; attendu 1", arcs.Blue.Num().Int64())
	}
	if arcs.Green.Num().Int64() != 0 {
		t.Errorf("Arc Green got %d; attendu 0", arcs.Green.Num().Int64())
	}
	if arcs.Red.Num().Int64() != 0 {
		t.Errorf("Arc Red got %d; attendu 0", arcs.Red.Num().Int64())
	}
	if arcs.Orange.Num().Int64() != 0 {
		t.Errorf("Arc Orange got %d; attendu 0", arcs.Orange.Num().Int64())
	}

}
