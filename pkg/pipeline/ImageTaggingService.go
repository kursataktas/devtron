package pipeline

import (
	"encoding/json"
	"errors"
	repository "github.com/devtron-labs/devtron/internal/sql/repository/imageTagging"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	repository2 "github.com/devtron-labs/devtron/pkg/cluster/repository"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"time"
)

const TagsKey = "tags"
const CommentKey = "comments"

type ImageTaggingResponseDTO struct {
	ImageReleaseTags []repository.ImageTag   `json:"imageReleaseTags"`
	AppReleaseTags   []repository.ImageTag   `json:"appReleaseTags"`
	ImageComment     repository.ImageComment `json:"imageComments"`
	ProdEnvExists    bool                    `json:"prodEnvExists"`
}

type ImageTaggingRequestDTO struct {
	CreateTags     []repository.ImageTag
	SoftDeleteTags []repository.ImageTag
	ImageComment   repository.ImageComment
	HardDeleteTags []repository.ImageTag
}

type ImageTaggingService interface {
	GetTagsData(ciPipelineId, appId, artifactId int) (*ImageTaggingResponseDTO, error)
	CreateUpdateImageTagging(ciPipelineId, appId, artifactId, userId int, imageTaggingRequest *ImageTaggingRequestDTO) (*ImageTaggingResponseDTO, error)
	GetProdEnvFromParentAndLinkedWorkflow(ciPipelineId int) (bool, error)
	GetProdEnvByCdPipelineId(pipelineId int) (bool, error)
	ValidateImageTaggingRequest(imageTaggingRequest *ImageTaggingRequestDTO) (bool, error)
	GetTagsByArtifactId(artifactId int) ([]repository.ImageTag, error)
	GetTagsByAppId(appId int) ([]repository.ImageTag, error)
	GetTaggingDataMapByAppId(appId int) (map[int]*ImageTaggingResponseDTO, error)
}

type ImageTaggingServiceImpl struct {
	imageTaggingRepo      repository.ImageTaggingRepository
	ciPipelineRepository  pipelineConfig.CiPipelineRepository
	cdPipelineRepository  pipelineConfig.PipelineRepository
	environmentRepository repository2.EnvironmentRepository
	logger                *zap.SugaredLogger
}

func NewImageTaggingServiceImpl(imageTaggingRepo repository.ImageTaggingRepository,
	ciPipelineRepository pipelineConfig.CiPipelineRepository,
	cdPipelineRepository pipelineConfig.PipelineRepository,
	environmentRepository repository2.EnvironmentRepository,
	logger *zap.SugaredLogger) *ImageTaggingServiceImpl {
	return &ImageTaggingServiceImpl{
		imageTaggingRepo:      imageTaggingRepo,
		ciPipelineRepository:  ciPipelineRepository,
		cdPipelineRepository:  cdPipelineRepository,
		environmentRepository: environmentRepository,
		logger:                logger,
	}
}

func (impl ImageTaggingServiceImpl) GetTagsData(ciPipelineId, appId, artifactId int) (*ImageTaggingResponseDTO, error) {
	resp := &ImageTaggingResponseDTO{}
	imageComment, err := impl.imageTaggingRepo.GetImageComment(artifactId)
	if err != nil && err != pg.ErrNoRows {
		//log error
		impl.logger.Errorw("error in fetching image comment using artifactId", "err", err, "artifactId", artifactId)
		return resp, err
	}
	appReleaseTags, err := impl.GetTagsByAppId(appId)
	if err != nil {
		//log error
		impl.logger.Errorw("error in fetching image tags using appId", "err", err, "appId", appId)
		return resp, err
	}
	imageReleaseTags, err := impl.GetTagsByArtifactId(artifactId)
	if err != nil {
		//log error
		impl.logger.Errorw("error in fetching image tags using artifactId", "err", err, "artifactId", artifactId)
		return resp, err
	}
	prodEnvExists, err := impl.GetProdEnvFromParentAndLinkedWorkflow(ciPipelineId)
	if err != nil {
		//log error
		impl.logger.Errorw("error in GetProdEnvFromParentAndLinkedWorkflow", "err", err, "ciPipelineId", ciPipelineId)
		return resp, err
	}
	resp.AppReleaseTags = appReleaseTags
	resp.ImageReleaseTags = imageReleaseTags
	resp.ImageComment = imageComment
	resp.ProdEnvExists = prodEnvExists
	return resp, err
}

