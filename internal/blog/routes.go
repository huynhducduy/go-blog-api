package blog

import (
	"encoding/json"
	"go-blog/pkg/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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

func RouterCreate(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	var newBlog Blog
	json.Unmarshal(reqBody, &newBlog)

	*newBlog.CreatedAt = int64(time.Now().Unix())

	id, err := Create(newBlog)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.ResponseCreated(w, id)
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

func RouterUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseMessage(w, http.StatusBadRequest, "Id must be an integer!")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	var updatedBlog Blog
	json.Unmarshal(reqBody, &updatedBlog)
	updatedBlog.Id = &id

	err = Update(updatedBlog)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}

func RouterDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ResponseMessage(w, http.StatusBadRequest, "Id must be an integer!")
		return
	}

	err = Delete(id)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.ResponseMessage(w, http.StatusOK, "Delete succeed!")
}
