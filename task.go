package worktile

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/xudai3/worktile/utils"
)

func (w *Worktile) GetTasksByIds(taskIds []string, fields []string) ([]TaskDetail, error) {
	accessToken, err := w.GetTenant()
	if err != nil {
		return nil, err
	}
	var taskFields string
	var rsp []TaskDetail
	if len(taskIds) == 0 {
		return nil, errors.New("GetTasksByIds>taskIds empty")
	}
	if len(fields) == 0 {
		taskFields = DefaultTaskFields
	} else {
		taskFields = strings.Join(fields, ",")
	}
	req := GetTaskDetailsReq{AccessToken: accessToken, Fields: taskFields, TaskIds: strings.Join(taskIds, ",")}
	bytes, err := w.Client.Get(ApiGetTaskDetail, utils.ConvertStructToMap(req), utils.BuildTokenHeaderOptions(accessToken))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (w *Worktile) GetTaskById(taskId string) (*TaskDetail, error) {
	accessToken, err := w.GetTenant()
	if err != nil {
		return nil, err
	}
	req := GetTaskDetailsReq{AccessToken: accessToken, Fields: DefaultTaskFields, TaskIds: taskId}
	var rsp []*TaskDetail
	bytes, err := w.Client.Get(ApiGetTaskDetail, utils.ConvertStructToMap(req), utils.BuildTokenHeaderOptions(accessToken))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		return nil, err
	}
	if len(rsp) == 0 {
		return nil, errors.New("result empty")
	}
	return rsp[0], nil
}

func (w *Worktile) GetRelationTasks(taskId string) ([]TaskDetail, error) {
	accessToken, err := w.GetTenant()
	if err != nil {
		return nil, err
	}
	req := GetTaskDetailReq{AccessToken: accessToken, Fields: DefaultTaskFields, TaskId: taskId}
	var rsp []TaskDetail
	bytes, err := w.Client.Get(ApiGetRelationTasks, utils.ConvertStructToMap(req), utils.BuildTokenHeaderOptions(accessToken))
	if err != nil {
		return nil, fmt.Errorf("GetRelationTasks>%s detail failed:%v", taskId, err)
	}
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		return nil, fmt.Errorf("GetRelationTasks>unmarshal %v error:%v", string(bytes), err)
	}
	if len(rsp) == 0 {
		return nil, fmt.Errorf("GetRelationTasks>taskId:%v result(%v) empty", taskId, rsp)
	}
	return rsp, nil
}

func (w *Worktile) GetAssigneeNameByTaskId(taskId string) (string, error) {
	task, err := w.GetTaskById(taskId)
	if err != nil {
		return "", err
	}
	assigneeUid := task.Properties.Assignee
	if assigneeUid == "" {
		return "", fmt.Errorf("GetAssigneeNameByTaskId>>task:%+v assignee is empty", task)
	}
	assignee, err := w.GetUserByUid(assigneeUid)
	if err != nil {
		return "", err
	}
	return assignee.DisplayName, nil
}

