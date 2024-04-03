package bean

import (
	"context"
	"github.com/devtron-labs/devtron/api/bean"
	"github.com/devtron-labs/devtron/enterprise/pkg/deploymentWindow"
	"github.com/devtron-labs/devtron/enterprise/pkg/resourceFilter"
	"github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	"github.com/devtron-labs/devtron/pkg/resourceQualifiers"
	"time"
)

const (
	ARGOCD_SYNC_ERROR = "error in syncing argoCD app"
)

type TriggerEvent struct {
	PerformChartPush           bool
	PerformDeploymentOnCluster bool
	GetManifestInResponse      bool
	DeploymentAppType          string
	ManifestStorageType        string
	TriggeredBy                int32
	TriggerdAt                 time.Time
}

type TriggerRequest struct {
	CdWf                   *pipelineConfig.CdWorkflow
	Pipeline               *pipelineConfig.Pipeline
	Artifact               *repository.CiArtifact
	ApplyAuth              bool
	TriggeredBy            int32
	RefCdWorkflowRunnerId  int
	RunStageInEnvNamespace string
	WorkflowType           bean.WorkflowType
	TriggerMessage         string
	DeploymentWindowState  *deploymentWindow.EnvironmentState
	TriggerContext
}

type TriggerContext struct {
	// Context is a context object to be passed to the pipeline trigger
	// +optional
	Context context.Context
	// ReferenceId is a unique identifier for the workflow runner
	// refer pipelineConfig.CdWorkflowRunner
	ReferenceId *string

	// manual or automatic
	TriggerType TriggerType
}

type TriggerType int

const (
	Automatic TriggerType = 1
	Manual    TriggerType = 2
)

func (context TriggerContext) IsAutoTrigger() bool {
	return context.TriggerType == Automatic
}

func (context TriggerContext) ToTriggerTypeString() string {
	if context.IsAutoTrigger() {
		return "AUTO"
	}
	return "MANUAL"
}

type DeploymentType = string

const (
	Helm                    DeploymentType = "helm"
	ArgoCd                  DeploymentType = "argo_cd"
	ManifestDownload        DeploymentType = "manifest_download"
	GitOpsWithoutDeployment DeploymentType = "git_ops_without_deployment"
	ManifestPush            DeploymentType = "manifest_push"
)

type TriggerRequirementRequestDto struct {
	Scope          resourceQualifiers.Scope
	TriggerRequest TriggerRequest
	Stage          resourceFilter.ReferenceType
}

type TriggerFeasibilityResponse struct {
	ApprovalRequestId int
	TriggerRequest    TriggerRequest
	FilterIdVsState   map[int]resourceFilter.FilterState
	Filters           []*resourceFilter.FilterMetaDataBean
}
