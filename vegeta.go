package main

import (
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type Target struct {
	Url      string        `json:"url"`
	Method   string        `json:"method"`
	Req      int           `json:"requests_per_seconds"`
	Duration time.Duration `json:"duration"`
}

func (t *Target) Attack() (rep *vegeta.Metrics) {
	rate := vegeta.Rate{Freq: t.Req, Per: time.Second}
	duration := t.Duration * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: t.Method,
		URL:    t.Url,
	})
	attacker := vegeta.NewAttacker()

	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		rep.Add(res)
	}
	rep.Close()

	return rep
}
