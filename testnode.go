package main

import (
	"fmt"
	"k8s.io/api/core/v1"
	"strconv"
)

//Testnode is just a test struct
type Testnode struct {
	name       string
	latencies  map[string]float64
	DeviceType string
}

//Testnode is just a test function
//func (node *Testnode) FillValues(NodeName string, offset int) {
//	node.name = NodeName
//
//	node.latencies = map[string]float64{
//		"HOME-EDGE":   64 + offset,
//		"EKETA-EDGE":  72 - offset,
//		"EKETA-CLOUD": 100 + offset,
//	}
//}

//TestNode for adding Latencies from the Labels
func (node *Testnode) FillLatencies(nodeInfo v1.Node) {

	homeEdgeLat, _ := strconv.ParseFloat(nodeInfo.Labels["HOME-EDGE"], 64)
	eketaCloudLat, _ := strconv.ParseFloat(nodeInfo.Labels["EKETA-CLOUD"], 64)
	eketaEdgeLat, _ := strconv.ParseFloat(nodeInfo.Labels["EKETA-EDGE"], 64)

	fmt.Println(nodeInfo.Labels)
	fmt.Println(homeEdgeLat)
	fmt.Println(eketaEdgeLat)
	fmt.Println(eketaCloudLat)

	node.latencies = map[string]float64{
		"HOME-EDGE":   homeEdgeLat,
		"EKETA-EDGE":  eketaEdgeLat,
		"EKETA-CLOUD": eketaCloudLat,
	}

}
