package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

// https://api.slack.com/slack-apps
// https://api.slack.com/internal-integrations
type envConfig struct {
	// Port is server port to be listened.
	Port string `envconfig:"PORT" default:"3000"`

	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// VerificationToken is used to validate interactive messages from slack.
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID" required:"true"`

	// ChannelID is slack channel ID where bot is working.
	// Bot responses to the mention in this channel.
	ChannelID      string `envconfig:"CHANNEL_ID" required:"true"`
	AwsRegion      string `envconfig:"AWS_REGION" required:"true"`
	VegetaQueueURL string `envconfig:"VEGETA_QUEUE_URL" required:"true"`
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	// Listening slack event and response
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(env.BotToken, slack.OptionDebug(true), slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)))
	slackListener := &SlackListener{
		client:         client,
		botID:          env.BotID,
		channelID:      env.ChannelID,
		AwsRegion:      env.AwsRegion,
		VegetaQueueURL: env.VegetaQueueURL,
	}
	go slackListener.ListenAndResponse()

	log.Printf("[INFO] Server listening on :%s", env.Port)
	if err := http.ListenAndServe(":"+env.Port, nil); err != nil {
		log.Printf("[ERROR] %s", err)
		return 1
	}

	return 0
}
