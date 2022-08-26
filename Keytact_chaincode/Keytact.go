package main


import (
	"encoding/json"
	"fmt"
	"time"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}


type Trade struct{
	Trade_ID string `json:"trade_id"`
	Product string `json:"product"`
	Address string `json:"address"`
	Seller_Name string `json:"seller_name"`
	Purchaser_Name string `json:"purchaser_name"`
	Status string `json:"status"`
}

type HistoryQueryResult struct { //  Product History
	Record    *Trade    `json:"record"`
	TxId     string    `json:"transaction_id"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

func (s *SmartContract) Register(ctx contractapi.TransactionContextInterface, trade_id string, product string, seller_name string) error {
	trade := Trade{
		Trade_ID: trade_id,
		Product: product,
		Seller_Name: seller_name,
		Status: "Regist",
	}
	tradeAsBytes, _ := json.Marshal(trade)
	return ctx.GetStub().PutState(trade_id, tradeAsBytes)
}

func (s *SmartContract) GetTrade(ctx contractapi.TransactionContextInterface, trade_id string) (*Trade, error) {
	tradeAsBytes, err := ctx.GetStub().GetState(trade_id)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if tradeAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", trade_id)
	}

	trade := new(Trade)
	_ = json.Unmarshal(tradeAsBytes, trade)

	return trade, nil
}


func(s *SmartContract) Request(ctx contractapi.TransactionContextInterface, trade_id string, product string, address string, purchaser_name string) error {
	// tradeAsBytes, err := ctx.GetStub().GetState(trade_id)
	tradeData, err := s.GetTrade(ctx, trade_id)
	if err != nil {
		fmt.Errorf("Failed to read from world state. %s", err.Error())
		return nil
	}

	if (tradeData.Status == "Buy Requset") || (tradeData.Status == "Requset Approve") {
		fmt.Printf("Trade ID : %s is over", trade_id)
		return nil
	}
	tradeData.Address = address
	tradeData.Purchaser_Name = purchaser_name
	tradeData.Status = "Buy Requset"

	tradeAsBytes, _ := json.Marshal(tradeData)
	return ctx.GetStub().PutState(trade_id, tradeAsBytes)
}

func (t *SmartContract) GetTradeHistory(ctx contractapi.TransactionContextInterface, product string) ([]HistoryQueryResult, error) {
	log.Printf("GetTradeHistory: Product Name %v", product)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(product)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var trade Trade
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &trade)
			if err != nil {
				return nil, err
			}
		} else {
			trade = Trade{
				Product: product,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &trade,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

/*
type Account struct { // user account 
	User_ID string `json:"user_id"` // id
	Active bool `json:"active"` // active
	Trade_Count int64 `json:"trade_count` // trade count
	Credit float64 `json:"credit"` // credit
	Cash int64 `json:"cash"` // money
}
*/
/*
type HistoryQueryResult struct { //  User History
	Record    *Account    `json:"record"`
	TxId     string    `json:"transaction_id"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}
*/

/*func(s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, seller_id string, trade_id string) error {
	tradeAsBytes, err := ctx.GetStub().GetState(trade_id)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if tradeAsBytes == nil{
		return nil, fmt.Errorf("%s does not exist", trade_id)
	}
	if tradeAsBytes.seller_id != seller_id {
		return nil, fmt.Printf("ID doesn't match")
	}
	tradeAsBytes.status = "Requset Approve"
	return ctx.GetStub().PutState(user_id, tradeAsBytes)
}
*/

/*
func (s *SmartContract) Get_User_Data(ctx contractapi.TransactionContextInterface, user_id string) (*Account, error) {
	AccountAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if AccountAsByte == nil {
		return nil, fmt.Errorf("%s does not exist", user_id)
	}

	account := new(Account)
	_ = json.Unmarshal(AccountAsByte, account)

	return account, nil
}

func (s *SmartContract) Set_User_Data(ctx contractapi.TransactionContextInterface, user_id string, active bool trade_count int64 credit float64 cash int64) error {
	account := Account{
		User_ID: user_id,
		Active: active,
		Trade_Count: trade_count,
		Credit: credit,
		Cash: cash,
	}

	AccountAsByte, _ := json.Marshal(account)

	return ctx.GetStub().PutState(user_id, AccountAsByte)
}
*/



func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create Keytact chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Keytact chaincode: %s", err.Error())
	}
}