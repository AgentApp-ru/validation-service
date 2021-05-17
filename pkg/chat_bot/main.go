package chat_bot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"validation_service/pkg/log"
)

var UptimeChatChannel = "https://rocketchat.b2bpolis.ru/hooks/RQPxhvZc3BFvGQxhC/bs6JK6qhmCJRrR54SJqLxQAytWuN6qpZhdSTMoaXFB4A9CCg"

func sendToChat(data map[string]string) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		log.Logger.Error(err)
	}
	resp, err := http.Post(UptimeChatChannel, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Logger.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Logger.Error("not OK from chat")
	}
}
