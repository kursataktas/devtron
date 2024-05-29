/*
 * Copyright (c) 2020-2024. Devtron Inc.
 */

package trigger

import (
	"github.com/devtron-labs/devtron/api/restHandler/app/pipeline/trigger"
	sse2 "github.com/devtron-labs/devtron/api/sse"
	"github.com/gorilla/mux"
	"github.com/juju/errors"
	_ "k8s.io/client-go/rest"
	"net/http"
	"strconv"
	"time"
)

var sse *sse2.SSE

type PipelineTriggerRouter interface {
	InitPipelineTriggerRouter(pipelineTriggerRouter *mux.Router)
}

func PollTopic(r *http.Request) (string, error) {
	parameters := mux.Vars(r)
	if parameters == nil {
		return "", errors.New("missing mandatory parameters")
	}
	name := parameters["name"]
	if name == "" {
		return "", errors.New("missing mandatory parameters")
	}
	return "/" + name, nil
}

func NewPipelineTriggerRouter(pipelineRestHandler trigger.PipelineTriggerRestHandler, sseChannel *sse2.SSE) *PipelineTriggerRouterImpl {
	routerImpl := &PipelineTriggerRouterImpl{restHandler: pipelineRestHandler}
	sse = sseChannel
	return routerImpl
}

type PipelineTriggerRouterImpl struct {
	restHandler trigger.PipelineTriggerRestHandler
}

func (router PipelineTriggerRouterImpl) InitPipelineTriggerRouter(pipelineTriggerRouter *mux.Router) {
	pipelineTriggerRouter.Path("/cd-pipeline/trigger").HandlerFunc(router.restHandler.OverrideConfig).Methods("POST")
	pipelineTriggerRouter.Path("/update-release-status").HandlerFunc(router.restHandler.ReleaseStatusUpdate).Methods("POST")
	pipelineTriggerRouter.Path("/rotate-pods").HandlerFunc(router.restHandler.RotatePods).Methods("POST")
	pipelineTriggerRouter.Path("/stop-start-app").HandlerFunc(router.restHandler.StartStopApp).Methods("POST")
	pipelineTriggerRouter.Path("/stop-start-dg").HandlerFunc(router.restHandler.StartStopDeploymentGroup).Methods("POST")
	pipelineTriggerRouter.Path("/release/").
		Handler(sse2.SubscribeHandler(sse.Broker, PollTopic, fetchReleaseData)).
		Methods("GET").
		Queries("name", "{name}")

	pipelineTriggerRouter.Path("/deployment-configuration/latest/saved/{appId}/{pipelineId}").HandlerFunc(router.restHandler.GetAllLatestDeploymentConfiguration).Methods("GET")
	pipelineTriggerRouter.Path("/manifest/download/{appId}/{envId}").Queries("runner", "{runner}").HandlerFunc(router.restHandler.DownloadManifest).Methods("GET")
	pipelineTriggerRouter.Path("/manifest/download/{appId}/{envId}/{cd_workflow_id}").HandlerFunc(router.restHandler.DownloadManifestForSpecificTrigger).Methods("GET")
}

func fetchReleaseData(r *http.Request, receive <-chan int, send chan<- int) {
	parameters := mux.Vars(r)
	name := parameters["name"]
	for i := 0; i <= 10; i++ {
		select {
		case <-receive:
			return
		default:
		}
		time.Sleep(1 * time.Second)
		data := []byte(time.Now().String() + "-" + strconv.Itoa(i))
		sse.OutboundChannel <- sse2.SSEMessage{"", data, "/" + name}
	}
	send <- 1
}
