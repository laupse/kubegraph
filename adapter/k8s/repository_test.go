package k8s

import "testing"

func TestComputeReplicasetArcAllready(t *testing.T) {
	arcs := computeReplicasetArc(2, 2, 2)
	if arcs.Blue != 0 {
		t.Errorf("Arc Blue got %f; attendu 0", arcs.Blue)
	}
	if arcs.Green != 1 {
		t.Errorf("Arc Green got %f; attendu 1", arcs.Green)
	}
	if arcs.Red != 0 {
		t.Errorf("Arc Red got %f; attendu 0", arcs.Red)
	}
	if arcs.Orange != 0.4 {
		t.Errorf("Arc Orange got %f; attendu 0", arcs.Orange)
	}
}

func TestComputeReplicasetArcNotAllready(t *testing.T) {
	arcs := computeReplicasetArc(5, 3, 5)
	if arcs.Blue != 0 {
		t.Errorf("Arc Blue got %f; attendu 0", arcs.Blue)
	}
	if arcs.Green != 0.6 {
		t.Errorf("Arc Green got %f; attendu 0.6", arcs.Green)
	}
	if arcs.Red != 0 {
		t.Errorf("Arc Red got %f; attendu 0", arcs.Red)
	}
	if arcs.Orange != 0.4 {
		t.Errorf("Arc Orange got %f; attendu 0.4", arcs.Orange)
	}
}

func TestComputeReplicasetArcNotAllreadyWithMissing(t *testing.T) {
	arcs := computeReplicasetArc(5, 3, 4)
	if arcs.Blue != 0 {
		t.Errorf("Arc Blue got %f; attendu 0", arcs.Blue)
	}
	if arcs.Green != 0.6 {
		t.Errorf("Arc Green got %f; attendu 0.6", arcs.Green)
	}
	if arcs.Red != 0.2 {
		t.Errorf("Arc Red got %f; attendu 0.2", arcs.Red)
	}
	if arcs.Orange != 0.2 {
		t.Errorf("Arc Orange got %f; attendu 0.2", arcs.Orange)
	}

}

func TestComputeReplicasetArcZeroReplica(t *testing.T) {
	arcs := computeReplicasetArc(0, 0, 0)
	if arcs.Blue != 1 {
		t.Errorf("Arc Blue got %f; attendu 1", arcs.Blue)
	}
	if arcs.Green != 0 {
		t.Errorf("Arc Green got %f; attendu 0", arcs.Green)
	}
	if arcs.Red != 0 {
		t.Errorf("Arc Red got %f; attendu 0", arcs.Red)
	}
	if arcs.Orange != 0 {
		t.Errorf("Arc Orange got %f; attendu 0", arcs.Orange)
	}

}
