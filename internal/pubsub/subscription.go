package pubsub

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kunalsinghdadhwal/nyx/internal/data"
	"github.com/kunalsinghdadhwal/nyx/pkg/logger"
)

type SubscriptionRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (s *SubscriptionRequest) GetRegex() *regexp.Regexp {
	pattern, err := regexp.Compile("^(block|(transaction(/(0x[a-zA-Z0-9]{40}|\\*)(/(0x[a-zA-Z0-9]{40}|\\*))?)?)|(event(/(0x[a-zA-Z0-9]{40}|\\*)(/(0x[a-zA-Z0-9]{64}|\\*)(/(0x[a-zA-Z0-9]{64}|\\*)(/(0x[a-zA-Z0-9]{64}|\\*)(/(0x[a-zA-Z0-9]{64}|\\*))?)?)?)?)?))$")
	if err != nil {
		logger.S().Warnln("failed to compile subscription regex:", err)
		return nil
	}
	return pattern
}

func (s *SubscriptionRequest) Topic() string {
	if strings.HasPrefix(s.Type, "block") {
		return "block"
	}

	if strings.HasPrefix(s.Type, "transaction") {
		return "transaction"
	}

	if strings.HasPrefix(s.Type, "event") {
		return "event"
	}

	return ""
}

func (s *SubscriptionRequest) GetLogEventFilters() []string {
	pattern := s.GetRegex()
	if pattern == nil {
		return nil
	}

	matches := pattern.FindStringSubmatch(s.Type)
	return []string{matches[9], matches[11], matches[13], matches[15], matches[17]}
}

func (s *SubscriptionRequest) DoesMatchWithPublishedEventData(event *data.Event) bool {
	matchTopicInEvent := func(topic string, x uint8) bool {
		if !(int(x) < len(event.Topics)) {
			return topic == "*" || topic == ""
		}

		status := false

		switch topic {
		case "*":
			status = true
		case "":
			status = true
		default:
			status = CheckSimilarity(topic, event.Topics[x])
		}
		return status
	}

	filters := s.GetLogEventFilters()
	if filters == nil {
		return false
	}

	status := false

	switch filters[0] {
	case "*", "":
		status = matchTopicInEvent(filters[1], 0) &&
			matchTopicInEvent(filters[2], 1) &&
			matchTopicInEvent(filters[3], 2) &&
			matchTopicInEvent(filters[4], 3)

	default:
		if CheckSimilarity(filters[0], event.Origin) {
			status = matchTopicInEvent(filters[1], 0) &&
				matchTopicInEvent(filters[2], 1) &&
				matchTopicInEvent(filters[3], 2) &&
				matchTopicInEvent(filters[4], 3)
		}
	}
	return status
}


func (s *SubscriptionRequest) GetTransactionFilters() []string {
	pattern := s.GetRegex()
	if pattern == nil {
		return nil
	}

	matches := pattern.FindStringSubmatch(s.Name)
	return []string{matches[4], matches[6]}
}

func CheckSimilarity(a string, b string) bool {
	req, err := regexp.Compile(fmt.Sprintf("(?i)^(%s)$", a))
	if err != nil {
		logger.S().Warnln("failed to compile similarity regex:", err)
		return false
	}
	return req.MatchString(b)
}

func (s *SubscriptionRequest) IsValidTopic() bool {
	pattern := s.GetRegex()
	if pattern == nil {
		return false
	}

	return pattern.MatchString(s.Name)
}
