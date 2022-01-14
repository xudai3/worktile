package worktile

const (
	ApiGetAuthCode      = "https://dev.worktile.com/api/oauth2/authorize"
	ApiGetTenant        = "https://dev.worktile.com/open-api/tenant-access-token"
	ApiGetTaskDetail    = "https://dev.worktile.com/open-api/mission/get-tasks-by-ids"
	ApiGetRelationTasks = "https://dev.worktile.com/open-api/mission/get-relation-tasks"
	ApiGetUserByUid     = "https://dev.worktile.com/open-api/contact/get-members-by-ids"
)

const (
	StatusInactivated   = "未激活"
	StatusInDesign      = "设计中"
	StatusInDevelop     = "研发中"
	StatusInTest        = "测试中"
	StatusAccepted      = "已验收"
	StatusPublished     = "已发布"
	StatusNotStarted    = "未开始"
	StatusInProgress    = "进行中"
	StatusToBeConfirmed = "待确认"
	StatusCompleted     = "已完成"
	StatusPending       = "挂起"
)

const (
	DefaultTaskFields = "assignee,participant,start,due,tag,priority,desc,attachment"
	FieldAssignee     = "assignee"
	FieldParticipant  = "participant"
	FieldStart        = "start"
	FieldDue          = "due"
	FieldTag          = "tag"
	FieldPriority     = "priority"
	FieldDesc         = "desc"
	FieldAttachment   = "attachment"
)
