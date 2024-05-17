package artifactPromotion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	bean3 "github.com/devtron-labs/devtron/api/bean"
	client "github.com/devtron-labs/devtron/client/events"
	"github.com/devtron-labs/devtron/enterprise/pkg/expressionEvaluators"
	"github.com/devtron-labs/devtron/enterprise/pkg/resourceFilter"
	"github.com/devtron-labs/devtron/internal/sql/models"
	repository2 "github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/internal/sql/repository/appWorkflow"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	"github.com/devtron-labs/devtron/internal/util"
	appWorkflow2 "github.com/devtron-labs/devtron/pkg/appWorkflow"
	bean4 "github.com/devtron-labs/devtron/pkg/appWorkflow/bean"
	"github.com/devtron-labs/devtron/pkg/auth/user"
	bean5 "github.com/devtron-labs/devtron/pkg/auth/user/bean"
	"github.com/devtron-labs/devtron/pkg/cluster"
	bean6 "github.com/devtron-labs/devtron/pkg/cluster/repository/bean"
	bean2 "github.com/devtron-labs/devtron/pkg/deployment/trigger/devtronApps/bean"
	"github.com/devtron-labs/devtron/pkg/pipeline"
	repository3 "github.com/devtron-labs/devtron/pkg/pipeline/repository"
	"github.com/devtron-labs/devtron/pkg/policyGovernance"
	"github.com/devtron-labs/devtron/pkg/policyGovernance/artifactPromotion/bean"
	"github.com/devtron-labs/devtron/pkg/policyGovernance/artifactPromotion/constants"
	"github.com/devtron-labs/devtron/pkg/policyGovernance/artifactPromotion/read"
	"github.com/devtron-labs/devtron/pkg/policyGovernance/artifactPromotion/repository"
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/devtron-labs/devtron/pkg/team"
	"github.com/devtron-labs/devtron/pkg/workflow/cd"
	"github.com/devtron-labs/devtron/pkg/workflow/dag"
	util3 "github.com/devtron-labs/devtron/util"
	util2 "github.com/devtron-labs/devtron/util/event"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type ApprovalRequestService interface {
	HandleArtifactPromotionRequest(ctx *util3.RequestCtx, request *bean.ArtifactPromotionRequest, authorizedEnvironments map[string]bool) ([]bean.EnvironmentPromotionMetaData, error)
	GetPromotionRequestById(promotionRequestId int) (*bean.ArtifactPromotionApprovalResponse, error)
	FetchWorkflowPromoteNodeList(ctx *util3.RequestCtx, workflowId int, artifactId int, rbacChecker func(token string, appName string, envNames []string) map[string]bool) (*bean.EnvironmentListingResponse, error)
	FetchApprovalAllowedEnvList(ctx *util3.RequestCtx, artifactId int, environmentName string, promotionApproverAuth func(*util3.RequestCtx, []string) map[string]bool) ([]bean.EnvironmentApprovalMetadata, error)
}

type ApprovalRequestServiceImpl struct {
	logger                               *zap.SugaredLogger
	ciPipelineService                    pipeline.CiPipelineConfigService
	cdPipelineService                    pipeline.CdPipelineConfigService
	pipelineStageService                 pipeline.PipelineStageService
	environmentService                   cluster.EnvironmentService
	userService                          user.UserService
	appWorkflowService                   appWorkflow2.AppWorkflowService
	cdWorkflowService                    cd.CdWorkflowCommonService
	resourceFilterConditionsEvaluator    resourceFilter.ResourceFilterEvaluator
	resourceFilterEvaluationAuditService resourceFilter.FilterEvaluationAuditService
	imageTaggingService                  pipeline.ImageTaggingService
	promotionPolicyDataReadService       read.ArtifactPromotionDataReadService
	workflowDagExecutor                  dag.WorkflowDagExecutor
	transactionManager                   sql.TransactionWrapper

	ciArtifactRepository                       repository2.CiArtifactRepository
	artifactPromotionApprovalRequestRepository repository.RequestRepository
	requestApprovalUserdataRepo                pipelineConfig.RequestApprovalUserdataRepository
	eventFactory                               client.EventFactory
	eventClient                                client.EventClient
	teamService                                team.TeamService
}

func NewApprovalRequestServiceImpl(
	logger *zap.SugaredLogger,
	ciPipelineService pipeline.CiPipelineConfigService,
	cdPipelineService pipeline.CdPipelineConfigService,
	userService user.UserService,
	appWorkflowService appWorkflow2.AppWorkflowService,
	cdWorkflowService cd.CdWorkflowCommonService,
	resourceFilterConditionsEvaluator resourceFilter.ResourceFilterEvaluator,
	imageTaggingService pipeline.ImageTaggingService,
	promotionPolicyService read.ArtifactPromotionDataReadService,
	workflowDagExecutor dag.WorkflowDagExecutor,
	policyEventNotifier PolicyEventNotifier,
	commonPoliyApplyEventNotifier policyGovernance.CommonPoliyApplyEventNotifier,
	pipelineStageService pipeline.PipelineStageService,
	environmentService cluster.EnvironmentService,
	resourceFilterEvaluationAuditService resourceFilter.FilterEvaluationAuditService,
	transactionManager sql.TransactionWrapper,

	ciArtifactRepository repository2.CiArtifactRepository,
	artifactPromotionApprovalRequestRepository repository.RequestRepository,
	requestApprovalUserdataRepo pipelineConfig.RequestApprovalUserdataRepository,
	eventFactory client.EventFactory,
	eventClient client.EventClient,
	teamService team.TeamService,
) *ApprovalRequestServiceImpl {

	artifactApprovalService := &ApprovalRequestServiceImpl{
		logger:                               logger,
		ciPipelineService:                    ciPipelineService,
		cdPipelineService:                    cdPipelineService,
		userService:                          userService,
		appWorkflowService:                   appWorkflowService,
		cdWorkflowService:                    cdWorkflowService,
		resourceFilterConditionsEvaluator:    resourceFilterConditionsEvaluator,
		imageTaggingService:                  imageTaggingService,
		promotionPolicyDataReadService:       promotionPolicyService,
		workflowDagExecutor:                  workflowDagExecutor,
		pipelineStageService:                 pipelineStageService,
		environmentService:                   environmentService,
		resourceFilterEvaluationAuditService: resourceFilterEvaluationAuditService,
		transactionManager:                   transactionManager,

		ciArtifactRepository:                       ciArtifactRepository,
		artifactPromotionApprovalRequestRepository: artifactPromotionApprovalRequestRepository,
		requestApprovalUserdataRepo:                requestApprovalUserdataRepo,
		eventFactory:                               eventFactory,
		eventClient:                                eventClient,
		teamService:                                teamService,
	}

	// register hooks
	policyEventNotifier.AddDeleteEventObserver(artifactApprovalService.onPolicyDelete)
	policyEventNotifier.AddUpdateEventObserver(artifactApprovalService.onPolicyUpdate)
	commonPoliyApplyEventNotifier.AddApplyEventObserver(policyGovernance.ImagePromotion, artifactApprovalService.onApplyPolicy)
	return artifactApprovalService
}

func (impl *ApprovalRequestServiceImpl) HandleArtifactPromotionRequest(ctx *util3.RequestCtx, request *bean.ArtifactPromotionRequest, authorizedEnvironments map[string]bool) ([]bean.EnvironmentPromotionMetaData, error) {

	switch request.Action {

	case constants.ACTION_PROMOTE:
		return impl.promoteArtifact(ctx, request, authorizedEnvironments)
	case constants.ACTION_APPROVE:
		return impl.approveArtifactPromotion(ctx, request, authorizedEnvironments)
	case constants.ACTION_CANCEL:
		_, err := impl.cancelPromotionApprovalRequest(ctx, request)
		return nil, err
	}
	return nil, errors.New(constants.UnknownRequestAction)
}

func (impl *ApprovalRequestServiceImpl) GetPromotionRequestById(promotionRequestId int) (*bean.ArtifactPromotionApprovalResponse, error) {
	promotionRequest, err := impl.artifactPromotionApprovalRequestRepository.FindById(promotionRequestId)
	if err != nil {
		impl.logger.Errorw("error in getting promotion request by id", "promotionRequestId", promotionRequestId, "err", err)
		if errors.Is(err, pg.ErrNoRows) {
			return nil, util.NewApiError().WithHttpStatusCode(http.StatusNotFound).WithUserMessage(constants.ArtifactPromotionRequestNotFoundErr)
		}
		return nil, err
	}
	artifactPromotionResponse := &bean.ArtifactPromotionApprovalResponse{
		Id:                      promotionRequest.Id,
		PolicyId:                promotionRequest.PolicyId,
		PolicyEvaluationAuditId: promotionRequest.PolicyEvaluationAuditId,
		ArtifactId:              promotionRequest.ArtifactId,
		SourceType:              promotionRequest.SourceType,
		SourcePipelineId:        promotionRequest.SourcePipelineId,
		DestinationPipelineId:   promotionRequest.DestinationPipelineId,
		Status:                  promotionRequest.Status,
	}
	return artifactPromotionResponse, nil
}

