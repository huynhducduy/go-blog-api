package tag

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"go-blog/pkg/utils"
	"io/ioutil"
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

func RouterUpdate(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "tag")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	var updated Tag
	json.Unmarshal(reqBody, &updated)
	updated.Tag = &t

	err = Update(updated)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	utils.Response(w, http.StatusOK, updated)
}

func RouterRead(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "tag")

	data, err := Read(t)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}
	utils.Response(w, 200, data)
}