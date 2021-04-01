package main

import (
	"errors"
	"log"

	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

//Predicate is the struct with name and func
type Predicate struct {
	Name string
	Func func(pod v1.Pod, node v1.Node) (bool, error)
}

func (p Predicate) Handler(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	pod := args.Pod
	canSchedule := make([]v1.Node, 0, len(args.Nodes.Items))
	canNotSchedule := make(map[string]string)

	log.Print("info:", pod.Labels)

	//	nodes := args.Nodes.Ite
	//our function  selectNode(pod, nodes)
	for _, node := range args.Nodes.Items {
		result, err := p.Func(*pod, node)
		if err != nil {
			canNotSchedule[node.Name] = err.Error()
		} else {
			if result {
				canSchedule = append(canSchedule, node)
			}
		}
	}

	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: canSchedule,
		},
		FailedNodes: canNotSchedule,
		Error:       "",
	}

	return &result
}

func getOurLocation(pod v1.Pod) (string, error) {

	if value, ok := pod.Labels["OurLocation"]; ok {
		return value, nil
	}

	return "", errors.New("No value OurLocation")
}

func getTypeOfComponent(pod v1.Pod) (string, error) {

	if value, found := pod.Labels["typeOfComponent"]; found {
		return value, nil
	}
	return "", errors.New("No value typeOfComponent")
}

///func selectNode(pod v1.Pod, nodes []v1.Node) v1.Node {
//	return nil
//}
