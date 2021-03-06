package worktile

import (
	"encoding/json"
	"github.com/xudai3/worktile/logger"
	"github.com/xudai3/worktile/utils"
	"regexp"
	"strings"
)

type TaskDetailsReq struct {
	AccessToken string `json:"access_token"`
	Fields      string `json:"fields"`
	TaskIds     string `json:"task_ids"`
}

type TaskDetail struct {
	Id         string       `json:"_id"`
	CreatedBy  string       `json:"created_by"`
	CreatedAt  int64        `json:"created_at"`
	IsArchived int          `json:"is_archived"` //任务是否被归档。0：未归档，1：已归档
	Title      string       `json:"title"`
	Identifier string       `json:"identifier"`
	ProjectId  string       `json:"project_id"`
	Properties TaskProperty `json:"properties"`
	ParentId   string       `json:"parent_id"`
	ParentIds  []string     `json:"parent_ids"`
	DerivedIds []string     `json:"derived_ids"`
	TaskState  TaskState    `json:"task_state"`
}

type TaskProperty struct {
	Assignee string `json:"assignee"`
}

type TaskState struct {
	Id          string `json:"_id"`
	Name        string `json:"name"`
	Type        int    `json:"type"`
	Description string `json:"description"`
}

const defaultTaskFields = "assignee,workload"

func (w *Worktile) GetTasksByIds(accessToken string, taskIds []string) []TaskDetail {
	req := TaskDetailsReq{AccessToken: accessToken, Fields: defaultTaskFields, TaskIds: strings.Join(taskIds, ",")}
	var rsp []TaskDetail
	bytes, err := w.Client.Get(ApiGetTaskDetail, utils.ConvertStructToMap(req), utils.BuildTokenHeaderOptions(accessToken))
	if err != nil {
		logger.Debugf("get task:%v detail failed:%v\n", taskIds, err)
		return nil
	}
	json.Unmarshal(bytes, &rsp)
	return rsp
}

func (w *Worktile) GetTaskById(accessToken string, taskId string) *TaskDetail {
	req := TaskDetailsReq{AccessToken: accessToken, Fields: defaultTaskFields, TaskIds: taskId}
	var rsp []*TaskDetail
	bytes, err := w.Client.Get(ApiGetTaskDetail, utils.ConvertStructToMap(req), utils.BuildTokenHeaderOptions(accessToken))
	if err != nil {
		logger.Errorf("get task:%s detail failed:%v\n", taskId, err)
		return nil
	}
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		logger.Errorf("unmarshal taskdetails error:%v", err)
		return nil
	}
	if len(rsp) > 0 {
		return rsp[0]
	} else {
		return nil
	}
}

func (w *Worktile) GetAssigneeNameByTaskId(accessToken string, taskId string) string {
	task := w.GetTaskById(accessToken, taskId)
	assigneeUid := task.Properties.Assignee
	if assigneeUid == "" {
		logger.Errorf("assignee is empty")
		return ""
	}
	assignee := w.GetUserByUid(accessToken, assigneeUid)
	return assignee.DisplayName
}

func (w *Worktile) GetSubTasks(accessToken string, taskId string) []TaskDetail {
	currentTask := w.GetTaskById(accessToken, taskId)
	tasks := make([]TaskDetail, 0)
	if currentTask == nil || len(currentTask.DerivedIds) == 0 {
		return tasks
	}
	taskIds := make([]string, 0)
	taskIds = append(taskIds, currentTask.DerivedIds...)
	tasks = w.GetTasksByIds(accessToken, taskIds)
	return tasks
}

func (w *Worktile) GetAllSubTasks(accessToken string, taskId string) []TaskDetail {
	taskIds := make([]string, 0)
	tasks := make([]TaskDetail, 0)
	currentTask := w.GetTaskById(accessToken, taskId)
	if currentTask == nil {
		return tasks
	}
	if len(currentTask.DerivedIds) == 0 {
		return tasks
	}
	taskIds = append(taskIds, currentTask.DerivedIds...)
	for len(taskIds) > 0 {
		subTasks := w.GetTasksByIds(accessToken, taskIds)
		taskIds = nil
		for _, subTask := range subTasks {
			logger.Debugf("subTask:%+v title:%s", subTask, subTask.Title)
			if len(subTask.DerivedIds) != 0 {
				taskIds = append(taskIds, subTask.DerivedIds...)
			}
			tasks = append(tasks, subTask)
		}
	}
	return tasks
}

func (w *Worktile) GetAllAssigneeUids(accessToken string, taskId string) []string {
	assignees := make([]string, 0)
	currentTask := w.GetTaskById(accessToken, taskId)
	if currentTask == nil {
		return assignees
	}
	tasks := w.GetAllSubTasks(accessToken, taskId)
	assignees = append(assignees, currentTask.Properties.Assignee)
	for _, task := range tasks {
		if task.Properties.Assignee != "" {
			assignees = append(assignees, task.Properties.Assignee)
		}
	}
	return assignees
}

func (w *Worktile) GetAllAssigneeNames(accessToken string, taskId string) []string {
	assigneeUids := w.GetAllAssigneeUids(accessToken, taskId)
	users := w.GetUsersByUids(accessToken, assigneeUids)
	names := make([]string, 0, len(users))
	for _, user := range users {
		names = append(names, user.DisplayName)
	}
	return names
}

func (w *Worktile) GetAllAssigneeUidsFilterByTitle(accessToken string, taskId string, filter string) []string {
	assignees := make([]string, 0)
	currentTask := w.GetTaskById(accessToken, taskId)
	if currentTask == nil {
		return assignees
	}
	tasks := w.GetAllSubTasks(accessToken, taskId)
	r := regexp.MustCompile(filter)
	if matched := r.MatchString(currentTask.Title); matched {
		assignees = append(assignees, currentTask.Properties.Assignee)
	}
	for _, task := range tasks {
		matched := r.MatchString(task.Title)
		if matched && task.Properties.Assignee != "" {
			assignees = append(assignees, task.Properties.Assignee)
		}
	}
	return assignees
}

func (w *Worktile) GetAllAssigneeNamesFilterByTitle(accessToken string, taskId string, filter string) []string {
	assignees := w.GetAllAssigneeUidsFilterByTitle(accessToken, taskId, filter)
	users := w.GetUsersByUids(accessToken, assignees)
	names := make([]string, 0, len(users))
	for _, user := range users {
		names = append(names, user.DisplayName)
	}
	return names
}

func (w *Worktile) GetMainTaskDetail(accessToken string, taskId string) *TaskDetail {
	var mainTaskId string
	mainTaskId = taskId
	for {
		mainTask := w.GetTaskById(accessToken, mainTaskId)
		if mainTask == nil || mainTask.ParentId == "" {
			break
		} else {
			mainTaskId = mainTask.ParentId
		}
	}
	return w.GetTaskById(accessToken, mainTaskId)
}
