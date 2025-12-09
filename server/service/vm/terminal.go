package vm

import (
	"NanoKVM-Server/proto"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

const (
	messageWait     = 10 * time.Second
	maxMessageSize  = 1024
	TeminalAuthLock = "/etc/kvm/terminal_auth"
)

var login bool = false

type WinSize struct {
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Service) TerminalAuth(ctx *gin.Context) {
	var rsp proto.Response
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.Header("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
		rsp.ErrRsp(ctx, -1, "authentication failed")
		return
	}
	authParts := strings.SplitN(authHeader, " ", 2)
	if len(authParts) != 2 || authParts[0] != "Basic" {
		rsp.ErrRsp(ctx, -1, "authentication failed")
		return
	}

	payload, err := base64.StdEncoding.DecodeString(authParts[1])
	if err != nil {
		rsp.ErrRsp(ctx, -1, "authentication failed")
		return
	}

	credentials := strings.SplitN(string(payload), ":", 2)
	if len(credentials) != 2 {
		rsp.ErrRsp(ctx, -1, "authentication failed")
		return
	}

	sshConfig := &ssh.ClientConfig{
		User: credentials[0],
		Auth: []ssh.AuthMethod{
			ssh.Password(credentials[1]),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		}),
	}

	sshClient, err := ssh.Dial("tcp", "127.0.0.1", sshConfig)
	if err != nil {
		log.Printf("SSH authentication failed for user %s password %s: %v", credentials[0], credentials[1], err)
		rsp.ErrRsp(ctx, -1, "authentication failed")
	}
	defer sshClient.Close()
	file, err := os.OpenFile(TeminalAuthLock, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("failed to create terminal auth lock file: %v", err)
		rsp.ErrRsp(ctx, -1, "authentication failed")
	}
	defer file.Close()
	rsp.OkRsp(ctx)
}

func (s *Service) Terminal(c *gin.Context) {
	// _, err := os.Stat(TeminalAuthLock)
	// if err != nil {
	// 	log.Printf("failed to stat terminal auth lock file: %v", err)
	// 	return
	// }

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("failed to init websocket: %s", err)
		return
	}
	defer func() {
		_ = ws.Close()
	}()

	cmd := exec.Command("/bin/sh")
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"COLORTERM=truecolor",
	)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Errorf("failed to start pty: %s", err)
		return
	}
	defer func() {
		_ = ptmx.Close()
		_ = cmd.Process.Kill()
		// Wait & reap sub-process; needed to prevent zombie processes
		_ = cmd.Wait()
	}()

	go wsWrite(ws, ptmx)
	wsRead(ws, ptmx)
}

// pty to ws
func wsWrite(ws *websocket.Conn, ptmx *os.File) {
	data := make([]byte, maxMessageSize)

	for {
		n, err := ptmx.Read(data)
		if err != nil {
			return
		}

		if n > 0 {
			_ = ws.SetWriteDeadline(time.Now().Add(messageWait))

			err = ws.WriteMessage(websocket.BinaryMessage, data[:n])
			if err != nil {
				log.Errorf("write ws message failed: %s", err)
				return
			}
		}
	}
}

// ws to pty
func wsRead(ws *websocket.Conn, ptmx *os.File) {
	var zeroTime time.Time
	_ = ws.SetReadDeadline(zeroTime)

	for {
		msgType, p, err := ws.ReadMessage()
		if err != nil {
			return
		}

		// resize message
		if msgType == websocket.BinaryMessage {
			var winSize WinSize
			if err := json.Unmarshal(p, &winSize); err == nil {
				_ = pty.Setsize(ptmx, &pty.Winsize{
					Rows: winSize.Rows,
					Cols: winSize.Cols,
				})
			}
			continue
		}

		_, err = ptmx.Write(p)
		if err != nil {
			log.Errorf("failed to write to pty: %s", err)
			return
		}
	}
}
