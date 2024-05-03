/*
 * Bitbucket Data Center
 *
 * This is the reference document for the Atlassian Bitbucket REST API. The REST API is for developers who want to:    - integrate Bitbucket with other applications;   - create scripts that interact with Bitbucket; or   - develop plugins that enhance the Bitbucket UI, using REST to interact with the backend.    You can read more about developing Bitbucket plugins in the [Bitbucket Developer Documentation](https://developer.atlassian.com/bitbucket/server/docs/latest/).
 *
 * API version: 8.19
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RepositoryHookDetails struct {
	Version         string   `json:"version,omitempty"`
	Name            string   `json:"name,omitempty"`
	Key             string   `json:"key,omitempty"`
	Type_           string   `json:"type,omitempty"`
	ConfigFormKey   string   `json:"configFormKey,omitempty"`
	SupportedScopes []string `json:"supportedScopes,omitempty"`
	Description     string   `json:"description,omitempty"`
}