func (impl *ApprovalRequestServiceImpl) FetchApprovalAllowedEnvList(ctx *util3.RequestCtx, artifactId int, environmentName string, promotionApproverAuth func(*util3.RequestCtx, []string) map[string]bool) ([]bean.EnvironmentApprovalMetadata, error) {

	environmentApprovalMetadata := make([]bean.EnvironmentApprovalMetadata, 0)

	artifact, err := impl.ciArtifactRepository.Get(artifactId)
	if err != nil {
		impl.logger.Errorw(constants.ARTIFACT_NOT_FOUND_ERR, "artifactId", artifactId, "err", err)
		return nil, util.NewApiError().WithHttpStatusCode(http.StatusUnprocessableEntity).WithUserMessage(constants.ARTIFACT_NOT_FOUND_ERR).WithInternalMessage(constants.ARTIFACT_NOT_FOUND_ERR)
	}

	promotionRequests, err := impl.artifactPromotionApprovalRequestRepository.FindRequestsByArtifactAndOptionalEnv(artifactId, environmentName, constants.AWAITING_APPROVAL)
	if err != nil {
		impl.logger.Errorw("error in finding promotion requests in awaiting state for given artifactId", "artifactId", artifactId, "err", err)
		return nil, err
	}

	if len(promotionRequests) == 0 {
		return environmentApprovalMetadata, nil
	}

	destinationPipelineIds := make([]int, len(promotionRequests))
	for i, request := range promotionRequests {
		destinationPipelineIds[i] = request.DestinationPipelineId
	}

	pipelineIdToDaoMapping, err := impl.getPipelineIdToDaoMapping(destinationPipelineIds)
	if err != nil {
		impl.logger.Errorw("error in getting pipelineId to Dao mapping", "destinationPipelineIds", destinationPipelineIds, "err", err)
		return environmentApprovalMetadata, err
	}

	envIds := make([]int, 0, len(pipelineIdToDaoMapping))
	for _, pipelineDao := range pipelineIdToDaoMapping {
		envIds = append(envIds, pipelineDao.EnvironmentId)
	}

	rbacObjects, pipelineIdToRbacObjMap := impl.getRbacObjects(pipelineIdToDaoMapping)
	rbacResults := promotionApproverAuth(ctx, rbacObjects)

	appId := pipelineIdToDaoMapping[promotionRequests[0].DestinationPipelineId].AppId

	policiesMap, err := impl.promotionPolicyDataReadService.GetPromotionPolicyByAppAndEnvIds(ctx, appId, envIds)
	if err != nil {
		impl.logger.Errorw("error in fetching policies by appId and envIds", "appId", appId, "envIds", envIds, "err", err)
		return nil, err
	}

	for _, request := range promotionRequests {

		pipelineDao := pipelineIdToDaoMapping[request.DestinationPipelineId]

		environmentMetadata := bean.EnvironmentApprovalMetadata{
			Name:            pipelineDao.Environment.Name,
			ApprovalAllowed: true,
			Reasons:         make([]string, 0),
		}

		policy := policiesMap[pipelineDao.Environment.Name]
		if policy == nil {
			environmentMetadata.ApprovalAllowed = false
			environmentMetadata.Reasons = append(environmentMetadata.Reasons, constants.USER_DOES_NOT_HAVE_ARTIFACT_PROMOTER_ACCESS)
			environmentApprovalMetadata = append(environmentApprovalMetadata, environmentMetadata)
			continue
		} else if policy.BlockImageBuilderToApprove(artifact.CreatedBy, ctx.GetUserId()) {
			environmentMetadata.ApprovalAllowed = false
			environmentMetadata.Reasons = append(environmentMetadata.Reasons, constants.BUILD_TRIGGER_USER_CANNOT_APPROVE_MSG)
		} else if policy.BlockPromotionRequesterToApprove(request.CreatedBy, ctx.GetUserId()) {
			environmentMetadata.ApprovalAllowed = false
			environmentMetadata.Reasons = append(environmentMetadata.Reasons, constants.PROMOTION_REQUESTED_BY_USER_CANNOT_APPROVE_MSG)
		}

		rbacObj := pipelineIdToRbacObjMap[request.DestinationPipelineId]
		if isAuthorized := rbacResults[rbacObj]; !isAuthorized {
			environmentMetadata.ApprovalAllowed = false
			environmentMetadata.Reasons = append(environmentMetadata.Reasons, constants.USER_DOES_NOT_HAVE_ARTIFACT_PROMOTER_ACCESS)
		}
		environmentApprovalMetadata = append(environmentApprovalMetadata, environmentMetadata)
	}
	return environmentApprovalMetadata, nil
}

func (impl *ApprovalRequestServiceImpl) FetchWorkflowPromoteNodeList(ctx *util3.RequestCtx, workflowId int, artifactId int, rbacChecker func(token string, appName string, envNames []string) map[string]bool) (*bean.EnvironmentListingResponse, error) {
	metadata, err := impl.fetchEnvMetaDataListingRequestMetadata(ctx.GetToken(), workflowId, artifactId, rbacChecker)
	if err != nil {
		impl.logger.Errorw("error in fetching envMetaDataListing request metadata", "workflowId", workflowId, "artifactId", artifactId, "err", err)
		return nil, err
	}
	envMap := metadata.GetActiveEnvironmentsMap()
	result := &bean.EnvironmentListingResponse{}
	result.CiSource = metadata.GetSourceMetaData().GetCiSourceMeta()
	policiesMap, err := impl.promotionPolicyDataReadService.GetPromotionPolicyByAppAndEnvIds(ctx, metadata.GetAppId(), metadata.GetActiveAuthorisedEnvIds())
	if err != nil {
		impl.logger.Errorw("error in getting the policies", "appId", metadata.GetAppId(), "envIds", metadata.GetActiveAuthorisedEnvIds(), "err", err)
		return nil, err
	}

	if artifactId != 0 {
		responses, err := impl.evaluatePoliciesOnArtifact(metadata, policiesMap)
		if err != nil {
			impl.logger.Errorw("error in evaluating policies on an ciArtifact", "ciArtifactId", artifactId, "policiesMap", policiesMap, "authorizedEnvironments", metadata.GetActiveAuthorisedEnvIds(), "err", err)
			return nil, err
		}
		result.Environments = responses
		return result, nil
	}

	responseMap := metadata.GetDefaultEnvironmentPromotionMetaDataResponseMap()
	for envName, resp := range responseMap {
		resp.IsVirtualEnvironment = envMap[envName].IsVirtualEnvironment
		responseMap[envName] = resp
	}

	for envName, policy := range policiesMap {
		approvalCount := 0
		if policy != nil {
			approvalCount = policy.ApprovalMetaData.ApprovalCount
		}
		responseMap[envName] = bean.EnvironmentPromotionMetaData{
			PromotionPossible:          true,
			Name:                       envName,
			ApprovalCount:              approvalCount,
			IsVirtualEnvironment:       envMap[envName].IsVirtualEnvironment,
			PromotionValidationMessage: constants.EMPTY,
		}
	}

	for envName, resp := range responseMap {
		if resp.PromotionValidationMessage == constants.EMPTY {
			resp.PromotionValidationMessage = constants.POLICY_NOT_CONFIGURED
			resp.PromotionPossible = false
		}
		responseMap[envName] = resp
	}

	responses := make([]bean.EnvironmentPromotionMetaData, 0, len(responseMap))
	for _, envResponse := range responseMap {
		responses = append(responses, envResponse)
	}

	result.Environments = responses
	return result, nil
}

func (impl *ApprovalRequestServiceImpl) getSourceInfoAndPipelineIds(workflowId int) (*bean.SourceMetaData, []int, error) {
	allAppWorkflowMappings, err := impl.appWorkflowService.FindAppWorkflowMapping(workflowId)
	if err != nil {
		impl.logger.Errorw("error in finding the app workflow mappings using appWorkflowId", "workflowId", workflowId, "err", err)
		return nil, nil, err
	}

	if len(allAppWorkflowMappings) == 0 {
		return nil, nil, util.NewApiError().WithHttpStatusCode(http.StatusConflict).WithUserMessage(constants.WorkflowNotFoundErr)
	}

	sourcePipelineMapping := bean4.AppWorkflowMappingDto{}
	pipelineIds := make([]int, 0, len(allAppWorkflowMappings))
	for _, mapping := range allAppWorkflowMappings {
		if mapping.Type == appWorkflow.CDPIPELINE {
			pipelineIds = append(pipelineIds, mapping.ComponentId)
		}
		if mapping.ParentId == 0 {
			sourcePipelineMapping = mapping
		}
	}

	sourceId := sourcePipelineMapping.ComponentId
	var sourceName string
	var sourceType constants.SourceTypeStr
	if sourcePipelineMapping.Type == appWorkflow.CIPIPELINE {
		ciPipeline, err := impl.ciPipelineService.GetCiPipelineById(sourceId)
		if err != nil {
			impl.logger.Errorw("error in fetching ci pipeline by id", "ciPipelineId", sourceId, "err", err)
			return nil, nil, err
		}
		sourceName = ciPipeline.Name
		sourceType = bean.GetSourceTypeFromPipelineType(ciPipeline.PipelineType)
	} else if sourcePipelineMapping.Type == appWorkflow.WEBHOOK {
		sourceType = constants.SOURCE_TYPE_WEBHOOK
	}

	// set source metadata
	sourceInfo := &bean.SourceMetaData{}
	sourceInfo = sourceInfo.WithName(sourceName).WithType(sourceType).WithId(sourceId).WithSourceWorkflowId(workflowId)
	return sourceInfo, pipelineIds, nil
}

