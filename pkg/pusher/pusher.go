package pusher

import (
	pushnotifications "github.com/pusher/push-notifications-go"
)

type IPusherBeamClient interface {
	PublishToUsers(users []string, request map[string]interface{}) (string, error)
}

type PusherBeamClient struct {
	pushNotification pushnotifications.PushNotifications
}

func NewPusherClient(instanceId string, secretKey string) *PusherBeamClient {
	pushNotification, err := pushnotifications.New(instanceId, secretKey)
	if err != nil {
		panic(err)
	}
	return &PusherBeamClient{
		pushNotification: pushNotification,
	}
}

func (p *PusherBeamClient) PublishToUsers(users []string, request map[string]interface{}) (string, error) {
	pubId, err := p.pushNotification.PublishToUsers(users, request)
	if err != nil {
		panic(err)
	}
	return pubId, err
}
