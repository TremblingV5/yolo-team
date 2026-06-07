package handler

import (
	"context"

	"yolo-team/yolo-cli/internal/common"
	"yolo-team/yolo-cli/internal/model"
	"yolo-team/yolo-cli/internal/repository"
)

type ProjectHandler struct {
	repo *repository.ProjectRepo
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{repo: repository.NewProjectRepo()}
}

type CreateProjectReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetProjectReq struct {
	Key string `uri:"key"`
}

type UpdateProjectReq struct {
	Key         string  `uri:"key"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type DeleteProjectReq struct {
	Key string `uri:"key"`
}

// ListProjects godoc
//
//	@Id				ListProjects
//	@Summary		List all projects
//	@Description	Get a list of all projects
//	@Tags			Projects
//	@Produce		json
//	@Success		200	{object}	common.Response{data=[]model.Project}
//	@Failure		500	{object}	common.Response
//	@Router			/api/v1/projects [get]
func (h *ProjectHandler) List(ctx context.Context, req struct{}) ([]model.Project, error) {
	return h.repo.List()
}

// CreateProject godoc
//
//	@Id				CreateProject
//	@Summary		Create a project
//	@Description	Create a new project
//	@Tags			Projects
//	@Accept			json
//	@Produce		json
//	@Param			project	body		CreateProjectReq	true	"Project info"
//	@Success		200		{object}	common.Response{data=model.Project}
//	@Failure		400		{object}	common.Response
//	@Router			/api/v1/projects [post]
func (h *ProjectHandler) Create(ctx context.Context, req CreateProjectReq) (*model.Project, error) {
	project := model.NewProject(req.Name, req.Description)
	if err := project.Validate(); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	return h.repo.Create(project)
}

// GetProject godoc
//
//	@Id				GetProject
//	@Summary		Get a project
//	@Description	Get project details by key
//	@Tags			Projects
//	@Produce		json
//	@Param			key	path		string	true	"Project key"
//	@Success		200	{object}	common.Response{data=model.Project}
//	@Failure		404	{object}	common.Response
//	@Router			/api/v1/projects/{key} [get]
func (h *ProjectHandler) Get(ctx context.Context, req GetProjectReq) (*model.Project, error) {
	project, err := h.repo.GetByKey(req.Key)
	if err != nil {
		return nil, common.NewAppError(40401, "project not found")
	}

	return project, nil
}

// UpdateProject godoc
//
//	@Id				UpdateProject
//	@Summary		Update a project
//	@Description	Update project name and/or description
//	@Tags			Projects
//	@Accept			json
//	@Produce		json
//	@Param			key		path		string				true	"Project key"
//	@Param			project	body		UpdateProjectReq	true	"Project update info"
//	@Success		200		{object}	common.Response{data=model.Project}
//	@Failure		400		{object}	common.Response
//	@Failure		404		{object}	common.Response
//	@Router			/api/v1/projects/{key} [put]
func (h *ProjectHandler) Update(ctx context.Context, req UpdateProjectReq) (*model.Project, error) {
	project, err := h.repo.GetByKey(req.Key)
	if err != nil {
		return nil, common.NewAppError(40401, "project not found")
	}

	project.ApplyUpdate(req.Name, req.Description)
	if err := project.Validate(); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	if err := h.repo.Save(project); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	return project, nil
}

// DeleteProject godoc
//
//	@Id				DeleteProject
//	@Summary		Delete a project
//	@Description	Delete a project by key
//	@Tags			Projects
//	@Produce		json
//	@Param			key	path		string	true	"Project key"
//	@Success		200	{object}	common.Response
//	@Failure		404	{object}	common.Response
//	@Router			/api/v1/projects/{key} [delete]
func (h *ProjectHandler) Delete(ctx context.Context, req DeleteProjectReq) (struct{}, error) {
	if err := h.repo.Delete(req.Key); err != nil {
		return struct{}{}, common.NewAppError(40401, "project not found")
	}

	return struct{}{}, nil
}


