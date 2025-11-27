package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Bayan2019/ai-hackathon-2025-api/configuration"
	"github.com/Bayan2019/ai-hackathon-2025-api/repositories/database"
	"github.com/Bayan2019/ai-hackathon-2025-api/views"
	"github.com/asafschers/goscore"
)

type TransactionsHandlers struct {
	DB    *database.Queries
	Model goscore.RandomForest
	// jwtSecret   string
	// email       string
	// appPassword string
}

func NewTransactionsHandlers(config configuration.ApiConfiguration) *TransactionsHandlers {
	return &TransactionsHandlers{
		DB:    config.DB,
		Model: config.Model,
	}
}

// GetClients godoc
// @Tags Transactions
// @Summary      Get User profile
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Success      200  {array} views.Transaction "OK"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't Get user"
// @Router       /v1/transactions [get]
// @Security Bearer
func (uh *TransactionsHandlers) GetTransactions(w http.ResponseWriter, r *http.Request, user views.User) {
	transactions, err := uh.DB.GetTransactions(r.Context())
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "couldn't get clients", err)
		return
	}
	views.RespondWithJSON(w, http.StatusOK, views.DatabaseGetTransactionsRows2viewTransactions(transactions))
}

// GetProbability godoc
// @Tags Transactions
// @Summary      Get Result
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Bearer AccessToken"
// @Param request body views.TransactionData true "Transaction data"
// @Success      200  {object} views.ClassificationResult "OK"
// @Failure   	 401  {object} views.ErrorResponse "No token Middleware"
// @Failure   	 404  {object} views.ErrorResponse "Not found User Middleware"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't Get user"
// @Router       /v1/transactions [post]
// @Security Bearer
func (uh *TransactionsHandlers) GetProbability(w http.ResponseWriter, r *http.Request, user views.User) {
	decoder := json.NewDecoder(r.Body)
	td := views.TransactionData{}

	err := decoder.Decode(&td)
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Error parsing JSON of TransactionData", err)
		return
	}

	// Prepare features for scoring
	features := map[string]interface{}{
		"amount":                       td.Amount,
		"monthly_os_changes":           td.MonthlyOsChanges,
		"monthly_phone_model_changes":  td.MonthlyPhoneModelChanges,
		"logins_last_7_days":           td.LoginsLast7Days,
		"logins_last_30_days":          td.LoginsLast30Days,
		"login_frequency_7d":           td.LoginFrequency7D,
		"freq_change_7d_vs_mean":       td.FreqDhange7DvsMean,
		"logins_7d_over_30d_ratio":     td.Logins7DOver30DRatio,
		"avg_login_interval_30d":       td.AvgLoginInterval30d,
		"std_login_interval_30d":       td.StdLoginInterval30d,
		"var_login_interval_30d":       td.VarLoginInterval30d,
		"ewm_login_interval_7d":        td.EwmLoginInterval7d,
		"burstiness_login_interval":    td.BurstinessLoginInterval,
		"fano_factor_login_interval":   td.FanoFactorLoginInterval,
		"zscore_avg_login_interval_7d": td.ZscoreAvgLoginInterval7d,
		"hour":                         td.Transdate.Hour(),
		"GIONEE":                       (td.LastPhoneModelCategorical == "GIONEE"),
		"Google":                       (td.LastPhoneModelCategorical == "Google"),
		"HONOR":                        (td.LastPhoneModelCategorical == "HONOR"),
		"HUAWEI":                       (td.LastPhoneModelCategorical == "HUAWEI"),
		"Honor":                        (td.LastPhoneModelCategorical == "Honor"),
		"Huawei":                       (td.LastPhoneModelCategorical == "Huawei"),
		"Meizu":                        (td.LastPhoneModelCategorical == "Meizu"),
		"Motorola":                     (td.LastPhoneModelCategorical == "Motorola"),
		"OPPO":                         (td.LastPhoneModelCategorical == "OPPO"),
		"OnePlus":                      (td.LastPhoneModelCategorical == "OnePlus"),
		"Oppo":                         (td.LastPhoneModelCategorical == "Oppo"),
		"Realme":                       (td.LastPhoneModelCategorical == "Realme"),
		"Samsung":                      (td.LastPhoneModelCategorical == "Samsung"),
		"TECNO":                        (td.LastPhoneModelCategorical == "TECNO"),
		"Tecno":                        (td.LastPhoneModelCategorical == "Tecno"),
		"Vivo":                         (td.LastPhoneModelCategorical == "Vivo"),
		"Xiaomi":                       (td.LastPhoneModelCategorical == "Xiaomi"),
		"iPhone":                       (td.LastPhoneModelCategorical == "iPhone"),
		"implyForteApp":                (td.LastPhoneModelCategorical == "implyForteApp"),
		"x":                            (td.LastPhoneModelCategorical == "HONOR"),
		"iOS":                          (td.LastOsCategorical == "iOS"),
		"mib":                          (td.LastOsCategorical == "mib"),
		"mibWebv3":                     (td.LastOsCategorical == "mibWebv3"),
		"year":                         td.Transdate.Year(),
		"month":                        td.Transdate.Month(),
		"day":                          td.Transdate.Day(),
		"dayofweek":                    td.Transdate.Weekday(),
		"timestamp_unix":               td.Transdate.UnixNano(),
	}

	// Score the features
	score, err := uh.Model.Score(features, "target") // Example: score for class "1"
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Error getting score JSON for TransactionData", err)
		return
	}
	// log.Printf("Prediction score for class 1: %f", score)
	views.RespondWithJSON(w, http.StatusOK, views.ClassificationResult{
		ProbabilityOfFraud: score,
	})
}
