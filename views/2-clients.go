package views

import "github.com/Bayan2019/ai-hackathon-2025-api/repositories/database"

type Client struct {
	CstDimID  int64  `json:"cst_dim_id"`
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ClientDetailed struct {
	CstDimID     int64         `json:"cst_dim_id"`
	Gender       string        `json:"gender"`
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	Transactions []Transaction `json:"transactions"`
}

//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////

func databaseClient2viewClient(dbC database.Client) Client {
	return Client{
		CstDimID:  dbC.CstDimID,
		Gender:    dbC.Gender,
		FirstName: dbC.FirstName,
		LastName:  dbC.LastName,
	}
}

func DatabaseClients2viewClients(dbCs []database.Client) []Client {
	clients := []Client{}
	for _, c := range dbCs {
		clients = append(clients, databaseClient2viewClient(c))
	}
	return clients
}
