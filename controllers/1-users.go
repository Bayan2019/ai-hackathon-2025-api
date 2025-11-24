package controllers

import (
	"net/http"

	"github.com/Bayan2019/ai-hackathon-2025-api/configuration"
	"github.com/Bayan2019/ai-hackathon-2025-api/repositories/database"
	"github.com/Bayan2019/ai-hackathon-2025-api/views"
)

type UsersHandlers struct {
	DB *database.Queries
	// jwtSecret   string
	// email       string
	// appPassword string
}

func NewUsersHandlers(config configuration.ApiConfiguration) *UsersHandlers {
	return &UsersHandlers{
		DB: config.DB,
	}
}

// GetProfile godoc
// @Tags Users
// @Summary      Get User profile
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Success      200  {object} views.User "OK"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't Get user"
// @Router       /v1/profile [get]
// @Security Bearer
func (uh *UsersHandlers) GetProfile(w http.ResponseWriter, r *http.Request, user views.User) {
	views.RespondWithJSON(w, http.StatusOK, user)
}
