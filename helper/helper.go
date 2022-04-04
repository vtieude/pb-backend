package helper

import (
	"database/sql"
	"time"
)

const DateTimeFormat = "YYYY-MM-DD"

func BeginningOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func EndOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, -1, time.UTC)
}

func ConvertToNullPointDateTime(date *string) *time.Time {
	if date == nil {
		return nil
	}
	t, err := time.Parse(*date, DateTimeFormat)

	if err != nil {
		return nil
	}
	return &t
}
func ConvertToNullPointSqlString(data *string) sql.NullString {
	if data == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *data, Valid: true}
}

func ConvertToString(data *sql.NullString) string {
	if data == nil || !data.Valid {
		return ""
	}
	return data.String
}
func ConvertToPoinerString(data string) *string {
	return &data
}

func GetRoleLabelByRole(roleName string) string {
	fullRolePermision := map[string]string{
		"admin":       "Quản lí",
		"super_admin": "Super Admin",
		"staff":       "Nhân Viên",
		"user":        "Người dùng",
	}
	if roleName == "" {
		return ""
	}
	if role, ok := fullRolePermision[roleName]; ok {
		return role
	}
	return ""
}

func GetPermissionByRole(roleName string) int {
	fullRolePermision := map[string]int{
		"admin":       7,
		"super_admin": 7,
		"staff":       3,
		"user":        1,
	}
	if roleName == "" {
		return 0
	}
	if role, ok := fullRolePermision[roleName]; ok {
		return role
	}
	return 0
}
