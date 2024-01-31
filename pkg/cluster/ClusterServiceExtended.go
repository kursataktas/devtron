package cluster

import (
	"context"
	"fmt"
	repository3 "github.com/devtron-labs/devtron/internal/sql/repository"
	auth "github.com/devtron-labs/devtron/pkg/auth/authorisation/globalConfig"
	"github.com/devtron-labs/devtron/pkg/auth/user"
	"github.com/devtron-labs/devtron/pkg/imageDigestPolicy"
	"net/http"
	"strings"
	"time"

	cluster3 "github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/devtron-labs/common-lib-private/utils/k8s"
	k8s2 "github.com/devtron-labs/common-lib/utils/k8s"
	repository4 "github.com/devtron-labs/devtron/pkg/auth/user/repository"
	"github.com/devtron-labs/devtron/pkg/k8s/informer"
	"github.com/go-pg/pg"

	cluster2 "github.com/devtron-labs/devtron/client/argocdServer/cluster"
	"github.com/devtron-labs/devtron/client/grafana"
	"github.com/devtron-labs/devtron/internal/constants"
	"github.com/devtron-labs/devtron/internal/util"
	appStoreBean "github.com/devtron-labs/devtron/pkg/appStore/bean"
	repository2 "github.com/devtron-labs/devtron/pkg/appStore/deployment/repository"
	"github.com/devtron-labs/devtron/pkg/cluster/repository"
	"go.uber.org/zap"
)

// ClusterServiceImplExtended extends ClusterServiceImpl and enhances method of ClusterService with full mode specific errors
type ClusterServiceImplExtended struct {
	environmentRepository    repository.EnvironmentRepository
	grafanaClient            grafana.GrafanaClient
	installedAppRepository   repository2.InstalledAppRepository
	clusterServiceCD         cluster2.ServiceClient
	K8sInformerFactory       informer.K8sInformerFactory
	gitOpsRepository         repository3.GitOpsConfigRepository
	sshTunnelWrapperService  k8s.SSHTunnelWrapperService
	imageDigestPolicyService imageDigestPolicy.ImageDigestPolicyService
	*ClusterServiceImpl
}

func NewClusterServiceImplExtended(repository repository.ClusterRepository, environmentRepository repository.EnvironmentRepository,
	grafanaClient grafana.GrafanaClient, logger *zap.SugaredLogger, installedAppRepository repository2.InstalledAppRepository,
	K8sUtil *k8s.K8sUtilExtended,
	clusterServiceCD cluster2.ServiceClient, K8sInformerFactory informer.K8sInformerFactory,
	gitOpsRepository repository3.GitOpsConfigRepository, userAuthRepository repository4.UserAuthRepository,
	userRepository repository4.UserRepository, roleGroupRepository repository4.RoleGroupRepository,
	sshTunnelWrapperService k8s.SSHTunnelWrapperService,
	globalAuthorisationConfigService auth.GlobalAuthorisationConfigService,
	userService user.UserService,
	imageDigestPolicyService imageDigestPolicy.ImageDigestPolicyService) *ClusterServiceImplExtended {
	clusterServiceExt := &ClusterServiceImplExtended{
		environmentRepository:    environmentRepository,
		grafanaClient:            grafanaClient,
		installedAppRepository:   installedAppRepository,
		clusterServiceCD:         clusterServiceCD,
		gitOpsRepository:         gitOpsRepository,
		sshTunnelWrapperService:  sshTunnelWrapperService,
		imageDigestPolicyService: imageDigestPolicyService,
		ClusterServiceImpl: &ClusterServiceImpl{
			clusterRepository:                repository,
			logger:                           logger,
			K8sUtil:                          K8sUtil,
			K8sInformerFactory:               K8sInformerFactory,
			userAuthRepository:               userAuthRepository,
			userRepository:                   userRepository,
			roleGroupRepository:              roleGroupRepository,
			globalAuthorisationConfigService: globalAuthorisationConfigService,
			ClusterRbacServiceImpl: &ClusterRbacServiceImpl{
				userService: userService,
				logger:      logger,
			},
		},
	}
	go clusterServiceExt.updateClusterConnectionMap()
	go clusterServiceExt.buildInformer()
	return clusterServiceExt
}

