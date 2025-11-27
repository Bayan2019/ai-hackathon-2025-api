package views

import (
	"time"

	"github.com/Bayan2019/ai-hackathon-2025-api/repositories/database"
)

type Transaction struct {
	FullNameClient string `json:"full_name_of_sender"`
	Transdatetime  string `json:"transdatetime"`
	Amount         int64  `json:"amount"`
	Destination    string `json:"destination"`
	Target         int64  `json:"target"`
}

type ClassificationResult struct {
	ProbabilityOfFraud float64 `json:"probability_of_fraud"`
}

type TransactionData struct {
	Amount                    int64     `json:"amount"`
	Transdate                 time.Time `json:"transdate"`
	LastPhoneModelCategorical string    `json:"last_phone_model_categorical"`
	LastOsCategorical         string    `json:"last_os_categorical"`
	MonthlyOsChanges          int32     `json:"monthly_os_changes"`
	MonthlyPhoneModelChanges  int32     `json:"monthly_phone_model_changes"`
	LoginsLast7Days           float64   `json:"logins_last_7_days"`
	LoginsLast30Days          float64   `json:"logins_last_30_days"`
	LoginFrequency7D          float64   `json:"login_frequency_7d"`
	FreqDhange7DvsMean        float64   `json:"freq_change_7d_vs_mean"`
	Logins7DOver30DRatio      float64   `json:"flogins_7d_over_30d_ratio"`
	AvgLoginInterval30d       float64   `json:"avg_login_interval_30d"`
	StdLoginInterval30d       float64   `json:"std_login_interval_30d"`
	VarLoginInterval30d       float64   `json:"var_login_interval_30d"`
	EwmLoginInterval7d        float64   `json:"ewm_login_interval_7d"`
	BurstinessLoginInterval   float64   `json:"burstiness_login_interval"`
	FanoFactorLoginInterval   float64   `json:"fano_factor_login_interval"`
	ZscoreAvgLoginInterval7d  float64   `json:"zscore_avg_login_interval_7d"`
}

//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////

func databaseGetTransactionsOfClientRow2viewTransaction(dbT database.GetTransactionsOfClientRow) Transaction {
	return Transaction{
		Transdatetime: dbT.Transdatetime,
		Amount:        dbT.Amount,
		Destination:   dbT.Direction,
		Target:        dbT.Target,
	}
}

func DatabaseGetTransactionsOfClientRows2viewTransactions(dbTs []database.GetTransactionsOfClientRow) []Transaction {
	transactions := []Transaction{}
	for _, t := range dbTs {
		transactions = append(transactions, databaseGetTransactionsOfClientRow2viewTransaction(t))
	}
	return transactions
}

func databaseGetTransactionsRow2viewTransaction(dbT database.GetTransactionsRow) Transaction {
	return Transaction{
		FullNameClient: dbT.FirstName.String + " " + dbT.LastName.String,
		Transdatetime:  dbT.Transdatetime,
		Amount:         dbT.Amount,
		Destination:    dbT.Direction,
		Target:         dbT.Target,
	}
}

func DatabaseGetTransactionsRows2viewTransactions(dbTs []database.GetTransactionsRow) []Transaction {
	transactions := []Transaction{}
	for _, t := range dbTs {
		transactions = append(transactions, databaseGetTransactionsRow2viewTransaction(t))
	}
	return transactions
}
