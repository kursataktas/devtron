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

package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/devtron-labs/devtron/pkg/chart"
	"github.com/devtron-labs/devtron/pkg/deployment/manifest/deployedAppMetrics"
	bean2 "github.com/devtron-labs/devtron/pkg/deployment/manifest/deployedAppMetrics/bean"
	"github.com/devtron-labs/devtron/pkg/deployment/manifest/deploymentTemplate"
	"github.com/devtron-labs/devtron/pkg/deployment/manifest/deploymentTemplate/chartRef"
	chartRefBean "github.com/devtron-labs/devtron/pkg/deployment/manifest/deploymentTemplate/chartRef/bean"
	"github.com/devtron-labs/devtron/pkg/pipeline/bean"
	"github.com/devtron-labs/devtron/pkg/resourceQualifiers"
	"github.com/devtron-labs/devtron/pkg/variables"
	repository5 "github.com/devtron-labs/devtron/pkg/variables/repository"
	"time"

	chartRepoRepository "github.com/devtron-labs/devtron/pkg/chartRepo/repository"
	repository2 "github.com/devtron-labs/devtron/pkg/cluster/repository"
	"github.com/devtron-labs/devtron/pkg/pipeline/history"
	"github.com/devtron-labs/devtron/pkg/sql"

	"github.com/devtron-labs/devtron/internal/sql/models"
	"github.com/devtron-labs/devtron/internal/sql/repository/chartConfig"
	"github.com/go-pg/pg"
	"github.com/juju/errors"
	"go.uber.org/zap"
)

type PropertiesConfigService interface {
	CreateEnvironmentProperties(appId int, propertiesRequest *bean.EnvironmentProperties) (*bean.EnvironmentProperties, error)
	UpdateEnvironmentProperties(appId int, propertiesRequest *bean.EnvironmentProperties, userId int32, checkForNamespace bool) (*bean.EnvironmentProperties, error)
	//create environment entry for each new environment
	CreateIfRequired(chart *chartRepoRepository.Chart, environmentId int, userId int32, manualReviewed bool, chartStatus models.ChartStatus, isOverride, isAppMetricsEnabled bool, namespace string, IsBasicViewLocked bool, CurrentViewEditor models.ChartsViewEditorType, tx *pg.Tx) (*chartConfig.EnvConfigOverride, bool, error)
	GetEnvironmentProperties(appId, environmentId int, chartRefId int) (environmentPropertiesResponse *bean.EnvironmentPropertiesResponse, err error)
	GetEnvironmentPropertiesById(environmentId int) ([]bean.EnvironmentProperties, error)

	GetAppIdByChartEnvId(chartEnvId int) (*chartConfig.EnvConfigOverride, error)
	GetLatestEnvironmentProperties(appId, environmentId int) (*bean.EnvironmentProperties, error)
	ResetEnvironmentProperties(id int) (bool, error)
	CreateEnvironmentPropertiesWithNamespace(appId int, propertiesRequest *bean.EnvironmentProperties) (*bean.EnvironmentProperties, error)

	FetchEnvProperties(appId, envId, chartRefId int) (*chartConfig.EnvConfigOverride, error)
	ProcessEnvConfigProperties(ctx context.Context, request chart.ChartRefChangeRequest, envMetrics bool, userId int32) (*bean.EnvironmentProperties, error, bool, bool, bool)
}
type PropertiesConfigServiceImpl struct {
	logger                              *zap.SugaredLogger
	envConfigRepo                       chartConfig.EnvConfigOverrideRepository
	chartRepo                           chartRepoRepository.ChartRepository
	environmentRepository               repository2.EnvironmentRepository
	deploymentTemplateHistoryService    history.DeploymentTemplateHistoryService
	scopedVariableManager               variables.ScopedVariableManager
	deployedAppMetricsService           deployedAppMetrics.DeployedAppMetricsService
	chartService                        chart.ChartService
	chartRefService                     chartRef.ChartRefService
	deploymentTemplateValidationService deploymentTemplate.DeploymentTemplateValidationService
}

