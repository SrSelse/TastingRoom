package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/centrifugal/gocent/v3"
)

type Message struct {
	Channel string
	Payload map[string]string
}

type CentrifugoProvider struct {
	client *gocent.Client
	logger *slog.Logger
}

func NewCentrifugoProvider(c *gocent.Client, logger *slog.Logger) *CentrifugoProvider {
	return &CentrifugoProvider{client: c, logger: logger}
}

func (p *CentrifugoProvider) HandleMessage(ctx context.Context, message Message) {
	err := p.Send(ctx, message)
	if err != nil {
		p.logger.Error("HandleMessage", "error", err.Error())
	}
}

func (p *CentrifugoProvider) CreateVoteMessage(
	beerId int,
	userId int,
	voteValue int,
	voteNote string,
	username string,
	reason string,
) Message {
	c := fmt.Sprintf("beers:beer-%d", beerId)
	return p.CreateMessageWithChannel(c, map[string]string{
		"beerId":    strconv.Itoa(beerId),
		"userId":    strconv.Itoa(userId),
		"voteValue": strconv.Itoa(voteValue),
		"voteNote":  voteNote,
		"username":  username,
		"reason":    reason,
	})
}

func (p *CentrifugoProvider) CreateBeerMessage(
	roomId int,
	reason string,
) Message {
	c := fmt.Sprintf("beers:room-%d", roomId)
	return p.CreateMessageWithChannel(c, map[string]string{
		"roomId": strconv.Itoa(roomId),
		"reason": reason,
	})

}

func (p *CentrifugoProvider) CreateMessageWithChannel(channel string, payload map[string]string) Message {
	return Message{
		Channel: channel,
		Payload: payload,
	}

}

func (p *CentrifugoProvider) Send(ctx context.Context, message Message) error {

	payload, err := json.Marshal(message.Payload)
	if err != nil {
		return err
	}
	_, err = p.client.Publish(ctx, message.Channel, payload)
	return err
}
