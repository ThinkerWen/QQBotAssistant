package reply

import (
	"QQBotAssistant/config"
	"github.com/spf13/viper"
	"strings"
)

type Sequence struct {
	Sender   int64
	Receiver int64
	State    AutoReplyState
}

type AutoReplyState struct {
	Step   int
	Ask    string
	Answer string
	Range  int64
}

func matchAsk(ask string, group int64) string {
	for _, v := range config.AutoReplyList {
		if int64(v["range"].(int)) != 0 && group != int64(v["range"].(int)) {
			continue
		}
		if matchPattern(ask, v["ask"].(string)) {
			return v["answer"].(string)
		}
	}
	return ""
}

func matchPattern(str, pattern string) bool {
	if pattern == "*" {
		return true
	} else if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
		substr := strings.TrimPrefix(pattern, "*")
		substr = strings.TrimSuffix(substr, "*")
		return strings.Contains(str, substr)
	} else if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(str, suffix)
	} else if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(str, prefix)
	} else {
		return str == pattern
	}
}

func checkReply(msg string, state *AutoReplyState) int {
	switch state.Step {
	case 0:
		state.Ask = msg
		state.Step = 1
	case 1:
		state.Answer = msg
		state.Step = 2
	case 2:
		if msg == "1" || msg == "2" {
			state.Step = 3
		}
	default:
		state.Step = 0
	}
	return state.Step
}

func saveReply(s Sequence) {
	config.AutoReplyList = append(config.AutoReplyList, map[string]interface{}{"ask": s.State.Ask, "answer": s.State.Answer, "range": s.State.Range})
	_ = viper.WriteConfig()
}

func removeSequence(sender, receiver int64) []Sequence {
	data := make([]Sequence, 0)
	for _, v := range sequence {
		if v.Sender == sender && v.Receiver == receiver {
			continue
		}
		data = append(data, v)
	}
	return data
}