func NewPropertiesConfigServiceImpl(logger *zap.SugaredLogger,
	envConfigRepo chartConfig.EnvConfigOverrideRepository,
	chartRepo chartRepoRepository.ChartRepository,
	environmentRepository repository2.EnvironmentRepository,
	deploymentTemplateHistoryService history.DeploymentTemplateHistoryService,
	scopedVariableManager variables.ScopedVariableManager,
	deployedAppMetricsService deployedAppMetrics.DeployedAppMetricsService, chartService chart.ChartService,
	chartRefService chartRef.ChartRefService, deploymentTemplateValidationService deploymentTemplate.DeploymentTemplateValidationService) *PropertiesConfigServiceImpl {
	return &PropertiesConfigServiceImpl{
		logger:                              logger,
		envConfigRepo:                       envConfigRepo,
		chartRepo:                           chartRepo,
		environmentRepository:               environmentRepository,
		deploymentTemplateHistoryService:    deploymentTemplateHistoryService,
		scopedVariableManager:               scopedVariableManager,
		deployedAppMetricsService:           deployedAppMetricsService,
		chartService:                        chartService,
		chartRefService:                     chartRefService,
		deploymentTemplateValidationService: deploymentTemplateValidationService,
	}

}

func (impl PropertiesConfigServiceImpl) GetEnvironmentProperties(appId, environmentId int, chartRefId int) (environmentPropertiesResponse *bean.EnvironmentPropertiesResponse, err error) {
	environmentPropertiesResponse = &bean.EnvironmentPropertiesResponse{}
	env, err := impl.environmentRepository.FindById(environmentId)
	if err != nil {
		return nil, err
	}
	if len(env.Namespace) > 0 {
		environmentPropertiesResponse.Namespace = env.Namespace
	}

	// step 1
	envOverride, err := impl.envConfigRepo.ActiveEnvConfigOverride(appId, environmentId)
	if err != nil {
		return nil, err
	}
	environmentProperties := &bean.EnvironmentProperties{}
	if envOverride.Id > 0 {
		r := json.RawMessage{}
		if envOverride.IsOverride {
			err = r.UnmarshalJSON([]byte(envOverride.EnvOverrideValues))
			if err != nil {
				return nil, err
			}
			environmentPropertiesResponse.IsOverride = true
		} else {
			err = r.UnmarshalJSON([]byte(envOverride.EnvOverrideValues))
			if err != nil {
				return nil, err
			}
		}
		environmentProperties = &bean.EnvironmentProperties{
			//Id:                envOverride.Id,
			Status:            envOverride.Status,
			EnvOverrideValues: r,
			ManualReviewed:    envOverride.ManualReviewed,
			Active:            envOverride.Active,
			Namespace:         env.Namespace,
			Description:       env.Description,
			EnvironmentId:     environmentId,
			EnvironmentName:   env.Name,
			Latest:            envOverride.Latest,
			//ChartRefId:        chartRefId,
			IsOverride:        envOverride.IsOverride,
			IsBasicViewLocked: envOverride.IsBasicViewLocked,
			CurrentViewEditor: envOverride.CurrentViewEditor,
		}

		if environmentPropertiesResponse.Namespace == "" {
			environmentPropertiesResponse.Namespace = envOverride.Namespace
		}
	}
	ecOverride, err := impl.envConfigRepo.FindChartByAppIdAndEnvIdAndChartRefId(appId, environmentId, chartRefId)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if errors.IsNotFound(err) {
		environmentProperties.Id = 0
		environmentProperties.ChartRefId = chartRefId
		environmentProperties.IsOverride = false
	} else {
		environmentProperties.Id = ecOverride.Id
		environmentProperties.Latest = ecOverride.Latest
		environmentProperties.IsOverride = ecOverride.IsOverride
		environmentProperties.ChartRefId = chartRefId
		environmentProperties.ManualReviewed = ecOverride.ManualReviewed
		environmentProperties.Status = ecOverride.Status
		environmentProperties.Namespace = ecOverride.Namespace
		environmentProperties.Active = ecOverride.Active
		environmentProperties.IsBasicViewLocked = ecOverride.IsBasicViewLocked
		environmentProperties.CurrentViewEditor = ecOverride.CurrentViewEditor
	}
	environmentPropertiesResponse.ChartRefId = chartRefId
	environmentPropertiesResponse.EnvironmentConfig = *environmentProperties

	//setting global config
	chart, err := impl.chartRepo.FindLatestChartForAppByAppId(appId)
	if err != nil {
		return nil, err
	}
	if chart != nil && chart.Id > 0 {
		globalOverride := []byte(chart.GlobalOverride)
		environmentPropertiesResponse.GlobalConfig = globalOverride
		environmentPropertiesResponse.GlobalChartRefId = chart.ChartRefId
		if !environmentPropertiesResponse.IsOverride {
			environmentPropertiesResponse.EnvironmentConfig.IsBasicViewLocked = chart.IsBasicViewLocked
			environmentPropertiesResponse.EnvironmentConfig.CurrentViewEditor = chart.CurrentViewEditor
		}
	}
	isAppMetricsEnabled, err := impl.deployedAppMetricsService.GetMetricsFlagForAPipelineByAppIdAndEnvId(appId, environmentId)
	if err != nil {
		impl.logger.Errorw("error, GetMetricsFlagForAPipelineByAppIdAndEnvId", "err", err, "appId", appId, "envId", environmentId)
		return nil, err
	}
	environmentPropertiesResponse.AppMetrics = &isAppMetricsEnabled
	return environmentPropertiesResponse, nil
}

