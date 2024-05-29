/*
 * Copyright (c) 2024. Devtron Inc.
 */

package cd

import (
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	"github.com/devtron-labs/devtron/pkg/workflow/cd/adapter"
	"github.com/devtron-labs/devtron/pkg/workflow/cd/bean"
	"github.com/devtron-labs/devtron/pkg/workflow/cd/util"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type CdWorkflowRunnerService interface {
	FindWorkflowRunnerById(wfrId int) (*bean.CdWorkflowRunnerDto, error)
	CheckIfWfrLatest(wfrId, pipelineId int) (isLatest bool, err error)
	CreateBulkCdWorkflowRunners(tx *pg.Tx, cdWorkflowRunnerDtos []*bean.CdWorkflowRunnerDto) (map[int]int, error)
}

type CdWorkflowRunnerServiceImpl struct {
	logger               *zap.SugaredLogger
	cdWorkflowRepository pipelineConfig.CdWorkflowRepository
}

func NewCdWorkflowRunnerServiceImpl(logger *zap.SugaredLogger,
	cdWorkflowRepository pipelineConfig.CdWorkflowRepository) *CdWorkflowRunnerServiceImpl {
	return &CdWorkflowRunnerServiceImpl{
		logger:               logger,
		cdWorkflowRepository: cdWorkflowRepository,
	}
}

func (impl *CdWorkflowRunnerServiceImpl) FindWorkflowRunnerById(wfrId int) (*bean.CdWorkflowRunnerDto, error) {
	cdWfr, err := impl.cdWorkflowRepository.FindWorkflowRunnerById(wfrId)
	if err != nil {
		impl.logger.Errorw("error in getting cd workflow runner by id", "err", err, "id", wfrId)
		return nil, err
	}
	return adapter.ConvertCdWorkflowRunnerDbObjToDto(cdWfr), nil

}

func (impl *CdWorkflowRunnerServiceImpl) CheckIfWfrLatest(wfrId, pipelineId int) (isLatest bool, err error) {
	isLatest, err = impl.cdWorkflowRepository.IsLatestCDWfr(wfrId, pipelineId)
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("err in checking latest cd workflow runner", "err", err)
		return false, err
	}
	return isLatest, nil
}

func (impl *CdWorkflowRunnerServiceImpl) CreateBulkCdWorkflowRunners(tx *pg.Tx, cdWorkflowRunnerDtos []*bean.CdWorkflowRunnerDto) (map[int]int, error) {
	cdWorkFlowRunners := make([]*pipelineConfig.CdWorkflowRunner, 0, len(cdWorkflowRunnerDtos))
	for _, dto := range cdWorkflowRunnerDtos {
		cdWorkFlowRunners = append(cdWorkFlowRunners, adapter.ConvertCdWorkflowRunnerDtoToDbObj(dto))
	}
	err := impl.cdWorkflowRepository.BulkSaveWorkflowRunners(tx, cdWorkFlowRunners)
	if err != nil {
		impl.logger.Errorw("error encountered in CreateBulkCdWorkflowRunners", "cdWorkFlowRunners", cdWorkFlowRunners, "err", err)
		return nil, err
	}
	return util.GetCdWorkflowIdVsRunnerIdMap(cdWorkFlowRunners), nil
}
