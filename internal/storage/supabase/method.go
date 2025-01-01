package supabase

import "bandicute-server/internal/storage/repository/connection"

const (
	GetMethod   = connection.DML("GET")
	PostMethod  = connection.DML("POST")
	PatchMethod = connection.DML("PATCH")
)
