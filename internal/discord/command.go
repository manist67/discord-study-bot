package discord

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func (gw *Gateway) Handshake() {
	identifyPayload, err := json.Marshal(IdentifyPayload{
		Token: os.Getenv("DISCORD_BOT_TOKEN"),
		Properties: struct {
			Os      string "json:\"os\""
			Browser string "json:\"browser\""
			Device  string "json:\"device\""
		}{
			Os:      runtime.GOOS,
			Browser: "discordbot",
			Device:  "discordbot",
		},
		Intents: 1<<7 | 1<<0,
	})
	if err != nil {
		log.Printf("Handshake marshal error %v", err)
		return
	}

	raw := json.RawMessage(identifyPayload)
	gw.EventChannel <- Event{
		Op: 2,
		D:  &raw,
		S:  &gw.seq,
	}
}

func (gw *Gateway) Reconnect() {
	resumePayload, err := json.Marshal(ResumePayload{
		Token:     os.Getenv("DISCORD_BOT_TOKEN"),
		SessionId: gw.sessionId,
		Seq:       gw.seq,
	})
	if err != nil {
		log.Fatal("Fail to Reconnect : ", err)
	}
	d := json.RawMessage(resumePayload)
	gw.EventChannel <- Event{
		Op: 6,
		S:  &gw.seq,
		D:  &d,
	}
}

func (gw *Gateway) SetResume(p json.RawMessage) {
	var payload ReadyPayload
	if err := json.Unmarshal(p, &payload); err != nil {
		log.Printf("Fail to unmarshal payload %s", string(p))
		return
	}

	gw.setResumeValue(payload.ResumeGatewayUrl, payload.SessionId)
}

func (c *Connection) SendHeartbeat() {
	c.gw.EventChannel <- Event{
		Op: 1,
		S:  &c.gw.seq,
	}
}

func (c *Connection) HandleInvalidSession(event Event) {
	var resumable bool
	if err := json.Unmarshal(*event.D, &resumable); err != nil {
		log.Printf("Handshake marshal error %v", err)
		return
	}
	time.Sleep(time.Duration(1+rand.Intn(5)) * time.Second)
	if resumable {
		c.gw.Reconnect()
	} else {
		c.gw.isReconnect = false
		c.errChan <- errors.New("Reconnect Fail : opcode 6")
	}
}

func (c *Connection) StartHeartbeat(ctx context.Context, event Event) {
	var handshakeEvent HandshakePayload
	if err := json.Unmarshal(*event.D, &handshakeEvent); err != nil {
		log.Printf("Handshake marshal error %v", err)
		return
	}

	log.Printf("Start heartbeat interval %d\n", handshakeEvent.HeartbeatInterval)
	c.Interval = handshakeEvent.HeartbeatInterval

	go func() {
		// send first heatbeat
		jitter := time.Duration(float64(c.Interval)*rand.Float64()) * time.Millisecond
		select {
		case <-time.After(jitter):
			c.SendHeartbeat()
		case <-ctx.Done():
			return
		}

		// send interval heatbeat
		t := time.NewTicker(time.Duration(c.Interval) * time.Millisecond)
		for {
			select {
			case <-t.C:
				if !c.lastAckReceive {
					log.Println("ack is dead.")
					c.errChan <- errors.New("ACK is dead")
					return
				}
				c.lastAckReceive = false
				log.Printf("heartbeat duration: %d", time.Duration(c.Interval))
				c.SendHeartbeat()
			case <-ctx.Done():
				t.Stop()
				return
			}
		}
	}()
}
