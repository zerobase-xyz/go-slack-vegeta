package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/nlopes/slack"
)

const (
	// action is used for slack attament action.
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)

type SlackListener struct {
	client         *slack.Client
	botID          string
	channelID      string
	AwsRegion      string
	VegetaQueueURL string
}

// LstenAndResponse listens slack events and response
// particular messages. It replies by slack message button.
func (s *SlackListener) ListenAndResponse() {
	rtm := s.client.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "GCTG5978Q"))
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

// handleMesageEvent handles message events.
func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	// Only response in specific channel. Ignore else.
	if ev.Channel != s.channelID {
		log.Printf("%s %s", ev.Channel, ev.Msg.Text)
		return nil
	}

	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	var resText string
	switch m[0] {
	case "attack":
		resText = "Wait until the attack is over"
		sess := session.Must(session.NewSession())
		svc := sqs.New(sess, aws.NewConfig().WithRegion(s.AwsRegion))

		params := &sqs.SendMessageInput{
			MessageBody:  aws.String(ev.Msg.Text),
			QueueUrl:     aws.String(s.VegetaQueueURL),
			DelaySeconds: aws.Int64(1),
		}

		_, err := svc.SendMessage(params)
		if err != nil {
			return err
		}
	default:
		resText = m[0]

	}

	if _, _, err := s.client.PostMessage(ev.Channel, slack.MsgOptionText(resText, false)); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}
