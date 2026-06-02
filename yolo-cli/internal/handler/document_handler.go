package handler

import (
	"context"
	"os"
	"time"

	"yolo-team/yolo-cli/internal/common"
	"yolo-team/yolo-cli/internal/model"
	"yolo-team/yolo-cli/internal/repository"
	"yolo-team/yolo-cli/internal/utils/fileutils"
)

type DocumentHandler struct {
	repo      *repository.DocumentRepo
	projRepo  *repository.ProjectRepo
	workspace string
}

func NewDocumentHandler(workspace string) *DocumentHandler {
	return &DocumentHandler{
		repo:      repository.NewDocumentRepo(),
		projRepo:  repository.NewProjectRepo(),
		workspace: workspace,
	}
}

type ListDocReq struct {
	ProjectKey string `uri:"key"`
}

type CreateDocReq struct {
	ProjectKey string `uri:"key"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Creator    string `json:"creator"`
}

type GetDocReq struct {
	Key string `uri:"key"`
}

type UpdateDocReq struct {
	Key     string  `uri:"key"`
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type DeleteDocReq struct {
	Key string `uri:"key"`
}

type DocumentResponse struct {
	ID        int64     `json:"id"`
	Key       string    `json:"key"`
	ProjectID int64     `json:"project_id"`
	Title     string    `json:"title"`
	Creator   string    `json:"creator"`
	Content   string    `json:"content"`
	FilePath  string    `json:"file_path"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListDocuments godoc
//
//	@Id				ListDocuments
//	@Summary		List documents
//	@Description	List documents in a project
//	@Tags			Documents
//	@Produce		json
//	@Param			key	path		string	true	"Project key"
//	@Success		200	{object}	common.Response{data=[]model.Document}
//	@Failure		404	{object}	common.Response
//	@Failure		500	{object}	common.Response
//	@Router			/api/v1/projects/{key}/documents [get]
func (h *DocumentHandler) List(ctx context.Context, req ListDocReq) ([]model.Document, error) {
	project, err := h.projRepo.GetByKey(req.ProjectKey)
	if err != nil {
		return nil, common.NewAppError(40401, "project not found")
	}

	return h.repo.ListByProject(project.ID)
}

// CreateDocument godoc
//
//	@Id				CreateDocument
//	@Summary		Create a document
//	@Description	Create a new document in a project
//	@Tags			Documents
//	@Accept			json
//	@Produce		json
//	@Param			key		path		string				true	"Project key"
//	@Param			document	body		CreateDocReq	true	"Document info"
//	@Success		200			{object}	common.Response{data=DocumentResponse}
//	@Failure		400			{object}	common.Response
//	@Failure		404			{object}	common.Response
//	@Router			/api/v1/projects/{key}/documents [post]
func (h *DocumentHandler) Create(ctx context.Context, req CreateDocReq) (DocumentResponse, error) {
	project, err := h.projRepo.GetByKey(req.ProjectKey)
	if err != nil {
		return DocumentResponse{}, common.NewAppError(40401, "project not found")
	}

	doc := model.NewDocument(project.ID, req.Title)
	if req.Creator != "" {
		doc.Creator = req.Creator
	}
	if err := doc.Validate(); err != nil {
		return DocumentResponse{}, common.NewAppError(40001, err.Error())
	}

	result, err := h.repo.Create(doc)
	if err != nil {
		return DocumentResponse{}, common.NewAppError(40001, err.Error())
	}

	filePath := fileutils.BuildDocPath(h.workspace, doc.ProjectID, req.Title)
	fileutils.Write(filePath, req.Content)
	return newDocumentResponse(result, req.Content, filePath), nil
}

// GetDocument godoc
//
//	@Id				GetDocument
//	@Summary		Get a document
//	@Description	Get document details by key
//	@Tags			Documents
//	@Produce		json
//	@Param			key	path		string	true	"Document key"
//	@Success		200	{object}	common.Response{data=DocumentResponse}
//	@Failure		404	{object}	common.Response
//	@Router			/api/v1/documents/{key} [get]
func (h *DocumentHandler) Get(ctx context.Context, req GetDocReq) (DocumentResponse, error) {
	doc, err := h.repo.GetByKey(req.Key)
	if err != nil {
		return DocumentResponse{}, common.NewAppError(40401, "document not found")
	}

	filePath := fileutils.BuildDocPath(h.workspace, doc.ProjectID, doc.Title)
	content, _ := fileutils.Read(filePath)
	return newDocumentResponse(doc, content, filePath), nil
}

// UpdateDocument godoc
//
//	@Id				UpdateDocument
//	@Summary		Update a document
//	@Description	Update document title and/or content
//	@Tags			Documents
//	@Accept			json
//	@Produce		json
//	@Param			key			path		string			true	"Document key"
//	@Param			document	body		UpdateDocReq	true	"Document update info"
//	@Success		200			{object}	common.Response{data=DocumentResponse}
//	@Failure		400			{object}	common.Response
//	@Failure		404			{object}	common.Response
//	@Router			/api/v1/documents/{key} [put]
func (h *DocumentHandler) Update(ctx context.Context, req UpdateDocReq) (DocumentResponse, error) {
	doc, err := h.repo.GetByKey(req.Key)
	if err != nil {
		return DocumentResponse{}, common.NewAppError(40401, "document not found")
	}

	oldPath := fileutils.BuildDocPath(h.workspace, doc.ProjectID, doc.Title)

	if req.Title != nil {
		doc.Title = *req.Title
		if err := doc.Validate(); err != nil {
			return DocumentResponse{}, common.NewAppError(40001, err.Error())
		}

		newPath := fileutils.BuildDocPath(h.workspace, doc.ProjectID, doc.Title)
		if oldPath != newPath {
			os.Rename(oldPath, newPath)
		}
	}

	if req.Content != nil {
		currPath := fileutils.BuildDocPath(h.workspace, doc.ProjectID, doc.Title)
		fileutils.Write(currPath, *req.Content)
	}

	if err := h.repo.Save(doc); err != nil {
		return DocumentResponse{}, common.NewAppError(40001, err.Error())
	}

	filePath := fileutils.BuildDocPath(h.workspace, doc.ProjectID, doc.Title)
	content, _ := fileutils.Read(filePath)
	return newDocumentResponse(doc, content, filePath), nil
}

// DeleteDocument godoc
//
//	@Id				DeleteDocument
//	@Summary		Delete a document
//	@Description	Delete a document by key
//	@Tags			Documents
//	@Produce		json
//	@Param			key	path		string	true	"Document key"
//	@Success		200	{object}	common.Response
//	@Failure		404	{object}	common.Response
//	@Router			/api/v1/documents/{key} [delete]
func (h *DocumentHandler) Delete(ctx context.Context, req DeleteDocReq) (struct{}, error) {
	doc, err := h.repo.GetByKey(req.Key)
	if err != nil {
		return struct{}{}, common.NewAppError(40401, "document not found")
	}

	filePath := fileutils.BuildDocPath(h.workspace, doc.ProjectID, doc.Title)
	os.Remove(filePath)
	h.repo.Delete(req.Key)
	return struct{}{}, nil
}

func newDocumentResponse(doc *model.Document, content, filePath string) DocumentResponse {
	return DocumentResponse{
		ID:        doc.ID,
		Key:       doc.Key,
		ProjectID: doc.ProjectID,
		Title:     doc.Title,
		Creator:   doc.Creator,
		Content:   content,
		FilePath:  filePath,
		SortOrder: doc.SortOrder,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}