func (impl PropertiesConfigServiceImpl) FetchEnvProperties(appId, envId, chartRefId int) (*chartConfig.EnvConfigOverride, error) {
	return impl.envConfigRepo.GetByAppIdEnvIdAndChartRefId(appId, envId, chartRefId)
}

func (impl PropertiesConfigServiceImpl) CreateEnvironmentProperties(appId int, environmentProperties *bean.EnvironmentProperties) (*bean.EnvironmentProperties, error) {
	chart, err := impl.chartRepo.FindChartByAppIdAndRefId(appId, environmentProperties.ChartRefId)
	if err != nil && pg.ErrNoRows != err {
		return nil, err
	}
	if pg.ErrNoRows == err {
		impl.logger.Errorw("create new chart set latest=false", "a", "b")
		return nil, fmt.Errorf("NOCHARTEXIST")
	}
	chart.GlobalOverride = string(environmentProperties.EnvOverrideValues)
	appMetrics := false
	if environmentProperties.AppMetrics != nil {
		appMetrics = *environmentProperties.AppMetrics
	}
	envOverride, appMetrics, err := impl.CreateIfRequired(chart, environmentProperties.EnvironmentId, environmentProperties.UserId, environmentProperties.ManualReviewed, models.CHARTSTATUS_SUCCESS, true, appMetrics, environmentProperties.Namespace, environmentProperties.IsBasicViewLocked, environmentProperties.CurrentViewEditor, nil)
	if err != nil {
		return nil, err
	}
	environmentProperties.AppMetrics = &appMetrics
	r := json.RawMessage{}
	err = r.UnmarshalJSON([]byte(envOverride.EnvOverrideValues))
	if err != nil {
		return nil, err
	}
	env, err := impl.environmentRepository.FindById(environmentProperties.EnvironmentId)
	if err != nil {
		return nil, err
	}
	environmentProperties = &bean.EnvironmentProperties{
		Id:                envOverride.Id,
		Status:            envOverride.Status,
		EnvOverrideValues: r,
		ManualReviewed:    envOverride.ManualReviewed,
		Active:            envOverride.Active,
		Namespace:         env.Namespace,
		EnvironmentId:     environmentProperties.EnvironmentId,
		EnvironmentName:   env.Name,
		Latest:            envOverride.Latest,
		ChartRefId:        environmentProperties.ChartRefId,
		IsOverride:        envOverride.IsOverride,
	}
	if err != nil {
		impl.logger.Errorw("chart version parsing", "err", err, "chartVersion", chart.ChartVersion)
		return nil, err
	}

	return environmentProperties, nil
}

