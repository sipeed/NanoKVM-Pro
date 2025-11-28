package ws

import (
	"NanoKVM-Server/service/hid"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	conn          *websocket.Conn
	hid           *hid.Hid
	keyboard      chan []int
	mouse         chan []int
	lastHeartbeat time.Time
	mutex         sync.Mutex
}

var (
	clientMap = make(map[*WsClient]bool)
	mutex     sync.Mutex
)

func NewClient(conn *websocket.Conn) {
	client := &WsClient{
		hid:           hid.GetHid(),
		conn:          conn,
		keyboard:      make(chan []int, 200),
		mouse:         make(chan []int, 200),
		lastHeartbeat: time.Time{},
	}

	client.hid.Open()

	AddClient(client)

	go client.Start()
	go client.keepAlive()
}

func AddClient(client *WsClient) {
	mutex.Lock()
	defer mutex.Unlock()
	clientMap[client] = true
}

func RemoveClient(client *WsClient) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(clientMap, client)
}

func GetClientMap() map[*WsClient]bool {
	return clientMap
}
