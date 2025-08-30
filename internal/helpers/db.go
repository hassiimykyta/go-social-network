package helpers

import "database/sql"

func DerefStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
