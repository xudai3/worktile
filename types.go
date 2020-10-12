package worktile

type HookUpdateState struct {
	Event string `json:"event"`
	Payload TaskUpdateState `json:"payload"`
}

type HookAssignee struct {
	Event string `json:"event"`
	Payload TaskAssignee `json:"payload"`
}

type TaskUpdateState struct {
	Task TaskDetail `json:"task"`
	UpdateBy UserDetail `json:"update_by"`
	From TaskState `json:"from"`
	To TaskState `json:"to"`
}

type TaskDetail struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Identifier string `json:"identifier"`
	Type TaskType `json:"type"`
	Project ProjectDetail `json:"project"`
	State StateDetail `json:"state"`
}

type UserDetail struct {
	Uid string `json:"id"`
	UserName string `json:"user_name"`
	DisplayName string `json:"display_name"`
}

type TaskState struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type TaskType struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type ProjectDetail struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type StateDetail struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type TaskAssignee struct {
	Task TaskDetail `json:"task"`
	UpdateBy UserDetail `json:"update_by"`
	From UserDetail `json:"from"`
	To UserDetail `json:"to"`
}
