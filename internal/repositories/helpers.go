package repositories

import (
	"database/sql"
	"fmt"
)

func lastInsertID(result sql.Result) int32 {
	id, _ := result.LastInsertId()
	fmt.Println(id)
	return int32(id)
}
