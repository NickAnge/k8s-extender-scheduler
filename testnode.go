package main

//Testnode is just a test struct
type Testnode struct {
	name       string
	latencies  map[string]int
	DeviceType string
}

//Testnode is just a test function
func (node *Testnode) FillValues(NodeName string, offset int) {
	node.name = NodeName

	node.latencies = map[string]int{
		"HOME-EDGE":   64 + offset,
		"EKETA-EDGE":  72 - offset,
		"EKETA-CLOUD": 100 + offset,
	}
}
