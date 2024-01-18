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

package appStoreDeploymentCommon

import (
	"encoding/json"
	"fmt"
	"github.com/devtron-labs/devtron/api/bean"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/appStore/adapter"
	appStoreBean "github.com/devtron-labs/devtron/pkg/appStore/bean"
	"github.com/devtron-labs/devtron/pkg/appStore/deployment/repository"
	appStoreDiscoverRepository "github.com/devtron-labs/devtron/pkg/appStore/discover/repository"
	util2 "github.com/devtron-labs/devtron/pkg/appStore/util"
	repository2 "github.com/devtron-labs/devtron/pkg/cluster/repository"
	commonBean "github.com/devtron-labs/devtron/pkg/deployment/gitOps/common/bean"
	"github.com/devtron-labs/devtron/pkg/deployment/gitOps/config"
	"github.com/devtron-labs/devtron/pkg/deployment/gitOps/remote"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"os"
	"path"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

type AppStoreDeploymentCommonService interface {
	// move to DB service
	GetInstalledAppByClusterNamespaceAndName(clusterId int, namespace string, appName string) (*appStoreBean.InstallAppVersionDTO, error)
	GetInstalledAppByInstalledAppId(installedAppId int) (*appStoreBean.InstallAppVersionDTO, error)

	// common for all gitops?
	CreateChartProxyAndGetPath(installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (*util.ChartCreateResponse, error)

	// keep here only
	GetValuesAndRequirementGitConfig(installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (*util.ChartConfig, *util.ChartConfig, error)
	GetGitCommitConfig(installAppVersionRequest *appStoreBean.InstallAppVersionDTO, fileString string, filename string) (*util.ChartConfig, error)
	GetValuesString(chartName, valuesOverrideYaml string) (string, error)
	GetRequirementsString(appStoreVersionId int) (string, error)

	// rethink
	GenerateManifest(installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (manifestResponse *AppStoreManifestResponse, err error)

	// gitops service
	ParseGitRepoErrorResponse(err error) (bool, error)
	CreateGitOpsRepoAndPushChart(installAppVersionRequest *appStoreBean.InstallAppVersionDTO, builtChartPath string, requirementsConfig *util.ChartConfig, valuesConfig *util.ChartConfig) (*commonBean.ChartGitAttribute, bool, string, error)
	GitOpsOperations(manifestResponse *AppStoreManifestResponse, installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (*AppStoreGitOpsResponse, error)
	GenerateManifestAndPerformGitOperations(installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (*AppStoreGitOpsResponse, error)
}

type AppStoreDeploymentCommonServiceImpl struct {
	logger                               *zap.SugaredLogger
	installedAppRepository               repository.InstalledAppRepository
	appStoreApplicationVersionRepository appStoreDiscoverRepository.AppStoreApplicationVersionRepository
	environmentRepository                repository2.EnvironmentRepository
	chartTemplateService                 util.ChartTemplateService
	gitFactory                           *util.GitFactory
	gitOpsConfigReadService              config.GitOpsConfigReadService
	gitOpsRemoteOperationService         remote.GitOpsRemoteOperationService
}

func NewAppStoreDeploymentCommonServiceImpl(
	logger *zap.SugaredLogger,
	installedAppRepository repository.InstalledAppRepository,
	appStoreApplicationVersionRepository appStoreDiscoverRepository.AppStoreApplicationVersionRepository,
	environmentRepository repository2.EnvironmentRepository,
	chartTemplateService util.ChartTemplateService,
	gitFactory *util.GitFactory,
	gitOpsConfigReadService config.GitOpsConfigReadService,
	gitOpsRemoteOperationService remote.GitOpsRemoteOperationService) *AppStoreDeploymentCommonServiceImpl {
	return &AppStoreDeploymentCommonServiceImpl{
		logger:                               logger,
		installedAppRepository:               installedAppRepository,
		appStoreApplicationVersionRepository: appStoreApplicationVersionRepository,
		environmentRepository:                environmentRepository,
		chartTemplateService:                 chartTemplateService,
		gitFactory:                           gitFactory,
		gitOpsConfigReadService:              gitOpsConfigReadService,
		gitOpsRemoteOperationService:         gitOpsRemoteOperationService,
	}
}

// TODO: move this from here

func ParseChartCreateRequest(installAppRequestDTO *appStoreBean.InstallAppVersionDTO, chartPath string) *util.ChartCreateRequest {
	return &util.ChartCreateRequest{ChartMetaData: &chart.Metadata{
		Name:    installAppRequestDTO.AppName,
		Version: "1.0.1",
	}, ChartPath: chartPath}
}

func ParseChartGitPushRequest(installAppRequestDTO *appStoreBean.InstallAppVersionDTO, repoURl string, tempRefChart string) *appStoreBean.PushChartToGitRequestDTO {
	return &appStoreBean.PushChartToGitRequestDTO{
		AppName:           installAppRequestDTO.AppName,
		EnvName:           installAppRequestDTO.Environment.Name,
		ChartAppStoreName: installAppRequestDTO.AppStoreName,
		RepoURL:           repoURl,
		TempChartRefDir:   tempRefChart,
		UserId:            installAppRequestDTO.UserId,
	}
}

type AppStoreManifestResponse struct {
	ChartResponse      *util.ChartCreateResponse
	ValuesConfig       *util.ChartConfig
	RequirementsConfig *util.ChartConfig
}

type AppStoreGitOpsResponse struct {
	ChartGitAttribute *commonBean.ChartGitAttribute
	GitHash           string
}

func (impl AppStoreDeploymentCommonServiceImpl) GetInstalledAppByClusterNamespaceAndName(clusterId int, namespace string, appName string) (*appStoreBean.InstallAppVersionDTO, error) {
	installedApp, err := impl.installedAppRepository.GetInstalledApplicationByClusterIdAndNamespaceAndAppName(clusterId, namespace, appName)
	if err != nil {
		if err == pg.ErrNoRows {
			impl.logger.Warnw("no installed apps found", "clusterId", clusterId)
			return nil, nil
		} else {
			impl.logger.Errorw("error while fetching installed apps", "clusterId", clusterId, "error", err)
			return nil, err
		}
	}

	if installedApp.Id > 0 {
		installedAppVersion, err := impl.installedAppRepository.GetInstalledAppVersionByInstalledAppIdAndEnvId(installedApp.Id, installedApp.EnvironmentId)
		if err != nil {
			return nil, err
		}
		return adapter.GenerateInstallAppVersionDTO(installedApp, installedAppVersion), nil
	}

	return nil, nil
}

func (impl AppStoreDeploymentCommonServiceImpl) GetInstalledAppByInstalledAppId(installedAppId int) (*appStoreBean.InstallAppVersionDTO, error) {
	installedAppVersion, err := impl.installedAppRepository.GetActiveInstalledAppVersionByInstalledAppId(installedAppId)
	if err != nil {
		return nil, err
	}
	installedApp := &installedAppVersion.InstalledApp
	return adapter.GenerateInstallAppVersionDTO(installedApp, installedAppVersion), nil
}

func (impl AppStoreDeploymentCommonServiceImpl) GetValuesString(chartName, valuesOverrideYaml string) (string, error) {

	ValuesOverrideByte, err := yaml.YAMLToJSON([]byte(valuesOverrideYaml))
	if err != nil {
		impl.logger.Errorw("")
	}

	var dat map[string]interface{}
	err = json.Unmarshal(ValuesOverrideByte, &dat)
	if err != nil {
		impl.logger.Errorw("error in unmarshalling values override byte", "err", err)
		return "", err
	}

	valuesMap := make(map[string]map[string]interface{})
	valuesMap[chartName] = dat
	valuesByte, err := json.Marshal(valuesMap)
	if err != nil {
		impl.logger.Errorw("error in marshaling", "err", err)
		return "", err
	}
	return string(valuesByte), nil
}

func (impl AppStoreDeploymentCommonServiceImpl) GetRequirementsString(appStoreVersionId int) (string, error) {
	appStoreAppVersion, err := impl.appStoreApplicationVersionRepository.FindById(appStoreVersionId)
	if err != nil {
		impl.logger.Errorw("fetching error", "err", err)
		return "", err
	}
	dependency := appStoreBean.Dependency{
		Name:    appStoreAppVersion.AppStore.Name,
		Version: appStoreAppVersion.Version,
	}
	if appStoreAppVersion.AppStore.ChartRepo != nil {
		dependency.Repository = appStoreAppVersion.AppStore.ChartRepo.Url
	}
	var dependencies []appStoreBean.Dependency
	dependencies = append(dependencies, dependency)
	requirementDependencies := &appStoreBean.Dependencies{
		Dependencies: dependencies,
	}
	requirementDependenciesByte, err := json.Marshal(requirementDependencies)
	if err != nil {
		return "", err
	}
	requirementDependenciesByte, err = yaml.JSONToYAML(requirementDependenciesByte)
	if err != nil {
		return "", err
	}
	return string(requirementDependenciesByte), nil
}

func (impl AppStoreDeploymentCommonServiceImpl) GetGitCommitConfig(installAppVersionRequest *appStoreBean.InstallAppVersionDTO, fileString string, filename string) (*util.ChartConfig, error) {
	appStoreAppVersion, err := impl.appStoreApplicationVersionRepository.FindById(installAppVersionRequest.AppStoreVersion)
	if err != nil {
		impl.logger.Errorw("fetching error", "err", err)
		return nil, err
	}
	environment, err := impl.environmentRepository.FindById(installAppVersionRequest.EnvironmentId)
	if err != nil {
		impl.logger.Errorw("fetching error", "err", err)
		return nil, err
	}

	argocdAppName := installAppVersionRequest.AppName + "-" + environment.Name
	gitOpsRepoName := impl.gitOpsConfigReadService.GetGitOpsRepoName(installAppVersionRequest.AppName)
	userEmailId, userName := impl.gitOpsConfigReadService.GetUserEmailIdAndNameForGitOpsCommit(installAppVersionRequest.UserId)
	YamlConfig := &util.ChartConfig{
		FileName:       filename,
		FileContent:    fileString,
		ChartName:      installAppVersionRequest.AppName,
		ChartLocation:  argocdAppName,
		ChartRepoName:  gitOpsRepoName,
		ReleaseMessage: fmt.Sprintf("release-%d-env-%d ", appStoreAppVersion.Id, environment.Id),
		UserEmailId:    userEmailId,
		UserName:       userName,
	}
	return YamlConfig, nil
}

func (impl AppStoreDeploymentCommonServiceImpl) GetValuesAndRequirementGitConfig(installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (*util.ChartConfig, *util.ChartConfig, error) {

	appStoreAppVersion, err := impl.appStoreApplicationVersionRepository.FindById(installAppVersionRequest.AppStoreVersion)
	if err != nil {
		impl.logger.Errorw("fetching error", "err", err)
		return nil, nil, err
	}
	values, err := impl.GetValuesString(appStoreAppVersion.AppStore.Name, installAppVersionRequest.ValuesOverrideYaml)
	if err != nil {
		impl.logger.Errorw("error in getting values fot installedAppVersionRequest", "err", err)
		return nil, nil, err
	}
	dependency, err := impl.GetRequirementsString(installAppVersionRequest.AppStoreVersion)
	if err != nil {
		impl.logger.Errorw("error in getting dependency array fot installedAppVersionRequest", "err", err)
		return nil, nil, err
	}
	valuesConfig, err := impl.GetGitCommitConfig(installAppVersionRequest, values, appStoreBean.VALUES_YAML_FILE)
	if err != nil {
		impl.logger.Errorw("error in creating values config for git", "err", err)
		return nil, nil, err
	}
	RequirementConfig, err := impl.GetGitCommitConfig(installAppVersionRequest, dependency, appStoreBean.REQUIREMENTS_YAML_FILE)
	if err != nil {
		impl.logger.Errorw("error in creating dependency config for git", "err", err)
		return nil, nil, err
	}
	return valuesConfig, RequirementConfig, nil
}

func (impl AppStoreDeploymentCommonServiceImpl) GenerateManifest(installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (manifestResponse *AppStoreManifestResponse, err error) {

	manifestResponse = &AppStoreManifestResponse{}

	ChartCreateResponse, err := impl.CreateChartProxyAndGetPath(installAppVersionRequest)
	if err != nil {
		impl.logger.Errorw("Error in building chart while generating manifest", "err", err)
		return manifestResponse, err
	}
	valuesConfig, dependencyConfig, err := impl.GetValuesAndRequirementGitConfig(installAppVersionRequest)
	if err != nil {
		impl.logger.Errorw("error in fetching values and requirements.yaml config while generating manifest", "err", err)
		return manifestResponse, err
	}

	manifestResponse.ChartResponse = ChartCreateResponse
	manifestResponse.ValuesConfig = valuesConfig
	manifestResponse.RequirementsConfig = dependencyConfig

	return manifestResponse, nil
}

// CreateGitOpsRepo creates a gitOps repo with readme
func (impl AppStoreDeploymentCommonServiceImpl) CreateGitOpsRepo(installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (string, bool, error) {
	if len(installAppVersionRequest.GitOpsRepoName) == 0 {
		//here gitops repo will be the app name, to breaking the mono repo structure
		gitOpsRepoName := impl.gitOpsConfigReadService.GetGitOpsRepoName(installAppVersionRequest.AppName)
		installAppVersionRequest.GitOpsRepoName = gitOpsRepoName
	}
	bitbucketMetadata, err := impl.gitOpsConfigReadService.GetBitbucketMetadata()
	if err != nil {
		impl.logger.Errorw("error in getting bitbucket metadata", "err", err)
		return "", false, err
	}
	//getting user name & emailId for commit author data
	gitRepoRequest := &bean.GitOpsConfigDto{
		GitRepoName:          installAppVersionRequest.GitOpsRepoName,
		Description:          "helm chart for " + installAppVersionRequest.GitOpsRepoName,
		BitBucketWorkspaceId: bitbucketMetadata.BitBucketWorkspaceId,
		BitBucketProjectKey:  bitbucketMetadata.BitBucketProjectKey,
	}
	repoUrl, isNew, err := impl.gitOpsRemoteOperationService.CreateRepository(gitRepoRequest, installAppVersionRequest.UserId)
	if err != nil {
		impl.logger.Errorw("error in creating git project", "name", installAppVersionRequest.GitOpsRepoName, "err", err)
		return "", false, err
	}
	return repoUrl, isNew, err
}

// AddConfigFileToChart will override requirements.yaml file in chart
func (impl AppStoreDeploymentCommonServiceImpl) AddConfigFileToChart(config *util.ChartConfig, dir string, clonedDir string) error {
	filePath := filepath.Join(clonedDir, config.FileName)
	file, err := os.Create(filePath)
	if err != nil {
		impl.logger.Errorw("error in creating file", "err", err, "fileName", config.FileName)
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(config.FileContent))
	if err != nil {
		impl.logger.Errorw("error in writing file content", "err", err, "fileName", config.FileName)
		return err
	}
	destinationFilePath := filepath.Join(dir, config.FileName)
	err = util2.MoveFileToDestination(filePath, destinationFilePath)
	if err != nil {
		impl.logger.Errorw("error in moving file from source to destination", "err", err)
		return err
	}
	return nil
}

// CreateChartProxyAndGetPath parse chart in local directory and returns path of local dir and values.yaml
func (impl AppStoreDeploymentCommonServiceImpl) CreateChartProxyAndGetPath(installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (*util.ChartCreateResponse, error) {

	ChartCreateResponse := &util.ChartCreateResponse{}
	template := appStoreBean.CHART_PROXY_TEMPLATE
	chartPath := path.Join(appStoreBean.RefChartProxyDirPath, template)
	valid, err := chartutil.IsChartDir(chartPath)
	if err != nil || !valid {
		impl.logger.Errorw("invalid base chart", "dir", chartPath, "err", err)
		return ChartCreateResponse, err
	}
	chartCreateRequest := ParseChartCreateRequest(installAppVersionRequest, chartPath)
	chartCreateResponse, err := impl.chartTemplateService.BuildChartProxyForHelmApps(chartCreateRequest)
	if err != nil {
		impl.logger.Errorw("Error in building chart proxy", "err", err)
		return chartCreateResponse, err
	}
	return chartCreateResponse, nil

}