func (impl ImageTaggingServiceImpl) GetTagsByArtifactId(artifactId int) ([]repository.ImageTag, error) {
	imageReleaseTags, err := impl.imageTaggingRepo.GetTagsByArtifactId(artifactId)
	if err != nil && err != pg.ErrNoRows {
		//log error
		impl.logger.Errorw("error in fetching image tags using artifactId", "err", err, "artifactId", artifactId)
		return imageReleaseTags, err
	}
	return imageReleaseTags, nil
}

func (impl ImageTaggingServiceImpl) GetTagsByAppId(appId int) ([]repository.ImageTag, error) {
	appReleaseTags, err := impl.imageTaggingRepo.GetTagsByAppId(appId)
	if err != nil && err != pg.ErrNoRows {
		//log error
		impl.logger.Errorw("error in fetching image tags using appId", "err", err, "appId", appId)
		return appReleaseTags, err
	}
	return appReleaseTags, nil
}
func (impl ImageTaggingServiceImpl) GetTaggingDataMapByAppId(appId int) (map[int]*ImageTaggingResponseDTO, error) {
	tags, err := impl.GetTagsByAppId(appId)
	if err != nil {
		impl.logger.Errorw("error occurred in getting image tags by appId", "appId", appId, "err", err)
		return nil, err
	}
	result := make(map[int]*ImageTaggingResponseDTO)
	for _, tag := range tags {
		if _, ok := result[tag.ArtifactId]; !ok {
			result[tag.ArtifactId] = &ImageTaggingResponseDTO{
				ImageReleaseTags: make([]repository.ImageTag, 0),
			}
		}
		result[tag.ArtifactId].ImageReleaseTags = append(result[tag.ArtifactId].ImageReleaseTags, tag)
	}

	imageComments, err := impl.imageTaggingRepo.GetImageCommentsByAppId(appId)
	if err != nil && err != pg.ErrNoRows {
		//log error
		impl.logger.Errorw("error in fetching imageComments using appId", "appId", appId)
		return nil, err
	}

	//it may be possible that there are no tags for a artifact,but comment exists
	for _, comment := range imageComments {
		if _, ok := result[comment.ArtifactId]; !ok {
			result[comment.ArtifactId] = &ImageTaggingResponseDTO{}
		}
		result[comment.ArtifactId].ImageComment = comment
	}
	return result, nil
}

func (impl ImageTaggingServiceImpl) ValidateImageTaggingRequest(imageTaggingRequest *ImageTaggingRequestDTO) (bool, error) {
	//validate create tags
	for _, tags := range imageTaggingRequest.CreateTags {
		if tags.Id != 0 || tags.AppId == 0 || tags.ArtifactId == 0 {
			return false, errors.New("bad request,create tags cannot contain id, missing artifactId or appId")
		}
		err := tagNameValidation(tags.TagName)
		if err != nil {
			return false, err
		}
	}
	//validate update tags
	for _, tags := range imageTaggingRequest.SoftDeleteTags {
		if tags.Id == 0 {
			return false, errors.New("bad request,tags requested to delete should contain id")
		}
		err := tagNameValidation(tags.TagName)
		if err != nil {
			return false, err
		}
	}

	for _, tags := range imageTaggingRequest.HardDeleteTags {
		if tags.Id == 0 {
			return false, errors.New("bad request,tags requested to delete should contain id")
		}
		err := tagNameValidation(tags.TagName)
		if err != nil {
			return false, err
		}
	}

	//TODO: validate comment, currently no validation on comment
	return true, nil
}

