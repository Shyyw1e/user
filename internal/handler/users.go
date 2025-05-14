package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Shyyw1e/user/internal/store"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store store.Store
}

func New(s store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	users, err := h.store.All(ctx)
	if err != nil {
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var u store.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	created, err := h.store.Create(ctx, u)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)

}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	u, err := h.store.Get(ctx, id)
	if err != nil {
		if err == store.ErrNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
		}else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)

}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var u store.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	updated, err := h.store.Update(ctx, id, u)
	if err != nil {
		if err == store.ErrNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)

}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	err := h.store.Delete(ctx, id)
	if err != nil {
		if err == store.ErrNotFound {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
