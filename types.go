package worktile

type HookCreateTask struct {
	Event string `json:"event"`
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
	Event string `json:"event"`
	Payload PayloadUpdateTitle `json:"payload"`
}

type HookUpdateDesc struct {
	Event string `json:"event"`
	Payload PayloadUpdateTitle `json:"payload"`
}

type PayloadCreateTask struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Identifier string `json:"identifier"`
	Type CommonInfo `json:"type"`
	Project CommonInfo `json:"project"`	
	Creator CommonUserInfo `json:"creator"`
	Assignee CommonUserInfo `json:"assignee"`
	Due DateInfo `json:"due"`
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
	From string `json:"from"`	
	To string `json:"to"`
	Task TaskInfo `json:"task"`
	UpdateBy CommonUserInfo `json:"update_by"`
}

type TaskInfo struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Identifier string `json:"identifier"`
	Type CommonInfo `json:"type"`
	Project CommonInfo `json:"project"`
	State CommonInfo `json:"state"`
}

type DateInfo struct {
	Date int64 `json:"date"`
	WithTime int `json:"with_time"`
}

type CommonInfo struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type CommonUserInfo struct {
	Uid string `json:"uid"`
	UserName string `json:"user_name"`
	DisplayName string `json:"display_name"`
}

