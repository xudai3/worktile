package worktile

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/xudai3/worktile/logger"
	"github.com/xudai3/worktile/utils"
)

func (w *Worktile) GetTasksByIds(accessToken string, taskIds []string, fields []string) ([]TaskDetail, error) {
	var taskFields string
	var rsp []TaskDetail
	if len(taskIds) == 0 {
		return nil, errors.New("taskIds empty")
	}
	if len(fields) == 0 {
		taskFields = DefaultTaskFields
	} else {
		taskFields = strings.Join(fields, ",")
	}
	req := GetTaskDetailsReq{AccessToken: accessToken, Fields: taskFields, TaskIds: strings.Join(taskIds, ",")}
	bytes, err := w.Client.Get(ApiGetTaskDetail, utils.ConvertStructToMap(req), utils.BuildTokenHeaderOptions(accessToken))
	if err != nil {
		logger.Debugf("get task:%v detail failed:%v\n", taskIds, err)
		return nil, err
	}
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (w *Worktile) GetTaskById(accessToken string, taskId string) (*TaskDetail, error) {
	req := GetTaskDetailsReq{AccessToken: accessToken, Fields: DefaultTaskFields, TaskIds: taskId}
	var rsp []*TaskDetail
	bytes, err := w.Client.Get(ApiGetTaskDetail, utils.ConvertStructToMap(req), utils.BuildTokenHeaderOptions(accessToken))
	if err != nil {
		logger.Errorf("get task:%s detail failed:%v\n", taskId, err)
		return nil, err
	}
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		logger.Errorf("unmarshal taskdetails error:%v", err)
		return nil, err
	}
	if len(rsp) == 0 {
		return nil, errors.New("result empty")
	}
	return rsp[0], nil
}

func (w *Worktile) GetRelationTasks(accessToken string, taskId string) ([]TaskDetail, error) {
	req := GetTaskDetailReq{AccessToken: accessToken, Fields: DefaultTaskFields, TaskId: taskId}
	var rsp []TaskDetail
	bytes, err := w.Client.Get(ApiGetRelationTasks, utils.ConvertStructToMap(req), utils.BuildTokenHeaderOptions(accessToken))
	if err != nil {
		logger.Errorf("GetRelationTasks>%s detail failed:%v\n", taskId, err)
		return nil, err
	}
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		logger.Errorf("GetRelationTasks>unmarshal taskdetails error:%v", err)
		return nil, err
	}
	if len(rsp) == 0 {
		return nil, errors.New("GetRelationTasks>result empty")
	}
	return rsp, nil
}

func (w *Worktile) GetAssigneeNameByTaskId(accessToken string, taskId string) (string, error) {
	task, err := w.GetTaskById(accessToken, taskId)
	if err != nil {
		return "", err
	}
	assigneeUid := task.Properties.Assignee
	if assigneeUid == "" {
		err = fmt.Errorf("assignee is empty")
		logger.Error(err.Error())
		return "", err
	}
	assignee, err := w.GetUserByUid(accessToken, assigneeUid)
	if err != nil {
		return "", err
	}
	return assignee.DisplayName, nil
}

