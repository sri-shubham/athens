package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sri-shubham/athens/models"
)

type GenericCrud interface {
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}

type IDResp struct {
	ID int64 `json:"id"`
}

type GenericCrudHelper[Model models.IdGetter] struct {
	db models.CRUD[Model]
}

func NewGenericCrudHelper[Model models.IdGetter](db models.CRUD[Model]) GenericCrud {
	return &GenericCrudHelper[Model]{
		db: db,
	}
}

// Create implements GenericCrud.
func (h *GenericCrudHelper[Model]) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody Model

	// Decode the JSON request body into the struct
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON request body", http.StatusBadRequest)
		return
	}

	id, err := h.db.Create(reqBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON request body", http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, IDResp{ID: id}, http.StatusOK)
}

// Delete implements GenericCrud.
func (h *GenericCrudHelper[Model]) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		w.Header().Set("code", fmt.Sprint(http.StatusBadRequest))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Failed to parse ID", http.StatusBadRequest)
		return
	}

	err = h.db.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	sendJSONResponse(w, IDResp{ID: id}, http.StatusOK)
}

// Read implements GenericCrud.
func (h *GenericCrudHelper[Model]) Read(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		w.Header().Set("code", fmt.Sprint(http.StatusBadRequest))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model, err := h.db.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, model, http.StatusOK)
}

// Update implements GenericCrud.
func (h *GenericCrudHelper[Model]) Update(w http.ResponseWriter, r *http.Request) {
	var reqBody Model

	// Decode the JSON request body into the struct
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON request body", http.StatusBadRequest)
		return
	}

	if reqBody.GetID() == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := h.db.Update(reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, IDResp{ID: id}, http.StatusOK)
}