func (impl *ApprovalRequestServiceImpl) fetchEnvMetaDataListingRequestMetadata(token string, workflowId int, artifactId int, rbacChecker func(token string, appName string, envNames []string) map[string]bool) (*bean.RequestMetaData, error) {

	sourceInfo, pipelineIds, err := impl.getSourceInfoAndPipelineIds(workflowId)
	if err != nil {
		impl.logger.Errorw("error in finding source info and pipelinesIds", "sourceInfo", sourceInfo, "pipelineIds", pipelineIds, "err", err)
		return nil, err
	}
	pipelines, err := impl.cdPipelineService.FindByIdsIn(pipelineIds)
	if err != nil {
		impl.logger.Errorw("error in finding pipelines", "pipelineIds", pipelineIds, "err", err)
		return nil, err
	}
	environmentNames := make([]string, 0, len(pipelines))
	appName := ""
	appId := 0
	environments := make([]*bean6.EnvironmentBean, 0, len(pipelines))
	for _, pipeline := range pipelines {
		environmentNames = append(environmentNames, pipeline.Environment.Name)
		environment := &bean6.EnvironmentBean{}
		environment.AdaptFromEnvironment(&pipeline.Environment)
		environments = append(environments, environment)
		appName = pipeline.App.AppName
		appId = pipeline.AppId
	}
	authorizedEnvironments := rbacChecker(token, appName, environmentNames)
	cdPipelines := make([]*pipelineConfig.Pipeline, 0, len(pipelines))
	for _, pipeline := range pipelines {
		if authorizedEnvironments[pipeline.Environment.Name] {
			cdPipelines = append(cdPipelines, pipeline)
		}
	}

	requestMetaData := &bean.RequestMetaData{}
	requestMetaData = requestMetaData.WithAppId(appId).WithWorkflowId(workflowId)
	requestMetaData.SetSourceMetaData(sourceInfo)
	requestMetaData.SetActiveEnvironments(environmentNames, authorizedEnvironments, environments)
	requestMetaData.SetDestinationPipelineMetaData(cdPipelines)
	if artifactId > 0 {
		ciArtifact, err := impl.ciArtifactRepository.Get(artifactId)
		if err != nil {
			impl.logger.Errorw("error in finding the artifact using id", "artifactId", artifactId, "err", err)
			errorResp := util.NewApiError().WithHttpStatusCode(http.StatusInternalServerError).WithInternalMessage(fmt.Sprintf(constants.ArtifactFindErr+" , err : %s", err.Error())).WithUserMessage(constants.ArtifactFindErr)
			if errors.Is(err, pg.ErrNoRows) {
				errorResp = errorResp.WithHttpStatusCode(http.StatusConflict).WithUserMessage(constants.ArtifactFindErr)
			}
			return nil, errorResp
		}
		requestMetaData = requestMetaData.WithCiArtifact(ciArtifact)
	}

	return requestMetaData, nil
}

func (impl *ApprovalRequestServiceImpl) computeFilterParams(ciArtifact *repository2.CiArtifact) ([]expressionEvaluators.ExpressionParam, error) {
	var ciMaterials []repository2.CiMaterialInfo
	err := json.Unmarshal([]byte(ciArtifact.MaterialInfo), &ciMaterials)
	if err != nil {
		impl.logger.Errorw("error in parsing ci artifact material info", "materialInfo", "")
		return nil, err
	}

	imageTags, err := impl.imageTaggingService.GetTagsByArtifactId(ciArtifact.Id)
	if err != nil {
		impl.logger.Errorw("error in fetching the image tags using artifact id", "artifactId", ciArtifact.Id, "err", err)
		return nil, err
	}

	releaseTags := make([]string, 0, len(imageTags))
	for _, imageTag := range imageTags {
		releaseTags = append(releaseTags, imageTag.TagName)
	}
	params, err := resourceFilter.GetParamsFromArtifact(ciArtifact.Image, releaseTags, ciMaterials)
	return params, err
}

func (impl *ApprovalRequestServiceImpl) evaluatePoliciesOnArtifact(metadata *bean.RequestMetaData, policiesMap map[string]*bean.PromotionPolicy) ([]bean.EnvironmentPromotionMetaData, error) {
	envMap := metadata.GetActiveEnvironmentsMap()
	responseMap := metadata.GetDefaultEnvironmentPromotionMetaDataResponseMap()
	for envName, resp := range responseMap {
		if env, ok := envMap[envName]; ok {
			resp.PromotionValidationMessage = constants.POLICY_NOT_CONFIGURED
			resp.IsVirtualEnvironment = env.IsVirtualEnvironment
			responseMap[envName] = resp
		}
	}

	if len(policiesMap) > 0 {
		// can be concurrent
		params, err := impl.computeFilterParams(metadata.GetCiArtifact())
		if err != nil {
			impl.logger.Errorw("error in finding the required CEL expression parameters for using ciArtifact", "err", err)
			return nil, err
		}
		for envName, policy := range policiesMap {
			evaluationResult, err := impl.resourceFilterConditionsEvaluator.EvaluateFilter(policy.Conditions, expressionEvaluators.ExpressionMetadata{Params: params})
			if err != nil {
				impl.logger.Errorw("policy evaluation failed with error", "policyConditions", policy.Conditions, "envName", envName, policy.Conditions, "params", params, "err", err)
				responseMap[envName] = bean.EnvironmentPromotionMetaData{
					Name:                       envName,
					ApprovalCount:              policy.ApprovalMetaData.ApprovalCount,
					PromotionPossible:          false,
					PromotionValidationMessage: constants.POLICY_EVALUATION_ERRORED,
				}
				continue
			}
			envResp := responseMap[envName]
			envResp.ApprovalCount = policy.ApprovalMetaData.ApprovalCount
			envResp.PromotionValidationMessage = constants.EMPTY
			envResp.PromotionPossible = evaluationResult
			// checks on metadata not needed as this is just an evaluation flow (kinda validation)
			if !evaluationResult {
				envResp.PromotionValidationMessage = constants.BLOCKED_BY_POLICY
			}
			responseMap[envName] = envResp
		}
	}

	result := make([]bean.EnvironmentPromotionMetaData, 0, len(responseMap))
	for _, envResponse := range responseMap {
		result = append(result, envResponse)
	}
	return result, nil
}

func (impl *ApprovalRequestServiceImpl) getMetadataForApproveRequest(request *bean.ArtifactPromotionRequest, authorizedEnvironments map[string]bool) (*bean.RequestMetaData, error) {
	metadata, err := impl.getMetadataObjectWithDestinationInfo(nil, request, authorizedEnvironments)
	if err != nil {
		impl.logger.Errorw("error in getting destination pipelines metadata for the request", "request", request, "err", err)
		return nil, err
	}

	metadata, err = impl.getMetadataObjectWithArtifactInfo(metadata, request.ArtifactId)
	if err != nil {
		impl.logger.Errorw("error in getting artifact metadata", "artifactId", request.ArtifactId, "err", err)
		return nil, err
	}
	return metadata, nil
}

func (impl *ApprovalRequestServiceImpl) approveArtifactPromotion(ctx *util3.RequestCtx, request *bean.ArtifactPromotionRequest, authorizedEnvironments map[string]bool) ([]bean.EnvironmentPromotionMetaData, error) {
	// get request and check if it is promoted already.
	// attempt approving this by creating new resource_approval_user_data, if unique constraint error ,current user already did something.
	// attempt success , then get the approval count and check no of approvals got
	//  promote if approvalCount > approvals received

	metadata, err := impl.getMetadataForApproveRequest(request, authorizedEnvironments)
	if err != nil {
		impl.logger.Errorw("error in getting metadata information for approve request", "request", request, "authorizedEnvironments", authorizedEnvironments, "err", err)
		return nil, err
	}
	responseMap := metadata.GetDefaultEnvironmentPromotionMetaDataResponseMap()

	promotionRequests, err := impl.artifactPromotionApprovalRequestRepository.FindByArtifactAndDestinationPipelineIds(request.ArtifactId, metadata.GetActiveAuthorisedPipelineIds())
	if err != nil {
		impl.logger.Errorw("error in getting artifact promotion requests object by pipeline ids", "pipelineIds", metadata.GetActiveAuthorisedPipelineIds(), "err", err)
		return nil, err
	}

	if len(promotionRequests) == 0 {
		result := make([]bean.EnvironmentPromotionMetaData, 0, len(responseMap))
		for _, resp := range responseMap {
			resp.PromotionValidationState = constants.ERROR
			resp.PromotionValidationMessage = constants.ArtifactPromotionRequestNotFoundErr
			result = append(result, resp)

		}
		return result, nil
	}
	// policies fetched form above policy ids
	policies, err := impl.promotionPolicyDataReadService.GetPromotionPolicyByAppAndEnvIds(ctx, request.AppId, metadata.GetActiveAuthorisedPipelineEnvIds())
	if err != nil {
		impl.logger.Errorw("error in finding the promotionPolicy by appId and envId names", "appid", request.AppId, "envIds", metadata.GetActiveAuthorisedPipelineEnvIds(), "err", err)
		return nil, err
	}

	// map the policies for O(1) access
	policyIdMap := make(map[int]*bean.PromotionPolicy)
	for _, policy := range policies {
		policyIdMap[policy.Id] = policy
	}

	environmentResponses, err := impl.initiateApprovalProcess(ctx, metadata, promotionRequests, responseMap, policyIdMap)
	if err != nil {
		impl.logger.Errorw("error in finding approving the artifact promotion requests", "promotionRequests", promotionRequests, "err", err)
		return nil, err
	}
	return environmentResponses, nil
}

func (impl *ApprovalRequestServiceImpl) approveRequests(ctx *util3.RequestCtx, metadata *bean.RequestMetaData, validRequestIds []int, policyIdMap map[int]*bean.PromotionPolicy, promotionRequests []*repository.ArtifactPromotionApprovalRequest, responses map[string]bean.EnvironmentPromotionMetaData) map[string]bean.EnvironmentPromotionMetaData {
	validRequestsMap := make(map[int]bool)
	for _, requestId := range validRequestIds {
		validRequestsMap[requestId] = true
	}

	pipelineIdVsEnvMap := metadata.GetActiveAuthorisedPipelineIdEnvMap()
	for _, promotionRequest := range promotionRequests {
		// skip the invalid requests
		if ok := validRequestsMap[promotionRequest.Id]; !ok {
			continue
		}
		resp := responses[pipelineIdVsEnvMap[promotionRequest.DestinationPipelineId]]

		if !policyIdMap[promotionRequest.PolicyId].CanApprove(promotionRequest.CreatedBy, metadata.GetCiArtifact().CreatedBy, ctx.GetUserId()) {
			resp.PromotionValidationMessage = constants.BLOCKED_BY_POLICY
			continue
		}
		promotionRequestApprovedUserData := &pipelineConfig.RequestApprovalUserData{
			ApprovalRequestId: promotionRequest.Id,
			RequestType:       models.ARTIFACT_PROMOTION_APPROVAL,
			UserId:            ctx.GetUserId(),
			UserResponse:      pipelineConfig.APPROVED,
		}
		// have to do this in loop as we have to ensure partial approval even in case of partial failure
		err := impl.requestApprovalUserdataRepo.SaveRequestApprovalUserData(promotionRequestApprovedUserData)
		if err != nil {
			impl.logger.Errorw("error in saving promotion approval user data", "promotionRequestId", promotionRequest.Id, "err", err)
			if strings.Contains(err.Error(), string(pipelineConfig.UNIQUE_USER_REQUEST_ACTION)) {
				resp.PromotionValidationMessage = constants.ALREADY_APPROVED

			} else {
				resp.PromotionValidationMessage = constants.ERRORED_APPROVAL
			}
			responses[pipelineIdVsEnvMap[promotionRequest.DestinationPipelineId]] = resp
			continue
		}

		resp.PromotionValidationMessage = constants.APPROVED
		responses[pipelineIdVsEnvMap[promotionRequest.DestinationPipelineId]] = resp
	}
	return responses
}

