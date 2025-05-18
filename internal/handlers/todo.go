package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo/internal/models"
	"todo/internal/repository"

	"github.com/go-chi/chi/v5"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	repo *repository.TodoRepository
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(repo *repository.TodoRepository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

// GetAll godoc
// @Summary Get all todos
// @Description Get all todos
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} models.Todo
// @Router /todos [get]
func (h *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	todos, err := h.repo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// GetByID godoc
// @Summary Get a todo by ID
// @Description Get a specific todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 404 {string} string "Todo not found"
// @Router /todos/{id} [get]
func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todo, err := h.repo.FindByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if todo == nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// Create godoc
// @Summary Create a new todo
// @Description Create a new todo with the provided data
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body models.TodoCreate true "Todo data"
// @Success 201 {object} models.Todo
// @Failure 400 {string} string "Bad request"
// @Router /todos [post]
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var todoCreate models.TodoCreate
	if err := json.NewDecoder(r.Body).Decode(&todoCreate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todo := models.Todo{
		Title:       todoCreate.Title,
		Description: todoCreate.Description,
	}

	if err := h.repo.Create(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// Update godoc
// @Summary Update a todo
// @Description Update a todo with the provided data
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body models.TodoUpdate true "Todo data to update"
// @Success 200 {object} models.Todo
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Todo not found"
// @Router /todos/{id} [put]
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todo, err := h.repo.FindByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if todo == nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	var todoUpdate models.TodoUpdate
	if err := json.NewDecoder(r.Body).Decode(&todoUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update fields if provided
	if todoUpdate.Title != nil {
		todo.Title = *todoUpdate.Title
	}
	if todoUpdate.Description != nil {
		todo.Description = *todoUpdate.Description
	}
	if todoUpdate.Completed != nil {
		todo.Completed = *todoUpdate.Completed
	}

	if err := h.repo.Update(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// Delete godoc
// @Summary Delete a todo
// @Description Delete a todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 204 {string} string "No Content"
// @Failure 404 {string} string "Todo not found"
// @Router /todos/{id} [delete]
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Check if todo exists
	todo, err := h.repo.FindByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if todo == nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// At the end of the file, add these helper functions

// respondWithError sends an error response with the specified status code
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON sends a JSON response with the specified status code
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
