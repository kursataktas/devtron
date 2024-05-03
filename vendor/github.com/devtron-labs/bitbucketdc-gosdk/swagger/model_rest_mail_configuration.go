/*
 * Bitbucket Data Center
 *
 * This is the reference document for the Atlassian Bitbucket REST API. The REST API is for developers who want to:    - integrate Bitbucket with other applications;   - create scripts that interact with Bitbucket; or   - develop plugins that enhance the Bitbucket UI, using REST to interact with the backend.    You can read more about developing Bitbucket plugins in the [Bitbucket Developer Documentation](https://developer.atlassian.com/bitbucket/server/docs/latest/).
 *
 * API version: 8.19
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RestMailConfiguration struct {
	Port            int32  `json:"port,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	Hostname        string `json:"hostname,omitempty"`
	RequireStartTls bool   `json:"requireStartTls,omitempty"`
	UseStartTls     bool   `json:"useStartTls,omitempty"`
	SenderAddress   string `json:"senderAddress,omitempty"`
	Password        string `json:"password,omitempty"`
	Username        string `json:"username,omitempty"`
}
