package gelf

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/jasonsoft/log"
)

// Gelf is an instance of the Gelf logger
type Gelf struct {
	conn net.Conn
	url  *url.URL
}

// New create a new Gelf instance
func New(connectionString string) *Gelf {
	url, err := url.Parse(connectionString)
	if err != nil {
		panic(fmt.Errorf("graylog connectionString is wrong: %v", err))
	}
	g := &Gelf{
		url: url,
	}
	g.manageConnections()
	return g
}

var empty byte

// Log handles the log entry
func (g *Gelf) Log(e log.Entry) error {
	if g.conn != nil {
		payload := entryToPayload(e)
		payload = append(payload, empty) // when we use tcp, we need to add null byte in the end.
		_, err := g.conn.Write(payload)
		if err != nil {
			_ = g.conn.Close()
			g.conn = nil
			return fmt.Errorf("send log to graylog failed: %w", err)
		}

		// msg := fmt.Sprintf("payload size: %d", size)
		// println(msg)

		// msg = fmt.Sprintf("payload body: %s", string(payload))
		// println(msg)
	}

	return nil
}

// Flush clear all buffer and close connection
func (g *Gelf) Flush() error {
	if g.conn != nil {
		g.conn.Close()
	}

	return nil
}

func (g *Gelf) manageConnections() {
	var err error
	if strings.EqualFold(g.url.Scheme, "tcp") {
		g.conn, err = net.Dial("tcp", g.url.Host)
		if err != nil {
			println("gelf tcp connection was failed:", err.Error())
		}
	} else {
		g.conn, err = net.Dial("udp", g.url.Host)
		if err != nil {
			println("gelf udp connection was failed:", err.Error())
		}
	}

	// check connection status every 1 second
	go func() {
		for {
			if g.conn == nil {
				// TODO: tcp is hard-code at the point, we need to remove that later
				newConn, err := net.Dial("tcp", g.url.Host)
				if err == nil {
					g.conn = newConn
					println("created new connection")
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func entryToPayload(e log.Entry) []byte {
	items := make(map[string]interface{})
	items["version"] = "1.1"
	items["level"] = toGelfLevel(e.Level)
	items["short_message"] = e.Message
	items["full_message"] = e.Message
	items["timestamp"] = float64(e.Timestamp.UnixNano()) / float64(time.Second)

	for key, value := range e.Fields {
		switch key {
		case "short_message", "host":
			items[key] = value
		default:
			items["_"+key] = value
		}
	}

	payload, _ := json.Marshal(items)
	return payload
}

func toGelfLevel(level log.Level) uint8 {
	switch level {
	case log.DebugLevel:
		return 7
	case log.InfoLevel:
		return 6
	case log.WarnLevel:
		return 4
	case log.ErrorLevel:
		return 3
	case log.FatalLevel:
		return 2
	default:
		return 1
	}
}
