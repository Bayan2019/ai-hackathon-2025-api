package views

import (
	"github.com/Bayan2019/ai-hackathon-2025-api/repositories/database"
)

type Transaction struct {
	FullNameClient string `json:"full_name_of_sender"`
	Transdatetime  string `json:"transdatetime"`
	Amount         int64  `json:"amount"`
	Destination    string `json:"destination"`
	Target         int64  `json:"target"`
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
