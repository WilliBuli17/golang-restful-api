package helper

import "database/sql"

func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	PanicIfError(err)
}