func (impl *ApprovalRequestServiceImpl) initiateApprovalProcess(ctx *util3.RequestCtx, metadata *bean.RequestMetaData, promotionRequests []*repository.ArtifactPromotionApprovalRequest, responses map[string]bean.EnvironmentPromotionMetaData, policyIdMap map[int]*bean.PromotionPolicy) ([]bean.EnvironmentPromotionMetaData, error) {

	if len(metadata.GetActiveAuthorisedEnvIds()) == 0 {
		return nil, util.NewApiError().WithHttpStatusCode(http.StatusUnauthorized).WithUserMessage(constants.NoApprovePermissionOnEnvsErr).WithInternalMessage(constants.NoApprovePermissionOnEnvsErr)
	}
	pipelineIdVsEnvMap := metadata.GetActiveAuthorisedPipelineIdEnvMap()
	staleRequestIds, validRequestIds, responses := impl.filterValidAndStaleRequests(promotionRequests, responses, pipelineIdVsEnvMap, policyIdMap)

	tx, err := impl.transactionManager.StartTx()
	if err != nil {
		impl.logger.Errorw("error in starting the transaction", "promotionRequests", promotionRequests, "err", err)
		return nil, err
	}
	defer impl.transactionManager.RollbackTx(tx)
	responses = impl.approveRequests(ctx, metadata, validRequestIds, policyIdMap, promotionRequests, responses)
	if len(staleRequestIds) > 0 {
		// attempt deleting the stale requests in bulk
		err = impl.artifactPromotionApprovalRequestRepository.MarkStaleByIds(tx, staleRequestIds)
		if err != nil {
			impl.logger.Errorw("error in deleting the request raised using a deleted promotion policy (stale requests)", "staleRequestIds", staleRequestIds, "err", err)
			// not returning by choice, don't interrupt the user flow
		}
	}
	if len(validRequestIds) > 0 {
		// fetch all the approved users data for the valid requestIds
		approvedUsersData, err := impl.requestApprovalUserdataRepo.FetchApprovalDataForRequests(validRequestIds, models.ARTIFACT_PROMOTION_APPROVAL)
		if err != nil && !errors.Is(err, pg.ErrNoRows) {
			impl.logger.Errorw("error in finding the approved users data for a artifact promotion request", "promotionRequestIds", validRequestIds, "err", err)
			return nil, err
		}

		// club the approved users(we just need count for now) per requestId
		promotionRequestIdVsApprovedUserCount := make(map[int]int)
		for _, _approvedUsersData := range approvedUsersData {
			count := promotionRequestIdVsApprovedUserCount[_approvedUsersData.ApprovalRequestId]
			promotionRequestIdVsApprovedUserCount[_approvedUsersData.ApprovalRequestId] = count + 1
		}

		validRequestMap := make(map[int]bool)
		for _, id := range validRequestIds {
			validRequestMap[id] = true
		}
		// filter out promotable requests.
		// we will promote if the current number approvals got for any request exceeds the current configured no of approvals in the policy
		promotableRequestIds := make([]int, 0, len(validRequestIds))
		for _, promotionRequest := range promotionRequests {
			if validRequestMap[promotionRequest.Id] {
				approvalCount := promotionRequestIdVsApprovedUserCount[promotionRequest.Id]
				if policyIdMap[promotionRequest.PolicyId].CanBePromoted(approvalCount) {
					promotableRequestIds = append(promotableRequestIds, promotionRequest.Id)
				}
			}
		}

		// promote the promotableRequestIds
		err = impl.artifactPromotionApprovalRequestRepository.MarkPromoted(tx, promotableRequestIds, ctx.GetUserId())
		if err != nil {
			impl.logger.Errorw("error in promoting the approval requests", "promotableRequestIds", promotableRequestIds, "err", err)
			return nil, err
		}

		promotionRequestIdToDaoMap := make(map[int]*repository.ArtifactPromotionApprovalRequest)
		for _, promotionRequest := range promotionRequests {
			promotionRequestIdToDaoMap[promotionRequest.Id] = promotionRequest
		}

		if len(promotableRequestIds) > 0 {
			err = impl.handleArtifactPromotionSuccess(promotableRequestIds, promotionRequestIdToDaoMap, metadata.GetActiveAuthorisedPipelineDaoMap())
			if err != nil {
				impl.logger.Errorw("error in handling the successful artifact promotion event for promotedRequests", "promotableRequestIds", promotableRequestIds, "err", err)
				return nil, err
			}
		}
	}
	err = impl.transactionManager.CommitTx(tx)
	if err != nil {
		impl.logger.Errorw("error in committing the transaction", "validRequestIds", validRequestIds, "err", err)
		return nil, err
	}

	result := make([]bean.EnvironmentPromotionMetaData, 0, len(responses))
	for _, resp := range responses {
		result = append(result, resp)
	}
	return result, nil
}

func (impl *ApprovalRequestServiceImpl) filterValidAndStaleRequests(promotionRequests []*repository.ArtifactPromotionApprovalRequest, responses map[string]bean.EnvironmentPromotionMetaData, pipelineIdVsEnvMap map[int]string, policyIdMap map[int]*bean.PromotionPolicy) ([]int, []int, map[string]bean.EnvironmentPromotionMetaData) {
	staleRequestIds := make([]int, 0)
	validRequestIds := make([]int, 0)
	for _, promotionRequest := range promotionRequests {
		resp := responses[pipelineIdVsEnvMap[promotionRequest.DestinationPipelineId]]
		_, policyFound := policyIdMap[promotionRequest.PolicyId]
		_, destinationPipelineFoud := pipelineIdVsEnvMap[promotionRequest.DestinationPipelineId]
		ok := policyFound && destinationPipelineFoud
		if ok {
			validRequestIds = append(validRequestIds, promotionRequest.Id)
		}
		if !ok {
			// policy is not found in the map, and the request is still in awaiting state.
			// although the policy is no longer governing the current pipeline
			// this is a stale case.
			// mark it stale
			staleRequestIds = append(staleRequestIds, promotionRequest.Id)

			// also set the response messages
			resp.PromotionPossible = false
			resp.PromotionValidationMessage = constants.PromotionRequestStale
		} else if promotionRequest.Status != constants.AWAITING_APPROVAL {
			resp.PromotionValidationMessage = constants.PromotionValidationMsg(fmt.Sprintf("artifact is in %s state", promotionRequest.Status.Status()))
		}
		responses[pipelineIdVsEnvMap[promotionRequest.DestinationPipelineId]] = resp
	}
	return staleRequestIds, validRequestIds, responses
}

func (impl *ApprovalRequestServiceImpl) handleArtifactPromotionSuccess(promotableRequestIds []int, promotionRequestIdToDaoMap map[int]*repository.ArtifactPromotionApprovalRequest, pipelineIdToDaoMap map[int]*pipelineConfig.Pipeline) error {
	promotedCiArtifactIds := make([]int, 0)
	for _, id := range promotableRequestIds {
		promotableRequest := promotionRequestIdToDaoMap[id]
		promotedCiArtifactIds = append(promotedCiArtifactIds, promotableRequest.ArtifactId)
	}

	artifacts, err := impl.ciArtifactRepository.GetByIds(promotedCiArtifactIds)
	if err != nil {
		impl.logger.Errorw("error in fetching the artifacts by ids", "artifactIds", promotedCiArtifactIds, "err", err)
		return err
	}

	artifactsMap := make(map[int]*repository2.CiArtifact)
	for _, artifact := range artifacts {
		artifactsMap[artifact.Id] = artifact
	}
	for _, id := range promotableRequestIds {
		promotableRequest := promotionRequestIdToDaoMap[id]
		pipelineDao := pipelineIdToDaoMap[promotableRequest.DestinationPipelineId]
		triggerRequest := bean2.TriggerRequest{
			CdWf:        nil,
			Pipeline:    pipelineDao,
			Artifact:    artifactsMap[promotableRequest.ArtifactId],
			TriggeredBy: 1,
			TriggerContext: bean2.TriggerContext{
				// send request context's context
				Context: context.Background(),
			},
		}
		impl.workflowDagExecutor.HandleArtifactPromotionEvent(triggerRequest)
	}
	return nil
}

func (impl *ApprovalRequestServiceImpl) getMetadataObjectWithSourceInfo(requestMetaData *bean.RequestMetaData, request *bean.ArtifactPromotionRequest) (*bean.RequestMetaData, error) {

	if requestMetaData == nil {
		requestMetaData = &bean.RequestMetaData{}
	}
	// set source metadata
	sourceMeta, err := impl.fetchSourceMeta(request.SourceName, request.SourceType, request.AppId, request.WorkflowId)
	if err != nil {
		impl.logger.Errorw("error in validating the request", "request", request, "err", err)
		return nil, err
	}
	requestMetaData.SetSourceMetaData(sourceMeta)
	return requestMetaData, nil
}

