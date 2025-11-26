package controllers

import (
	"net/http"

	"github.com/Bayan2019/ai-hackathon-2025-api/views"
)

// GetClients godoc
// @Tags Transactions
// @Summary      Get User profile
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Success      200  {array} views.Client "OK"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't Get user"
// @Router       /v1/transactions [get]
// @Security Bearer
func (uh *UsersHandlers) GetTransactions(w http.ResponseWriter, r *http.Request, user views.User) {
	transactions, err := uh.DB.GetTransactions(r.Context())
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "couldn't get clients", err)
		return
	}
	views.RespondWithJSON(w, http.StatusOK, views.DatabaseGetTransactionsRows2viewTransactions(transactions))
}
