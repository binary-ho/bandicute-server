package api

import "net/http"

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /posts", app.WriteAllMembersPost)
	mux.HandleFunc("POST /studies/{studyId}/posts", app.WriteByStudy)
	mux.HandleFunc("POST /members/{memberId}/posts", app.WriteByMember)
	return mux
}
