package supabase

import "bandicute-server/internal/storage/repository/connection"

const (
	getMethod   = connection.DML("GET")
	postMethod  = connection.DML("POST")
	patchMethod = connection.DML("PATCH")
)