func (impl *ApprovalRequestServiceImpl) getMetadataObjectWithArtifactInfo(requestMetaData *bean.RequestMetaData, artifactId int) (*bean.RequestMetaData, error) {
	if artifactId > 0 {
		if requestMetaData == nil {
			requestMetaData = &bean.RequestMetaData{}
		}
		ciArtifact, err := impl.ciArtifactRepository.Get(artifactId)
		if err != nil {
			impl.logger.Errorw("error in finding the artifact using id", "artifactId", artifactId, "err", err)
			errorResp := util.NewApiError().WithHttpStatusCode(http.StatusInternalServerError).WithInternalMessage(fmt.Sprintf("error in finding artifact , err : %s", err.Error())).WithUserMessage("error in finding artifact")
			if errors.Is(err, pg.ErrNoRows) {
				errorResp = errorResp.WithHttpStatusCode(http.StatusConflict).WithUserMessage(constants.ArtifactNotFoundErr)
			}
			return nil, errorResp
		}
		requestMetaData = requestMetaData.WithCiArtifact(ciArtifact)
	}
	return requestMetaData, nil
}

func (impl *ApprovalRequestServiceImpl) getMetadataObjectWithDestinationInfo(requestMetaData *bean.RequestMetaData, request *bean.ArtifactPromotionRequest, authorizedEnvironments map[string]bool) (*bean.RequestMetaData, error) {
	if requestMetaData == nil {
		requestMetaData = &bean.RequestMetaData{}
	}
	// set environment metadata
	environments, err := impl.environmentService.FindByNames(request.EnvironmentNames)
	if err != nil {
		impl.logger.Errorw("error in fetching the environment details", "environmentNames", request.EnvironmentNames, "err", err)
		return nil, err
	}
	requestMetaData.SetActiveEnvironments(request.EnvironmentNames, authorizedEnvironments, environments)

	// set destination pipelines metadata
	cdPipelines, err := impl.cdPipelineService.FindActiveByAppIdAndEnvNames(request.AppId, requestMetaData.GetActiveAuthorisedEnvNames())
	if err != nil {
		impl.logger.Errorw("error in finding the cd pipelines using appID and environment names", "appId", request.AppId, "envNames", requestMetaData.GetActiveAuthorisedEnvNames(), "err", err)
		return nil, err
	}
	requestMetaData.SetDestinationPipelineMetaData(cdPipelines)
	return requestMetaData, nil
}

func (impl *ApprovalRequestServiceImpl) fetchSourceMeta(sourceName string, sourceType constants.SourceTypeStr, appId int, workflowId int) (*bean.SourceMetaData, error) {
	sourceInfo := &bean.SourceMetaData{}
	sourceInfo = sourceInfo.WithName(sourceName).WithType(sourceType)
	if sourceType == constants.SOURCE_TYPE_CD {
		// source type will be cd and source name will be envName.
		// get pipeline using appId and env name and get the workflowMapping
		pipelines, err := impl.cdPipelineService.FindActiveByAppIdAndEnvNames(appId, []string{sourceName})
		if err != nil {
			impl.logger.Errorw("error in getting the pipelines using appId and source environment name ", "workflowId", workflowId, "appId", appId, "source", sourceName, "err", err)
			return nil, err
		}
		if len(pipelines) == 0 {
			return nil, util.NewApiError().WithHttpStatusCode(http.StatusConflict).WithUserMessage(constants.WorkflowAndSourceMisMatchErr).WithInternalMessage(constants.WorkflowAndSourceMisMatchErr)
		}

		pipeline := pipelines[0]
		appWorkflowMapping, err := impl.appWorkflowService.FindWFMappingByComponent(appWorkflow.CDPIPELINE, pipeline.Id)
		if err != nil {
			impl.logger.Errorw("error in getting the app workflow mapping using workflow id and cd component id", "workflowId", workflowId, "appId", appId, "pipelineId", pipeline.Id, "err", err)
			if errors.Is(err, pg.ErrNoRows) {
				return nil, util.NewApiError().WithHttpStatusCode(http.StatusConflict).WithUserMessage(constants.WorkflowAndSourceMisMatchErr).WithInternalMessage(constants.WorkflowAndSourceMisMatchErr)
			}
			return nil, err
		}
		sourceInfo = sourceInfo.WithId(pipeline.Id).WithSourceWorkflowId(appWorkflowMapping.AppWorkflowId).WithCdPipeline(pipeline)
	} else {
		appWorkflowMapping, err := impl.appWorkflowService.FindByCiSourceWorkflowMappingById(workflowId)
		if err != nil {
			impl.logger.Errorw("error in getting the workflow mapping of ci-source/webhook using workflow id", "workflowId", workflowId, "err", err)
			if errors.Is(err, pg.ErrNoRows) {
				return nil, util.NewApiError().WithHttpStatusCode(http.StatusConflict).WithUserMessage(constants.WorkflowAndSourceMisMatchErr).WithInternalMessage(constants.WorkflowAndSourceMisMatchErr)
			}
			return nil, err
		}
		sourceInfo = sourceInfo.WithId(appWorkflowMapping.ComponentId).WithSourceWorkflowId(appWorkflowMapping.AppWorkflowId)
	}
	return sourceInfo, nil
}

func (impl *ApprovalRequestServiceImpl) validatePromoteAction(requestedWorkflowId int, metadata *bean.RequestMetaData) (map[string]bean.EnvironmentPromotionMetaData, error) {
	if requestedWorkflowId != metadata.GetWorkflowId() {
		return nil, util.NewApiError().WithHttpStatusCode(http.StatusConflict).WithUserMessage(constants.WorkflowAndSourceMisMatchErr).WithInternalMessage(constants.WorkflowAndSourceMisMatchErr)
	}

	allAppWorkflowMappings, err := impl.appWorkflowService.FindAppWorkflowMapping(metadata.GetWorkflowId())
	if err != nil {
		impl.logger.Errorw("error in finding the app workflow mappings", "workflowId", metadata.GetWorkflowId(), "err", err)
		return nil, err
	}

	respMap := crossWorkflowPromotionValidation(metadata, allAppWorkflowMappings)
	if metadata.GetSourceTypeStr() == constants.SOURCE_TYPE_CD {

		// if sourcePipelineId is 0, then the source pipeline given by user is not found in the workflow.
		if metadata.GetSourcePipelineId() == 0 {
			errMsg := fmt.Sprintf("no pipeline found against given source environment %s", metadata.GetSourceName())
			return nil, util.NewApiError().WithHttpStatusCode(http.StatusBadRequest).WithUserMessage(errMsg).WithInternalMessage(errMsg)
		}

		deployed, err := impl.checkIfDeployedAtSource(metadata.GetCiArtifactId(), metadata.GetSourceCdPipeline())
		if err != nil {
			impl.logger.Errorw("error in checking if artifact is available for promotion at source pipeline", "ciArtifactId", metadata.GetCiArtifactId(), "sourcePipelineId", metadata.GetSourcePipelineId(), "err", err)
			return nil, err
		}

		if !deployed {
			errMsg := fmt.Sprintf("artifact is not deployed on the source environment %s", metadata.GetSourceName())
			return nil, util.NewApiError().WithHttpStatusCode(http.StatusConflict).WithUserMessage(errMsg).WithInternalMessage(errMsg)
		}
		topologyValidationRespMap := runSourceAndDestinationTopologyValidations(metadata, allAppWorkflowMappings)
		for env, resp := range topologyValidationRespMap {
			if _, ok := respMap[env]; !ok {
				respMap[env] = resp
			}
		}
	} else {
		artifact := metadata.GetCiArtifact()
		switch metadata.GetSourceTypeStr() {
		case constants.SOURCE_TYPE_CI, constants.SOURCE_TYPE_LINKED_CI, constants.SOURCE_TYPE_LINKED_CD, constants.SOURCE_TYPE_JOB_CI:
			if artifact.PipelineId != metadata.GetSourcePipelineId() {
				return nil, util.NewApiError().WithHttpStatusCode(http.StatusConflict).WithUserMessage(constants.ArtifactSourceMisMatch).WithInternalMessage(constants.ArtifactSourceMisMatch)
			}
		case constants.SOURCE_TYPE_WEBHOOK:
			if artifact.ExternalCiPipelineId != metadata.GetSourcePipelineId() {
				return nil, util.NewApiError().WithHttpStatusCode(http.StatusConflict).WithUserMessage(constants.ArtifactSourceMisMatch).WithInternalMessage(constants.ArtifactSourceMisMatch)
			}
		}

	}
	return respMap, nil
}

func crossWorkflowPromotionValidation(metadata *bean.RequestMetaData, allAppWorkflowMappings []bean4.AppWorkflowMappingDto) map[string]bean.EnvironmentPromotionMetaData {
	responseMap := make(map[string]bean.EnvironmentPromotionMetaData)
	for _, pipelineId := range metadata.GetActiveAuthorisedPipelineIds() {
		// if pipeline is not found in the requested workflow, we should reject the request for the particular pipeline
		if !util3.Contains(allAppWorkflowMappings, func(appWorkflowMapping bean4.AppWorkflowMappingDto) bool {
			return appWorkflowMapping.Type == appWorkflow.CDPIPELINE && appWorkflowMapping.ComponentId == pipelineId
		}) {
			envName := metadata.GetActiveAuthorisedPipelineIdEnvMap()[pipelineId]
			resp := bean.EnvironmentPromotionMetaData{
				Name:                       envName,
				PromotionValidationMessage: constants.PIPELINE_NOT_FOUND,
			}
			cdPipeline := metadata.GetActiveAuthorisedPipelineDaoMap()[pipelineId]
			if cdPipeline != nil {
				resp.IsVirtualEnvironment = cdPipeline.Environment.IsVirtualEnvironment
			}
			responseMap[envName] = resp
		}
	}
	return responseMap
}

