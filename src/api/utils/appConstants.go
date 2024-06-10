package utils

var USER_AUTHORITIES = map[string]string{
	"ADMIN_USER":  "ADMIN_USER",
	"NORMAL_USER": "NORMAL_USER",
}

var USER_STATUS = map[string]string{
	"VERIFIED":     "VERIFIED",
	"NOT_VERIFIED": "NOT_VERIFIED",
	"DEACTIVATED":  "DEACTIVATED",
}

var TASK_STATUS = map[string]string{
	"COMPLETED": "COMPLETED",
	"ON-GOING":  "ON-GOING",
	"DROPPED":   "DROPPED",
}
