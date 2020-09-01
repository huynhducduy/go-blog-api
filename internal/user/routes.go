package user

import (
	"go-blog/pkg/utils"
	"net/http"
)

func RouterList(w http.ResponseWriter, r *http.Request) {
	data, err := List()
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, 200, data)
}

func RouterMe(w http.ResponseWriter, r *http.Request) {
	employee := r.Context().Value("user").(*User)
	utils.Response(w, 200, employee)
}