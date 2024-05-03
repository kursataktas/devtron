/*
 * Bitbucket Data Center
 *
 * This is the reference document for the Atlassian Bitbucket REST API. The REST API is for developers who want to:    - integrate Bitbucket with other applications;   - create scripts that interact with Bitbucket; or   - develop plugins that enhance the Bitbucket UI, using REST to interact with the backend.    You can read more about developing Bitbucket plugins in the [Bitbucket Developer Documentation](https://developer.atlassian.com/bitbucket/server/docs/latest/).
 *
 * API version: 8.19
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RestChangesetToCommit struct {
	Message            string                          `json:"message,omitempty"`
	CommitterTimestamp int64                           `json:"committerTimestamp,omitempty"`
	Committer          *RestChangesetToCommitCommitter `json:"committer,omitempty"`
	AuthorTimestamp    int64                           `json:"authorTimestamp,omitempty"`
	Parents            []RestMinimalCommit             `json:"parents,omitempty"`
	Author             *RestChangesetToCommitCommitter `json:"author,omitempty"`
	Id                 string                          `json:"id,omitempty"`
	DisplayId          string                          `json:"displayId,omitempty"`
}
