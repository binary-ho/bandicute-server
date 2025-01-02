package request

import (
	"bandicute-server/internal/storage/repository/post"
	"context"
)

type ParsePostByMemberId struct {
	Context  context.Context
	MemberId string
}

type OpenPullRequest struct {
	Context    context.Context
	Post       *post.Model
	Repository string
	MemberName string
	Summary    string
	StudyId    string
	FilePath   string
}

type Summarize struct {
	Context context.Context
	Post    *post.Model
}
