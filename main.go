package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	slackToken = os.Getenv("SLACK_TOKEN")
	channelID  = os.Getenv("ChannelID")
)

type Target struct {
	URL      string        `json:"url"`
	Method   string        `json:"method"`
	Req      int           `json:"requests_per_seconds"`
	Duration time.Duration `json:"duration"`
}

func (t *Target) Attack() (rep vegeta.Metrics) {
	rate := vegeta.Rate{Freq: t.Req, Per: time.Second}
	duration := t.Duration * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: t.Method,
		URL:    t.URL,
	})
	attacker := vegeta.NewAttacker()

	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		rep.Add(res)
	}
	rep.Close()

	return rep
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		form := strings.Split(message.Body, " ")[2:]
		req, _ := strconv.Atoi(form[2])
		duration, _ := strconv.Atoi(form[3])

		t := Target{
			URL:      form[0][1 : len(form[0])-1],
			Method:   form[1],
			Req:      req,
			Duration: time.Duration(duration),
		}
		m := t.Attack()

		const fmtstr = "Requests\t[total, rate]\t%d, %.2f\n" +
			"Duration\t[total, attack, wait]\t%s, %s, %s\n" +
			"Latencies\t[mean, 50, 95, 99, max]\t%s, %s, %s, %s, %s\n" +
			"Bytes In\t[total, mean]\t%d, %.2f\n" +
			"Bytes Out\t[total, mean]\t%d, %.2f\n" +
			"Success\t[ratio]\t%.2f%%\n" +
			"Status Codes\t[code:count]\t"

		text := fmt.Sprintf(fmtstr,
			m.Requests, m.Rate,
			m.Duration+m.Wait, m.Duration, m.Wait,
			m.Latencies.Mean, m.Latencies.P50, m.Latencies.P95, m.Latencies.P99, m.Latencies.Max,
			m.BytesIn.Total, m.BytesIn.Mean,
			m.BytesOut.Total, m.BytesOut.Mean,
			m.Success*100,
		)

		attachment := slack.Attachment{
			Text: text + form[0],
		}
		attachments := []slack.Attachment{attachment}

		param := slack.PostMessageParameters{Attachments: attachments}
		api := slack.New(slackToken)
		_, _, _ = api.PostMessage(channelID, "Result", param)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
