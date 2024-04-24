package devtronResource

import (
	"github.com/devtron-labs/devtron/pkg/devtronResource"
	"github.com/devtron-labs/devtron/pkg/devtronResource/audit"
	"github.com/devtron-labs/devtron/pkg/devtronResource/in"
	"github.com/devtron-labs/devtron/pkg/devtronResource/read"
	"github.com/devtron-labs/devtron/pkg/devtronResource/repository"
	"github.com/google/wire"
)

var DevtronResourceWireSet = wire.NewSet(
	//old bindings, migrated from wire.go
	read.NewDevtronResourceSearchableKeyServiceImpl,
	wire.Bind(new(read.DevtronResourceSearchableKeyService), new(*read.DevtronResourceSearchableKeyServiceImpl)),
	repository.NewDevtronResourceSearchableKeyRepositoryImpl,
	wire.Bind(new(repository.DevtronResourceSearchableKeyRepository), new(*repository.DevtronResourceSearchableKeyRepositoryImpl)),

	NewDevtronResourceRouterImpl,
	wire.Bind(new(DevtronResourceRouter), new(*DevtronResourceRouterImpl)),
	NewDevtronResourceRestHandlerImpl,
	wire.Bind(new(DevtronResourceRestHandler), new(*DevtronResourceRestHandlerImpl)),

	in.NewInternalProcessingServiceImpl,
	wire.Bind(new(in.InternalProcessingService), new(*in.InternalProcessingServiceImpl)),
	read.NewReadServiceImpl,
	wire.Bind(new(read.ReadService), new(*read.ReadServiceImpl)),
	devtronResource.NewDevtronResourceServiceImpl,
	wire.Bind(new(devtronResource.DevtronResourceService), new(*devtronResource.DevtronResourceServiceImpl)),
	audit.NewObjectAuditServiceImpl,
	wire.Bind(new(audit.ObjectAuditService), new(*audit.ObjectAuditServiceImpl)),

	repository.NewDevtronResourceRepositoryImpl,
	wire.Bind(new(repository.DevtronResourceRepository), new(*repository.DevtronResourceRepositoryImpl)),
	repository.NewDevtronResourceSchemaRepositoryImpl,
	wire.Bind(new(repository.DevtronResourceSchemaRepository), new(*repository.DevtronResourceSchemaRepositoryImpl)),
	repository.NewDevtronResourceObjectRepositoryImpl,
	wire.Bind(new(repository.DevtronResourceObjectRepository), new(*repository.DevtronResourceObjectRepositoryImpl)),
	repository.NewDevtronResourceSchemaAuditRepositoryImpl,
	wire.Bind(new(repository.DevtronResourceSchemaAuditRepository), new(*repository.DevtronResourceSchemaAuditRepositoryImpl)),
	repository.NewDevtronResourceObjectAuditRepositoryImpl,
	wire.Bind(new(repository.DevtronResourceObjectAuditRepository), new(*repository.DevtronResourceObjectAuditRepositoryImpl)),
)

var DevtronResourceWireSetEA = wire.NewSet(
	devtronResource.NewDevtronResourceServiceImpl,
	wire.Bind(new(devtronResource.DevtronResourceService), new(*devtronResource.DevtronResourceServiceImpl)),
	repository.NewDevtronResourceRepositoryImpl,
	wire.Bind(new(repository.DevtronResourceRepository), new(*repository.DevtronResourceRepositoryImpl)),
	repository.NewDevtronResourceSchemaRepositoryImpl,
	wire.Bind(new(repository.DevtronResourceSchemaRepository), new(*repository.DevtronResourceSchemaRepositoryImpl)),
	repository.NewDevtronResourceObjectRepositoryImpl,
	wire.Bind(new(repository.DevtronResourceObjectRepository), new(*repository.DevtronResourceObjectRepositoryImpl)),
	repository.NewDevtronResourceSchemaAuditRepositoryImpl,
	wire.Bind(new(repository.DevtronResourceSchemaAuditRepository), new(*repository.DevtronResourceSchemaAuditRepositoryImpl)),
	repository.NewDevtronResourceObjectAuditRepositoryImpl,
	wire.Bind(new(repository.DevtronResourceObjectAuditRepository), new(*repository.DevtronResourceObjectAuditRepositoryImpl)),
)
