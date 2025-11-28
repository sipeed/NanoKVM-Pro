package application

import (
	"NanoKVM-Server/service/ws"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type Message struct {
	Step     string `json:"step"`
	Progress int    `json:"progress"`
}

func sendMessage(step string, progress int) error {
	message := &Message{
		Step:     step,
		Progress: progress,
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	clientMap := ws.GetClientMap()
	for client := range clientMap {
		if err := client.Write("update", string(data)); err != nil {
			log.Errorf("send message to client error: %s", err.Error())
		}
	}

	return nil
}

func progressHandler(total int64, downloaded int64) {
	progress := float64(downloaded) / float64(total) * 100
	_ = sendMessage("download", int(progress))
}
