package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/apex/gateway"
	"github.com/nlopes/slack"
)

var (
	slackToken = os.Getenv("SLACK_TOKEN")
)

func main() {
	ApexGatewayDisabled := os.Getenv("APEX_GATEWAY_DISABLED")
	http.HandleFunc("/attack", AttackRequest)

	if ApexGatewayDisabled == "true" {
		log.Fatal(http.ListenAndServe(":3000", nil))
	} else {
		log.Fatal(gateway.ListenAndServe(":3000", nil))
	}
}

func AttackRequest(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if s.Token != slackToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/attack":
		form := strings.Split(s.Text, " ")
		if len(form) != 4 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		req, _ := strconv.Atoi(form[2])
		duration, _ := strconv.Atoi(form[3])

		t := Target{
			Url:      form[0],
			Method:   form[1],
			Req:      req,
			Duration: time.Duration(duration),
		}

		params := &slack.Msg{Text: "Start Attack"}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
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
			Text: text,
		}
		attachments := []slack.Attachment{attachment}

		param := slack.PostMessageParameters{Attachments: attachments}
		api := slack.New(slackToken)
		_, _, _ = api.PostMessage(s.ChannelID, "攻撃完了", param)

		return

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
