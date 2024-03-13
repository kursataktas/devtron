/*
 * Copyright (c) 2020 Devtron Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package bean

import (
	"encoding/json"
	"github.com/devtron-labs/devtron/internal/sql/models"
	"github.com/devtron-labs/devtron/pkg/pipeline/repository"
	"github.com/devtron-labs/devtron/util"
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
}

type BulkCdDeployEvent struct {
	ValuesOverrideRequest *ValuesOverrideRequest `json:"valuesOverrideRequest"`
	UserId                int32                  `json:"userId"`
}

type ReleaseStatusUpdateRequest struct {
	RequestId string             `json:"requestId"`
	NewStatus models.ChartStatus `json:"newStatus"`
}

type ArtifactsListFilterOptions struct {
	//list filter data
	Limit        int
	Offset       int
	SearchString string
	Order        string

	//self stage data
	PipelineId int
	StageType  WorkflowType

	// CiPipelineId is id of ci-pipeline present in the same app-workflow of PipelineId
	CiPipelineId int

	//parent satge data
	ParentCdId      int
	ParentId        int
	ParentStageType WorkflowType

	//excludeArtifactIds
	ExcludeArtifactIds []int

	//excludeWfRunners
	ExcludeWfrIds []int

	//ApprovalNode data
	ApprovalNodeConfigured bool
	ApproversCount         int

	//pluginStage
	PluginStage string

	// UseCdStageQueryV2 is to set query version
	UseCdStageQueryV2 bool
}

type CiNodePromotionArtifactsRequest struct {
	CiPipelineId   int
	ListingOptions util.ListingFilterOptions
}

type ExtCiNodePromotionArtifactsRequest struct {
	ExternalCiPipelineId int
	ListingOptions       util.ListingFilterOptions
}

type CdNodePromotionArtifactsRequest struct {
	ResourceCdPipelineId int
	ListingOptions       util.ListingFilterOptions
}

type ArtifacPromotionPendingNodeRequest struct {
	ResourceCdPipelineId int
	ListingOptions       util.ListingFilterOptions
}

type ArtifactPromotionPendingForCurrentUserRequest struct {
	ImagePromoterAccessCdPipelineIds []int
	ListingOptions                   util.ListingFilterOptions
}