func runSourceAndDestinationTopologyValidations(metadata *bean.RequestMetaData, allAppWorkflowMappings []bean4.AppWorkflowMappingDto) map[string]bean.EnvironmentPromotionMetaData {
	tree := make(map[int][]int)
	for _, appWorkflowMapping := range allAppWorkflowMappings {
		// create the tree from the DAG excluding the ci source
		if appWorkflowMapping.Type == appWorkflow.CDPIPELINE && appWorkflowMapping.ParentType == appWorkflow.CDPIPELINE {
			tree[appWorkflowMapping.ParentId] = append(tree[appWorkflowMapping.ParentId], appWorkflowMapping.ComponentId)
		}
	}

	responseMap := make(map[string]bean.EnvironmentPromotionMetaData)
	for _, pipelineId := range metadata.GetActiveAuthorisedPipelineIds() {
		if !util.IsAncestor(tree, metadata.GetSourcePipelineId(), pipelineId) {
			envName := metadata.GetActiveAuthorisedPipelineIdEnvMap()[pipelineId]
			resp := bean.EnvironmentPromotionMetaData{
				Name:                       envName,
				PromotionValidationMessage: constants.SOURCE_AND_DESTINATION_PIPELINE_MISMATCH,
			}
			cdPipeline := metadata.GetActiveAuthorisedPipelineDaoMap()[pipelineId]
			if cdPipeline != nil {
				resp.IsVirtualEnvironment = cdPipeline.Environment.IsVirtualEnvironment
			}
			responseMap[envName] = resp
		}
	}
	return responseMap
}

func (impl *ApprovalRequestServiceImpl) getPromotionMetadata(request *bean.ArtifactPromotionRequest, authorizedEnvironments map[string]bool) (*bean.RequestMetaData, error) {
	metadata, err := impl.getMetadataObjectWithSourceInfo(nil, request)
	if err != nil {
		impl.logger.Errorw("error in getting source metadata for the request", "request", request, "err", err)
		return nil, err
	}

	metadata, err = impl.getMetadataObjectWithDestinationInfo(metadata, request, authorizedEnvironments)
	if err != nil {
		impl.logger.Errorw("error in getting destination environments/pipelines metadata for the request", "request", request, "authorizedEnvironments", authorizedEnvironments, "err", err)
		return nil, err
	}

	metadata, err = impl.getMetadataObjectWithArtifactInfo(metadata, request.ArtifactId)
	if err != nil {
		impl.logger.Errorw("error in getting artifact metadata", "artifactId", request.ArtifactId, "err", err)
		return nil, err
	}
	return metadata, nil
}

func (impl *ApprovalRequestServiceImpl) promoteArtifact(ctx *util3.RequestCtx, request *bean.ArtifactPromotionRequest, authorizedEnvironments map[string]bool) ([]bean.EnvironmentPromotionMetaData, error) {
	// 	step1: validate if artifact is deployed/created at the source pipeline.
	//      step1: if source is cd , check if this artifact is deployed on these environments
	//  step2: check if destination pipeline is topologically downwards from the source pipeline and also source and destination are on the same subtree.
	// 	step3: check if promotion request for this artifact on this destination pipeline has already been raised.
	//  step4: check if this artifact on this destination pipeline has already been promoted
	//  step5: raise request.

	// fetch artifact

	metadata, err := impl.getPromotionMetadata(request, authorizedEnvironments)
	if err != nil {
		impl.logger.Errorw("error in getting metadata info for promotion request", "request", request, "authorizedEnvironments", authorizedEnvironments, "err", err)
		return nil, err
	}

	responseMap := metadata.GetDefaultEnvironmentPromotionMetaDataResponseMap()

	validationResponseMap, err := impl.validatePromoteAction(request.WorkflowId, metadata)
	if err != nil {
		impl.logger.Errorw("error in validating the workflowPromotion request", "metadata", metadata, "err", err)
		return nil, err
	}

	for envName, resp := range validationResponseMap {
		responseMap[envName] = resp
	}

	policiesMap, err := impl.promotionPolicyDataReadService.GetPromotionPolicyByAppAndEnvIds(ctx, request.AppId, metadata.GetActiveAuthorisedEnvIds())
	if err != nil {
		impl.logger.Errorw("error in getting policies for some environments in an app", "appName", request.AppName, "envNames", metadata.GetActiveAuthorisedEnvNames(), "err", err)
		return nil, err
	}

	promotableEnvs := make([]string, 0, len(responseMap))
	for _, resp := range responseMap {
		if resp.NoValidationIssue() {
			promotableEnvs = append(promotableEnvs, resp.Name)
		}
	}

	if len(promotableEnvs) > 0 {
		promoteResponseMap, err := impl.raisePromoteRequestHelper(ctx, policiesMap, metadata.WithPromotableEnvs(promotableEnvs))
		if err != nil {
			impl.logger.Errorw("error in promoting the artifact on to destination pipelines", "artifactId", metadata.GetCiArtifactId(), "destinationPipelineIds", metadata.GetActiveAuthorisedPipelineIds(), "err", err)
			return nil, err
		}
		for envName, resp := range promoteResponseMap {
			responseMap[envName] = resp
		}
	}

	envResponses := make([]bean.EnvironmentPromotionMetaData, 0, len(responseMap))
	for _, resp := range responseMap {
		envResponses = append(envResponses, resp)
	}
	return envResponses, nil
}

func (impl *ApprovalRequestServiceImpl) raisePromoteRequestHelper(ctx *util3.RequestCtx, policiesMap map[string]*bean.PromotionPolicy, metadata *bean.RequestMetaData) (map[string]bean.EnvironmentPromotionMetaData, error) {
	responseMap := make(map[string]bean.EnvironmentPromotionMetaData)
	promotedCountPerPipeline, pendingCountPerPipeline, err := impl.fetchPendingAndPromotedRequests(metadata.GetPromotablePipelineIds(), metadata.GetCiArtifactId())
	if err != nil {
		impl.logger.Errorw("error in getting the pending and promoted requests using destination pipelines ids for an artifact", "destinationPipelineIds", metadata.GetActiveAuthorisedPipelineIds(), "artifactId", metadata.GetCiArtifactId(), "err", err)
		return nil, err
	}

	pipelineIdVsEnvNameMap := metadata.GetActiveAuthorisedPipelineIdEnvMap()
	pipelineIdToDaoMap := metadata.GetActiveAuthorisedPipelineDaoMap()

	for _, pipelineId := range metadata.GetPromotablePipelineIds() {
		policy := policiesMap[pipelineIdVsEnvNameMap[pipelineId]]
		cdPipeline := pipelineIdToDaoMap[pipelineId]

		envResp := impl.raiseRequestAndGetResponse(ctx, metadata, cdPipeline, policy, promotedCountPerPipeline, pendingCountPerPipeline)

		envResp.Name = pipelineIdVsEnvNameMap[pipelineId]
		envResp.IsVirtualEnvironment = pipelineIdToDaoMap[pipelineId].Environment.IsVirtualEnvironment
		responseMap[pipelineIdVsEnvNameMap[pipelineId]] = envResp
	}

	return responseMap, nil
}

func (impl *ApprovalRequestServiceImpl) raiseRequestAndGetResponse(ctx *util3.RequestCtx, metadata *bean.RequestMetaData, cdPipeline *pipelineConfig.Pipeline, promotionPolicy *bean.PromotionPolicy, promotedCountPerPipeline, pendingCountPerPipeline map[int]int) bean.EnvironmentPromotionMetaData {

	envResponse := bean.EnvironmentPromotionMetaData{}

	if metadata.GetSourceTypeStr() == constants.SOURCE_TYPE_CD && metadata.GetSourcePipelineId() == cdPipeline.Id {
		envResponse.PromotionValidationMessage = constants.SameSourceAndDestinationErr
		return envResponse
	}
	if promotedCountPerPipeline[cdPipeline.Id] > 0 {
		envResponse.PromotionValidationMessage = constants.ARTIFACT_ALREADY_PROMOTED
		return envResponse
	}

	if pendingCountPerPipeline[cdPipeline.Id] >= 1 {
		envResponse.PromotionValidationMessage = constants.ALREADY_REQUEST_RAISED
		return envResponse
	}

	policy := promotionPolicy
	if policy == nil {
		envResponse.PromotionValidationMessage = constants.POLICY_NOT_CONFIGURED
	} else {
		state, err := impl.raisePromoteRequest(ctx, policy, cdPipeline, metadata)
		if err != nil {
			impl.logger.Errorw("error in raising promotion request for the pipeline", "pipelineId", cdPipeline.Id, "artifactId", metadata.GetCiArtifactId(), "err", err)
			envResponse.PromotionValidationMessage = constants.ERRORED
		}
		envResponse.PromotionPossible = true
		envResponse.PromotionValidationMessage = state
	}
	return envResponse
}

func (impl *ApprovalRequestServiceImpl) fetchPendingAndPromotedRequests(destinationPipelineIds []int, artifactId int) (map[int]int, map[int]int, error) {
	requests, err := impl.artifactPromotionApprovalRequestRepository.FindRequestsByStatusesForDestinationPipelines(destinationPipelineIds, artifactId, []constants.ArtifactPromotionRequestStatus{constants.AWAITING_APPROVAL, constants.PROMOTED})
	if err != nil {
		impl.logger.Errorw("error in getting the pending and promoted requests using destination pipelines ids for an artifact", "destinationPipelineIds", destinationPipelineIds, "artifactId", artifactId, "err", err)
		return nil, nil, err
	}

	promotedCountPerPipeline := make(map[int]int)
	pendingCountPerPipeline := make(map[int]int)

	for _, request := range requests {
		if request.Status == constants.PROMOTED {
			promotedCountPerPipeline[request.DestinationPipelineId] = promotedCountPerPipeline[request.DestinationPipelineId] + 1
		}
		if request.Status == constants.AWAITING_APPROVAL {
			pendingCountPerPipeline[request.DestinationPipelineId] = pendingCountPerPipeline[request.DestinationPipelineId] + 1
		}
	}

	return promotedCountPerPipeline, pendingCountPerPipeline, nil

}

