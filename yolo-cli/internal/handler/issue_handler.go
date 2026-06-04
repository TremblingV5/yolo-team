package handler

import (
	"context"

	"yolo-team/yolo-cli/internal/common"
	"yolo-team/yolo-cli/internal/model"
	"yolo-team/yolo-cli/internal/repository"
	"yolo-team/yolo-cli/internal/utils/timeutils"
)

type IssueHandler struct {
	repo     *repository.IssueRepo
	execRepo *repository.ExecutorRepo
	docRepo  *repository.DocumentRepo
	taskRepo *repository.TaskRepo
}

func NewIssueHandler() *IssueHandler {
	return &IssueHandler{
		repo:     repository.NewIssueRepo(),
		execRepo: repository.NewExecutorRepo(),
		docRepo:  repository.NewDocumentRepo(),
		taskRepo: repository.NewTaskRepo(),
	}
}

type ListIssueReq struct {
	ProjectID  *int64  `form:"project_id"`
	Status     *string `form:"status"`
	ExecutorID *int64  `form:"executor_id"`
}

type CreateIssueReq struct {
	ProjectID   int64  `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	ExecutorID  *int64 `json:"executor_id"`
	Deadline    string `json:"deadline"`
}

type GetIssueReq struct {
	Key string `uri:"key"`
}

type UpdateIssueReq struct {
	Key string `uri:"key"`
	model.UpdateIssueRequest
}

type DeleteIssueReq struct {
	Key string `uri:"key"`
}

type TodoIssueReq struct {
	ExecutorID int64 `form:"executor_id"`
	Limit      int   `form:"limit"`
}

type LinkDocReq struct {
	IssueKey   string `uri:"key"`
	DocumentID int64  `json:"document_id"`
}

type UnlinkDocReq struct {
	IssueKey   string `uri:"key"`
	DocumentID int64  `form:"document_id"`
}

// Task request/response types
type ListTaskReq struct {
	IssueKey string `uri:"key"`
}

type CreateTaskReq struct {
	IssueKey    string `uri:"key"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskReq struct {
	TaskKey     string  `uri:"task_key"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

type DeleteTaskReq struct {
	TaskKey string `uri:"task_key"`
}

// ListIssues godoc
//
//	@Id				ListIssues
//	@Summary		List issues
//	@Description	List issues with optional filters
//	@Tags			Issues
//	@Produce		json
//	@Param			project_id	query		int		false	"Filter by project ID"
//	@Param			status		query		string	false	"Filter by status"
//	@Param			executor_id	query		int		false	"Filter by executor ID"
//	@Success		200			{object}	common.Response{data=[]model.Issue}
//	@Failure		500			{object}	common.Response
//	@Router			/api/v1/issues [get]
func (h *IssueHandler) List(ctx context.Context, req ListIssueReq) ([]model.Issue, error) {
	filter := repository.IssueFilter{}
	filter.ProjectID = req.ProjectID
	filter.Status = req.Status
	filter.ExecutorID = req.ExecutorID

	return h.repo.List(filter)
}

// CreateIssue godoc
//
//	@Id				CreateIssue
//	@Summary		Create an issue
//	@Description	Create a new issue
//	@Tags			Issues
//	@Accept			json
//	@Produce		json
//	@Param			issue	body		CreateIssueReq	true	"Issue info"
//	@Success		200		{object}	common.Response{data=model.Issue}
//	@Failure		400		{object}	common.Response
//	@Router			/api/v1/issues [post]
func (h *IssueHandler) Create(ctx context.Context, req CreateIssueReq) (*model.Issue, error) {
	issue := model.NewIssue(req.ProjectID, req.Title)
	issue.Description = req.Description
	issue.ExecutorID = req.ExecutorID
	if req.Priority != "" {
		issue.Priority = req.Priority
	}

	if req.Deadline != "" {
		t, err := timeutils.Parse(req.Deadline)
		if err == nil {
			issue.Deadline = &t
		}
	}

	if err := issue.Validate(); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	if err := h.repo.Create(issue); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	return h.repo.GetByKey(issue.Key)
}

// GetIssue godoc
//
//	@Id				GetIssue
//	@Summary		Get an issue
//	@Description	Get issue details by key, including tasks
//	@Tags			Issues
//	@Produce		json
//	@Param			key	path		string	true	"Issue key"
//	@Success		200	{object}	common.Response{data=model.Issue}
//	@Failure		404	{object}	common.Response
//	@Router			/api/v1/issues/{key} [get]
func (h *IssueHandler) Get(ctx context.Context, req GetIssueReq) (*model.Issue, error) {
	issue, err := h.repo.GetByKey(req.Key)
	if err != nil {
		return nil, common.NewAppError(40401, "issue not found")
	}

	return issue, nil
}

// UpdateIssue godoc
//
//	@Id				UpdateIssue
//	@Summary		Update an issue
//	@Description	Update issue fields
//	@Tags			Issues
//	@Accept			json
//	@Produce		json
//	@Param			key		path		string						true	"Issue key"
//	@Param			issue	body		model.UpdateIssueRequest	true	"Issue update info"
//	@Success		200		{object}	common.Response{data=model.Issue}
//	@Failure		400		{object}	common.Response
//	@Failure		404		{object}	common.Response
//	@Router			/api/v1/issues/{key} [put]
func (h *IssueHandler) Update(ctx context.Context, req UpdateIssueReq) (*model.Issue, error) {
	issue, err := h.repo.GetByKey(req.Key)
	if err != nil {
		return nil, common.NewAppError(40401, "issue not found")
	}

	issue.ApplyUpdate(&req.UpdateIssueRequest)
	if err := issue.Validate(); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	// Prevent moving to done if any task is not completed
	if issue.Status == model.StatusDone {
		tasks, err := h.taskRepo.ListByIssue(issue.ID)
		if err != nil {
			return nil, common.NewAppError(50001, err.Error())
		}
		for _, t := range tasks {
			if t.Status != model.TaskStatusDone {
				return nil, common.NewAppError(40001, "all tasks must be completed before closing the issue")
			}
		}
	}

	if err := h.repo.Save(issue); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	return h.repo.GetByKey(req.Key)
}

// DeleteIssue godoc
//
//	@Id				DeleteIssue
//	@Summary		Delete an issue
//	@Description	Delete an issue by key
//	@Tags			Issues
//	@Produce		json
//	@Param			key	path		string	true	"Issue key"
//	@Success		200	{object}	common.Response
//	@Failure		404	{object}	common.Response
//	@Router			/api/v1/issues/{key} [delete]
func (h *IssueHandler) Delete(ctx context.Context, req DeleteIssueReq) (struct{}, error) {
	if err := h.repo.Delete(req.Key); err != nil {
		return struct{}{}, common.NewAppError(40401, "issue not found")
	}

	return struct{}{}, nil
}

// TodoIssues godoc
//
//	@Id				TodoIssues
//	@Summary		Get todo issues
//	@Description	Get todo list for an executor
//	@Tags			Issues
//	@Produce		json
//	@Param			executor_id	query		int	true	"Executor ID"
//	@Param			limit		query		int	false	"Max results (default 20)"
//	@Success		200			{object}	common.Response{data=[]model.Issue}
//	@Failure		400			{object}	common.Response
//	@Failure		500			{object}	common.Response
//	@Router			/api/v1/issues/todo [get]
func (h *IssueHandler) Todo(ctx context.Context, req TodoIssueReq) ([]model.Issue, error) {
	if req.ExecutorID == 0 {
		return nil, common.NewAppError(40001, "executor_id is required")
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	return h.repo.GetTodo(req.ExecutorID, req.Limit)
}

// LinkDocument godoc
//
//	@Id				LinkDocument
//	@Summary		Link document to issue
//	@Description	Associate a document with an issue
//	@Tags			Issues
//	@Accept			json
//	@Produce		json
//	@Param			key			path		string		true	"Issue key"
//	@Param			document	body		LinkDocReq	true	"Document info"
//	@Success		200			{object}	common.Response
//	@Failure		404			{object}	common.Response
//	@Router			/api/v1/issues/{key}/documents [post]
func (h *IssueHandler) LinkDocument(ctx context.Context, req LinkDocReq) (struct{}, error) {
	issue, err := h.repo.GetByKey(req.IssueKey)
	if err != nil {
		return struct{}{}, common.NewAppError(40401, "issue not found")
	}

	if err := h.repo.LinkDocument(issue.ID, req.DocumentID); err != nil {
		return struct{}{}, common.NewAppError(40001, err.Error())
	}

	return struct{}{}, nil
}

// UnlinkDocument godoc
//
//	@Id				UnlinkDocument
//	@Summary		Unlink document from issue
//	@Description	Remove a document association from an issue
//	@Tags			Issues
//	@Produce		json
//	@Param			key			path		string	true	"Issue key"
//	@Param			document_id	query		int		true	"Document ID"
//	@Success		200			{object}	common.Response
//	@Failure		404			{object}	common.Response
//	@Router			/api/v1/issues/{key}/documents [delete]
func (h *IssueHandler) UnlinkDocument(ctx context.Context, req UnlinkDocReq) (struct{}, error) {
	issue, err := h.repo.GetByKey(req.IssueKey)
	if err != nil {
		return struct{}{}, common.NewAppError(40401, "issue not found")
	}

	if err := h.repo.UnlinkDocument(issue.ID, req.DocumentID); err != nil {
		return struct{}{}, common.NewAppError(40001, err.Error())
	}

	return struct{}{}, nil
}

// --- Task handlers ---

// ListTasks godoc
//
//	@Id				ListTasks
//	@Summary		List tasks for an issue
//	@Description	List all tasks belonging to an issue
//	@Tags			Tasks
//	@Produce		json
//	@Param			key	path		string	true	"Issue key"
//	@Success		200	{object}	common.Response{data=[]model.Task}
//	@Failure		404	{object}	common.Response
//	@Router			/api/v1/issues/{key}/tasks [get]
func (h *IssueHandler) ListTasks(ctx context.Context, req ListTaskReq) ([]model.Task, error) {
	issue, err := h.repo.GetByKey(req.IssueKey)
	if err != nil {
		return nil, common.NewAppError(40401, "issue not found")
	}

	return h.taskRepo.ListByIssue(issue.ID)
}

// CreateTask godoc
//
//	@Id				CreateTask
//	@Summary		Create a task for an issue
//	@Description	Create a new task under an issue
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			key		path		string			true	"Issue key"
//	@Param			task	body		CreateTaskReq	true	"Task info"
//	@Success		200		{object}	common.Response{data=model.Task}
//	@Failure		400		{object}	common.Response
//	@Failure		404		{object}	common.Response
//	@Router			/api/v1/issues/{key}/tasks [post]
func (h *IssueHandler) CreateTask(ctx context.Context, req CreateTaskReq) (*model.Task, error) {
	issue, err := h.repo.GetByKey(req.IssueKey)
	if err != nil {
		return nil, common.NewAppError(40401, "issue not found")
	}

	task := model.NewTask(issue.ID, req.Title)
	task.Description = req.Description

	if err := task.Validate(); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	if err := h.taskRepo.Create(task); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	return h.taskRepo.GetByKey(task.Key)
}

// UpdateTask godoc
//
//	@Id				UpdateTask
//	@Summary		Update a task
//	@Description	Update task fields or status
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			key		path		string			true	"Issue key"
//	@Param			task_key	path		string			true	"Task key"
//	@Param			task		body		UpdateTaskReq	true	"Task update info"
//	@Success		200			{object}	common.Response{data=model.Task}
//	@Failure		400			{object}	common.Response
//	@Failure		404			{object}	common.Response
//	@Router			/api/v1/issues/{key}/tasks/{task_key} [put]
func (h *IssueHandler) UpdateTask(ctx context.Context, req UpdateTaskReq) (*model.Task, error) {
	task, err := h.taskRepo.GetByKey(req.TaskKey)
	if err != nil {
		return nil, common.NewAppError(40401, "task not found")
	}

	if req.Title != nil {
		task.Title = *req.Title
	}

	if req.Description != nil {
		task.Description = *req.Description
	}

	if req.Status != nil {
		task.Status = *req.Status
	}

	if err := task.Validate(); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	if err := h.taskRepo.Update(task); err != nil {
		return nil, common.NewAppError(40001, err.Error())
	}

	return h.taskRepo.GetByKey(task.Key)
}

// DeleteTask godoc
//
//	@Id				DeleteTask
//	@Summary		Delete a task
//	@Description	Delete a task by ID
//	@Tags			Tasks
//	@Produce		json
//	@Param			key		path		string	true	"Issue key"
//	@Param			task_key	path		string			true	"Task key"
//	@Success		200			{object}	common.Response
//	@Failure		404			{object}	common.Response
//	@Router			/api/v1/issues/{key}/tasks/{task_key}/delete [post]
func (h *IssueHandler) DeleteTask(ctx context.Context, req DeleteTaskReq) (struct{}, error) {
	_, err := h.taskRepo.GetByKey(req.TaskKey)
	if err != nil {
		return struct{}{}, common.NewAppError(40401, "task not found")
	}

	if err := h.taskRepo.DeleteByKey(req.TaskKey); err != nil {
		return struct{}{}, common.NewAppError(40001, err.Error())
	}

	return struct{}{}, nil
}