func (impl PropertiesConfigServiceImpl) UpdateEnvironmentProperties(appId int, propertiesRequest *bean.EnvironmentProperties, userId int32, checkForNamespace bool) (*bean.EnvironmentProperties, error) {
	//check if exists
	oldEnvOverride, err := impl.envConfigRepo.Get(propertiesRequest.Id)
	if err != nil {
		return nil, err
	}
	overrideByte, err := propertiesRequest.EnvOverrideValues.MarshalJSON()
	if err != nil {
		return nil, err
	}
	env, err := impl.environmentRepository.FindById(oldEnvOverride.TargetEnvironment)
	if err != nil {
		return nil, err
	}
	if !checkForNamespace {
		propertiesRequest.Namespace = env.Namespace
	} else {
		//FIXME add check for restricted NS also like (kube-system, devtron, monitoring, etc)
		if env.Namespace != "" && env.Namespace != propertiesRequest.Namespace {
			return nil, fmt.Errorf("enviremnt is restricted to namespace: %s only, cant deploy to: %s", env.Namespace, propertiesRequest.Namespace)
		}
	}

	if !oldEnvOverride.Latest {
		envOverrideExisting, err := impl.envConfigRepo.FindLatestChartForAppByAppIdAndEnvId(appId, oldEnvOverride.TargetEnvironment)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
		if envOverrideExisting != nil {
			envOverrideExisting.Latest = false
			envOverrideExisting.IsOverride = false
			envOverrideExisting.UpdatedOn = time.Now()
			envOverrideExisting.UpdatedBy = userId
			envOverrideExisting, err = impl.envConfigRepo.Update(envOverrideExisting)
			if err != nil {
				return nil, err
			}
		}
	}

	override := &chartConfig.EnvConfigOverride{
		Active:            propertiesRequest.Active,
		Id:                propertiesRequest.Id,
		ChartId:           oldEnvOverride.ChartId,
		EnvOverrideValues: string(overrideByte),
		Status:            propertiesRequest.Status,
		ManualReviewed:    propertiesRequest.ManualReviewed,
		Namespace:         propertiesRequest.Namespace,
		TargetEnvironment: propertiesRequest.EnvironmentId,
		IsBasicViewLocked: propertiesRequest.IsBasicViewLocked,
		CurrentViewEditor: propertiesRequest.CurrentViewEditor,
		AuditLog:          sql.AuditLog{UpdatedBy: propertiesRequest.UserId, UpdatedOn: time.Now()},
	}

	override.Latest = true
	override.IsOverride = true
	impl.logger.Debugw("updating environment override ", "value", override)
	err = impl.envConfigRepo.UpdateProperties(override)

	if oldEnvOverride.Namespace != override.Namespace {
		return nil, fmt.Errorf("namespace name update not supported")
	}

	if err != nil {
		impl.logger.Errorw("chart version parsing", "err", err)
		return nil, err
	}

	isAppMetricsEnabled := false
	if propertiesRequest.AppMetrics != nil {
		isAppMetricsEnabled = *propertiesRequest.AppMetrics
	}
	envLevelMetricsUpdateReq := &bean2.DeployedAppMetricsRequest{
		EnableMetrics: isAppMetricsEnabled,
		AppId:         appId,
		EnvId:         oldEnvOverride.TargetEnvironment,
		ChartRefId:    oldEnvOverride.Chart.ChartRefId,
		UserId:        propertiesRequest.UserId,
	}
	err = impl.deployedAppMetricsService.CreateOrUpdateAppOrEnvLevelMetrics(context.Background(), envLevelMetricsUpdateReq)
	if err != nil {
		impl.logger.Errorw("error, CheckAndUpdateAppOrEnvLevelMetrics", "err", err, "req", envLevelMetricsUpdateReq)
		return nil, err
	}

	//creating history
	err = impl.deploymentTemplateHistoryService.CreateDeploymentTemplateHistoryFromEnvOverrideTemplate(override, nil, isAppMetricsEnabled, 0)
	if err != nil {
		impl.logger.Errorw("error in creating entry for env deployment template history", "err", err, "envOverride", override)
		return nil, err
	}
	//VARIABLE_MAPPING_UPDATE
	err = impl.scopedVariableManager.ExtractAndMapVariables(override.EnvOverrideValues, override.Id, repository5.EntityTypeDeploymentTemplateEnvLevel, override.CreatedBy, nil)
	if err != nil {
		return nil, err
	}

	return propertiesRequest, err
}