func (impl *ClusterServiceImplExtended) updateClusterConnectionMap() {
	//getting all clusters which have SSH Tunnel configured
	clusters, err := impl.clusterRepository.GetAllSSHTunnelConfiguredClusters()
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Error("error in getting all sshTunnel configured clusters", "err", err)
		return
	}
	for _, cluster := range clusters {
		clusterBean := GetClusterBean(*cluster)
		clusterConfig := clusterBean.GetClusterConfig()
		_, err = impl.sshTunnelWrapperService.StartUpdateConnectionForCluster(clusterConfig)
		if err != nil {
			impl.logger.Warnw("error in connecting with cluster through SSH tunnel", "err", err, "clusterId", cluster.Id)
			//not returning here as at startup it might be possible that the cluster connection status is not updated,
			//and we might be trying on an invalid cluster, so with one try we will move on
			continue
		}
	}
	return
}
func (impl *ClusterServiceImplExtended) FindAllWithoutConfig() ([]*ClusterBean, error) {
	beans, err := impl.FindAll()
	if err != nil {
		return nil, err
	}
	for _, bean := range beans {
		bean.Config = map[string]string{k8s2.BearerToken: ""}
		if bean.SSHTunnelConfig != nil {
			if len(bean.SSHTunnelConfig.Password) > 0 {
				bean.SSHTunnelConfig.Password = SecretDataObfuscatePlaceholder
			}
			if len(bean.SSHTunnelConfig.AuthKey) > 0 {
				bean.SSHTunnelConfig.AuthKey = SecretDataObfuscatePlaceholder
			}
		}
	}
	return beans, nil
}

func (impl *ClusterServiceImplExtended) GetClusterFullModeDTO(beans []*ClusterBean) ([]*ClusterBean, error) {
	//devtron full mode logic
	var clusterIds []int
	for _, cluster := range beans {
		clusterIds = append(clusterIds, cluster.Id)
	}
	clusterComponentsMap := make(map[int][]*repository2.InstalledAppVersions)
	charts, err := impl.installedAppRepository.GetInstalledAppVersionByClusterIdsV2(clusterIds)
	if err != nil {
		impl.logger.Errorw("error on fetching installed apps for cluster ids", "err", err, "clusterIds", clusterIds)
		return nil, err
	}
	for _, item := range charts {
		if _, ok := clusterComponentsMap[item.InstalledApp.Environment.ClusterId]; !ok {
			var charts []*repository2.InstalledAppVersions
			charts = append(charts, item)
			clusterComponentsMap[item.InstalledApp.Environment.ClusterId] = charts
		} else {
			charts := clusterComponentsMap[item.InstalledApp.Environment.ClusterId]
			charts = append(charts, item)
			clusterComponentsMap[item.InstalledApp.Environment.ClusterId] = charts
		}
	}

	for _, item := range beans {
		defaultClusterComponents := make([]*DefaultClusterComponent, 0)
		if _, ok := clusterComponentsMap[item.Id]; ok {
			charts := clusterComponentsMap[item.Id]
			failed := false
			chartLen := 0
			chartPass := 0
			if len(charts) > 0 {
				chartLen = len(charts)
			}
			for _, chart := range charts {
				defaultClusterComponent := &DefaultClusterComponent{}
				defaultClusterComponent.AppId = chart.InstalledApp.AppId
				defaultClusterComponent.InstalledAppId = chart.InstalledApp.Id
				defaultClusterComponent.EnvId = chart.InstalledApp.EnvironmentId
				defaultClusterComponent.EnvName = chart.InstalledApp.Environment.Name
				defaultClusterComponent.ComponentName = chart.AppStoreApplicationVersion.AppStore.Name
				defaultClusterComponent.Status = chart.InstalledApp.Status.String()
				defaultClusterComponents = append(defaultClusterComponents, defaultClusterComponent)
				if chart.InstalledApp.Status == appStoreBean.QUE_ERROR || chart.InstalledApp.Status == appStoreBean.TRIGGER_ERROR ||
					chart.InstalledApp.Status == appStoreBean.DEQUE_ERROR || chart.InstalledApp.Status == appStoreBean.GIT_ERROR ||
					chart.InstalledApp.Status == appStoreBean.ACD_ERROR {
					failed = true
				}
				if chart.InstalledApp.Status == appStoreBean.DEPLOY_SUCCESS {
					chartPass = chartPass + 1
				}
			}
			if chartPass == chartLen {
				item.AgentInstallationStage = 2
			} else if failed {
				item.AgentInstallationStage = 3
			} else {
				item.AgentInstallationStage = 1
			}
		}
		if item.Id == 1 {
			item.AgentInstallationStage = -1
		}
		item.DefaultClusterComponent = defaultClusterComponents
	}
	return beans, nil
}

