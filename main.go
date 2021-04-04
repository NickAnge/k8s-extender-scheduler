package main

import (
	//	"fmt"
	"errors"
	"log"
	"math"

	"net/http"
	"os"
	"strings"

	"github.com/comail/colog"
	"github.com/julienschmidt/httprouter"

	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

const (
	versionPath      = "/version"
	apiPrefix        = "/scheduler"
	bindPath         = apiPrefix + "/bind"
	preemptionPath   = apiPrefix + "/preemption"
	predicatesPrefix = apiPrefix + "/predicates"
	prioritiesPrefix = apiPrefix + "/priorities"

	LabelTraining = "Training"
	LabelFetching = "Fetching"
	LabelPredict  = "Predict"

	LabelEketaEdge  = "EKETA-EDGE"
	LabelHomeEdge   = "HOME-EDGE"
	LabelEketaCloud = "EKETA-CLOUD"
)

var (
	version string // injected via ldflags at build time

	TruePredicate = Predicate{
		Name: "nikos_true",
		Func: func(pod v1.Pod, node v1.Node) (bool, error) {
			return true, nil
		},
	}

	ZeroPriority = Prioritize{
		Name: "nikos_priority",
		Func: func(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error) {
			var priorityList schedulerapi.HostPriorityList

			typeOfComponent, err := GetTypeOfComponent(pod)

			if err != nil {
				return &priorityList, err
			}
			if typeOfComponent == LabelTraining {
				return &priorityList, errors.New("Type of training, let the default scheduler do the priority")
			} else if typeOfComponent == LabelFetching || typeOfComponent == LabelPredict {
				ourLocation, err := GetOurLocation(pod)
				var arrayNodes []Testnode

				log.Print("info:", ourLocation)

				if err != nil {
					return &priorityList, err
				}
				//TEST PURPOSES-FILL NODE VALUES
				for i, node := range nodes {

					newNode := &Testnode{
						name: node.Name,
					}
					newNode.FillValues(node.Name, (2 + i))
					arrayNodes = append(arrayNodes, *newNode)
				}

				for _, node := range arrayNodes {
					log.Print("info:", node)
				}
				latencyMap := make(map[float64]Testnode)
				min := math.Inf(+1)

				for _, node := range arrayNodes {
					currentLatency := GetLatency(node, ourLocation)

					if currentLatency < min {
						min = currentLatency
					}
					latencyMap[currentLatency] = node
				}

				priorityList = make([]schedulerapi.HostPriority, 1)

				priorityList[0] = schedulerapi.HostPriority{
					Host:  latencyMap[min].name,
					Score: 1,
				}

				return &priorityList, nil
			} else {
				log.Print("info:", "scheduling left to DS")
				return nil, errors.New("Ds")
			}
		},
	}
	//	priorityList = make([]schedulerapi.HostPriority, len(nodes))

	//	priorityList[0] = schedulerapi.HostPriority{
	//		Host:  nodes[0].Name,
	//		Score: 91,
	//	}

	//	priorityList[1] = schedulerapi.HostPriority{
	//		Host:  nodes[1].Name,
	//		Score: 100,
	//	}
	//	//for _, node := range nodes {
	//	//	priorityList[0] = schedulerapi.HostPriority{
	//	//		Host:  node.Name,
	//	//		Score: 0,
	//	//			}
	//	//}
	//	return &priorityList, nil

)

func StringToLevel(levelStr string) colog.Level {
	switch level := strings.ToUpper(levelStr); level {
	case "TRACE":
		return colog.LTrace
	case "DEBUG":
		return colog.LDebug
	case "INFO":
		return colog.LInfo
	case "WARNING":
		return colog.LWarning
	case "ERROR":
		return colog.LError
	case "ALERT":
		return colog.LAlert
	default:
		log.Printf("warning: LOG_LEVEL=\"%s\" is empty or invalid, fallling back to \"INFO\".\n", level)
		return colog.LInfo
	}
}

func main() {
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()
	level := StringToLevel(os.Getenv("LOG_LEVEL"))
	log.Print("Log level was set to ", strings.ToUpper(level.String()))
	colog.SetMinLevel(level)

	router := httprouter.New()
	AddVersion(router)

	predicates := []Predicate{TruePredicate}

	//AddPredicateRoute(router, TruePredicate)
	for _, p := range predicates {
		log.Print("Trying to run predicate:", p.Name)
		AddPredicateRoute(router, p)
	}

	priorities := []Prioritize{ZeroPriority}
	for _, p := range priorities {
		AddPrioritize(router, p)
	}

	//AddBind(router, NoBind)

	log.Print("info: server starting on the port :80")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
