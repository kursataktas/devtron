// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/devtron-labs/authenticator/apiToken"
	"github.com/devtron-labs/authenticator/client"
	"github.com/devtron-labs/authenticator/middleware"
	"github.com/devtron-labs/common-lib-private/utils/k8s"
	apiToken2 "github.com/devtron-labs/devtron/api/apiToken"
	chartProvider2 "github.com/devtron-labs/devtron/api/appStore/chartProvider"
	"github.com/devtron-labs/devtron/api/appStore/deployment"
	"github.com/devtron-labs/devtron/api/appStore/discover"
	"github.com/devtron-labs/devtron/api/appStore/values"
	chartRepo2 "github.com/devtron-labs/devtron/api/chartRepo"
	cluster2 "github.com/devtron-labs/devtron/api/cluster"
	"github.com/devtron-labs/devtron/api/connector"
	"github.com/devtron-labs/devtron/api/dashboardEvent"
	externalLink2 "github.com/devtron-labs/devtron/api/externalLink"
	client3 "github.com/devtron-labs/devtron/api/helm-app"
	application2 "github.com/devtron-labs/devtron/api/k8s/application"
	capacity2 "github.com/devtron-labs/devtron/api/k8s/capacity"
	module2 "github.com/devtron-labs/devtron/api/module"
	"github.com/devtron-labs/devtron/api/restHandler"
	"github.com/devtron-labs/devtron/api/router"
	server2 "github.com/devtron-labs/devtron/api/server"
	sso2 "github.com/devtron-labs/devtron/api/sso"
	team2 "github.com/devtron-labs/devtron/api/team"
	terminal2 "github.com/devtron-labs/devtron/api/terminal"
	user2 "github.com/devtron-labs/devtron/api/user"
	webhookHelm2 "github.com/devtron-labs/devtron/api/webhook/helm"
	"github.com/devtron-labs/devtron/client/dashboard"
	"github.com/devtron-labs/devtron/client/telemetry"
	repository3 "github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/internal/sql/repository/app"
	"github.com/devtron-labs/devtron/internal/sql/repository/appStatus"
	repository5 "github.com/devtron-labs/devtron/internal/sql/repository/dockerRegistry"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	"github.com/devtron-labs/devtron/internal/sql/repository/security"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/apiToken"
	app2 "github.com/devtron-labs/devtron/pkg/app"
	"github.com/devtron-labs/devtron/pkg/appStore/bean"
	"github.com/devtron-labs/devtron/pkg/appStore/chartProvider"
	"github.com/devtron-labs/devtron/pkg/appStore/deployment/common"
	repository4 "github.com/devtron-labs/devtron/pkg/appStore/deployment/repository"
	service3 "github.com/devtron-labs/devtron/pkg/appStore/deployment/service"
	"github.com/devtron-labs/devtron/pkg/appStore/deployment/tool"
	"github.com/devtron-labs/devtron/pkg/appStore/discover/repository"
	"github.com/devtron-labs/devtron/pkg/appStore/discover/service"
	"github.com/devtron-labs/devtron/pkg/appStore/values/repository"
	service2 "github.com/devtron-labs/devtron/pkg/appStore/values/service"
	"github.com/devtron-labs/devtron/pkg/attributes"
	"github.com/devtron-labs/devtron/pkg/auth"
	"github.com/devtron-labs/devtron/pkg/chartRepo"
	"github.com/devtron-labs/devtron/pkg/chartRepo/repository"
	"github.com/devtron-labs/devtron/pkg/cluster"
	repository2 "github.com/devtron-labs/devtron/pkg/cluster/repository"
	"github.com/devtron-labs/devtron/pkg/clusterTerminalAccess"
	delete2 "github.com/devtron-labs/devtron/pkg/delete"
	casbin2 "github.com/devtron-labs/devtron/pkg/enterprise/user/casbin"
	"github.com/devtron-labs/devtron/pkg/externalLink"
	"github.com/devtron-labs/devtron/pkg/genericNotes"
	repository6 "github.com/devtron-labs/devtron/pkg/genericNotes/repository"
	k8s2 "github.com/devtron-labs/devtron/pkg/k8s"
	"github.com/devtron-labs/devtron/pkg/k8s/application"
	"github.com/devtron-labs/devtron/pkg/k8s/capacity"
	"github.com/devtron-labs/devtron/pkg/k8s/informer"
	"github.com/devtron-labs/devtron/pkg/kubernetesResourceAuditLogs"
	repository7 "github.com/devtron-labs/devtron/pkg/kubernetesResourceAuditLogs/repository"
	"github.com/devtron-labs/devtron/pkg/module"
	"github.com/devtron-labs/devtron/pkg/module/repo"
	"github.com/devtron-labs/devtron/pkg/module/store"
	"github.com/devtron-labs/devtron/pkg/pipeline"
	"github.com/devtron-labs/devtron/pkg/server"
	"github.com/devtron-labs/devtron/pkg/server/config"
	"github.com/devtron-labs/devtron/pkg/server/store"
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/devtron-labs/devtron/pkg/sso"
	"github.com/devtron-labs/devtron/pkg/team"
	"github.com/devtron-labs/devtron/pkg/terminal"
	"github.com/devtron-labs/devtron/pkg/user"
	"github.com/devtron-labs/devtron/pkg/user/casbin"
	client2 "github.com/devtron-labs/devtron/pkg/user/casbin/client"
	"github.com/devtron-labs/devtron/pkg/user/repository"
	util3 "github.com/devtron-labs/devtron/pkg/util"
	"github.com/devtron-labs/devtron/pkg/webhook/helm"
	util2 "github.com/devtron-labs/devtron/util"
	"github.com/devtron-labs/devtron/util/argo"
	"github.com/devtron-labs/devtron/util/rbac"
)