func (impl *ClusterServiceImplExtended) FindAll() ([]*ClusterBean, error) {
	beans, err := impl.ClusterServiceImpl.FindAll()
	if err != nil {
		return nil, err
	}
	return impl.GetClusterFullModeDTO(beans)
}

func (impl *ClusterServiceImplExtended) FindAllExceptVirtual() ([]*ClusterBean, error) {
	beans, err := impl.ClusterServiceImpl.FindAll()
	if err != nil {
		return nil, err
	}
	return impl.GetClusterFullModeDTO(beans)
}

func (impl *ClusterServiceImplExtended) Update(ctx context.Context, bean *ClusterBean, userId int32) (*ClusterBean, error) {
	isGitOpsConfigured, err1 := impl.gitOpsRepository.IsGitOpsConfigured()
	if err1 != nil {
		return nil, err1
	}

	bean, err := impl.ClusterServiceImpl.Update(ctx, bean, userId)
	if err != nil {
		return nil, err
	}

	envs, err := impl.environmentRepository.FindByClusterId(bean.Id)
	if err != nil && !util.IsErrNoRows(err) {
		impl.logger.Error(err)
		return nil, err
	}

	// TODO: Can be called in goroutines if performance issue
	for _, env := range envs {
		if len(bean.PrometheusUrl) > 0 && env.GrafanaDatasourceId == 0 {
			grafanaDatasourceId, _ := impl.CreateGrafanaDataSource(bean, env)
			if grafanaDatasourceId == 0 {
				impl.logger.Errorw("unable to create data source for environment which doesn't exists", "env", env)
				continue
			}
			env.GrafanaDatasourceId = grafanaDatasourceId
		}
		//if the request doesn't have a non-empty prometheus url, and we don't have a GrafanaDataSourceId defined yet, no point in
		//going to grafana client and trying to get data source
		if bean.PrometheusUrl != "" && env.GrafanaDatasourceId != 0 {
			promDatasource, err := impl.grafanaClient.GetDatasource(env.GrafanaDatasourceId)
			if err != nil {
				impl.logger.Errorw("error on getting data source", "err", err)
				return nil, err
			}

			updateDatasourceReq := grafana.UpdateDatasourceRequest{
				Id:                env.GrafanaDatasourceId,
				OrgId:             promDatasource.OrgId,
				Name:              promDatasource.Name,
				Type:              promDatasource.Type,
				Url:               bean.PrometheusUrl,
				Access:            promDatasource.Access,
				BasicAuth:         promDatasource.BasicAuth,
				BasicAuthUser:     promDatasource.BasicAuthUser,
				BasicAuthPassword: promDatasource.BasicAuthPassword,
				JsonData:          promDatasource.JsonData,
			}

			if bean.PrometheusAuth != nil {
				secureJsonData := &grafana.SecureJsonData{}
				if len(bean.PrometheusAuth.UserName) > 0 {
					updateDatasourceReq.BasicAuthUser = bean.PrometheusAuth.UserName
					updateDatasourceReq.BasicAuthPassword = bean.PrometheusAuth.Password
					secureJsonData.BasicAuthPassword = bean.PrometheusAuth.Password
				}
				if len(bean.PrometheusAuth.TlsClientCert) > 0 {
					secureJsonData.TlsClientCert = bean.PrometheusAuth.TlsClientCert
					secureJsonData.TlsClientKey = bean.PrometheusAuth.TlsClientKey
					updateDatasourceReq.BasicAuth = false

					jsonData := &grafana.JsonData{
						HttpMethod: http.MethodGet,
						TlsAuth:    true,
					}
					updateDatasourceReq.JsonData = *jsonData
				}
				updateDatasourceReq.SecureJsonData = secureJsonData
			}
			_, err = impl.grafanaClient.UpdateDatasource(updateDatasourceReq, env.GrafanaDatasourceId)
			if err != nil {
				impl.logger.Errorw("Error while updating the datasource", "Datasource id : ", env.GrafanaDatasourceId, "error", err)
				return nil, err
			}
		}

	}

	// if git-ops configured and no proxy is configured, then only update cluster in ACD, otherwise ignore
	if isGitOpsConfigured && len(bean.ProxyUrl) == 0 && !bean.ToConnectWithSSHTunnel {
		configMap := bean.Config
		serverUrl := bean.ServerUrl
		bearerToken := ""
		if configMap[k8s2.BearerToken] != "" {
			bearerToken = configMap[k8s2.BearerToken]
		}

		tlsConfig := v1alpha1.TLSClientConfig{
			Insecure: bean.InsecureSkipTLSVerify,
		}
		if !bean.InsecureSkipTLSVerify {
			tlsConfig.KeyData = []byte(configMap[k8s2.TlsKey])
			tlsConfig.CertData = []byte(configMap[k8s2.CertData])
			tlsConfig.CAData = []byte(configMap[k8s2.CertificateAuthorityData])
		}

		cdClusterConfig := v1alpha1.ClusterConfig{
			BearerToken:     bearerToken,
			TLSClientConfig: tlsConfig,
		}

		cl := &v1alpha1.Cluster{
			Name:   bean.ClusterName,
			Server: serverUrl,
			Config: cdClusterConfig,
		}

		_, err = impl.clusterServiceCD.Update(ctx, &cluster3.ClusterUpdateRequest{Cluster: cl})

		if err != nil {
			impl.logger.Errorw("service err, Update", "error", err, "payload", cl)
			userMsg := "failed to update on cluster via ACD"
			if strings.Contains(err.Error(), k8s2.DefaultClusterUrl) {
				userMsg = fmt.Sprintf("%s, %s", err.Error(), ", successfully updated in ACD")
			}
			err = &util.ApiError{
				Code:            constants.ClusterUpdateACDFailed,
				InternalMessage: err.Error(),
				UserMessage:     userMsg,
			}
			return nil, err
		}
	}

	if bean.HasConfigOrUrlChanged {
		impl.ClusterServiceImpl.SyncNsInformer(bean)
	}

	return bean, err
}