func (w *Worktile) GetSubTasks(accessToken string, taskId string) ([]TaskDetail, error) {
	currentTask, err := w.GetTaskById(accessToken, taskId)
	if err != nil {
		return nil, err
	}
	var tasks []TaskDetail
	if currentTask == nil {
		return nil, errors.New("task empty")
	}
	if len(currentTask.DerivedIds) == 0 {
		return nil, errors.New("derivedIds empty")
	}
	taskIds := make([]string, 0)
	taskIds = append(taskIds, currentTask.DerivedIds...)
	tasks, err = w.GetTasksByIds(accessToken, taskIds, nil)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (w *Worktile) GetAllSubTasks(accessToken string, taskId string) ([]TaskDetail, error) {
	taskIds := make([]string, 0)
	tasks := make([]TaskDetail, 0)
	currentTask, err := w.GetTaskById(accessToken, taskId)
	if err != nil {
		return nil, err
	}
	if currentTask == nil {
		return nil, errors.New("task empty")
	}
	if len(currentTask.DerivedIds) == 0 {
		return nil, errors.New("derivedIds empty")
	}
	taskIds = append(taskIds, currentTask.DerivedIds...)
	for len(taskIds) > 0 {
		subTasks, err := w.GetTasksByIds(accessToken, taskIds, nil)
		if err != nil {
			logger.Error(err)
			break
		}
		taskIds = nil
		for _, subTask := range subTasks {
			logger.Debugf("subTask:%+v title:%s", subTask, subTask.Title)
			if len(subTask.DerivedIds) != 0 {
				taskIds = append(taskIds, subTask.DerivedIds...)
			}
			tasks = append(tasks, subTask)
		}
	}
	if len(tasks) == 0 {
		return nil, errors.New("tasks empty")
	}
	return tasks, nil
}

func (w *Worktile) GetAllAssigneeUids(accessToken string, taskId string) ([]string, error) {
	assignees := make([]string, 0)
	currentTask, err := w.GetTaskById(accessToken, taskId)
	if err != nil {
		return nil, err
	}
	if currentTask == nil {
		return assignees, errors.New("current task empty")
	}
	tasks, err := w.GetAllSubTasks(accessToken, taskId)
	if err != nil {
		return nil, err
	}
	assignees = append(assignees, currentTask.Properties.Assignee)
	for _, task := range tasks {
		if task.Properties.Assignee != "" {
			assignees = append(assignees, task.Properties.Assignee)
		}
	}
	return assignees, nil
}

func (w *Worktile) GetAllAssigneeNames(accessToken string, taskId string) ([]string, error) {
	assigneeUids, err := w.GetAllAssigneeUids(accessToken, taskId)
	if err != nil {
		return nil, err
	}
	users, err := w.GetUsersByUids(accessToken, assigneeUids)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(users))
	for _, user := range users {
		names = append(names, user.DisplayName)
	}
	return names, nil
}

func (w *Worktile) GetAllAssigneeUidsFilterByTitle(accessToken string, taskId string, filter string) ([]string, error) {
	assignees := make([]string, 0)
	currentTask, err := w.GetTaskById(accessToken, taskId)
	if err != nil {
		return nil, err
	}
	if currentTask == nil {
		return assignees, errors.New("current task empty")
	}
	tasks, err := w.GetAllSubTasks(accessToken, taskId)
	if err != nil {
		return nil, err
	}
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
	return assignees, nil
}

func (w *Worktile) GetAllAssigneeNamesFilterByTitle(accessToken string, taskId string, filter string) ([]string, error) {
	assignees, err := w.GetAllAssigneeUidsFilterByTitle(accessToken, taskId, filter)
	if err != nil {
		return nil, err
	}
	users, err := w.GetUsersByUids(accessToken, assignees)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(users))
	for _, user := range users {
		names = append(names, user.DisplayName)
	}
	return names, nil
}

func (w *Worktile) GetMainTaskDetail(accessToken string, taskId string) (*TaskDetail, error) {
	var mainTaskId string
	mainTaskId = taskId
	for {
		mainTask, err := w.GetTaskById(accessToken, mainTaskId)
		if err != nil {
			return nil, err
		}
		if mainTask == nil || mainTask.ParentId == "" {
			break
		} else {
			mainTaskId = mainTask.ParentId
		}
	}
	return w.GetTaskById(accessToken, mainTaskId)
}

func (w *Worktile) GetParticipantUids(accessToken string, taskId string) ([]string, error) {
	taskDetail, err := w.GetTaskById(accessToken, taskId)
	if err != nil {
		return nil, err
	}
	return taskDetail.Properties.Participant, nil
}

func (w *Worktile) GetParticipantNames(accessToken string, taskId string) ([]string, error) {
	Uids, err := w.GetParticipantUids(accessToken, taskId)
	if err != nil {
		return nil, err
	}
	userNames := make([]string, 0)
	for _, uid := range Uids {
		user, err := w.GetUserByUid(accessToken, uid)
		if err != nil {
			logger.Error(err)
			continue
		}
		userNames = append(userNames, user.DisplayName)
	}
	return userNames, nil
}