func (impl PropertiesConfigServiceImpl) CreateIfRequired(chart *chartRepoRepository.Chart, environmentId int, userId int32, manualReviewed bool, chartStatus models.ChartStatus, isOverride, isAppMetricsEnabled bool, namespace string, IsBasicViewLocked bool, CurrentViewEditor models.ChartsViewEditorType, tx *pg.Tx) (*chartConfig.EnvConfigOverride, bool, error) {
	env, err := impl.environmentRepository.FindById(environmentId)
	if err != nil {
		return nil, isAppMetricsEnabled, err
	}

	if env != nil && len(env.Namespace) > 0 {
		namespace = env.Namespace
	}

	if isOverride { //case of override, to do app metrics operation
		envLevelMetricsUpdateReq := &bean2.DeployedAppMetricsRequest{
			EnableMetrics: isAppMetricsEnabled,
			AppId:         chart.AppId,
			EnvId:         environmentId,
			ChartRefId:    chart.ChartRefId,
			UserId:        userId,
		}
		err = impl.deployedAppMetricsService.CreateOrUpdateAppOrEnvLevelMetrics(context.Background(), envLevelMetricsUpdateReq)
		if err != nil {
			impl.logger.Errorw("error, CheckAndUpdateAppOrEnvLevelMetrics", "err", err, "req", envLevelMetricsUpdateReq)
			return nil, isAppMetricsEnabled, err
		}
		//updating metrics flag because it might be possible that the chartRef used was not supported and that could have override the metrics flag got in request
		isAppMetricsEnabled = envLevelMetricsUpdateReq.EnableMetrics
	}

	envOverride, err := impl.envConfigRepo.GetByChartAndEnvironment(chart.Id, environmentId)
	if err != nil && !errors.IsNotFound(err) {
		return nil, isAppMetricsEnabled, err
	}
	if errors.IsNotFound(err) {
		if isOverride {
			// before creating new entry, remove previous one from latest tag
			envOverrideExisting, err := impl.envConfigRepo.FindLatestChartForAppByAppIdAndEnvId(chart.AppId, environmentId)
			if err != nil && !errors.IsNotFound(err) {
				return nil, isAppMetricsEnabled, err
			}
			if envOverrideExisting != nil {
				envOverrideExisting.Latest = false
				envOverrideExisting.UpdatedOn = time.Now()
				envOverrideExisting.UpdatedBy = userId
				envOverrideExisting.IsOverride = isOverride
				envOverrideExisting.IsBasicViewLocked = IsBasicViewLocked
				envOverrideExisting.CurrentViewEditor = CurrentViewEditor
				//maintaining backward compatibility for while
				if tx != nil {
					envOverrideExisting, err = impl.envConfigRepo.UpdateWithTxn(envOverrideExisting, tx)
				} else {
					envOverrideExisting, err = impl.envConfigRepo.Update(envOverrideExisting)
				}
				if err != nil {
					return nil, isAppMetricsEnabled, err
				}
			}
		}

		impl.logger.Debugw("env config not found creating new ", "chart", chart.Id, "env", environmentId)
		//create new
		envOverride = &chartConfig.EnvConfigOverride{
			Active:            true,
			ManualReviewed:    manualReviewed,
			Status:            chartStatus,
			TargetEnvironment: environmentId,
			ChartId:           chart.Id,
			AuditLog:          sql.AuditLog{UpdatedBy: userId, UpdatedOn: time.Now(), CreatedOn: time.Now(), CreatedBy: userId},
			Namespace:         namespace,
			IsOverride:        isOverride,
			IsBasicViewLocked: IsBasicViewLocked,
			CurrentViewEditor: CurrentViewEditor,
		}
		if isOverride {
			envOverride.EnvOverrideValues = chart.GlobalOverride
			envOverride.Latest = true
		} else {
			envOverride.EnvOverrideValues = "{}"
		}
		//maintaining backward compatibility for while
		if tx != nil {
			err = impl.envConfigRepo.SaveWithTxn(envOverride, tx)
		} else {
			err = impl.envConfigRepo.Save(envOverride)
		}
		if err != nil {
			impl.logger.Errorw("error in creating envconfig", "data", envOverride, "error", err)
			return nil, isAppMetricsEnabled, err
		}
		err = impl.deploymentTemplateHistoryService.CreateDeploymentTemplateHistoryFromEnvOverrideTemplate(envOverride, tx, isAppMetricsEnabled, 0)
		if err != nil {
			impl.logger.Errorw("error in creating entry for env deployment template history", "err", err, "envOverride", envOverride)
			return nil, isAppMetricsEnabled, err
		}

		//VARIABLE_MAPPING_UPDATE
		if envOverride.EnvOverrideValues != "{}" {
			err = impl.scopedVariableManager.ExtractAndMapVariables(envOverride.EnvOverrideValues, envOverride.Id, repository5.EntityTypeDeploymentTemplateEnvLevel, envOverride.CreatedBy, tx)
			if err != nil {
				return nil, isAppMetricsEnabled, err
			}
		}
	}
	return envOverride, isAppMetricsEnabled, nil
}

