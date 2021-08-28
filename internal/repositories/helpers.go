package repositories

import "database/sql"

func lastInsertID(result sql.Result) int32 {
	id, _ := result.LastInsertId()
	return int32(id)
}
