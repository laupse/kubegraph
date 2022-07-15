package utils

import (
	"math"
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	newVal := math.Round(val*ratio) / ratio
	return newVal
}

func ComputeReplicasetArc(replicas, readyReplicas, currentReplicas int32) (ready, notready, missing float64) {
	if replicas == 0 {
		return 0, 1, 0
	}
	ready = roundFloat(float64(readyReplicas)/float64(replicas), 2)
	notready = roundFloat(float64((currentReplicas-readyReplicas))/float64(replicas), 2)
	missing = roundFloat(1-(ready+notready), 2)
	return
}

func ComputePodArc(state string, isready bool) (ready, running, failed float64) {
	if state != "Pending" && state != "Running" {
		return 0, 0, 1
	}
	if isready {
		return 1, 0, 0
	}
	return 0, 1, 0

}