func (impl PropertiesConfigServiceImpl) GetEnvironmentPropertiesById(envId int) ([]bean.EnvironmentProperties, error) {

	var envProperties []bean.EnvironmentProperties
	envOverrides, err := impl.envConfigRepo.GetByEnvironment(envId)
	if err != nil {
		impl.logger.Error("error fetching override config", "err", err)
		return nil, err
	}

	for _, envOverride := range envOverrides {
		envProperties = append(envProperties, bean.EnvironmentProperties{
			Id:             envOverride.Id,
			Status:         envOverride.Status,
			ManualReviewed: envOverride.ManualReviewed,
			Active:         envOverride.Active,
			Namespace:      envOverride.Namespace,
		})
	}

	return envProperties, nil
}

func (impl PropertiesConfigServiceImpl) GetAppIdByChartEnvId(chartEnvId int) (*chartConfig.EnvConfigOverride, error) {

	envOverride, err := impl.envConfigRepo.Get(chartEnvId)
	if err != nil {
		impl.logger.Error("error fetching override config", "err", err)
		return nil, err
	}

	return envOverride, nil
}

func (impl PropertiesConfigServiceImpl) GetLatestEnvironmentProperties(appId, environmentId int) (environmentProperties *bean.EnvironmentProperties, err error) {
	env, err := impl.environmentRepository.FindById(environmentId)
	if err != nil {
		return nil, err
	}
	// step 1
	envOverride, err := impl.envConfigRepo.ActiveEnvConfigOverride(appId, environmentId)
	if err != nil {
		return nil, err
	}
	if envOverride.Id == 0 {
		//return nil, errors.New("No env config exists with tag latest for given appId and envId")
		impl.logger.Warnw("No env config exists with tag latest for given appId and envId", "envId", environmentId)
	} else {
		r := json.RawMessage{}
		err = r.UnmarshalJSON([]byte(envOverride.EnvOverrideValues))
		if err != nil {
			return nil, err
		}

		environmentProperties = &bean.EnvironmentProperties{
			Id:                envOverride.Id,
			EnvOverrideValues: r,
			Status:            envOverride.Status,
			ManualReviewed:    envOverride.ManualReviewed,
			Active:            envOverride.Active,
			Namespace:         env.Namespace,
			EnvironmentId:     environmentId,
			EnvironmentName:   env.Name,
			Latest:            envOverride.Latest,
			IsOverride:        envOverride.IsOverride,
			IsBasicViewLocked: envOverride.IsBasicViewLocked,
			CurrentViewEditor: envOverride.CurrentViewEditor,
			ChartRefId:        envOverride.Chart.ChartRefId,
			ClusterId:         env.ClusterId,
		}
	}

	return environmentProperties, nil
}

