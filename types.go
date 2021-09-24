package worktile

type HookCreateTask struct {
	Event   string            `json:"event"`
	Payload PayloadCreateTask `json:"payload"`
}

type HookUpdateState struct {
	Event   string             `json:"event"`
	Payload PayloadUpdateState `json:"payload"`
}

type HookAssignee struct {
	Event   string          `json:"event"`
	Payload PayloadAssignee `json:"payload"`
}

type HookUpdateTitle struct {
	Event   string             `json:"event"`
	Payload PayloadUpdateTitle `json:"payload"`
}

type HookUpdateDesc struct {
	Event   string             `json:"event"`
	Payload PayloadUpdateTitle `json:"payload"`
}

type PayloadCreateTask struct {
	Id         string         `json:"id"`
	Title      string         `json:"title"`
	Identifier string         `json:"identifier"`
	Type       CommonInfo     `json:"type"`
	Project    CommonInfo     `json:"project"`
	Creator    CommonUserInfo `json:"creator"`
	Assignee   CommonUserInfo `json:"assignee"`
	Due        DateInfo       `json:"due"`
}

type PayloadUpdateState struct {
	Task     TaskInfo       `json:"task"`
	UpdateBy CommonUserInfo `json:"update_by"`
	From     CommonInfo     `json:"from"`
	To       CommonInfo     `json:"to"`
}

type PayloadAssignee struct {
	Task     TaskInfo       `json:"task"`
	UpdateBy CommonUserInfo `json:"update_by"`
	From     CommonUserInfo `json:"from"`
	To       CommonUserInfo `json:"to"`
}

type PayloadUpdateTitle struct {
	From     string         `json:"from"`
	To       string         `json:"to"`
	Task     TaskInfo       `json:"task"`
	UpdateBy CommonUserInfo `json:"update_by"`
}

type TaskInfo struct {
	Id         string     `json:"id"`
	Title      string     `json:"title"`
	Identifier string     `json:"identifier"`
	Type       CommonInfo `json:"type"`
	Project    CommonInfo `json:"project"`
	State      CommonInfo `json:"state"`
}

type DateInfo struct {
	Date     int64 `json:"date"`
	WithTime int   `json:"with_time"`
}

type CommonInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CommonUserInfo struct {
	Uid         string `json:"uid"`
	UserName    string `json:"user_name"`
	DisplayName string `json:"display_name"`
}

// task types

type GetTaskDetailsReq struct {
	AccessToken string `json:"access_token"`
	Fields      string `json:"fields"`
	TaskIds     string `json:"task_ids"`
}

type TaskDetail struct {
	Id          string       `json:"_id,omitempty"`
	CreatedBy   string       `json:"created_by,omitempty"`
	CreatedAt   int64        `json:"created_at,omitempty"`
	UpdatedAt   int          `json:"updated_at,omitempty"`
	IsArchived  int          `json:"is_archived,omitempty"` //任务是否被归档。0：未归档，1：已归档
	Title       string       `json:"title,omitempty"`
	ProjectId   string       `json:"project_id,omitempty"`
	Properties  TaskProperty `json:"properties,omitempty"`
	StateType   int          `json:"state_type,omitempty"`
	Identifier  string       `json:"identifier,omitempty"`
	CompletedAt int          `json:"completed_at,omitempty"`
	DerivedIds  []string     `json:"derived_ids,omitempty"`
	TaskState   TaskState    `json:"task_state,omitempty"`
	TaskType    TaskType     `json:"task_type,omitempty"`
	ParentId    string       `json:"parent_id,omitempty"`
	ParentIds   []string     `json:"parent_ids,omitempty"`
}

type TaskProperty struct {
	Assignee    string    `json:"assignee,omitempty"`
	Start       DateInfo  `json:"start,omitempty"`
	Due         DateInfo  `json:"due,omitempty"`
	Tag         []TaskTag `json:"tag,omitempty"`
	Desc        string    `json:"desc,omitempty"`
	Priority    string    `json:"priority,omitempty"`
	Participant []string  `json:"participant,omitempty"` // 参与人
	Attachment  []string  `json:"attachment,omitempty"`  // 附件
}

type TaskState struct {
	Id          string `json:"_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        int    `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

type TaskType struct {
	Id          string `json:"_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// user types

type GetUserDetailReq struct {
	AccessToken string `json:"access_token"`
	Uids        string `json:"uids"`
}

type UserDetail struct {
	Uid               string `json:"uid"`
	Name              string `json:"name"`
	DisplayName       string `json:"display_name"`
	Avatar            string `json:"avatar"`
	Status            int    `json:"status"`
	DisplayNamePinyin string `json:"display_name_pinyin"`
}

type TaskTag struct {
	Id     string `json:"_id,omitempty"`
	Name   string `json:"name,omitempty"`
	Color  string `json:"color,omitempty"`
	ModeId string `json:"mode_id,omitempty"`
}
