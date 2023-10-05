package resourceFilter

import (
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

// FilterTargetObject represents the targeted areas where filters are evaluated
type FilterTargetObject int

const (
	DeploymentPipeline FilterTargetObject = 0
	BuildPipeline      FilterTargetObject = 1
)

type ResourceConditionType int

const (
	FAIL ResourceConditionType = iota
	PASS
)

type ResourceFilter struct {
	tableName           struct{}            `sql:"resource_filter" pg:",discard_unknown_columns"`
	Id                  int                 `sql:"id"`
	Name                string              `sql:"name"`
	Description         string              `sql:"description"`
	TargetObject        *FilterTargetObject `sql:"target_object"`
	ConditionExpression string              `sql:"condition_expression"`
	Deleted             *bool               `sql:"deleted"`
	sql.AuditLog
}

func (o *ResourceFilter) GetMetaData() {

}

type ResourceFilterRepository interface {
	//transaction util funcs
	sql.TransactionWrapper
	GetConnection() *pg.DB
	CreateResourceFilter(tx *pg.Tx, filter *ResourceFilter) (*ResourceFilter, error)
	UpdateFilter(tx *pg.Tx, filter *ResourceFilter) error
	ListAll() ([]*ResourceFilter, error)
	GetById(id int) (*ResourceFilter, error)
	GetByIds(ids []int) ([]*ResourceFilter, error)
	GetDistinctFilterNames() ([]string, error)
}

type ResourceFilterRepositoryImpl struct {
	logger       *zap.SugaredLogger
	dbConnection *pg.DB
	*sql.TransactionUtilImpl
}

func NewResourceFilterRepositoryImpl(logger *zap.SugaredLogger,
	dbConnection *pg.DB) *ResourceFilterRepositoryImpl {
	return &ResourceFilterRepositoryImpl{
		logger:              logger,
		dbConnection:        dbConnection,
		TransactionUtilImpl: sql.NewTransactionUtilImpl(dbConnection),
	}
}

func (repo *ResourceFilterRepositoryImpl) GetConnection() *pg.DB {
	return repo.dbConnection
}
func (repo *ResourceFilterRepositoryImpl) CreateResourceFilter(tx *pg.Tx, filter *ResourceFilter) (*ResourceFilter, error) {
	err := tx.Insert(filter)
	return filter, err
}

func (repo *ResourceFilterRepositoryImpl) UpdateFilter(tx *pg.Tx, filter *ResourceFilter) error {
	err := tx.Update(filter)
	return err
}

func (repo *ResourceFilterRepositoryImpl) GetById(id int) (*ResourceFilter, error) {
	filter := &ResourceFilter{}
	err := repo.dbConnection.Model(filter).Where("id = ? and deleted = false", id).Select()
	return filter, err
}

func (repo *ResourceFilterRepositoryImpl) GetByIds(ids []int) ([]*ResourceFilter, error) {
	var resourceFilters []*ResourceFilter
	if len(ids) == 0 {
		return resourceFilters, nil
	}
	err := repo.dbConnection.Model(&resourceFilters).
		Where("id IN (?)", pg.In(ids)).
		Where("deleted = ?", false).
		Select()
	if err != nil {
		repo.logger.Errorw("error occurred while fetching filter", "ids", ids, "err", err)
		if err == pg.ErrNoRows {
			err = nil
		}
	}
	return resourceFilters, err
}

func (repo *ResourceFilterRepositoryImpl) ListAll() ([]*ResourceFilter, error) {
	list := make([]*ResourceFilter, 0)
	err := repo.dbConnection.
		Model(&list).
		Where("deleted=?", false).
		Order("name").
		Select()
	return list, err
}

func (repo *ResourceFilterRepositoryImpl) GetDistinctFilterNames() ([]string, error) {
	list := make([]string, 0)
	query := "SELECT DISTINCT name " +
		" FROM resource_filter " +
		" WHERE deleted = false"
	_, err := repo.dbConnection.Query(&list, query)
	if err != nil {
		return nil, err
	}
	return list, err
}
