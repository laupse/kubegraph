package utils

import "testing"

func TestComputeReplicasetArcAllready(t *testing.T) {
	ready, notready, missing := ComputeReplicasetArc(2, 2, 2)
	if ready != 1 {
		t.Errorf("got %f; attendu 1", ready)
	}
	if notready != 0 {
		t.Errorf("got %f; attendu 0", ready)
	}
	if missing != 0 {
		t.Errorf("got %f; attendu 0", ready)
	}
}

func TestComputeReplicasetArcNotAllready(t *testing.T) {
	ready, notready, missing := ComputeReplicasetArc(5, 3, 5)
	if ready != 0.6 {
		t.Errorf("got %f; attendu 0.6", ready)
	}
	if notready != 0.4 {
		t.Errorf("got %f; attendu 0.4", ready)
	}
	if missing != 0 {
		t.Errorf("got %f; attendu 0", ready)
	}
}

func TestComputeReplicasetArcNotAllreadyWithMissing(t *testing.T) {
	ready, notready, missing := ComputeReplicasetArc(5, 3, 4)
	if ready != 0.6 {
		t.Errorf("got %f; attendu 0.6", ready)
	}
	if notready != 0.2 {
		t.Errorf("got %f; attendu 0.2", ready)
	}
	if missing != 0.2 {
		t.Errorf("got %f; attendu 0.2", ready)
	}
}

func TestComputeReplicasetArcZeroReplica(t *testing.T) {
	ready, notready, missing := ComputeReplicasetArc(0, 0, 0)
	if ready != 0 {
		t.Errorf("got %f; attendu 1", ready)
	}
	if notready != 1 {
		t.Errorf("got %f; attendu 0", ready)
	}
	if missing != 0 {
		t.Errorf("got %f; attendu 0", ready)
	}

}
