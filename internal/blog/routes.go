package blog

import (
	"github.com/go-chi/chi"
	"go-blog/pkg/utils"
	"net/http"
	"strconv"
)

func RouterList(w http.ResponseWriter, r *http.Request) {

	cursors, ok := r.URL.Query()["cursor"]

	cursor := 0
	var err error

	if ok && len(cursors) > 0 {
		cursor, err = strconv.Atoi(cursors[0])
		if err != nil {
			utils.ResponseMessage(w, http.StatusBadRequest, "Cursor must be a number!")
			return
		}
	}

	data, err := List(cursor)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, 200, data)
}

func RouterRead(w http.ResponseWriter, r *http.Request) {
	blogIdStr := chi.URLParam(r, "id")
	blogIdNum, err := strconv.Atoi(blogIdStr)
	if err != nil {
		utils.ResponseMessage(w, http.StatusBadRequest, "Blog ID must be a number!")
		return
	}

	data, err := Read(blogIdNum)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}
	utils.Response(w, 200, data)
}