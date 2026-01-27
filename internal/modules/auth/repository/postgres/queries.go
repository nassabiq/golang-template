package postgres

import _ "embed"

//go:embed queries.sql
var rawQuery string