func (w *Worktile) GetSubTasks(taskId string) ([]TaskDetail, error) {
	currentTask, err := w.GetTaskById(taskId)
	if err != nil {
		return nil, err
	}
	var tasks []TaskDetail
	if currentTask == nil {
		return nil, errors.New("GetSubTasks>task empty")
	}
	if len(currentTask.DerivedIds) == 0 {
		return nil, errors.New("GetSubTasks>derivedIds empty")
	}
	taskIds := make([]string, 0)
	taskIds = append(taskIds, currentTask.DerivedIds...)
	tasks, err = w.GetTasksByIds(taskIds, nil)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (w *Worktile) GetAllSubTasks(taskId string) ([]TaskDetail, error) {
	taskIds := make([]string, 0)
	tasks := make([]TaskDetail, 0)
	currentTask, err := w.GetTaskById(taskId)
	if err != nil {
		return nil, err
	}
	if currentTask == nil {
		return nil, fmt.Errorf("GetAllSubTasks>taskId:%v empty", taskId)
	}
	if len(currentTask.DerivedIds) == 0 {
		return nil, fmt.Errorf("GetAllSubTasks>task:%+v derivedIds empty", currentTask)
	}
	taskIds = append(taskIds, currentTask.DerivedIds...)
	for len(taskIds) > 0 {
		subTasks, err := w.GetTasksByIds(taskIds, nil)
		if err != nil {
			return nil, fmt.Errorf("GetAllSubTasks>taskIds:%v err:%v", taskIds, err)
		}
		taskIds = nil
		for _, subTask := range subTasks {
			if len(subTask.DerivedIds) != 0 {
				taskIds = append(taskIds, subTask.DerivedIds...)
			}
			tasks = append(tasks, subTask)
		}
	}
	if len(tasks) == 0 {
		return nil, errors.New("GetAllSubTasks>tasks empty")
	}
	return tasks, nil
}

func (w *Worktile) GetAllAssigneeUids(taskId string) ([]string, error) {
	assignees := make([]string, 0)
	currentTask, err := w.GetTaskById(taskId)
	if err != nil {
		return nil, err
	}
	if currentTask == nil {
		return assignees, errors.New("GetAllAssigneeUids>current task empty")
	}
	tasks, err := w.GetAllSubTasks(taskId)
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

func (w *Worktile) GetAllAssigneeNames(taskId string) ([]string, error) {
	assigneeUids, err := w.GetAllAssigneeUids(taskId)
	if err != nil {
		return nil, err
	}
	users, err := w.GetUsersByUids(assigneeUids)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(users))
	for _, user := range users {
		names = append(names, user.DisplayName)
	}
	return names, nil
}

func (w *Worktile) GetAllAssigneeUidsFilterByTitle(taskId string, filter string) ([]string, error) {
	assignees := make([]string, 0)
	currentTask, err := w.GetTaskById(taskId)
	if err != nil {
		return nil, err
	}
	if currentTask == nil {
		return assignees, errors.New("GetAllAssigneeUidsFilterByTitle>current task empty")
	}
	tasks, err := w.GetAllSubTasks(taskId)
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

func (w *Worktile) GetAllAssigneeNamesFilterByTitle(taskId string, filter string) ([]string, error) {
	assignees, err := w.GetAllAssigneeUidsFilterByTitle(taskId, filter)
	if err != nil {
		return nil, err
	}
	users, err := w.GetUsersByUids(assignees)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(users))
	for _, user := range users {
		names = append(names, user.DisplayName)
	}
	return names, nil
}

func (w *Worktile) GetMainTaskDetail(taskId string) (*TaskDetail, error) {
	var mainTaskId string
	mainTaskId = taskId
	for {
		mainTask, err := w.GetTaskById(mainTaskId)
		if err != nil {
			return nil, err
		}
		if mainTask == nil || mainTask.ParentId == "" {
			break
		} else {
			mainTaskId = mainTask.ParentId
		}
	}
	return w.GetTaskById(mainTaskId)
}

func (w *Worktile) GetParticipantUids(taskId string) ([]string, error) {
	taskDetail, err := w.GetTaskById(taskId)
	if err != nil {
		return nil, err
	}
	return taskDetail.Properties.Participant, nil
}

func (w *Worktile) GetParticipantNames(taskId string) ([]string, error) {
	Uids, err := w.GetParticipantUids(taskId)
	if err != nil {
		return nil, err
	}
	userNames := make([]string, 0)
	for _, uid := range Uids {
		user, err := w.GetUserByUid(uid)
		if err != nil {
			return nil, err
		}
		userNames = append(userNames, user.DisplayName)
	}
	return userNames, nil
}

func (w *Worktile) GetDescRTFInfo(descs []DescRTF) string {
	var descString string
	for i, desc := range descs {
		if i != 0 {
			descString += "\n"
		}
		if desc.Type == "paragraph" {
			for _, child := range desc.Children {
				descString += child.Text
			}
		} else if desc.Type == "image" {
			descString += desc.OriginUrl
		}
	}
	return descString
}