func (impl *ClusterServiceImplExtended) CreateGrafanaDataSource(clusterBean *ClusterBean, env *repository.Environment) (int, error) {
	grafanaDatasourceId := env.GrafanaDatasourceId
	if grafanaDatasourceId == 0 {
		//starts grafana creation
		createDatasourceReq := grafana.CreateDatasourceRequest{
			Name:      "Prometheus-" + env.Name,
			Type:      "prometheus",
			Url:       clusterBean.PrometheusUrl,
			Access:    "proxy",
			BasicAuth: true,
		}

		if clusterBean.PrometheusAuth != nil {
			secureJsonData := &grafana.SecureJsonData{}
			if len(clusterBean.PrometheusAuth.UserName) > 0 {
				createDatasourceReq.BasicAuthUser = clusterBean.PrometheusAuth.UserName
				createDatasourceReq.BasicAuthPassword = clusterBean.PrometheusAuth.Password
				secureJsonData.BasicAuthPassword = clusterBean.PrometheusAuth.Password
			}
			if len(clusterBean.PrometheusAuth.TlsClientCert) > 0 {
				secureJsonData.TlsClientCert = clusterBean.PrometheusAuth.TlsClientCert
				secureJsonData.TlsClientKey = clusterBean.PrometheusAuth.TlsClientKey

				jsonData := &grafana.JsonData{
					HttpMethod: http.MethodGet,
					TlsAuth:    true,
				}
				createDatasourceReq.JsonData = jsonData
			}
			createDatasourceReq.SecureJsonData = secureJsonData
		}

		grafanaResp, err := impl.grafanaClient.CreateDatasource(createDatasourceReq)
		if err != nil {
			impl.logger.Errorw("error on create grafana datasource", "err", err)
			return 0, err
		}
		//ends grafana creation
		grafanaDatasourceId = grafanaResp.Id
	}
	env.GrafanaDatasourceId = grafanaDatasourceId
	err := impl.environmentRepository.Update(env)
	if err != nil {
		impl.logger.Errorw("error in updating environment", "err", err)
		return 0, err
	}
	return grafanaDatasourceId, nil
}