func (impl *ApprovalRequestServiceImpl) raisePromoteRequest(ctx *util3.RequestCtx, promotionPolicy *bean.PromotionPolicy, cdPipeline *pipelineConfig.Pipeline, metadata *bean.RequestMetaData) (constants.PromotionValidationMsg, error) {

	params, err := impl.computeFilterParams(metadata.GetCiArtifact())
	if err != nil {
		impl.logger.Errorw("error in finding the required CEL expression parameters for using ciArtifact", "err", err)
		return constants.POLICY_EVALUATION_ERRORED, err
	}

	evaluationResult, err := impl.resourceFilterConditionsEvaluator.EvaluateFilter(promotionPolicy.Conditions, expressionEvaluators.ExpressionMetadata{Params: params})
	if err != nil {
		impl.logger.Errorw("evaluation failed with error", "policyConditions", promotionPolicy.Conditions, "pipelineId", cdPipeline.Id, promotionPolicy.Conditions, "params", params, "err", err)
		return constants.POLICY_EVALUATION_ERRORED, err
	}

	if !evaluationResult {
		return constants.BLOCKED_BY_POLICY, nil
	}

	evaluationAuditJsonString, err := impl.evaluationJsonString(evaluationResult, promotionPolicy)
	if err != nil {
		return constants.ERRORED, err
	}

	tx, err := impl.transactionManager.StartTx()
	if err != nil {
		impl.logger.Errorw("error in starting the transaction", "evaluationResult", evaluationResult, "promotionPolicy", promotionPolicy, "err", err)
		return constants.ERRORED, err
	}
	defer impl.transactionManager.RollbackTx(tx)

	// save evaluation audit
	evaluationAuditEntry, err := impl.resourceFilterEvaluationAuditService.SaveFilterEvaluationAudit(tx, resourceFilter.Artifact, metadata.GetCiArtifactId(), cdPipeline.Id, resourceFilter.Pipeline, ctx.GetUserId(), evaluationAuditJsonString, resourceFilter.ARTIFACT_PROMOTION_POLICY)
	if err != nil {
		impl.logger.Errorw("error in saving policy evaluation audit data", "evaluationAuditEntry", evaluationAuditEntry, "err", err)
		return constants.ERRORED, err
	}
	promotionRequest := &repository.ArtifactPromotionApprovalRequest{
		SourceType:              metadata.GetSourceTypeStr().GetSourceType(),
		SourcePipelineId:        metadata.GetSourcePipelineId(),
		DestinationPipelineId:   cdPipeline.Id,
		Status:                  constants.AWAITING_APPROVAL,
		ArtifactId:              metadata.GetCiArtifactId(),
		PolicyId:                promotionPolicy.Id,
		PolicyEvaluationAuditId: evaluationAuditEntry.Id,
		AuditLog:                sql.NewDefaultAuditLog(ctx.GetUserId()),
	}

	status := constants.SENT_FOR_APPROVAL
	if promotionPolicy.CanBePromoted(0) {
		promotionRequest.Status = constants.PROMOTED
		status = constants.PROMOTION_SUCCESSFUL
	}
	_, err = impl.artifactPromotionApprovalRequestRepository.Create(tx, promotionRequest)
	if err != nil {
		impl.logger.Errorw("error in finding the pending promotion request using pipelineId and artifactId", "pipelineId", cdPipeline.Id, "artifactId", metadata.GetCiArtifactId())
		return constants.ERRORED, err
	}

	err = impl.transactionManager.CommitTx(tx)
	if err != nil {
		impl.logger.Errorw("error in committing the db transaction", "pipelineId", cdPipeline.Id, "artifactId", metadata.GetCiArtifactId(), "err", err)
		return constants.ERRORED, err
	}
	if promotionRequest.Status == constants.PROMOTED {
		triggerRequest := bean2.TriggerRequest{
			CdWf:        nil,
			Pipeline:    cdPipeline,
			Artifact:    metadata.GetCiArtifact(),
			TriggeredBy: 1,
			TriggerContext: bean2.TriggerContext{
				Context: context.Background(),
			},
		}
		err := impl.workflowDagExecutor.HandleArtifactPromotionEvent(triggerRequest)
		if err != nil {
			impl.logger.Errorw("auto trigger error on artifact promotion", "cdPipelineId", cdPipeline.Id, "artifactId", metadata.GetCiArtifactId(), "err", err)
		}
	} else if promotionRequest.Status == constants.AWAITING_APPROVAL {
		impl.logger.Infow("sending email notification for artifact promotion request", "cdPipelineId", cdPipeline.Id, "artifactId", metadata.GetCiArtifactId())
		go impl.sendPromotionRequestNotification(cdPipeline.Id, metadata, ctx.GetUserId())
	}
	return status, nil

}

func (impl *ApprovalRequestServiceImpl) sendPromotionRequestNotification(pipelineId int, metadata *bean.RequestMetaData, userId int32) {

	pipelineIdToDapMapping := metadata.GetActiveAuthorisedPipelineDaoMap()

	pipeline := pipelineIdToDapMapping[pipelineId]

	event := impl.eventFactory.Build(util2.ArtifactPromotionApproval, nil, pipeline.AppId, &pipeline.EnvironmentId, "")

	artifactPromotionNotificationRequest, err := impl.buildArtifactPromotionNotificationRequest(pipeline, metadata, userId)
	if err != nil {
		impl.logger.Errorw("error in building artifact promotion notification request", "pipelineId", pipeline.Id, "err", err)
		return
	}
	events := impl.eventFactory.BuildExtraArtifactPromotionData(event, artifactPromotionNotificationRequest)
	for _, evnt := range events {
		_, evtErr := impl.eventClient.WriteNotificationEvent(evnt)
		if evtErr != nil {
			impl.logger.Errorw("unable to send approval event", "error", evtErr)
		}
	}
}

func (impl *ApprovalRequestServiceImpl) buildArtifactPromotionNotificationRequest(pipeline *pipelineConfig.Pipeline, metadata *bean.RequestMetaData, userId int32) (client.ArtifactPromotionNotificationRequest, error) {

	team, err := impl.teamService.FetchOne(pipeline.App.TeamId)
	if err != nil {
		impl.logger.Errorw("error in fetching team by id", "teamId", pipeline.App.TeamId, "err", err)
		return client.ArtifactPromotionNotificationRequest{}, err
	}

	imagePromoterEmails, err := impl.userService.GetUsersByEnvAndAction(pipeline.App.AppName, pipeline.Environment.EnvironmentIdentifier, team.Name, bean5.ArtifactPromoter)
	if err != nil {
		impl.logger.Errorw("error in finding image promoter access emails", "appName", pipeline.App.AppName, "envName", pipeline.Environment.Name, "team", team.Name, "err", err)
		return client.ArtifactPromotionNotificationRequest{}, err
	}

	artifactId := metadata.GetCiArtifactId()

	imageComment, imageTagNames, err := impl.imageTaggingService.GetImageTagsAndComment(artifactId)
	if err != nil {
		impl.logger.Errorw("error in fetching tags and comment", "artifactId", artifactId, "err", err)
		return client.ArtifactPromotionNotificationRequest{}, err
	}
	metadata = metadata.WithImageComment(imageComment.Comment).WithImageTags(imageTagNames)

	artifactPromotionNotificationRequest := parseArtifactPromotionRequest(pipeline, metadata, imagePromoterEmails, userId)

	return artifactPromotionNotificationRequest, nil
}

func parseArtifactPromotionRequest(pipeline *pipelineConfig.Pipeline, metadata *bean.RequestMetaData, imagePromoterEmails []string, userId int32) client.ArtifactPromotionNotificationRequest {
	artifactPromotionNotificationRequest := client.ArtifactPromotionNotificationRequest{
		CDPipelineId:            pipeline.Id,
		AppId:                   pipeline.AppId,
		AppName:                 pipeline.App.AppName,
		EnvId:                   pipeline.EnvironmentId,
		EnvName:                 pipeline.Environment.Name,
		ArtifactId:              metadata.GetCiArtifactId(),
		UserId:                  userId,
		ImageTags:               metadata.GetImageTags(),
		ImageComment:            metadata.GetImageComment(),
		ImagePath:               metadata.GetCiArtifact().Image,
		ArtifactPromotionSource: metadata.GetSourceDisplayName(),
		PromoterAccessEmailIds:  imagePromoterEmails,
		WorkflowId:              metadata.GetWorkflowId(),
	}
	return artifactPromotionNotificationRequest
}

func (impl *ApprovalRequestServiceImpl) evaluationJsonString(evaluationResult bool, promotionPolicy *bean.PromotionPolicy) (string, error) {
	state := expressionEvaluators.ALLOW
	if !evaluationResult {
		state = expressionEvaluators.BLOCK
	}

	histories, err := impl.promotionPolicyDataReadService.GetPolicyHistoryIdsByPolicyIds([]int{promotionPolicy.Id})
	if err != nil {
		impl.logger.Errorw("error in fetching the latest promotion policy audit history", "err", err, "policyId", promotionPolicy.Id)
		return "", err
	}
	if len(histories) == 0 {
		return "", errors.New(fmt.Sprintf("no audit history found for the policy : %s", promotionPolicy.Name))
	}

	evaluationAudit := resourceFilter.FilterHistoryObject{
		State:           state,
		FilterHistoryId: histories[0],
	}
	evaluationAuditJsonBytes, err := json.Marshal(&evaluationAudit)
	if err != nil {
		return "", err
	}
	return string(evaluationAuditJsonBytes), nil
}

func (impl *ApprovalRequestServiceImpl) checkIfDeployedAtSource(ciArtifactId int, pipeline *pipelineConfig.Pipeline) (bool, error) {
	if pipeline == nil {
		return false, errors.New("invalid cd pipeline")
	}
	postStage, err := impl.pipelineStageService.GetCdStageByCdPipelineIdAndStageType(pipeline.Id, repository3.PIPELINE_STAGE_TYPE_POST_CD)
	if err != nil && !errors.Is(err, pg.ErrNoRows) {
		impl.logger.Errorw("error in finding the post-cd existence for the pipeline", "pipelineId", pipeline.Id, "err", err)
		return false, err
	}
	workflowType := bean3.CD_WORKFLOW_TYPE_DEPLOY
	if len(pipeline.PostStageConfig) > 0 || (postStage != nil && postStage.Id > 0) {
		workflowType = bean3.CD_WORKFLOW_TYPE_POST
	}

	deployed, err := impl.cdWorkflowService.IsArtifactDeployedOnStage(ciArtifactId, pipeline.Id, workflowType)
	if err != nil {
		impl.logger.Errorw("error in finding if the artifact is successfully deployed on a pipeline", "ciArtifactId", ciArtifactId, "pipelineId", pipeline.Id, "workflowType", workflowType, "err", err)
		return false, err
	}
	return deployed, nil
}

