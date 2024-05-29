/*
 * Copyright (c) 2020-2024. Devtron Inc.
 */

package bean

import (
	"encoding/json"
	"github.com/devtron-labs/devtron/internal/sql/models"
	"github.com/devtron-labs/devtron/pkg/pipeline/repository"
)

type WorkflowType string

type DeploymentConfigurationType string

const (
	CD_WORKFLOW_TYPE_PRE              WorkflowType                = "PRE"
	CD_WORKFLOW_TYPE_POST             WorkflowType                = "POST"
	CD_WORKFLOW_TYPE_DEPLOY           WorkflowType                = "DEPLOY"
	CI_WORKFLOW_TYPE                  WorkflowType                = "CI"
	WEBHOOK_WORKFLOW_TYPE             WorkflowType                = "WEBHOOK"
	DEPLOYMENT_CONFIG_TYPE_LAST_SAVED DeploymentConfigurationType = "LAST_SAVED_CONFIG"
	//latest trigger is not being used because this is being handled at FE and we anyhow identify latest trigger as
	//last deployed wfr which is also a specific trigger
	DEPLOYMENT_CONFIG_TYPE_LATEST_TRIGGER   DeploymentConfigurationType = "LATEST_TRIGGER_CONFIG"
	DEPLOYMENT_CONFIG_TYPE_SPECIFIC_TRIGGER DeploymentConfigurationType = "SPECIFIC_TRIGGER_CONFIG"
)

func (workflowType WorkflowType) WorkflowTypeToStageType() repository.PipelineStageType {
	switch workflowType {
	case CD_WORKFLOW_TYPE_PRE:
		return repository.PIPELINE_STAGE_TYPE_PRE_CD
	case CD_WORKFLOW_TYPE_POST:
		return repository.PIPELINE_STAGE_TYPE_POST_CD
	default:
		return ""
	}
}

func (workflowType WorkflowType) GetDeploymentStageType() models.DeploymentType {
	switch workflowType {
	case CD_WORKFLOW_TYPE_PRE:
		return models.DEPLOYMENTTYPE_PRE
	case CD_WORKFLOW_TYPE_POST:
		return models.DEPLOYMENTTYPE_POST
	case CD_WORKFLOW_TYPE_DEPLOY:
		return models.DEPLOYMENTTYPE_DEPLOY
	default:
		return models.DEPLOYMENTTYPE_UNKNOWN
	}
}

type ValuesOverrideRequest struct {
	PipelineId                            int                         `json:"pipelineId" validate:"required"`
	AppId                                 int                         `json:"appId" validate:"required"`
	CiArtifactId                          int                         `json:"ciArtifactId" validate:"required"`
	AdditionalOverride                    json.RawMessage             `json:"additionalOverride,omitempty"`
	TargetDbVersion                       int                         `json:"targetDbVersion"`
	ForceTrigger                          bool                        `json:"forceTrigger,notnull"`
	DeploymentTemplate                    string                      `json:"strategy,omitempty"` // validate:"oneof=BLUE-GREEN ROLLING"`
	DeploymentWithConfig                  DeploymentConfigurationType `json:"deploymentWithConfig"`
	WfrIdForDeploymentWithSpecificTrigger int                         `json:"wfrIdForDeploymentWithSpecificTrigger"`
	CdWorkflowType                        WorkflowType                `json:"cdWorkflowType,notnull"`
	WfrId                                 int                         `json:"wfrId,notnull"`
	CdWorkflowId                          int                         `json:"cdWorkflowId"`
	PipelineOverrideId                    int                         `json:"pipelineOverrideId"` //required for async install/upgrade event;
	DeploymentType                        models.DeploymentType       `json:"deploymentType"`     //required for async install/upgrade handling; previously if was used internally
	CdWorkflowRunnerId                    int                         `json:"cdWorkflowRunnerId"` // used for deployment triggered through release, will use existing cd workflow runner if CdWorkflowRunnerId is found.
	UserId                                int32                       `json:"-"`
	EnvId                                 int                         `json:"-"`
	EnvName                               string                      `json:"-"`
	ClusterId                             int                         `json:"-"`
	AppName                               string                      `json:"-"`
	PipelineName                          string                      `json:"-"`
	DeploymentAppType                     string                      `json:"-"`
	ProjectId                             int                         `json:"-"`
	IsProdEnv                             bool                        `json:"-"`
	Image                                 string                      `json:"-"`
	TriggerMetadata                       string                      `json:"-"`
}

func (v ValuesOverrideRequest) IsDeployDeploymentType() bool {
	return v.DeploymentType == models.DEPLOYMENTTYPE_DEPLOY
}
func (v ValuesOverrideRequest) IsUnknownDeploymentType() bool {
	return v.DeploymentType == models.DEPLOYMENTTYPE_UNKNOWN
}

type BulkCdDeployEvent struct {
	ValuesOverrideRequest *ValuesOverrideRequest `json:"valuesOverrideRequest"`
	UserId                int32                  `json:"userId"`
}

type ReleaseStatusUpdateRequest struct {
	RequestId string             `json:"requestId"`
	NewStatus models.ChartStatus `json:"newStatus"`
}

func CheckIfDeploymentTypePrePostOrDeployOrUnknown(deploymentType models.DeploymentType) bool {
	if deploymentType == models.DEPLOYMENTTYPE_DEPLOY || deploymentType == models.DEPLOYMENTTYPE_UNKNOWN || deploymentType == models.DEPLOYMENTTYPE_PRE || deploymentType == models.DEPLOYMENTTYPE_POST {
		return true
	}
	return false
}

type CdWorkflowRunnerIds struct {
	PreCdWorkflowRunnerId  int
	CdWorkflowRunnerId     int
	PostCdWorkflowRunnerId int
}