func (impl *ClusterServiceImplExtended) Save(ctx context.Context, bean *ClusterBean, userId int32) (*ClusterBean, error) {
	isGitOpsConfigured, err := impl.gitOpsRepository.IsGitOpsConfigured()
	if err != nil {
		return nil, err
	}

	clusterBean, err := impl.ClusterServiceImpl.Save(ctx, bean, userId)
	if err != nil {
		return nil, err
	}

	// if git-ops configured and no proxy or ssh tunnel is configured, then only add cluster in ACD, otherwise ignore
	if isGitOpsConfigured && len(clusterBean.ProxyUrl) == 0 && !clusterBean.ToConnectWithSSHTunnel {
		//create it into argo cd as well
		cl := impl.ConvertClusterBeanObjectToCluster(bean)

		_, err = impl.clusterServiceCD.Create(ctx, &cluster3.ClusterCreateRequest{Upsert: true, Cluster: cl})
		if err != nil {
			impl.logger.Errorw("service err, Save", "err", err, "payload", cl)
			err1 := impl.ClusterServiceImpl.Delete(bean, userId) //FIXME nishant call local
			if err1 != nil {
				impl.logger.Errorw("service err, Save, delete on rollback", "err", err, "payload", bean)
				err = &util.ApiError{
					Code:            constants.ClusterDBRollbackFailed,
					InternalMessage: err.Error(),
					UserMessage:     "failed to rollback cluster from db as it has failed in registering on ACD",
				}
				return nil, err

			}
			err = &util.ApiError{
				Code:            constants.ClusterCreateACDFailed,
				InternalMessage: err.Error(),
				UserMessage:     "failed to register on ACD, rollback completed from db",
			}
			return nil, err
		}
	}

	//on successful creation of new cluster, update informer cache for namespace group by cluster
	impl.SyncNsInformer(bean)
	return clusterBean, nil
}

func (impl ClusterServiceImplExtended) DeleteFromDb(bean *ClusterBean, userId int32) error {
	existingCluster, err := impl.clusterRepository.FindById(bean.Id)
	if err != nil {
		impl.logger.Errorw("No matching entry found for delete.", "id", bean.Id)
		return err
	}
	deleteReq := existingCluster
	deleteReq.UpdatedOn = time.Now()
	deleteReq.UpdatedBy = userId
	err = impl.clusterRepository.MarkClusterDeleted(deleteReq)
	if err != nil {
		impl.logger.Errorw("error in deleting cluster", "id", bean.Id, "err", err)
		return err
	}
	k8sClient, err := impl.ClusterServiceImpl.K8sUtil.GetCoreV1ClientInCluster()
	if err != nil {
		impl.logger.Errorw("error in creating k8s client set", "err", err, "clusterName", bean.ClusterName)
	}
	secretName := fmt.Sprintf("%s-%v", "cluster-event", bean.Id)
	err = impl.K8sUtil.DeleteSecret("default", secretName, k8sClient)
	impl.logger.Errorw("error in deleting secret", "error", err)
	return nil
}

func (impl ClusterServiceImplExtended) IsPolicyConfiguredForCluster(envId, clusterId int) (bool, error) {

	digestConfigurationRequest := imageDigestPolicy.DigestPolicyConfigurationRequest{ClusterId: clusterId, EnvironmentId: envId}
	digestPolicyConfigurations, err := impl.imageDigestPolicyService.GetDigestPolicyConfigurations(digestConfigurationRequest)
	if err != nil {
		impl.logger.Errorw("error in checking if isImageDigestPolicyConfiguredForPipeline", "err", err, "clusterId", clusterId, "envId", envId)
		return false, err
	}
	return digestPolicyConfigurations.DigestConfiguredForEnvOrCluster, nil

}