func tagNameValidation(tag string) error {
	err := errors.New("tag name should be max of 128 characters long,tag name should not start with '.' and '-'")
	if len(tag) > 128 || len(tag) == 0 || tag[0] == '.' || tag[1] == '-' {
		return err
	}
	return nil
}
func (impl ImageTaggingServiceImpl) CreateUpdateImageTagging(ciPipelineId, appId, artifactId, userId int, imageTaggingRequest *ImageTaggingRequestDTO) (*ImageTaggingResponseDTO, error) {

	db := impl.imageTaggingRepo.GetDbConnection()
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		//add logs
		impl.logger.Errorw("error in creating transaction", "err", err)
		return nil, err
	}

	//first perform delete and then perform create operation.
	//case : user can delete existing tag and then create a new tag with same name, this is a valid request

	//soft delete tags
	softDeleteAuditTags := make([]string, len(imageTaggingRequest.SoftDeleteTags))
	for i, tag := range imageTaggingRequest.SoftDeleteTags {
		tag.AppId = appId
		tag.Active = true
		tag.ArtifactId = artifactId
		err := impl.imageTaggingRepo.UpdateReleaseTag(tx, &tag)
		if err != nil {
			//log
			impl.logger.Errorw("error in updating releaseTag", "err", err, "payLoad", tag)
			return nil, err
		}
		softDeleteAuditTags[i] = tag.TagName
	}

	//hard delete tags
	hardDeleteAuditTags := make([]string, len(imageTaggingRequest.HardDeleteTags))
	for i, tag := range imageTaggingRequest.HardDeleteTags {
		tag.AppId = appId
		tag.ArtifactId = artifactId
		err := impl.imageTaggingRepo.DeleteReleaseTag(tx, &tag)
		if err != nil {
			//log
			impl.logger.Errorw("error in deleting releaseTag", "err", err, "payLoad", tag)
			return nil, err
		}
		hardDeleteAuditTags[i] = tag.TagName
	}

	//save release tags
	createAuditTags := make([]string, len(imageTaggingRequest.HardDeleteTags))
	for i, tag := range imageTaggingRequest.CreateTags {
		tag.AppId = appId
		tag.ArtifactId = artifactId
		err := impl.imageTaggingRepo.SaveReleaseTag(tx, &tag)
		if err != nil {
			//log
			impl.logger.Errorw("error in saving releaseTag", "err", err, "releaseTag", tag)
			return nil, err
		}
		createAuditTags[i] = tag.TagName
	}

	imageTaggingRequest.ImageComment.ArtifactId = artifactId
	imageTaggingRequest.ImageComment.UserId = userId
	//save or update comment
	if imageTaggingRequest.ImageComment.Id > 0 {
		savedComment, err := impl.imageTaggingRepo.GetImageComment(artifactId)
		if err != nil {
			impl.logger.Errorw("error in getting imageComment by artifactId", "err", err, "artifactId", artifactId)
			return nil, err
		}
		//update only if the comment is different from saved comment
		if savedComment.Comment != imageTaggingRequest.ImageComment.Comment {
			err = impl.imageTaggingRepo.UpdateImageComment(tx, &imageTaggingRequest.ImageComment)
			if err != nil {
				//log
				impl.logger.Errorw("error in updating imageComment ", "err", err, "ImageComment", imageTaggingRequest.ImageComment)
				return nil, err
			}
			//save comment audit
			err = impl.saveImageCommentAudit(tx, imageTaggingRequest.ImageComment.Comment, userId, artifactId, repository.ActionEdit)
			if err != nil {
				//log
				impl.logger.Errorw("error in saving image tagging audit", "err", err, "comment", imageTaggingRequest.ImageComment.Comment, "action", "editAction")
				return nil, err
			}
		}
	} else {
		err := impl.imageTaggingRepo.SaveImageComment(tx, &imageTaggingRequest.ImageComment)
		if err != nil {
			//log
			impl.logger.Errorw("error in saving imageComment ", "err", err, "ImageComment", imageTaggingRequest.ImageComment)
			return nil, err
		}
		//save comment audit
		err = impl.saveImageCommentAudit(tx, imageTaggingRequest.ImageComment.Comment, userId, artifactId, repository.ActionSave)
		if err != nil {
			//log
			impl.logger.Errorw("error in saving image tagging audit", "err", err, "comment", imageTaggingRequest.ImageComment.Comment, "action", "editAction")
			return nil, err
		}
	}

	//save tags audit
	err = impl.saveImageTagAudit(tx, softDeleteAuditTags, hardDeleteAuditTags, createAuditTags, userId, artifactId)
	if err != nil {
		//log
		impl.logger.Errorw("error in saving image tagging audit", "err", err)
		return nil, err
	}

	//commit transaction
	err = tx.Commit()
	if err != nil {
		//log
		impl.logger.Errorw("error in committing transaction", "err", err)
		return nil, err
	}
	return impl.GetTagsData(ciPipelineId, appId, artifactId)
}

