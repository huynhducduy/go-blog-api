package user

import (
	"github.com/go-chi/chi"
	"go-blog/pkg/utils"
	"net/http"
	"strconv"
)

func RouterList(w http.ResponseWriter, r *http.Request) {
	data, err := List()
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, 200, data)
}

func RouterRead(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idNum, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ResponseMessage(w, http.StatusBadRequest, "User's ID must be a number!")
		return
	}

	data, err := Read(idNum)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}
	utils.Response(w, 200, data)
}