func (impl PropertiesConfigServiceImpl) ResetEnvironmentProperties(id int) (bool, error) {
	envOverride, err := impl.envConfigRepo.Get(id)
	if err != nil {
		return false, err
	}
	envOverride.EnvOverrideValues = "{}"
	envOverride.IsOverride = false
	envOverride.Latest = false
	impl.logger.Infow("reset environment override ", "value", envOverride)
	err = impl.envConfigRepo.UpdateProperties(envOverride)
	if err != nil {
		impl.logger.Warnw("error in update envOverride", "envOverrideId", id)
	}
	err = impl.deployedAppMetricsService.DeleteEnvLevelMetricsIfPresent(envOverride.Chart.AppId, envOverride.TargetEnvironment)
	if err != nil {
		impl.logger.Errorw("error, DeleteEnvLevelMetricsIfPresent", "err", err, "appId", envOverride.Chart.AppId, "envId", envOverride.TargetEnvironment)
		return false, err
	}
	//VARIABLES
	err = impl.scopedVariableManager.RemoveMappedVariables(envOverride.Id, repository5.EntityTypeDeploymentTemplateEnvLevel, envOverride.UpdatedBy, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (impl PropertiesConfigServiceImpl) CreateEnvironmentPropertiesWithNamespace(appId int, environmentProperties *bean.EnvironmentProperties) (*bean.EnvironmentProperties, error) {
	chart, err := impl.chartRepo.FindChartByAppIdAndRefId(appId, environmentProperties.ChartRefId)
	if err != nil && pg.ErrNoRows != err {
		return nil, err
	}
	if pg.ErrNoRows == err {
		impl.logger.Warnw("no chart found this ref id", "refId", environmentProperties.ChartRefId)
		chart, err = impl.chartRepo.FindLatestChartForAppByAppId(appId)
		if err != nil && pg.ErrNoRows != err {
			return nil, err
		}
	}

	var envOverride *chartConfig.EnvConfigOverride
	if environmentProperties.Id == 0 {
		chart.GlobalOverride = "{}"
		appMetrics := false
		if environmentProperties.AppMetrics != nil {
			appMetrics = *environmentProperties.AppMetrics
		}
		envOverride, appMetrics, err = impl.CreateIfRequired(chart, environmentProperties.EnvironmentId, environmentProperties.UserId, environmentProperties.ManualReviewed, models.CHARTSTATUS_SUCCESS, false, appMetrics, environmentProperties.Namespace, environmentProperties.IsBasicViewLocked, environmentProperties.CurrentViewEditor, nil)
		if err != nil {
			return nil, err
		}
		environmentProperties.AppMetrics = &appMetrics
	} else {
		envOverride, err = impl.envConfigRepo.Get(environmentProperties.Id)
		if err != nil {
			impl.logger.Errorw("error in fetching envOverride", "err", err)
		}
		envOverride.Namespace = environmentProperties.Namespace
		envOverride.UpdatedBy = environmentProperties.UserId
		envOverride.IsBasicViewLocked = environmentProperties.IsBasicViewLocked
		envOverride.CurrentViewEditor = environmentProperties.CurrentViewEditor
		envOverride.UpdatedOn = time.Now()
		impl.logger.Debugw("updating environment override ", "value", envOverride)
		err = impl.envConfigRepo.UpdateProperties(envOverride)
	}

	r := json.RawMessage{}
	err = r.UnmarshalJSON([]byte(envOverride.EnvOverrideValues))
	if err != nil {
		return nil, err
	}
	env, err := impl.environmentRepository.FindById(environmentProperties.EnvironmentId)
	if err != nil {
		return nil, err
	}
	environmentProperties = &bean.EnvironmentProperties{
		Id:                envOverride.Id,
		Status:            envOverride.Status,
		EnvOverrideValues: r,
		ManualReviewed:    envOverride.ManualReviewed,
		Active:            envOverride.Active,
		Namespace:         env.Namespace,
		EnvironmentId:     environmentProperties.EnvironmentId,
		EnvironmentName:   env.Name,
		Latest:            envOverride.Latest,
		ChartRefId:        environmentProperties.ChartRefId,
		IsOverride:        envOverride.IsOverride,
		ClusterId:         env.ClusterId,
	}
	return environmentProperties, nil
}

// ProcessEnvConfigProperties processes the environment configuration properties based on the provided parameters.
// Steps performed:
// 1. Fetch the latest environment properties for the given `appId` and `envId`.
// 2. Handle error if fetching environment properties fails.
// 3. If environment properties are not found:
//   - Set `isEnvironmentOverridden` to false.
//   - If `allowEnvOverride` flag is not set in the request:
//   - Log an error indicating that environment properties are not found and advise to set the `allowEnvOverride` flag to true.
//   - Return nil, error, and flags indicating the absence of environment properties and updateOverride.
//   - Otherwise, create an empty environment properties object.
//
// 4. If environment properties are found:
//   - Check if the environment is overridden.
//   - Handle error if `isOverride` is not true.
//   - Check compatibility between old and new chart versions.
//   - Patch environment overrides based on old and new chart types.
//   - If the new chart type is rollout:
//   - Check if Flagger canary is enabled; if so, log an error and return.
//   - Fetch environment metrics for the given `appId` and `envId`.
//   - Handle error if fetching environment metrics fails.
//   - Set environment metrics in the environment properties.
//   - Define scope for validation including `appId`, `envId`, and `clusterId`.
//   - Validate deployment template with the provided environment override values, chart reference ID, and scope.
//   - Handle error if validation fails.
//
// 5. Set `userId`, `chartRefId`, and `environmentId` in the environment properties.
// 6. If environment is not overridden:
//   - Set `isOverride` to true.
//   - Fetch previous environment properties with the default chart reference ID.
//   - If no error occurs and previous environment properties are found:
//   - Set environment override values and activate the environment.
//   - Set ID of the previous environment properties.
//   - Otherwise, set `updateOverride` flag to false.
//   - Fetch app override for the default template with the target chart reference ID.
//   - Handle error if fetching app override fails.
//   - Set environment override values to the app override.
//
// 7. Return the processed environment properties, nil error, and flags indicating environment status and update behavior.
func (impl PropertiesConfigServiceImpl) ProcessEnvConfigProperties(ctx context.Context, request chart.ChartRefChangeRequest, envMetrics bool, userId int32) (*bean.EnvironmentProperties, error, bool, bool, bool) {
	// Default values
	isEnvironmentOverridden := true
	updateOverride := true

	// Retrieve latest environment properties
	envConfigProperties, err := impl.GetLatestEnvironmentProperties(request.AppId, request.EnvId)
	if err != nil {
		impl.logger.Errorw("error in finding env properties ChangeChartRef", "err", err, "payload", request)
		return nil, err, false, false, false
	}

	// Handle case when environment properties are not found
	if envConfigProperties == nil {
		isEnvironmentOverridden = false
		if !request.AllowEnvOverride {
			impl.logger.Errorw("env properties not found, ChangeChartRef, to override environment add flag allowEnvOverride = true", "err", err, "payload", request)
			return nil, err, false, false, false
		}
		// Create empty environment properties if environment override is allowed
		envConfigProperties = &bean.EnvironmentProperties{}
	} else {
		// Check if environment is overridden
		if !envConfigProperties.IsOverride {
			impl.logger.Errorw("isOverride is not true, ChangeChartRef", "err", err, "payload", request)
			return nil, errors.New("specific environment is not overridden"), false, false, false
		}

		// Check compatibility of chart references
		compatible, oldChartType, newChartType := impl.chartRefService.ChartRefIdsCompatible(envConfigProperties.ChartRefId, request.TargetChartRefId)
		if !compatible {
			return nil, errors.New("chart version is not compatible"), false, false, false
		}

		// Patch environment overrides based on chart changes
		envConfigProperties.EnvOverrideValues, err = impl.chartService.PatchEnvOverrides(envConfigProperties.EnvOverrideValues, oldChartType, newChartType)
		if err != nil {
			return nil, err, false, false, false
		}

		// Handle rollout chart type
		if newChartType == chartRefBean.RolloutChartType {
			enabled, err := impl.deploymentTemplateValidationService.FlaggerCanaryEnabled(envConfigProperties.EnvOverrideValues)
			if err != nil || enabled {
				impl.logger.Errorw("rollout charts do not support flaggerCanary, ChangeChartRef", "err", err, "payload", request)
				return nil, err, false, false, false
			}
		}

		// Retrieve environment metrics
		envMetrics, err = impl.deployedAppMetricsService.GetMetricsFlagByAppIdAndEnvId(request.AppId, request.EnvId)
		if err != nil {
			impl.logger.Errorw("could not find envMetrics for, ChangeChartRef", "err", err, "payload", request)
			return nil, err, false, false, false
		}
		envConfigProperties.AppMetrics = &envMetrics

		// Validate deployment template
		scope := resourceQualifiers.Scope{
			AppId:     request.AppId,
			EnvId:     request.EnvId,
			ClusterId: envConfigProperties.ClusterId,
		}
		validate, err := impl.deploymentTemplateValidationService.DeploymentTemplateValidate(ctx, envConfigProperties.EnvOverrideValues, envConfigProperties.ChartRefId, scope)
		if !validate {
			impl.logger.Errorw("validation err, UpdateAppOverride", "err", err, "payload", request)
			return nil, err, false, false, false
		}
	}

	// Update environment properties
	envConfigProperties.UserId = userId
	envConfigProperties.ChartRefId = request.TargetChartRefId
	envConfigProperties.EnvironmentId = request.EnvId

	// Handle case when environment is not overridden
	if !isEnvironmentOverridden {
		envConfigProperties.IsOverride = true

		// Fetch existing environment properties if available
		envConfigPropertiesOld, err := impl.FetchEnvProperties(request.AppId, request.EnvId, request.TargetChartRefId)
		if err == nil {
			if envConfigPropertiesOld.Chart != nil {
				envConfigProperties.EnvOverrideValues = json.RawMessage(envConfigPropertiesOld.Chart.Values)
				envConfigProperties.Active = true
			}
			envConfigProperties.Id = envConfigPropertiesOld.Id
		} else {
			updateOverride = false

			// Retrieve default chart override
			_, appOverride, err := impl.chartRefService.GetAppOverrideForDefaultTemplate(request.TargetChartRefId)
			if err != nil {
				impl.logger.Errorw("service err, GetAppOverrideForDefaultTemplate", "err", err, "payload", envConfigProperties)
				return nil, err, false, false, updateOverride
			}
			// Set default chart override
			envConfigProperties.EnvOverrideValues = json.RawMessage(appOverride)
		}
	}

	return envConfigProperties, nil, isEnvironmentOverridden, envMetrics, updateOverride
}
