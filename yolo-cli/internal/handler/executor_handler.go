package handler

import (
	"context"

	"yolo-team/yolo-cli/internal/common"
	"yolo-team/yolo-cli/internal/model"
	"yolo-team/yolo-cli/internal/repository"
	"yolo-team/yolo-cli/internal/utils/dbutils"
)

type ExecutorHandler struct {
	repo *repository.ExecutorRepo
}

func NewExecutorHandler() *ExecutorHandler {
	return &ExecutorHandler{repo: repository.NewExecutorRepo()}
}

type CreateExecutorReq struct {
	Name string `json:"name"`
	Soul string `json:"soul"`
}

type DeleteExecutorReq struct {
	Name string `uri:"name"`
}

type UpdateExecutorReq struct {
	Name string `uri:"name"`
	Role string `json:"role"`
	Soul string `json:"soul"`
}

// ListExecutors godoc
//
//	@Id				ListExecutors
//	@Summary		List all executors
//	@Description	Get a list of all executors
//	@Tags			Executors
//	@Produce		json
//	@Success		200	{object}	common.Response{data=[]model.Executor}
//	@Failure		500	{object}	common.Response
//	@Router			/api/v1/executors [get]
func (h *ExecutorHandler) List(ctx context.Context, req struct{}) ([]model.Executor, error) {
	return h.repo.List()
}

// CreateExecutor godoc
//
//	@Id				CreateExecutor
//	@Summary		Create an executor
//	@Description	Create a new executor
//	@Tags			Executors
//	@Accept			json
//	@Produce		json
//	@Param			executor	body		CreateExecutorReq	true	"Executor info"
//	@Success		200			{object}	common.Response{data=model.Executor}
//	@Failure		400			{object}	common.Response
//	@Failure		409			{object}	common.Response
//	@Router			/api/v1/executors [post]
func (h *ExecutorHandler) Create(ctx context.Context, req CreateExecutorReq) (*model.Executor, error) {
	executor := model.NewExecutor(req.Name, req.Soul)
	if err := executor.Validate(); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	result, err := h.repo.Create(executor)
	if err != nil {
		if dbutils.IsDuplicate(err) {
			return nil, common.NewAppError(40901, "executor name already exists")
		}

		return nil, common.NewAppError(40001, err.Error())
	}

	return result, nil
}

// UpdateExecutor godoc
//
//	@Id				UpdateExecutor
//	@Summary		Update an executor
//	@Description	Update executor role and/or soul
//	@Tags			Executors
//	@Accept			json
//	@Produce		json
//	@Param			name		path		string				true	"Executor name"
//	@Param			executor	body		UpdateExecutorReq	true	"Executor update info"
//	@Success		200			{object}	common.Response{data=model.Executor}
//	@Failure		400			{object}	common.Response
//	@Failure		404			{object}	common.Response
//	@Router			/api/v1/executors/{name} [put]
func (h *ExecutorHandler) Update(ctx context.Context, req UpdateExecutorReq) (*model.Executor, error) {
	executor, err := h.repo.GetByName(req.Name)
	if err != nil {
		return nil, common.NewAppError(40401, "executor not found")
	}

	if req.Role != "" {
		executor.SetRole(req.Role)
		if err := executor.Validate(); err != nil {
			return nil, common.NewAppError(40001, err.Error())
		}
	}

	if req.Soul != "" {
		executor.Soul = req.Soul
	}

	if err := h.repo.UpdateRole(executor.ID, executor.Role); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	if req.Soul != "" {
		h.repo.UpdateSoul(executor.ID, req.Soul)
	}

	return h.repo.GetByID(executor.ID)
}

// DeleteExecutor godoc
//
//	@Id				DeleteExecutor
//	@Summary		Delete an executor
//	@Description	Delete an executor by name
//	@Tags			Executors
//	@Produce		json
//	@Param			name	path		string	true	"Executor name"
//	@Success		200		{object}	common.Response
//	@Failure		404		{object}	common.Response
//	@Router			/api/v1/executors/{name} [delete]
func (h *ExecutorHandler) Delete(ctx context.Context, req DeleteExecutorReq) (struct{}, error) {
	if err := h.repo.DeleteByName(req.Name); err != nil {
		return struct{}{}, common.NewAppError(40401, "executor not found")
	}

	return struct{}{}, nil
}
