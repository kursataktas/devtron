/*
 * Bitbucket Data Center
 *
 * This is the reference document for the Atlassian Bitbucket REST API. The REST API is for developers who want to:    - integrate Bitbucket with other applications;   - create scripts that interact with Bitbucket; or   - develop plugins that enhance the Bitbucket UI, using REST to interact with the backend.    You can read more about developing Bitbucket plugins in the [Bitbucket Developer Documentation](https://developer.atlassian.com/bitbucket/server/docs/latest/).
 *
 * API version: 8.19
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RestSecretScanningRule struct {
	// Human readable name for the rule
	Name string `json:"name,omitempty"`
	// The ID of the rule
	Id int64 `json:"id,omitempty"`
	// If present, regular expression for matching a secret on a code line
	LineRegex string `json:"lineRegex,omitempty"`
	// If present, regular expression matching file names
	PathRegex string                       `json:"pathRegex,omitempty"`
	Scope     *RestSecretScanningRuleScope `json:"scope,omitempty"`
}
