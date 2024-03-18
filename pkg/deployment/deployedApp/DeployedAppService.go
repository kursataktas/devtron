package deployedApp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	util5 "github.com/devtron-labs/common-lib/utils/k8s"
	bean2 "github.com/devtron-labs/devtron/api/bean"
	"github.com/devtron-labs/devtron/enterprise/pkg/deploymentWindow"
	"github.com/devtron-labs/devtron/internal/sql/models"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	"github.com/devtron-labs/devtron/pkg/cluster/repository"
	"github.com/devtron-labs/devtron/pkg/deployment/deployedApp/bean"
	"github.com/devtron-labs/devtron/pkg/deployment/trigger/devtronApps"
	bean3 "github.com/devtron-labs/devtron/pkg/deployment/trigger/devtronApps/bean"
	"github.com/devtron-labs/devtron/pkg/k8s"
	"go.uber.org/zap"
	"time"
)

type DeployedAppService interface {
	StopStartApp(ctx context.Context, stopRequest *bean.StopAppRequest) (int, error)
	RotatePods(ctx context.Context, podRotateRequest *bean.PodRotateRequest) (*k8s.RotatePodResponse, error)
}

type DeployedAppServiceImpl struct {
	logger                  *zap.SugaredLogger
	k8sCommonService        k8s.K8sCommonService
	cdTriggerService        devtronApps.TriggerService
	envRepository           repository.EnvironmentRepository
	pipelineRepository      pipelineConfig.PipelineRepository
	cdWorkflowRepository    pipelineConfig.CdWorkflowRepository
	deploymentWindowService deploymentWindow.DeploymentWindowService
}

func NewDeployedAppServiceImpl(logger *zap.SugaredLogger,
	k8sCommonService k8s.K8sCommonService,
	cdTriggerService devtronApps.TriggerService,
	envRepository repository.EnvironmentRepository,
	pipelineRepository pipelineConfig.PipelineRepository,
	cdWorkflowRepository pipelineConfig.CdWorkflowRepository,
	deploymentWindowService deploymentWindow.DeploymentWindowService,
) *DeployedAppServiceImpl {
	return &DeployedAppServiceImpl{
		logger:                  logger,
		k8sCommonService:        k8sCommonService,
		cdTriggerService:        cdTriggerService,
		envRepository:           envRepository,
		pipelineRepository:      pipelineRepository,
		cdWorkflowRepository:    cdWorkflowRepository,
		deploymentWindowService: deploymentWindowService,
	}
}

func (impl *DeployedAppServiceImpl) StopStartApp(ctx context.Context, stopRequest *bean.StopAppRequest) (int, error) {
	pipelines, err := impl.pipelineRepository.FindActiveByAppIdAndEnvironmentId(stopRequest.AppId, stopRequest.EnvironmentId)
	if err != nil {
		impl.logger.Errorw("error in fetching pipeline", "app", stopRequest.AppId, "env", stopRequest.EnvironmentId, "err", err)
		return 0, err
	}
	if len(pipelines) == 0 {
		return 0, fmt.Errorf("no pipeline found")
	}
	pipeline := pipelines[0]

	// find pipeline with default
	var pipelineIds []int
	for _, p := range pipelines {
		impl.logger.Debugw("adding pipelineId", "pipelineId", p.Id)
		pipelineIds = append(pipelineIds, p.Id)
		// FIXME
	}
	wf, err := impl.cdWorkflowRepository.FindLatestCdWorkflowByPipelineId(pipelineIds)
	if err != nil {
		impl.logger.Errorw("error in fetching latest release", "err", err)
		return 0, err
	}
	stopTemplate := `{"replicaCount":0,"autoscaling":{"MinReplicas":0,"MaxReplicas":0 ,"enabled": false} }`
	latestArtifactId := wf.CiArtifactId
	cdPipelineId := pipeline.Id
	if pipeline.ApprovalNodeConfigured() {
		return 0, errors.New("application deployment requiring approval cannot be hibernated")
	}
	overrideRequest := &bean2.ValuesOverrideRequest{
		PipelineId:     cdPipelineId,
		AppId:          stopRequest.AppId,
		CiArtifactId:   latestArtifactId,
		UserId:         stopRequest.UserId,
		CdWorkflowType: bean2.CD_WORKFLOW_TYPE_DEPLOY,
	}
	if stopRequest.RequestType == bean.STOP {
		overrideRequest.AdditionalOverride = json.RawMessage([]byte(stopTemplate))
		overrideRequest.DeploymentType = models.DEPLOYMENTTYPE_STOP
	} else if stopRequest.RequestType == bean.START {
		overrideRequest.DeploymentType = models.DEPLOYMENTTYPE_START
	} else {
		return 0, fmt.Errorf("unsupported operation %s", stopRequest.RequestType)
	}
	triggerContext := bean3.TriggerContext{
		Context: ctx,
	}
	id, _, err := impl.cdTriggerService.ManualCdTrigger(triggerContext, overrideRequest)
	if err != nil {
		impl.logger.Errorw("error in stopping app", "err", err, "appId", stopRequest.AppId, "envId", stopRequest.EnvironmentId)
		return 0, err
	}
	return id, err
}

func (impl *DeployedAppServiceImpl) checkForDeploymentWindow(podRotateRequest *bean.PodRotateRequest) (*bean.PodRotateRequest, error) {
	actionState, envstate, err := impl.deploymentWindowService.GetStateForAppEnv(time.Now(), podRotateRequest.AppId, podRotateRequest.EnvironmentId, podRotateRequest.UserId)
	if err != nil {
		return podRotateRequest, fmt.Errorf("error in getting deployment window state %v", err)
	}
	if !actionState.IsActionAllowedWithBypass() {
		return podRotateRequest, deploymentWindow.GetActionBlockedError(actionState.GetBypassActionMessageForProfileAndState(envstate))
	}
	return podRotateRequest, nil
}

func (impl *DeployedAppServiceImpl) RotatePods(ctx context.Context, podRotateRequest *bean.PodRotateRequest) (*k8s.RotatePodResponse, error) {

	impl.logger.Infow("rotate pod request", "payload", podRotateRequest)
	podRotateRequest, err := impl.checkForDeploymentWindow(podRotateRequest)
	if err != nil {
		return nil, err
	}

	//extract cluster id and namespace from env id
	environmentId := podRotateRequest.EnvironmentId
	environment, err := impl.envRepository.FindById(environmentId)
	if err != nil {
		impl.logger.Errorw("error occurred while fetching env details", "envId", environmentId, "err", err)
		return nil, err
	}
	var resourceIdentifiers []util5.ResourceIdentifier
	for _, resourceIdentifier := range podRotateRequest.ResourceIdentifiers {
		resourceIdentifier.Namespace = environment.Namespace
		resourceIdentifiers = append(resourceIdentifiers, resourceIdentifier)
	}
	rotatePodRequest := &k8s.RotatePodRequest{
		ClusterId: environment.ClusterId,
		Resources: resourceIdentifiers,
	}
	response, err := impl.k8sCommonService.RotatePods(ctx, rotatePodRequest)
	if err != nil {
		return nil, err
	}
	//TODO KB: make entry in cd workflow runner
	return response, nil
}
