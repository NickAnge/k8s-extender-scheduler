package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

func checkBody(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		http.Error(res, "Please send a request body", 400)
		return
	}
}

func PredicateRoute(predicate Predicate) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		checkBody(res, req)

		var buf bytes.Buffer
		body := io.TeeReader(req.Body, &buf)

		//		log.Print("info: ", predicate.Name, " ExtenderArgs = " buf.String())
		log.Print("info: ", predicate.Name, ": Inside Route Handler")
		var extenderArgs schedulerapi.ExtenderArgs
		var extenderFilterResult *schedulerapi.ExtenderFilterResult

		if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
			extenderFilterResult = &schedulerapi.ExtenderFilterResult{
				Nodes:       nil,
				FailedNodes: nil,
				Error:       err.Error(),
			}
		} else {
			log.Print("info: ", "Call of predicate Handler")
			extenderFilterResult = predicate.Handler(extenderArgs)
		}

		if resultBody, err := json.Marshal(extenderFilterResult); err != nil {
			panic(err)
		} else {
			//log.Print("info: ", predicate.Name, " extenderFilterResult = ", string(resultBody))
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusOK)
			res.Write(resultBody)
		}
	}
}

func PrioritizeRoute(prioritize Prioritize) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		checkBody(w, r)

		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)
		log.Print("info: ", prioritize.Name, " ExtenderArgs = ", buf.String())

		var extenderArgs schedulerapi.ExtenderArgs
		var hostPriorityList *schedulerapi.HostPriorityList

		if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
			panic(err)
		}

		if list, err := prioritize.Handler(extenderArgs); err != nil {
			panic(err)
		} else {
			hostPriorityList = list
		}

		if resultBody, err := json.Marshal(hostPriorityList); err != nil {
			panic(err)
		} else {
			log.Print("info: ", prioritize.Name, " hostPriorityList = ", string(resultBody))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resultBody)
		}
	}
}

//
//func BindRoute(bind Bind) httprouter.Handle {
//	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//		checkBody(w, r)
//
//		var buf bytes.Buffer
//		body := io.TeeReader(r.Body, &buf)
//		log.Print("info: extenderBindingArgs = ", buf.String())
//
//		var extenderBindingArgs schedulerapi.ExtenderBindingArgs
//		var extenderBindingResult *schedulerapi.ExtenderBindingResult
//
//		if err := json.NewDecoder(body).Decode(&extenderBindingArgs); err != nil {
//			extenderBindingResult = &schedulerapi.ExtenderBindingResult{
//				Error: err.Error(),
//			}
//		} else {
//			extenderBindingResult = bind.Handler(extenderBindingArgs)
//		}
//
//		if resultBody, err := json.Marshal(extenderBindingResult); err != nil {
//			panic(err)
//		} else {
//			log.Print("info: extenderBindingResult = ", string(resultBody))
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(http.StatusOK)
//			w.Write(resultBody)
//		}
//	}
//}
//
//func PreemptionRoute(preemption Preemption) httprouter.Handle {
//	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//		checkBody(w, r)
//
//		var buf bytes.Buffer
//		body := io.TeeReader(r.Body, &buf)
//		log.Print("info: extenderPreemptionArgs = ", buf.String())
//
//		var extenderPreemptionArgs schedulerapi.ExtenderPreemptionArgs
//		var extenderPreemptionResult *schedulerapi.ExtenderPreemptionResult
//
//		if err := json.NewDecoder(body).Decode(&extenderPreemptionArgs); err != nil {
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(http.StatusBadRequest)
//		} else {
//			extenderPreemptionResult = preemption.Handler(extenderPreemptionArgs)
//		}
//
//		if resultBody, err := json.Marshal(extenderPreemptionResult); err != nil {
//			panic(err)
//		} else {
//			log.Print("info: extenderPreemptionResult = ", string(resultBody))
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(http.StatusOK)
//			w.Write(resultBody)
//		}
//	}
//}

func VersionRoute(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Fprint(res, fmt.Sprint(version))
}

func AddVersion(router *httprouter.Router) {
	router.GET(versionPath, DebugLogging(VersionRoute, versionPath))
}

func DebugLogging(h httprouter.Handle, path string) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		log.Print("debug: ", path, " request body = ", req.Body)
		h(res, req, params)
		log.Print("debug: ", path, " response=", res)
	}
}

func AddPredicateRoute(router *httprouter.Router, predicate Predicate) {
	path := predicatesPrefix + "/" + predicate.Name
	router.POST(path, DebugLogging(PredicateRoute(predicate), path))
}

func AddPrioritize(router *httprouter.Router, prioritize Prioritize) {
	path := prioritiesPrefix + "/" + prioritize.Name
	router.POST(path, DebugLogging(PrioritizeRoute(prioritize), path))
}

//
//func AddBind(router *httprouter.Router, bind Bind) {
//	if handle, _, _ := router.Lookup("POST", bindPath); handle != nil {
//		log.Print("warning: AddBind was called more then once!")
//	} else {
//		router.POST(bindPath, DebugLogging(BindRoute(bind), bindPath))
//	}
//}
//
//func AddPreemption(router *httprouter.Router, preemption Preemption) {
//	if handle, _, _ := router.Lookup("POST", preemptionPath); handle != nil {
//		log.Print("warning: AddPreemption was called more then once!")
//	} else {
//		router.POST(preemptionPath, DebugLogging(PreemptionRoute(preemption), preemptionPath))
//	}
//}
