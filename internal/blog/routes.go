package blog

import (
	"encoding/json"
	"go-blog/internal/user"
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

	if ok && len(cursors) > 0 && cursors[0] != ""{
		cursor, err = strconv.Atoi(cursors[0])
		if err != nil {
			utils.ResponseMessage(w, http.StatusBadRequest, "Cursor must be a number!")
			return
		}
	}

	var filter BlogFilter
	var sortMethod BlogSortMethod

	tags, ok := r.URL.Query()["tag"]
	if ok && len(tags) > 0 && tags[0] != "" {
		filter.Tags = &tags[0]
	}
	userIds, ok := r.URL.Query()["user_id"]
	if ok && len(userIds) > 0 && userIds[0] != "" {
		if userId, err := strconv.ParseInt(userIds[0], 10, 64); err == nil {
			filter.UserId = userId
		} else {
			utils.ResponseMessage(w, http.StatusBadRequest, "User's ID must be a number!")
		}

	}

	data, err := List(cursor, filter, sortMethod)
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

	now := int64(time.Now().Unix())
	newBlog.CreatedAt = &now

	usr := r.Context().Value("user").(*user.User)
	newBlog.UserId = usr.Id

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