func (impl *ApprovalRequestServiceImpl) cancelPromotionApprovalRequest(ctx *util3.RequestCtx, request *bean.ArtifactPromotionRequest) (*bean.ArtifactPromotionRequest, error) {
	rowsUpdated, err := impl.artifactPromotionApprovalRequestRepository.MarkCancel(request.PromotionRequestId, ctx.GetUserId())
	if err != nil {
		impl.logger.Errorw("error in canceling promotion approval request for given id", "promotionRequestId", request.PromotionRequestId, "err", err)
		return nil, err
	}
	if err != nil {
		impl.logger.Errorw("error in fetching artifact promotion request by id", "artifactPromotionRequestId", request.PromotionRequestId, "err", err)
		return nil, err
	}
	if rowsUpdated == 0 {
		return nil, util.NewApiError().WithHttpStatusCode(http.StatusUnprocessableEntity).WithInternalMessage(constants.UserCannotCancelRequest).WithUserMessage(constants.UserCannotCancelRequest)
	}
	return nil, nil
}

func (impl *ApprovalRequestServiceImpl) getRbacObjects(pipelineIdToDaoMapping map[int]*pipelineConfig.Pipeline) ([]string, map[int]string) {
	rbacObjects := make([]string, len(pipelineIdToDaoMapping))
	pipelineIdToRbacObjMap := make(map[int]string)
	for _, pipelineDao := range pipelineIdToDaoMapping {
		teamRbacObj := fmt.Sprintf("%s/%s/%s", pipelineDao.App.Team.Name, pipelineDao.Environment.EnvironmentIdentifier, pipelineDao.App.AppName)
		rbacObjects = append(rbacObjects, teamRbacObj)
		pipelineIdToRbacObjMap[pipelineDao.Id] = teamRbacObj
	}
	return rbacObjects, pipelineIdToRbacObjMap
}

func (impl *ApprovalRequestServiceImpl) getPipelineIdToDaoMapping(destinationPipelineIds []int) (map[int]*pipelineConfig.Pipeline, error) {
	pipelines, err := impl.cdPipelineService.FindAppAndEnvironmentAndProjectByPipelineIds(destinationPipelineIds)
	if err != nil {
		impl.logger.Errorw("error in fetching pipelines by ids", "pipelineIds", destinationPipelineIds, "err", err)
		return nil, err
	}
	pipelineIdToDaoMapping := make(map[int]*pipelineConfig.Pipeline)
	for _, pipelineDao := range pipelines {
		pipelineIdToDaoMapping[pipelineDao.Id] = pipelineDao
	}
	return pipelineIdToDaoMapping, err
}

func (impl *ApprovalRequestServiceImpl) onPolicyDelete(tx *pg.Tx, policyId int) error {
	err := impl.artifactPromotionApprovalRequestRepository.MarkStaleByPolicyId(tx, policyId)
	if err != nil {
		impl.logger.Errorw("error in marking artifact promotion requests stale", "policyId", policyId, "err", err)
	}
	return err
}

func (impl *ApprovalRequestServiceImpl) onPolicyUpdate(tx *pg.Tx, policy *bean.PromotionPolicy) error {
	// get all the requests whose id is policy.id
	existingRequests, err := impl.artifactPromotionApprovalRequestRepository.FindAwaitedRequestByPolicyId(policy.Id)
	if err != nil {
		impl.logger.Errorw("error in fetching the awaiting artifact promotion requests using policy Id", "policyId", policy.Id, "err", err)
		return err
	}
	artifactIds := make([]int, 0, len(existingRequests))
	for _, request := range existingRequests {
		artifactIds = append(artifactIds, request.ArtifactId)
	}
	if len(artifactIds) == 0 {
		impl.logger.Debugw("no awaiting requests found for the policy", "policyId", policy.Id)
		return nil
	}

	artifacts, err := impl.ciArtifactRepository.GetByIds(artifactIds)
	if err != nil {
		impl.logger.Errorw("error in fetching the artifacts by ids", "artifactIds", artifactIds, "err", err)
		return err
	}

	artifactsMap := make(map[int]*repository2.CiArtifact)
	for _, artifact := range artifacts {
		artifactsMap[artifact.Id] = artifact
	}

	// get all the artifacts using request.artifactId
	// re-evaluate the artifacts using the policy

	requestsToBeUpdatedAsStaled, err := impl.reEvaluatePolicyAndUpdateRequests(tx, policy, artifactsMap, existingRequests)
	if err != nil {
		impl.logger.Errorw("error in re evaluating artifact promotion requests ", "existingRequests", existingRequests, "err", err)
		return err
	}

	err = impl.artifactPromotionApprovalRequestRepository.UpdateInBulk(tx, requestsToBeUpdatedAsStaled)
	if err != nil {
		impl.logger.Errorw("error in marking artifact promotion requests stale", "policyId", policy.Id, "err", err)
		return err
	}

	unStaledRequestsIds := make([]int, 0, len(existingRequests))
	for _, existingRequest := range existingRequests {
		found := false
		for _, staleRequest := range requestsToBeUpdatedAsStaled {
			if staleRequest.Id == existingRequest.Id {
				found = true
			}
		}
		if !found {
			unStaledRequestsIds = append(unStaledRequestsIds, existingRequest.Id)
		}
	}

	approvbleRequestIds := make([]int, 0, len(unStaledRequestsIds))
	approvedUserData, err := impl.requestApprovalUserdataRepo.FetchApprovalDataForRequests(unStaledRequestsIds, models.ARTIFACT_PROMOTION_APPROVAL)
	if err != nil {
		impl.logger.Errorw("error in fetching approved user data for some artifact promotion requestIds", "unStaledRequestsIds", unStaledRequestsIds, "err", err)
		return err
	}
	approverCountForRequests := make(map[int]int)
	for _, unStaledRequestsId := range unStaledRequestsIds {
		approverCountForRequests[unStaledRequestsId] = 0
	}

	for _, userData := range approvedUserData {
		count := approverCountForRequests[userData.ApprovalRequestId]
		approverCountForRequests[userData.ApprovalRequestId] = count + 1
	}

	for _, unStaledRequestsId := range unStaledRequestsIds {
		if approverCountForRequests[unStaledRequestsId] >= policy.ApprovalMetaData.ApprovalCount {
			approvbleRequestIds = append(approvbleRequestIds, unStaledRequestsId)
		}
	}

	err = impl.artifactPromotionApprovalRequestRepository.MarkPromoted(tx, approvbleRequestIds, 1)
	if err != nil {
		impl.logger.Errorw("error in marking status of artifact promotion requests to approved as these requests already got approvals that the updated policy count", "policyId", policy.Id, "newApproverCount", policy.ApprovalMetaData.ApprovalCount, "requestIds", approvbleRequestIds, "err", err)
	}
	return err
}

func (impl *ApprovalRequestServiceImpl) reEvaluatePolicyAndUpdateRequests(tx *pg.Tx, policy *bean.PromotionPolicy, artifactsMap map[int]*repository2.CiArtifact, existingRequests []*repository.ArtifactPromotionApprovalRequest) ([]*repository.ArtifactPromotionApprovalRequest, error) {
	requestsToBeUpdatedAsStaled := make([]*repository.ArtifactPromotionApprovalRequest, 0, len(existingRequests))
	for _, request := range existingRequests {
		artifact := artifactsMap[request.ArtifactId]
		params, err := impl.computeFilterParams(artifact)
		if err != nil {
			continue
		}

		evaluationResult, err := impl.resourceFilterConditionsEvaluator.EvaluateFilter(policy.Conditions, expressionEvaluators.ExpressionMetadata{Params: params})
		if err != nil {
			impl.logger.Errorw("evaluation failed with error", "policyConditions", policy.Conditions, "pipelineId", request.DestinationPipelineId, "policyConditions", policy.Conditions, "params", params, "err", err)
			continue
		}

		// policy is blocking the request, so need to update these as staled requests
		if !evaluationResult {
			evaluationAuditJsonString, err := impl.evaluationJsonString(evaluationResult, policy)
			if err != nil {
				impl.logger.Errorw("error in creating evaluation audit json for a policy ", "policy", policy, "err", err)
				continue
			}

			// save evaluation audit
			evaluationAuditEntry, err := impl.resourceFilterEvaluationAuditService.SaveFilterEvaluationAudit(tx, resourceFilter.Artifact, request.ArtifactId, request.DestinationPipelineId, resourceFilter.Pipeline, 1, evaluationAuditJsonString, resourceFilter.ARTIFACT_PROMOTION_POLICY)
			if err != nil {
				impl.logger.Errorw("error in saving policy evaluation audit data", "evaluationAuditEntry", evaluationAuditEntry, "err", err)
				continue
			}
			request.UpdatedOn = time.Now()
			request.Status = constants.STALE
			request.PolicyEvaluationAuditId = evaluationAuditEntry.Id
			requestsToBeUpdatedAsStaled = append(requestsToBeUpdatedAsStaled, request)
		}

	}

	return requestsToBeUpdatedAsStaled, nil

}

func (impl *ApprovalRequestServiceImpl) onApplyPolicy(tx *pg.Tx, commaSeperatedAppEnvIds [][]int) error {
	return impl.artifactPromotionApprovalRequestRepository.MarkStaleByAppEnvIds(tx, commaSeperatedAppEnvIds)
}