func (impl ImageTaggingServiceImpl) saveImageTagAudit(tx *pg.Tx, softDeleteTags, hardDeleteTags, createTags []string, userId, artifactId int) error {

	if len(softDeleteTags) > 0 {
		dataMap := make(map[string]interface{})
		dataMap[TagsKey] = softDeleteTags
		dataBytes, err := json.Marshal(&dataMap)
		if err != nil {
			impl.logger.Errorw("error in marshaling imageTagging data", "error", err, "data", dataMap)
			return err
		}
		auditLog := &repository.ImageTaggingAudit{
			Data:       string(dataBytes),
			DataType:   repository.TagType,
			UpdatedBy:  userId,
			UpdatedOn:  time.Now(),
			ArtifactId: artifactId,
			Action:     repository.ActionSoftDelete,
		}
		err = impl.imageTaggingRepo.SaveAuditLog(tx, auditLog)
		if err != nil {
			impl.logger.Errorw("error occured in saving image tagging audit", "err", err, "auditLog", auditLog)
			return err
		}
	}

	if len(hardDeleteTags) > 0 {
		dataMap := make(map[string]interface{})
		dataMap[TagsKey] = hardDeleteTags
		dataBytes, err := json.Marshal(&dataMap)
		if err != nil {
			impl.logger.Errorw("error in marshaling imageTagging data", "error", err, "data", dataMap)
			return err
		}
		auditLog := &repository.ImageTaggingAudit{
			Data:       string(dataBytes),
			DataType:   repository.TagType,
			UpdatedBy:  userId,
			UpdatedOn:  time.Now(),
			ArtifactId: artifactId,
			Action:     repository.ActionHardDelete,
		}
		err = impl.imageTaggingRepo.SaveAuditLog(tx, auditLog)
		if err != nil {
			impl.logger.Errorw("error occurred in saving image tagging audit", "err", err, "auditLog", auditLog)
			return err
		}
	}

	if len(createTags) > 0 {
		dataMap := make(map[string]interface{})
		dataMap[TagsKey] = createTags
		dataBytes, err := json.Marshal(&dataMap)
		if err != nil {
			impl.logger.Errorw("error in marshaling imageTagging data", "error", err, "data", dataMap)
			return err
		}
		auditLog := &repository.ImageTaggingAudit{
			Data:       string(dataBytes),
			DataType:   repository.TagType,
			UpdatedBy:  userId,
			UpdatedOn:  time.Now(),
			ArtifactId: artifactId,
			Action:     repository.ActionSave,
		}
		err = impl.imageTaggingRepo.SaveAuditLog(tx, auditLog)
		if err != nil {
			impl.logger.Errorw("error occurred in saving image tagging audit", "err", err, "auditLog", auditLog)
			return err
		}
	}

	return nil

}

func (impl ImageTaggingServiceImpl) saveImageCommentAudit(tx *pg.Tx, imageComment string, userId, artifactId int, action repository.ImageTaggingAction) error {

	dataMap := make(map[string]string)
	dataMap[CommentKey] = imageComment
	dataBytes, err := json.Marshal(&dataMap)
	if err != nil {
		impl.logger.Errorw("error in marshaling imageTagging data", "error", err, "data", dataMap)
		return err
	}
	auditLog := &repository.ImageTaggingAudit{
		Data:       string(dataBytes),
		DataType:   repository.CommentType,
		UpdatedBy:  userId,
		UpdatedOn:  time.Now(),
		ArtifactId: artifactId,
		Action:     action,
	}
	err = impl.imageTaggingRepo.SaveAuditLog(tx, auditLog)
	if err != nil {
		impl.logger.Errorw("error occurred in saving image tagging audit", "err", err, "auditLog", auditLog)
		return err
	}

	return nil
}

func (impl ImageTaggingServiceImpl) GetProdEnvFromParentAndLinkedWorkflow(ciPipelineId int) (bool, error) {
	prodEnvExists := false
	pipelines, err := impl.ciPipelineRepository.FindByParentCiPipelineId(ciPipelineId)
	if err != nil {
		//add log
		impl.logger.Errorw("error in getting all linked ciPipelineIds", "err", err, "ciPipelineId", ciPipelineId)
		return prodEnvExists, err
	}

	//get all the pipeline ids liked with the requested ciPipelineId
	pipelineIds := make([]int, len(pipelines)+1)
	pipelineIds[0] = ciPipelineId
	for i := 0; i < len(pipelines); i++ {
		pipelineIds[i+1] = pipelines[i].Id
	}

	envs, err := impl.environmentRepository.FindEnvLinkedWithCiPipelines(pipelineIds)
	if err != nil {
		//add log
		impl.logger.Errorw("error in getting envs using ciPipelineIds", "err", err, "ciPipelineIds", pipelineIds)
		return prodEnvExists, err
	}

	for _, env := range envs {
		//env id prod ,return true
		if env.Default {
			prodEnvExists = true
			break
		}
	}

	return prodEnvExists, nil

}

func (impl ImageTaggingServiceImpl) GetProdEnvByCdPipelineId(pipelineId int) (bool, error) {
	pipeline, err := impl.cdPipelineRepository.FindById(pipelineId)
	if err != nil {
		impl.logger.Errorw("error occurred in fetching cdPipeline with pipelineId", "err", err, "pipelineId", pipelineId)
		return false, err
	}
	if pipeline.Environment.Default {
		return true, nil
	}

	//CiPipelineId will be zero for external webhook ci
	if pipeline.CiPipelineId > 0 {
		return impl.GetProdEnvFromParentAndLinkedWorkflow(pipeline.CiPipelineId)
	}

	return false, nil

}
