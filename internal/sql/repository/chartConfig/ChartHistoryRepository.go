package chartConfig

import (
	"github.com/devtron-labs/devtron/internal/sql/models"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"time"
)

type ChartsGlobalHistory struct {
	tableName               struct{}  `sql:"charts_global_history" pg:",discard_unknown_columns"`
	Id                      int       `sql:"id,pk"`
	ChartsId                int       `sql:"charts_id"`
	Values                  string    `sql:"values_yaml"`
	GlobalOverride          string    `sql:"global_override"`
	ReleaseOverride         string    `sql:"release_override"`
	PipelineOverride        string    `sql:"pipeline_override"`
	ImageDescriptorTemplate string    `sql:"image_descriptor_template"`
	ChartRefId              int       `sql:"chart_ref_id"`
	Latest                  bool      `sql:"latest,notnull"`
	//TODO : confirm if deployment details are needed here, since for every environment we will have
	// a separate entry in env history table
	Deployed                bool      `sql:"deployed"`
	DeployedOn              time.Time `sql:"deployed_on"`
	DeployedBy              int32     `sql:"deployed_by"`
	models.AuditLog
}

type ChartHistoryRepository interface {
	CreateGlobalHistory(chart *ChartsGlobalHistory) (*ChartsGlobalHistory, error)
	UpdateGlobalHistory(chart *ChartsGlobalHistory) (*ChartsGlobalHistory, error)

	CreateEnvHistory(chart *ChartsEnvHistory) (*ChartsEnvHistory, error)
	UpdateEnvHistory(chart *ChartsEnvHistory) (*ChartsEnvHistory, error)
}

type ChartHistoryRepositoryImpl struct {
	dbConnection *pg.DB
	logger       *zap.SugaredLogger
}

func NewChartHistoryRepositoryImpl(logger *zap.SugaredLogger, dbConnection *pg.DB) *ConfigMapHistoryRepositoryImpl {
	return &ConfigMapHistoryRepositoryImpl{dbConnection: dbConnection, logger: logger}
}

func(impl ChartHistoryRepositoryImpl) CreateGlobalHistory(chart *ChartsGlobalHistory) (*ChartsGlobalHistory, error){
	err := impl.dbConnection.Insert(chart)
	if err != nil {
		impl.logger.Errorw("err in creating global chart history entry", "err", err)
		return chart, err
	}
	return chart, nil
}

func(impl ChartHistoryRepositoryImpl) UpdateGlobalHistory(chart *ChartsGlobalHistory) (*ChartsGlobalHistory, error){
	err := impl.dbConnection.Update(chart)
	if err != nil {
		impl.logger.Errorw("err in updating global chart history entry", "err", err)
		return chart, err
	}
	return chart, nil
}




//---------------------------------------------------

type ChartsEnvHistory struct {
	tableName               struct{}  `sql:"charts_env_history" pg:",discard_unknown_columns"`
	Id                      int       `sql:"id,pk"`
	TargetEnvironment 		int       `sql:"target_environment"`
	EnvOverride          	string    `sql:"env_override"`
	Latest                  bool      `sql:"latest,notnull"`
	Deployed                bool      `sql:"deployed"`
	DeployedOn              time.Time `sql:"deployed_on"`
	DeployedBy              int32     `sql:"deployed_by"`
	models.AuditLog
}


func(impl ChartHistoryRepositoryImpl) CreateEnvHistory(chart *ChartsEnvHistory) (*ChartsEnvHistory, error){
	err := impl.dbConnection.Insert(chart)
	if err != nil {
		impl.logger.Errorw("err in creating env chart history entry", "err", err)
		return chart, err
	}
	return chart, nil
}
func(impl ChartHistoryRepositoryImpl) UpdateEnvHistory(chart *ChartsEnvHistory) (*ChartsEnvHistory, error){
	err := impl.dbConnection.Update(chart)
	if err != nil {
		impl.logger.Errorw("err in updating env chart history entry", "err", err)
		return chart, err
	}
	return chart, nil
}