// Injectors from wire.go:

func InitializeApp() (*App, error) {
	config, err := sql.GetConfig()
	if err != nil {
		return nil, err
	}
	sugaredLogger, err := util.NewSugardLogger()
	if err != nil {
		return nil, err
	}
	db, err := sql.NewDbConnection(config, sugaredLogger)
	if err != nil {
		return nil, err
	}
	runtimeConfig, err := client.GetRuntimeConfig()
	if err != nil {
		return nil, err
	}
	k8sClient, err := client.NewK8sClient(runtimeConfig)
	if err != nil {
		return nil, err
	}
	dexConfig, err := client.BuildDexConfig(k8sClient)
	if err != nil {
		return nil, err
	}
	settings, err := client.GetSettings(dexConfig)
	if err != nil {
		return nil, err
	}
	apiTokenSecretStore := apiTokenAuth.InitApiTokenSecretStore()
	sessionManager := middleware.NewSessionManager(settings, dexConfig, apiTokenSecretStore)
	validate, err := util.IntValidator()
	if err != nil {
		return nil, err
	}
	syncedEnforcer := casbin.Create()
	casbinSyncedEnforcer := casbin.CreateV2()
	casbinClientConfig, err := client2.GetConfig()
	if err != nil {
		return nil, err
	}
	casbinClientImpl := client2.NewCasbinClientImpl(sugaredLogger, casbinClientConfig)
	casbinServiceImpl := casbin.NewCasbinServiceImpl(sugaredLogger, casbinClientImpl)
	enterpriseEnforcerImpl, err := casbin2.NewEnterpriseEnforcerImpl(syncedEnforcer, casbinSyncedEnforcer, sessionManager, sugaredLogger, casbinServiceImpl)
	if err != nil {
		return nil, err
	}
	defaultAuthPolicyRepositoryImpl := repository.NewDefaultAuthPolicyRepositoryImpl(db, sugaredLogger)
	defaultAuthRoleRepositoryImpl := repository.NewDefaultAuthRoleRepositoryImpl(db, sugaredLogger)
	userAuthRepositoryImpl := repository.NewUserAuthRepositoryImpl(db, sugaredLogger, defaultAuthPolicyRepositoryImpl, defaultAuthRoleRepositoryImpl)
	userRepositoryImpl := repository.NewUserRepositoryImpl(db, sugaredLogger)
	roleGroupRepositoryImpl := repository.NewRoleGroupRepositoryImpl(db, sugaredLogger)
	rbacPolicyDataRepositoryImpl := repository.NewRbacPolicyDataRepositoryImpl(sugaredLogger, db)
	rbacRoleDataRepositoryImpl := repository.NewRbacRoleDataRepositoryImpl(sugaredLogger, db)
	rbacDataCacheFactoryImpl := repository.NewRbacDataCacheFactoryImpl(sugaredLogger, rbacPolicyDataRepositoryImpl, rbacRoleDataRepositoryImpl)
	userCommonServiceImpl := user.NewUserCommonServiceImpl(userAuthRepositoryImpl, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl, sessionManager, rbacDataCacheFactoryImpl)
	userAuditRepositoryImpl := repository.NewUserAuditRepositoryImpl(db)
	userAuditServiceImpl := user.NewUserAuditServiceImpl(sugaredLogger, userAuditRepositoryImpl)
	userServiceImpl := user.NewUserServiceImpl(userAuthRepositoryImpl, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl, sessionManager, userCommonServiceImpl, userAuditServiceImpl)
	ssoLoginRepositoryImpl := sso.NewSSOLoginRepositoryImpl(db)
	sshTunnelWrapperServiceImpl, err := k8s.NewSSHTunnelWrapperServiceImpl(sugaredLogger)
	if err != nil {
		return nil, err
	}
	k8sUtil := k8s.NewK8sUtil(sugaredLogger, runtimeConfig, sshTunnelWrapperServiceImpl)
	devtronSecretConfig, err := util2.GetDevtronSecretName()
	if err != nil {
		return nil, err
	}
	selfRegistrationRolesRepositoryImpl := repository.NewSelfRegistrationRolesRepositoryImpl(db, sugaredLogger)
	selfRegistrationRolesServiceImpl := user.NewSelfRegistrationRolesServiceImpl(sugaredLogger, selfRegistrationRolesRepositoryImpl, userServiceImpl)
	userAuthOidcHelperImpl, err := auth.NewUserAuthOidcHelperImpl(sugaredLogger, selfRegistrationRolesServiceImpl, dexConfig, settings, sessionManager)
	if err != nil {
		return nil, err
	}
	ssoLoginServiceImpl := sso.NewSSOLoginServiceImpl(sugaredLogger, ssoLoginRepositoryImpl, k8sUtil, devtronSecretConfig, userAuthOidcHelperImpl)
	ssoLoginRestHandlerImpl := sso2.NewSsoLoginRestHandlerImpl(validate, sugaredLogger, enterpriseEnforcerImpl, userServiceImpl, ssoLoginServiceImpl)
	ssoLoginRouterImpl := sso2.NewSsoLoginRouterImpl(ssoLoginRestHandlerImpl)
	teamRepositoryImpl := team.NewTeamRepositoryImpl(db)
	loginService := middleware.NewUserLogin(sessionManager, k8sClient)
	userAuthServiceImpl := user.NewUserAuthServiceImpl(userAuthRepositoryImpl, sessionManager, loginService, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl, userServiceImpl)
	teamServiceImpl := team.NewTeamServiceImpl(sugaredLogger, teamRepositoryImpl, userAuthServiceImpl)
	clusterRepositoryImpl := repository2.NewClusterRepositoryImpl(db, sugaredLogger)
	v := informer.NewGlobalMapClusterNamespace()
	k8sInformerFactoryImpl := informer.NewK8sInformerFactoryImpl(sugaredLogger, v, runtimeConfig, k8sUtil)
	clusterServiceImpl := cluster.NewClusterServiceImpl(clusterRepositoryImpl, sugaredLogger, k8sUtil, k8sInformerFactoryImpl, userAuthRepositoryImpl, userRepositoryImpl, roleGroupRepositoryImpl)
	appStatusRepositoryImpl := appStatus.NewAppStatusRepositoryImpl(db, sugaredLogger)
	environmentRepositoryImpl := repository2.NewEnvironmentRepositoryImpl(db, sugaredLogger, appStatusRepositoryImpl)
	attributesRepositoryImpl := repository3.NewAttributesRepositoryImpl(db)
	environmentServiceImpl := cluster.NewEnvironmentServiceImpl(environmentRepositoryImpl, clusterServiceImpl, sugaredLogger, k8sUtil, k8sInformerFactoryImpl, userAuthServiceImpl, attributesRepositoryImpl)
	chartRepoRepositoryImpl := chartRepoRepository.NewChartRepoRepositoryImpl(db)
	acdAuthConfig, err := util3.GetACDAuthConfig()
	if err != nil {
		return nil, err
	}
	httpClient := util.NewHttpClient()
	serverEnvConfigServerEnvConfig, err := serverEnvConfig.ParseServerEnvConfig()
	if err != nil {
		return nil, err
	}
	chartRepositoryServiceImpl := chartRepo.NewChartRepositoryServiceImpl(sugaredLogger, chartRepoRepositoryImpl, k8sUtil, clusterServiceImpl, acdAuthConfig, httpClient, serverEnvConfigServerEnvConfig)
	installedAppRepositoryImpl := repository4.NewInstalledAppRepositoryImpl(sugaredLogger, db)
	helmClientConfig, err := client3.GetConfig()
	if err != nil {
		return nil, err
	}
	helmAppClientImpl := client3.NewHelmAppClientImpl(sugaredLogger, helmClientConfig)
	pumpImpl := connector.NewPumpImpl(sugaredLogger)
	appRepositoryImpl := app.NewAppRepositoryImpl(db, sugaredLogger)
	enforcerUtilHelmImpl := rbac.NewEnforcerUtilHelmImpl(sugaredLogger, clusterRepositoryImpl, teamRepositoryImpl, appRepositoryImpl, environmentRepositoryImpl, installedAppRepositoryImpl)
	serverDataStoreServerDataStore := serverDataStore.InitServerDataStore()
	appStoreApplicationVersionRepositoryImpl := appStoreDiscoverRepository.NewAppStoreApplicationVersionRepositoryImpl(sugaredLogger, db)
	pipelineRepositoryImpl := pipelineConfig.NewPipelineRepositoryImpl(db, sugaredLogger)
	helmReleaseConfig, err := client3.GetHelmReleaseConfig()
	if err != nil {
		return nil, err
	}
	helmAppServiceImpl := client3.NewHelmAppServiceImpl(sugaredLogger, clusterServiceImpl, helmAppClientImpl, pumpImpl, enforcerUtilHelmImpl, serverDataStoreServerDataStore, serverEnvConfigServerEnvConfig, appStoreApplicationVersionRepositoryImpl, environmentServiceImpl, pipelineRepositoryImpl, installedAppRepositoryImpl, appRepositoryImpl, clusterRepositoryImpl, k8sUtil, helmReleaseConfig)
	dockerArtifactStoreRepositoryImpl := repository5.NewDockerArtifactStoreRepositoryImpl(db)
	dockerRegistryIpsConfigRepositoryImpl := repository5.NewDockerRegistryIpsConfigRepositoryImpl(db)
	ociRegistryConfigRepositoryImpl := repository5.NewOCIRegistryConfigRepositoryImpl(db)
	dockerRegistryConfigImpl := pipeline.NewDockerRegistryConfigImpl(sugaredLogger, helmAppServiceImpl, dockerArtifactStoreRepositoryImpl, dockerRegistryIpsConfigRepositoryImpl, ociRegistryConfigRepositoryImpl)
	deleteServiceImpl := delete2.NewDeleteServiceImpl(sugaredLogger, teamServiceImpl, clusterServiceImpl, environmentServiceImpl, chartRepositoryServiceImpl, installedAppRepositoryImpl, dockerRegistryConfigImpl, dockerArtifactStoreRepositoryImpl)
	teamRestHandlerImpl := team2.NewTeamRestHandlerImpl(sugaredLogger, teamServiceImpl, userServiceImpl, enterpriseEnforcerImpl, validate, userAuthServiceImpl, deleteServiceImpl)
	teamRouterImpl := team2.NewTeamRouterImpl(teamRestHandlerImpl)
	userAuthHandlerImpl := user2.NewUserAuthHandlerImpl(userAuthServiceImpl, validate, sugaredLogger, enterpriseEnforcerImpl)
	userAuthRouterImpl := user2.NewUserAuthRouterImpl(sugaredLogger, userAuthHandlerImpl, userAuthOidcHelperImpl)
	roleGroupServiceImpl := user.NewRoleGroupServiceImpl(userAuthRepositoryImpl, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl, userCommonServiceImpl)
	policiesCleanUpRepositoryImpl := repository.NewPoliciesCleanUpRepositoryImpl(db, sugaredLogger)
	cleanUpPoliciesServiceImpl := user.NewCleanUpPoliciesServiceImpl(userAuthRepositoryImpl, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl, policiesCleanUpRepositoryImpl)
	userRestHandlerImpl := user2.NewUserRestHandlerImpl(userServiceImpl, validate, sugaredLogger, enterpriseEnforcerImpl, roleGroupServiceImpl, userCommonServiceImpl, cleanUpPoliciesServiceImpl)
	userRouterImpl := user2.NewUserRouterImpl(userRestHandlerImpl)
	genericNoteRepositoryImpl := repository6.NewGenericNoteRepositoryImpl(db)
	genericNoteHistoryRepositoryImpl := repository6.NewGenericNoteHistoryRepositoryImpl(db)
	genericNoteHistoryServiceImpl := genericNotes.NewGenericNoteHistoryServiceImpl(genericNoteHistoryRepositoryImpl, sugaredLogger)
	genericNoteServiceImpl := genericNotes.NewGenericNoteServiceImpl(genericNoteRepositoryImpl, genericNoteHistoryServiceImpl, userRepositoryImpl, sugaredLogger)
	clusterDescriptionRepositoryImpl := repository2.NewClusterDescriptionRepositoryImpl(db, sugaredLogger)
	clusterDescriptionServiceImpl := cluster.NewClusterDescriptionServiceImpl(clusterDescriptionRepositoryImpl, userRepositoryImpl, sugaredLogger)
	helmUserServiceImpl, err := argo.NewHelmUserServiceImpl(sugaredLogger)
	if err != nil {
		return nil, err
	}
	clusterRbacServiceImpl := cluster.NewClusterRbacServiceImpl(environmentServiceImpl, enterpriseEnforcerImpl, clusterServiceImpl, sugaredLogger, userServiceImpl)
	clusterRestHandlerImpl := cluster2.NewClusterRestHandlerImpl(clusterServiceImpl, genericNoteServiceImpl, clusterDescriptionServiceImpl, sugaredLogger, userServiceImpl, validate, enterpriseEnforcerImpl, deleteServiceImpl, helmUserServiceImpl, environmentServiceImpl, clusterRbacServiceImpl)
	clusterRouterImpl := cluster2.NewClusterRouterImpl(clusterRestHandlerImpl)
	dashboardConfig, err := dashboard.GetConfig()
	if err != nil {
		return nil, err
	}
	dashboardRouterImpl := dashboard.NewDashboardRouterImpl(sugaredLogger, dashboardConfig)
	chartWorkingDir := _wireChartWorkingDirValue
	gitOpsConfigRepositoryImpl := repository3.NewGitOpsConfigRepositoryImpl(sugaredLogger, db)
	gitCliUtil := util.NewGitCliUtil(sugaredLogger)
	gitFactory, err := util.NewGitFactory(sugaredLogger, gitOpsConfigRepositoryImpl, gitCliUtil)
	if err != nil {
		return nil, err
	}
	globalEnvVariables, err := util2.GetGlobalEnvVariables()
	if err != nil {
		return nil, err
	}
	chartRepositoryImpl := chartRepoRepository.NewChartRepository(db)
	chartTemplateServiceImpl := util.NewChartTemplateServiceImpl(sugaredLogger, chartWorkingDir, httpClient, gitFactory, globalEnvVariables, gitOpsConfigRepositoryImpl, userRepositoryImpl, chartRepositoryImpl)
	refChartProxyDir := _wireRefChartProxyDirValue
	appStoreDeploymentCommonServiceImpl := appStoreDeploymentCommon.NewAppStoreDeploymentCommonServiceImpl(sugaredLogger, installedAppRepositoryImpl, appStoreApplicationVersionRepositoryImpl, environmentRepositoryImpl, chartTemplateServiceImpl, refChartProxyDir, gitFactory, gitOpsConfigRepositoryImpl)
	attributesServiceImpl := attributes.NewAttributesServiceImpl(sugaredLogger, attributesRepositoryImpl)
	helmAppRestHandlerImpl := client3.NewHelmAppRestHandlerImpl(sugaredLogger, helmAppServiceImpl, enterpriseEnforcerImpl, clusterServiceImpl, enforcerUtilHelmImpl, appStoreDeploymentCommonServiceImpl, userServiceImpl, attributesServiceImpl, serverEnvConfigServerEnvConfig)
	helmAppRouterImpl := client3.NewHelmAppRouterImpl(helmAppRestHandlerImpl)
	k8sResourceHistoryRepositoryImpl := repository7.NewK8sResourceHistoryRepositoryImpl(db, sugaredLogger)
	k8sResourceHistoryServiceImpl := kubernetesResourceAuditLogs.Newk8sResourceHistoryServiceImpl(k8sResourceHistoryRepositoryImpl, sugaredLogger, appRepositoryImpl, environmentRepositoryImpl)
	k8sCommonServiceImpl := k8s2.NewK8sCommonServiceImpl(sugaredLogger, k8sUtil, helmAppServiceImpl, k8sResourceHistoryServiceImpl, clusterServiceImpl)
	environmentRestHandlerImpl := cluster2.NewEnvironmentRestHandlerImpl(environmentServiceImpl, sugaredLogger, userServiceImpl, validate, enterpriseEnforcerImpl, deleteServiceImpl, k8sUtil, k8sCommonServiceImpl)
	environmentRouterImpl := cluster2.NewEnvironmentRouterImpl(environmentRestHandlerImpl)
	ephemeralContainersRepositoryImpl := repository2.NewEphemeralContainersRepositoryImpl(db)
	ephemeralContainerServiceImpl := cluster.NewEphemeralContainerServiceImpl(ephemeralContainersRepositoryImpl, sugaredLogger)
	terminalSessionHandlerImpl := terminal.NewTerminalSessionHandlerImpl(environmentServiceImpl, clusterServiceImpl, sugaredLogger, k8sUtil, ephemeralContainerServiceImpl)
	k8sApplicationServiceImpl, err := application.NewK8sApplicationServiceImpl(sugaredLogger, clusterServiceImpl, pumpImpl, helmAppServiceImpl, k8sUtil, acdAuthConfig, k8sResourceHistoryServiceImpl, k8sCommonServiceImpl, terminalSessionHandlerImpl, ephemeralContainerServiceImpl, ephemeralContainersRepositoryImpl)
	if err != nil {
		return nil, err
	}
	ciPipelineRepositoryImpl := pipelineConfig.NewCiPipelineRepositoryImpl(db, sugaredLogger)
	enforcerUtilImpl := rbac.NewEnforcerUtilImpl(sugaredLogger, teamRepositoryImpl, appRepositoryImpl, environmentRepositoryImpl, pipelineRepositoryImpl, ciPipelineRepositoryImpl, clusterRepositoryImpl, enterpriseEnforcerImpl)
	k8sApplicationRestHandlerImpl := application2.NewK8sApplicationRestHandlerImpl(sugaredLogger, k8sApplicationServiceImpl, pumpImpl, terminalSessionHandlerImpl, enterpriseEnforcerImpl, enforcerUtilHelmImpl, enforcerUtilImpl, helmAppServiceImpl, userServiceImpl, k8sCommonServiceImpl, validate)
	k8sApplicationRouterImpl := application2.NewK8sApplicationRouterImpl(k8sApplicationRestHandlerImpl)
	chartRefRepositoryImpl := chartRepoRepository.NewChartRefRepositoryImpl(db)
	refChartDir := _wireRefChartDirValue
	chartRepositoryRestHandlerImpl := chartRepo2.NewChartRepositoryRestHandlerImpl(sugaredLogger, userServiceImpl, chartRepositoryServiceImpl, enterpriseEnforcerImpl, validate, deleteServiceImpl, chartRefRepositoryImpl, refChartDir, attributesServiceImpl)
	chartRepositoryRouterImpl := chartRepo2.NewChartRepositoryRouterImpl(chartRepositoryRestHandlerImpl)
	appStoreServiceImpl := service.NewAppStoreServiceImpl(sugaredLogger, appStoreApplicationVersionRepositoryImpl)
	appStoreRestHandlerImpl := appStoreDiscover.NewAppStoreRestHandlerImpl(sugaredLogger, userServiceImpl, appStoreServiceImpl, enterpriseEnforcerImpl)
	appStoreDiscoverRouterImpl := appStoreDiscover.NewAppStoreDiscoverRouterImpl(appStoreRestHandlerImpl)
	appStoreVersionValuesRepositoryImpl := appStoreValuesRepository.NewAppStoreVersionValuesRepositoryImpl(sugaredLogger, db)
	appStoreValuesServiceImpl := service2.NewAppStoreValuesServiceImpl(sugaredLogger, appStoreApplicationVersionRepositoryImpl, installedAppRepositoryImpl, appStoreVersionValuesRepositoryImpl, userServiceImpl)
	appStoreValuesRestHandlerImpl := appStoreValues.NewAppStoreValuesRestHandlerImpl(sugaredLogger, userServiceImpl, appStoreValuesServiceImpl)
	appStoreValuesRouterImpl := appStoreValues.NewAppStoreValuesRouterImpl(appStoreValuesRestHandlerImpl)
	chartGroupDeploymentRepositoryImpl := repository4.NewChartGroupDeploymentRepositoryImpl(db, sugaredLogger)
	clusterInstalledAppsRepositoryImpl := repository4.NewClusterInstalledAppsRepositoryImpl(db, sugaredLogger)
	appStoreDeploymentHelmServiceImpl := appStoreDeploymentTool.NewAppStoreDeploymentHelmServiceImpl(sugaredLogger, helmAppServiceImpl, appStoreApplicationVersionRepositoryImpl, environmentRepositoryImpl, helmAppClientImpl, installedAppRepositoryImpl, appStoreDeploymentCommonServiceImpl, ociRegistryConfigRepositoryImpl)
	installedAppVersionHistoryRepositoryImpl := repository4.NewInstalledAppVersionHistoryRepositoryImpl(sugaredLogger, db)
	deploymentServiceTypeConfig, err := service3.GetDeploymentServiceTypeConfig()
	if err != nil {
		return nil, err
	}
	appStoreDeploymentServiceImpl := service3.NewAppStoreDeploymentServiceImpl(sugaredLogger, installedAppRepositoryImpl, chartGroupDeploymentRepositoryImpl, appStoreApplicationVersionRepositoryImpl, environmentRepositoryImpl, clusterInstalledAppsRepositoryImpl, appRepositoryImpl, appStoreDeploymentHelmServiceImpl, appStoreDeploymentHelmServiceImpl, environmentServiceImpl, clusterServiceImpl, helmAppServiceImpl, appStoreDeploymentCommonServiceImpl, globalEnvVariables, installedAppVersionHistoryRepositoryImpl, gitOpsConfigRepositoryImpl, attributesServiceImpl, deploymentServiceTypeConfig, chartTemplateServiceImpl)
	appStoreDeploymentRestHandlerImpl := appStoreDeployment.NewAppStoreDeploymentRestHandlerImpl(sugaredLogger, userServiceImpl, enterpriseEnforcerImpl, enforcerUtilImpl, enforcerUtilHelmImpl, appStoreDeploymentServiceImpl, validate, helmAppServiceImpl, appStoreDeploymentCommonServiceImpl, helmUserServiceImpl, attributesServiceImpl)
	appStoreDeploymentRouterImpl := appStoreDeployment.NewAppStoreDeploymentRouterImpl(appStoreDeploymentRestHandlerImpl)
	chartProviderServiceImpl := chartProvider.NewChartProviderServiceImpl(sugaredLogger, chartRepoRepositoryImpl, chartRepositoryServiceImpl, dockerArtifactStoreRepositoryImpl, ociRegistryConfigRepositoryImpl)
	chartProviderRestHandlerImpl := chartProvider2.NewChartProviderRestHandlerImpl(sugaredLogger, userServiceImpl, validate, chartProviderServiceImpl, enterpriseEnforcerImpl)
	chartProviderRouterImpl := chartProvider2.NewChartProviderRouterImpl(chartProviderRestHandlerImpl)
	dockerRegRestHandlerImpl := restHandler.NewDockerRegRestHandlerImpl(dockerRegistryConfigImpl, sugaredLogger, chartProviderServiceImpl, userServiceImpl, validate, enterpriseEnforcerImpl, teamServiceImpl, deleteServiceImpl)
	dockerRegRouterImpl := router.NewDockerRegRouterImpl(dockerRegRestHandlerImpl)
	posthogClient, err := telemetry.NewPosthogClient(sugaredLogger)
	if err != nil {
		return nil, err
	}
	moduleRepositoryImpl := moduleRepo.NewModuleRepositoryImpl(db)
	telemetryEventClientImpl, err := telemetry.NewTelemetryEventClientImpl(sugaredLogger, httpClient, clusterServiceImpl, k8sUtil, acdAuthConfig, userServiceImpl, attributesRepositoryImpl, ssoLoginServiceImpl, posthogClient, moduleRepositoryImpl, serverDataStoreServerDataStore, userAuditServiceImpl, helmAppClientImpl, installedAppRepositoryImpl)
	if err != nil {
		return nil, err
	}
	dashboardTelemetryRestHandlerImpl := dashboardEvent.NewDashboardTelemetryRestHandlerImpl(sugaredLogger, telemetryEventClientImpl)
	dashboardTelemetryRouterImpl := dashboardEvent.NewDashboardTelemetryRouterImpl(dashboardTelemetryRestHandlerImpl)
	commonDeploymentRestHandlerImpl := appStoreDeployment.NewCommonDeploymentRestHandlerImpl(sugaredLogger, userServiceImpl, enterpriseEnforcerImpl, enforcerUtilImpl, enforcerUtilHelmImpl, appStoreDeploymentServiceImpl, validate, helmAppServiceImpl, appStoreDeploymentCommonServiceImpl, helmAppRestHandlerImpl)
	commonDeploymentRouterImpl := appStoreDeployment.NewCommonDeploymentRouterImpl(commonDeploymentRestHandlerImpl)
	externalLinkMonitoringToolRepositoryImpl := externalLink.NewExternalLinkMonitoringToolRepositoryImpl(db)
	externalLinkIdentifierMappingRepositoryImpl := externalLink.NewExternalLinkIdentifierMappingRepositoryImpl(db)
	externalLinkRepositoryImpl := externalLink.NewExternalLinkRepositoryImpl(db)
	externalLinkServiceImpl := externalLink.NewExternalLinkServiceImpl(sugaredLogger, externalLinkMonitoringToolRepositoryImpl, externalLinkIdentifierMappingRepositoryImpl, externalLinkRepositoryImpl)
	externalLinkRestHandlerImpl := externalLink2.NewExternalLinkRestHandlerImpl(sugaredLogger, externalLinkServiceImpl, userServiceImpl, enterpriseEnforcerImpl, enforcerUtilImpl)
	externalLinkRouterImpl := externalLink2.NewExternalLinkRouterImpl(externalLinkRestHandlerImpl)
	moduleActionAuditLogRepositoryImpl := module.NewModuleActionAuditLogRepositoryImpl(db)
	serverCacheServiceImpl := server.NewServerCacheServiceImpl(sugaredLogger, serverEnvConfigServerEnvConfig, serverDataStoreServerDataStore, helmAppServiceImpl)
	moduleEnvConfig, err := module.ParseModuleEnvConfig()
	if err != nil {
		return nil, err
	}
	moduleCacheServiceImpl := module.NewModuleCacheServiceImpl(sugaredLogger, k8sUtil, moduleEnvConfig, serverEnvConfigServerEnvConfig, serverDataStoreServerDataStore, moduleRepositoryImpl, teamServiceImpl)
	moduleServiceHelperImpl := module.NewModuleServiceHelperImpl(serverEnvConfigServerEnvConfig)
	moduleResourceStatusRepositoryImpl := moduleRepo.NewModuleResourceStatusRepositoryImpl(db)
	moduleDataStoreModuleDataStore := moduleDataStore.InitModuleDataStore()
	moduleCronServiceImpl, err := module.NewModuleCronServiceImpl(sugaredLogger, moduleEnvConfig, moduleRepositoryImpl, serverEnvConfigServerEnvConfig, helmAppServiceImpl, moduleServiceHelperImpl, moduleResourceStatusRepositoryImpl, moduleDataStoreModuleDataStore)
	if err != nil {
		return nil, err
	}
	scanToolMetadataRepositoryImpl := security.NewScanToolMetadataRepositoryImpl(db, sugaredLogger)
	moduleServiceImpl := module.NewModuleServiceImpl(sugaredLogger, serverEnvConfigServerEnvConfig, moduleRepositoryImpl, moduleActionAuditLogRepositoryImpl, helmAppServiceImpl, serverDataStoreServerDataStore, serverCacheServiceImpl, moduleCacheServiceImpl, moduleCronServiceImpl, moduleServiceHelperImpl, moduleResourceStatusRepositoryImpl, scanToolMetadataRepositoryImpl)
	moduleRestHandlerImpl := module2.NewModuleRestHandlerImpl(sugaredLogger, moduleServiceImpl, userServiceImpl, enterpriseEnforcerImpl, validate)
	moduleRouterImpl := module2.NewModuleRouterImpl(moduleRestHandlerImpl)
	serverActionAuditLogRepositoryImpl := server.NewServerActionAuditLogRepositoryImpl(db)
	serverServiceImpl := server.NewServerServiceImpl(sugaredLogger, serverActionAuditLogRepositoryImpl, serverDataStoreServerDataStore, serverEnvConfigServerEnvConfig, helmAppServiceImpl, moduleRepositoryImpl)
	serverRestHandlerImpl := server2.NewServerRestHandlerImpl(sugaredLogger, serverServiceImpl, userServiceImpl, enterpriseEnforcerImpl, validate)
	serverRouterImpl := server2.NewServerRouterImpl(serverRestHandlerImpl)
	apiTokenSecretServiceImpl, err := apiToken.NewApiTokenSecretServiceImpl(sugaredLogger, attributesServiceImpl, apiTokenSecretStore)
	if err != nil {
		return nil, err
	}
	apiTokenRepositoryImpl := apiToken.NewApiTokenRepositoryImpl(db)
	apiTokenServiceImpl := apiToken.NewApiTokenServiceImpl(sugaredLogger, apiTokenSecretServiceImpl, userServiceImpl, userAuditServiceImpl, apiTokenRepositoryImpl)
	apiTokenRestHandlerImpl := apiToken2.NewApiTokenRestHandlerImpl(sugaredLogger, apiTokenServiceImpl, userServiceImpl, enterpriseEnforcerImpl, validate)
	apiTokenRouterImpl := apiToken2.NewApiTokenRouterImpl(apiTokenRestHandlerImpl)
	clusterCronServiceImpl, err := cluster.NewClusterCronServiceImpl(sugaredLogger, clusterServiceImpl)
	if err != nil {
		return nil, err
	}
	k8sCapacityServiceImpl := capacity.NewK8sCapacityServiceImpl(sugaredLogger, clusterServiceImpl, k8sApplicationServiceImpl, k8sUtil, k8sCommonServiceImpl, clusterCronServiceImpl)
	k8sCapacityRestHandlerImpl := capacity2.NewK8sCapacityRestHandlerImpl(sugaredLogger, k8sCapacityServiceImpl, userServiceImpl, enterpriseEnforcerImpl, clusterServiceImpl, environmentServiceImpl, clusterRbacServiceImpl)
	k8sCapacityRouterImpl := capacity2.NewK8sCapacityRouterImpl(k8sCapacityRestHandlerImpl)
	webhookHelmServiceImpl := webhookHelm.NewWebhookHelmServiceImpl(sugaredLogger, helmAppServiceImpl, clusterServiceImpl, chartRepositoryServiceImpl, attributesServiceImpl)
	webhookHelmRestHandlerImpl := webhookHelm2.NewWebhookHelmRestHandlerImpl(sugaredLogger, webhookHelmServiceImpl, userServiceImpl, enterpriseEnforcerImpl, validate)
	webhookHelmRouterImpl := webhookHelm2.NewWebhookHelmRouterImpl(webhookHelmRestHandlerImpl)
	userAttributesRepositoryImpl := repository3.NewUserAttributesRepositoryImpl(db)
	userAttributesServiceImpl := attributes.NewUserAttributesServiceImpl(sugaredLogger, userAttributesRepositoryImpl)
	userAttributesRestHandlerImpl := restHandler.NewUserAttributesRestHandlerImpl(sugaredLogger, enterpriseEnforcerImpl, userServiceImpl, userAttributesServiceImpl)
	userAttributesRouterImpl := router.NewUserAttributesRouterImpl(userAttributesRestHandlerImpl)
	telemetryRestHandlerImpl := restHandler.NewTelemetryRestHandlerImpl(sugaredLogger, telemetryEventClientImpl, enterpriseEnforcerImpl, userServiceImpl)
	telemetryRouterImpl := router.NewTelemetryRouterImpl(sugaredLogger, telemetryRestHandlerImpl)
	terminalAccessRepositoryImpl := repository3.NewTerminalAccessRepositoryImpl(db, sugaredLogger)
	userTerminalSessionConfig, err := clusterTerminalAccess.GetTerminalAccessConfig()
	if err != nil {
		return nil, err
	}
	userTerminalAccessServiceImpl, err := clusterTerminalAccess.NewUserTerminalAccessServiceImpl(sugaredLogger, terminalAccessRepositoryImpl, userTerminalSessionConfig, k8sCommonServiceImpl, terminalSessionHandlerImpl, k8sCapacityServiceImpl, k8sUtil)
	if err != nil {
		return nil, err
	}
	userTerminalAccessRestHandlerImpl := terminal2.NewUserTerminalAccessRestHandlerImpl(sugaredLogger, userTerminalAccessServiceImpl, enterpriseEnforcerImpl, userServiceImpl, validate)
	userTerminalAccessRouterImpl := terminal2.NewUserTerminalAccessRouterImpl(userTerminalAccessRestHandlerImpl)
	attributesRestHandlerImpl := restHandler.NewAttributesRestHandlerImpl(sugaredLogger, enterpriseEnforcerImpl, userServiceImpl, attributesServiceImpl)
	attributesRouterImpl := router.NewAttributesRouterImpl(attributesRestHandlerImpl)
	appLabelRepositoryImpl := pipelineConfig.NewAppLabelRepositoryImpl(db)
	materialRepositoryImpl := pipelineConfig.NewMaterialRepositoryImpl(db)
	appCrudOperationServiceImpl := app2.NewAppCrudOperationServiceImpl(appLabelRepositoryImpl, sugaredLogger, appRepositoryImpl, userRepositoryImpl, installedAppRepositoryImpl, teamRepositoryImpl, genericNoteServiceImpl, materialRepositoryImpl)
	appRestHandlerImpl := restHandler.NewAppRestHandlerImpl(sugaredLogger, appCrudOperationServiceImpl, userServiceImpl, validate, enforcerUtilImpl, enterpriseEnforcerImpl, helmAppServiceImpl, enforcerUtilHelmImpl, genericNoteServiceImpl)
	appRouterImpl := router.NewAppRouterImpl(sugaredLogger, appRestHandlerImpl)
	rbacPolicyResourceDetailRepositoryImpl := repository.NewRbacPolicyResourceDetailRepositoryImpl(sugaredLogger, db)
	rbacRoleResourceDetailRepositoryImpl := repository.NewRbacRoleResourceDetailRepositoryImpl(sugaredLogger, db)
	rbacRoleServiceImpl := user.NewRbacRoleServiceImpl(sugaredLogger, rbacPolicyResourceDetailRepositoryImpl, rbacRoleResourceDetailRepositoryImpl, rbacRoleDataRepositoryImpl, rbacPolicyDataRepositoryImpl, rbacDataCacheFactoryImpl, userAuthRepositoryImpl, userCommonServiceImpl)
	defaultRbacRoleDataRepositoryImpl := repository.NewDefaultRbacRoleDataRepositoryImpl(sugaredLogger, db)
	defaultRbacRoleServiceImpl := user.NewDefaultRbacRoleServiceImpl(sugaredLogger, defaultRbacRoleDataRepositoryImpl, rbacRoleServiceImpl)
	rbacRoleRestHandlerImpl := user2.NewRbacRoleHandlerImpl(sugaredLogger, validate, rbacRoleServiceImpl, userServiceImpl, enterpriseEnforcerImpl, enforcerUtilImpl, defaultRbacRoleServiceImpl)
	rbacRoleRouterImpl := user2.NewRbacRoleRouterImpl(sugaredLogger, validate, rbacRoleRestHandlerImpl)
	muxRouter := NewMuxRouter(sugaredLogger, ssoLoginRouterImpl, teamRouterImpl, userAuthRouterImpl, userRouterImpl, clusterRouterImpl, dashboardRouterImpl, helmAppRouterImpl, environmentRouterImpl, k8sApplicationRouterImpl, chartRepositoryRouterImpl, appStoreDiscoverRouterImpl, appStoreValuesRouterImpl, appStoreDeploymentRouterImpl, chartProviderRouterImpl, dockerRegRouterImpl, dashboardTelemetryRouterImpl, commonDeploymentRouterImpl, externalLinkRouterImpl, moduleRouterImpl, serverRouterImpl, apiTokenRouterImpl, k8sCapacityRouterImpl, webhookHelmRouterImpl, userAttributesRouterImpl, telemetryRouterImpl, userTerminalAccessRouterImpl, attributesRouterImpl, appRouterImpl, rbacRoleRouterImpl)
	mainApp := NewApp(db, sessionManager, muxRouter, telemetryEventClientImpl, posthogClient, sugaredLogger)
	return mainApp, nil
}

var (
	_wireChartWorkingDirValue  = util.ChartWorkingDir("/tmp/charts/")
	_wireRefChartProxyDirValue = appStoreBean.RefChartProxyDir("scripts/devtron-reference-helm-charts")
	_wireRefChartDirValue      = chartRepoRepository.RefChartDir("scripts/devtron-reference-helm-charts")
)
