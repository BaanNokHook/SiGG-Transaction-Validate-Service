package usecase

import (
	"encoding/json"
	"fmt"
	"nextclan/transaction-gateway/transaction-validate-service/internal/entity"
	"nextclan/transaction-gateway/transaction-validate-service/pkg/logger"

	"github.com/streadway/amqp"
)

//Use cases include:
/*
Given verified transaction When submit with validatedTransaction command Then the transaction should complete without errors.
Given verified transaction When receive transaction with verifiedTransaction command Then the transaction should publish into rabbitMQ without error.
*/

//Receive Verified Txn
//Publish to RabbitMQ

type ReceiveRawTransactionFromQueueUseCase struct {
	log logger.Interface
}

func NewReceiveRawTransaction(l logger.Interface) *ReceiveRawTransactionFromQueueUseCase {
	return &ReceiveRawTransactionFromQueueUseCase{log: l}
}

func (txn *ReceiveRawTransactionFromQueueUseCase) Handle(d amqp.Delivery) {
	body := d.Body
	transaction := &entity.ValidateTransaction{}
	err := json.Unmarshal(body, transaction)
	if err != nil {
		txn.log.Debug("Problem parsing txnRecieve: %v", err.Error())
		panic(err)
	} else {
		txn.log.Debug("%s Receive Message %s", d.ConsumerTag, transaction.TransactionId)
		publishRequest := createMessage(transaction.TransactionId)
		userIds := getValidateUsers()
		pubId, err := PushNotificationClient.PublishToUsers(userIds, publishRequest)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Publish Id:", pubId)
		}
	}
}

func createMessage(transactionId string) map[string]interface{} {
	return map[string]interface{}{
		"web": map[string]interface{}{
			"notification": map[string]interface{}{
				"title": "Request Transaction Validate",
				"body":  transactionId,
			},
		},
	}
}

func getValidateUsers() []string {
	userId := "alice2"
	return []string{userId}
}
