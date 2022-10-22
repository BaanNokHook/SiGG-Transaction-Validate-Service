package usecase

import (
	"nextclan/transaction-gateway/transaction-validate-service/pkg/pusher"
	messaging "nextclan/transaction-gateway/transaction-validate-service/pkg/rabbitmq"

	"github.com/streadway/amqp"
)

type (
	ReceiveRawTransactionFromQueue interface {
		Handle(d amqp.Delivery)
	}
)

var MessagingClient messaging.IMessagingClient
var PushNotificationClient pusher.IPusherBeamClient
