package worktile

const (
	ApiGetAuthCode   = "https://dev.worktile.com/api/oauth2/authorize"
	ApiGetTenant     = "https://dev.worktile.com/open-api/tenant-access-token"
	ApiGetTaskDetail = "https://dev.worktile.com/open-api/mission/get-tasks-by-ids"
	ApiGetUserByUid = "https://dev.worktile.com/open-api/contact/get-members-by-ids"
)

const (
	StatusNotStarted = "未开始"
	StatusInProgress = "进行中"
	StatusToBeConfirmed = "待确认"
	StatusCompleted = "已完成"
	StatusPending = "挂起"
)