/*
 * Bitbucket Data Center
 *
 * This is the reference document for the Atlassian Bitbucket REST API. The REST API is for developers who want to:    - integrate Bitbucket with other applications;   - create scripts that interact with Bitbucket; or   - develop plugins that enhance the Bitbucket UI, using REST to interact with the backend.    You can read more about developing Bitbucket plugins in the [Bitbucket Developer Documentation](https://developer.atlassian.com/bitbucket/server/docs/latest/).
 *
 * API version: 8.19
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type RestUserReactionComment struct {
	Version            int32                                        `json:"version,omitempty"`
	Parent             *RestUserReactionCommentParent               `json:"parent,omitempty"`
	Id                 int64                                        `json:"id,omitempty"`
	State              string                                       `json:"state,omitempty"`
	Severity           string                                       `json:"severity,omitempty"`
	ThreadResolvedDate int64                                        `json:"threadResolvedDate,omitempty"`
	ThreadResolver     *RestUserReactionCommentParentThreadResolver `json:"threadResolver,omitempty"`
	// Indicates if this comment thread has been marked as resolved or not
	ThreadResolved bool                                         `json:"threadResolved,omitempty"`
	CreatedDate    int64                                        `json:"createdDate,omitempty"`
	ResolvedDate   int64                                        `json:"resolvedDate,omitempty"`
	Resolver       *RestUserReactionCommentParentThreadResolver `json:"resolver,omitempty"`
	UpdatedDate    int64                                        `json:"updatedDate,omitempty"`
	Comments       []RestComment                                `json:"comments,omitempty"`
	Text           string                                       `json:"text,omitempty"`
	Anchor         *RestUserReactionCommentParentAnchor         `json:"anchor,omitempty"`
	Author         *RestUserReactionCommentParentThreadResolver `json:"author,omitempty"`
	Html           string                                       `json:"html,omitempty"`
	Anchored       bool                                         `json:"anchored,omitempty"`
	Pending        bool                                         `json:"pending,omitempty"`
	Reply          bool                                         `json:"reply,omitempty"`
	Properties     *interface{}                                 `json:"properties,omitempty"`
}
