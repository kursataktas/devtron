package repository

import (
	"github.com/glebarez/sqlite"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type EphemeralContainerEntity struct {
	Id                  int
	Name                string
	ClusterId           int
	Namespace           string
	PodName             string
	TargetContainer     string
	Config              string
	IsExternallyCreated bool
}

type EphemeralContainerActionEntity struct {
	Id                   int
	EphemeralContainerId int
	ActionType           ContainerAction
	PerformedBy          int32
	PerformedAt          time.Time
}

type EphemeralContainerFileBasedRepositoryImpl struct {
	logger       *zap.SugaredLogger
	dbConnection *gorm.DB
}

func NewEphemeralContainerFileBasedRepository(logger *zap.SugaredLogger) *EphemeralContainerFileBasedRepositoryImpl {
	err, dbPath := createOrCheckClusterDbPath(logger)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	//db, err := sql.Open("sqlite3", "./cluster.db")
	if err != nil {
		logger.Fatal("error occurred while opening db connection", "error", err)
	}
	migrator := db.Migrator()
	ephemeralContainerEntity := &EphemeralContainerEntity{}
	hasTable := migrator.HasTable(ephemeralContainerEntity)
	if !hasTable {
		err = migrator.CreateTable(ephemeralContainerEntity)
		if err != nil {
			logger.Fatal("error occurred while creating ephemeral container table", "error", err)
		}
	}
	ephemeralContainerActionEntity := &EphemeralContainerActionEntity{}
	hasTable = migrator.HasTable(ephemeralContainerActionEntity)
	if !hasTable {
		err = migrator.CreateTable(ephemeralContainerActionEntity)
		if err != nil {
			logger.Fatal("error occurred while creating ephemeral container action table", "error", err)
		}
	}
	logger.Debugw("ephemeralContainer repository file based initialized")
	return &EphemeralContainerFileBasedRepositoryImpl{logger, db}
}

func (impl EphemeralContainerFileBasedRepositoryImpl) StartTx() (*pg.Tx, error) {
	return nil, nil
}

func (impl EphemeralContainerFileBasedRepositoryImpl) RollbackTx(tx *pg.Tx) error {
	return nil
}

func (impl EphemeralContainerFileBasedRepositoryImpl) CommitTx(tx *pg.Tx) error {
	return nil
}

func (impl EphemeralContainerFileBasedRepositoryImpl) SaveEphemeralContainerData(tx *pg.Tx, model *EphemeralContainerBean) error {
	//containerEntity := impl.convertToEntity(model)
	result := impl.dbConnection.Create(model)
	return result.Error
}

func (impl EphemeralContainerFileBasedRepositoryImpl) SaveEphemeralContainerActionAudit(tx *pg.Tx, model *EphemeralContainerAction) error {
	//auditEntity := impl.convertToAuditEntity(model)
	result := impl.dbConnection.Create(model)
	return result.Error
}

func (impl EphemeralContainerFileBasedRepositoryImpl) FindContainerByName(clusterID int, namespace, podName, name string) (*EphemeralContainerBean, error) {
	ephemeralContainerEntity := &EphemeralContainerBean{}
	result := impl.dbConnection.
		Where("cluster_id = ?", clusterID).
		Where("namespace = ?", namespace).
		Where("pod_name = ?", podName).
		Where("name = ?", name).
		Find(ephemeralContainerEntity)
	err := result.Error
	if err != nil {
		impl.logger.Errorw("error occurred while finding ephemeral container data ", "cluster_id", clusterID, "err", err)
		return nil, errors.New("failed to fetch ephemeral container")
	}
	//model := impl.convertToModel(ephemeralContainerEntity)
	return ephemeralContainerEntity, nil
}

func (impl EphemeralContainerFileBasedRepositoryImpl) convertToEntity(ephemeralContainerBean *EphemeralContainerBean) *EphemeralContainerEntity {
	entity := &EphemeralContainerEntity{
		Id:                  ephemeralContainerBean.Id,
		Name:                ephemeralContainerBean.Name,
		ClusterId:           ephemeralContainerBean.ClusterId,
		Namespace:           ephemeralContainerBean.Namespace,
		PodName:             ephemeralContainerBean.PodName,
		TargetContainer:     ephemeralContainerBean.TargetContainer,
		Config:              ephemeralContainerBean.Config,
		IsExternallyCreated: ephemeralContainerBean.IsExternallyCreated,
	}
	return entity
}
