package controllers

import (
	"net/http"
	"strconv"

	"github.com/Bayan2019/ai-hackathon-2025-api/views"
	"github.com/go-chi/chi/v5"
)

// GetClients godoc
// @Tags Clients
// @Summary      Get User profile
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Success      200  {array} views.Client "OK"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't Get user"
// @Router       /v1/clients [get]
// @Security Bearer
func (uh *UsersHandlers) GetClients(w http.ResponseWriter, r *http.Request, user views.User) {
	clients, err := uh.DB.GetClients(r.Context())
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "couldn't get clients", err)
		return
	}
	views.RespondWithJSON(w, http.StatusOK, views.DatabaseClients2viewClients(clients))
}

// GetClients godoc
// @Tags Clients
// @Summary      Get User profile
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Param cst_dim_id path int true "cst_dim_id"
// @Success      200  {object} views.ClientDetailed "OK"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't Get user"
// @Router       /v1/clients/{cst_dim_id} [get]
// @Security Bearer
func (uh *UsersHandlers) GetClient(w http.ResponseWriter, r *http.Request, user views.User) {
	cstDimID, err := strconv.ParseInt(chi.URLParam(r, "cst_dim_id"), 10, 64)
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Invalid id", err) // 400
		return
	}
	client, err := uh.DB.GetClientByCstDimId(r.Context(), cstDimID)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "couldn't get client", err)
		return
	}

	transactions, err := uh.DB.GetTransactionsOfClient(r.Context(), cstDimID)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "couldn't get transactions", err)
		return
	}

	behaviors, err := uh.DB.GetBehaviorsOfClient(r.Context(), cstDimID)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "couldn't get transactions", err)
		return
	}

	views.RespondWithJSON(w, http.StatusOK, views.ClientDetailed{
		FirstName:    client.FirstName,
		LastName:     client.LastName,
		CstDimID:     client.CstDimID,
		Gender:       client.Gender,
		Transactions: views.DatabaseGetTransactionsOfClientRows2viewTransactions(transactions),
		Behaviors:    behaviors,
	})
}
