package worktile

import (
	"encoding/json"
	"github.com/xudai3/worktile/utils"
)

type TaskDetailsReq struct {
	//AccessToken string `json:"access_token"`
	TaskId string `json:"task_id"`
}

type TaskDetailsRsp struct {
	Id string `json:"_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	IsDeleted int `json:"is_deleted"` //任务是否被删除。0：未删除，1：已删除
	IsArchived int `json:"is_archived"` //任务是否被归档。0：未归档，1：已归档
	Type int `json:"type"` //任务类型。0：普通任务，1：任务模板
	Priority int `json:"priority"` //任务优先级。0：无优先级，1：低，2：中，3：高
	Visibility int `json:"visibility"` //任务可见性。0：公开，1：私有
	Entry string `json:"entry"` //任务所在的列表ID，无列表时无此属性
	EntryName string `json:"entry_name"` //任务所在列表的名字（如无列表则无此属性）
	Children []interface{} `json:"children"`
	Watchers []interface{} `json:"watchers"`
	DueDate DateInfo `json:"due_date"` //任务截止日期（无截止日期则无此属性）
	CreatedBy CommonUserInfo `json:"created_by"`
	UpdatedBy CommonUserInfo `json:"updated_by"`
}

func (w *Worktile) GetTask(taskId string, accessToken string) *TaskDetailsRsp {
	req := &TaskDetailsReq{TaskId:taskId}
	rsp := &TaskDetailsRsp{}
	bytes, err := w.Client.GetWithParam(ApiGetTaskDetail, req.TaskId, utils.BuildTokenHeaderOptions(accessToken))
	if err != nil {
		return nil
	}
	json.Unmarshal(bytes, rsp)
	return rsp
}