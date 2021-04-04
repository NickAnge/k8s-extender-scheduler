package main

import (
	"errors"
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

type Prioritize struct {
	Name string
	Func func(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error)
}

func (p Prioritize) Handler(args schedulerapi.ExtenderArgs) (*schedulerapi.HostPriorityList, error) {
	return p.Func(*args.Pod, args.Nodes.Items)
}

//OurLocation pod label
func GetOurLocation(pod v1.Pod) (string, error) {

	if value, ok := pod.Labels["OurLocation"]; ok {
		return value, nil
	}

	return "", errors.New("No value OurLocation")
}

//typeOfComponent pod label
func GetTypeOfComponent(pod v1.Pod) (string, error) {

	if value, found := pod.Labels["typeOfComponent"]; found {
		return value, nil
	}
	return "", errors.New("No value typeOfComponent")
}

func GetLatency(node Testnode, Location string) float64 {
	return float64(node.latencies[Location])
}
