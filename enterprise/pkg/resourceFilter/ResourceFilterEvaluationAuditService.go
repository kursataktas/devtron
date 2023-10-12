package resourceFilter

import (
	"fmt"
	"github.com/devtron-labs/devtron/internal/sql/repository/helper"
	"github.com/devtron-labs/devtron/pkg/sql"
	"go.uber.org/zap"
	"time"
)

type FilterHistoryObject struct {
	FilterHistoryId int         `json:"filter_history_id"`
	State           FilterState `json:"state"`
	Message         string      `json:"message"`
}

type FilterEvaluationAuditService interface {
	CreateFilterEvaluation(subjectType SubjectType, subjectIds []int, refType ReferenceType, refId int, filters []*FilterMetaDataBean, filterIdVsState map[int]FilterState) (*ResourceFilterEvaluationAudit, error)
	UpdateFilterEvaluationAuditRef(id int, refType ReferenceType, refId int) error
	GetFilterEvaluationAudits()
}

type FilterEvaluationAuditServiceImpl struct {
	logger                    *zap.SugaredLogger
	filterEvaluationAuditRepo FilterEvaluationAuditRepository
	filterAuditRepo           FilterAuditRepository
}

func NewFilterEvaluationAuditServiceImpl(logger *zap.SugaredLogger,
	filterEvaluationAuditRepo FilterEvaluationAuditRepository,
	filterAuditRepo FilterAuditRepository) *FilterEvaluationAuditServiceImpl {
	return &FilterEvaluationAuditServiceImpl{
		logger:                    logger,
		filterEvaluationAuditRepo: filterEvaluationAuditRepo,
		filterAuditRepo:           filterAuditRepo,
	}
}

func (impl *FilterEvaluationAuditServiceImpl) CreateFilterEvaluation(subjectType SubjectType, subjectIds []int, refType ReferenceType, refId int, filters []*FilterMetaDataBean, filterIdVsState map[int]FilterState) (*ResourceFilterEvaluationAudit, error) {
	filterHistoryObjectsStr, err := impl.extractFilterHistoryObjects(filters, filterIdVsState)
	if err != nil {
		impl.logger.Errorw("error in extracting filter history objects", "err", err, "filters", filters, "filterIdVsState", filterIdVsState)
		return nil, err
	}
	currentTime := time.Now()
	filterEvaluationAudit := &ResourceFilterEvaluationAudit{
		SubjectType:   &subjectType,
		SubjectIds:    helper.GetCommaSepratedString(subjectIds),
		ReferenceType: &refType,
		ReferenceId:   refId,
		AuditLog: sql.AuditLog{
			CreatedOn: currentTime,
			UpdatedOn: currentTime,
			//TODO: created or updated by
		},
		FilterHistoryObjects: filterHistoryObjectsStr,
	}
	filterEvaluationAudit, err = impl.filterEvaluationAuditRepo.Create(filterEvaluationAudit)
	if err != nil {
		impl.logger.Errorw("error in saving resource filter evaluation result in resource_filter_evaluation_audit table", "err", err, "filterEvaluationAudit", filterEvaluationAudit)
		return filterEvaluationAudit, err
	}
	return filterEvaluationAudit, nil
}

func (impl *FilterEvaluationAuditServiceImpl) UpdateFilterEvaluationAuditRef(id int, refType ReferenceType, refId int) error {
	return impl.filterEvaluationAuditRepo.UpdateRefTypeAndRefId(id, refType, refId)
}

func (impl *FilterEvaluationAuditServiceImpl) GetFilterEvaluationAudits() {

}

func (impl *FilterEvaluationAuditServiceImpl) extractFilterHistoryObjects(filters []*FilterMetaDataBean, filterIdVsState map[int]FilterState) (string, error) {
	filterIds := make([]int, 0)
	filterHistoryObjectMap := make(map[int]*FilterHistoryObject)
	for _, filter := range filters {
		filterIds = append(filterIds, filter.Id)
		message := ""
		for _, condition := range filter.Conditions {
			message = fmt.Sprintf("\n%s conditionType : %v , errorMsg : %v", message, condition.ConditionType, condition.ErrorMsg)
		}
		filterHistoryObjectMap[filter.Id] = &FilterHistoryObject{
			State:   filterIdVsState[filter.Id],
			Message: message,
		}
	}

	resourceFilterEvaluationAudits, err := impl.filterAuditRepo.GetLatestResourceFilterAuditByFilterIds(filterIds)
	if err != nil {
		impl.logger.Errorw("error in getting latest resource filter audits for given filter id's", "filterIds", filterIds, "err", err)
		return "", err
	}

	for _, resourceFilterEvaluationAudit := range resourceFilterEvaluationAudits {
		filterHistoryObject := filterHistoryObjectMap[resourceFilterEvaluationAudit.FilterId]
		if filterHistoryObject != nil {
			filterHistoryObject.FilterHistoryId = resourceFilterEvaluationAudit.Id
		}
	}

	filterHistoryObjects := make([]*FilterHistoryObject, 0, len(filterHistoryObjectMap))
	for _, val := range filterHistoryObjectMap {
		filterHistoryObjects = append(filterHistoryObjects, val)
	}
	jsonStr, err := getJsonStringFromFilterHistoryObjects(filterHistoryObjects)
	if err != nil {
		impl.logger.Errorw("error in getting json string for filter history objects", "filterHistoryObjects", filterHistoryObjects, "err", err)
		return "", err
	}
	return jsonStr, err

}
