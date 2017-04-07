package gelf

import (
	"encoding/json"
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

func New(connectionString string) *Gelf {
	url, err := url.Parse(connectionString)
	if err != nil {
		panic(err)
	}
	return &Gelf{
		url: url,
	}
}

// Run starts the logger consuming on the returned channed
func (g *Gelf) Run() chan<- *log.Entry {
	// in a big high traffic app, set a higher buffer
	ch := make(chan *log.Entry, 30000)
	g.manageConnections()
	go func(entries <-chan *log.Entry) {
		var empty byte
		var e *log.Entry
		for e = range entries {
			if g.conn != nil {
				payload := entryToPayload(e)
				payload = append(payload, empty) // when we use tcp, we need to add null byte in the end.
				_, err := g.conn.Write(payload)
				if err != nil {
					println("failed to write: %v", err)
					g.conn.Close()
					g.conn = nil
				} else {
					//msg := fmt.Sprintf("payload size: %d", size)
					//println(msg)
				}
			}
			e.Consumed()
		}
	}(ch)

	return ch
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

func entryToPayload(e *log.Entry) []byte {
	items := make(map[string]interface{})
	items["version"] = "1.1"
	items["host"] = e.Host
	items["level"] = toGelfLevel(e.Level)
	items["short_message"] = e.Message
	items["full_message"] = e.Message
	items["timestamp"] = float64(e.Timestamp.UnixNano()) / float64(time.Second)
	items["_app_id"] = e.AppID
	items["_file"] = e.File

	if e.Line > 0 {
		items["_line"] = e.Line
	}

	for k, v := range e.Fields {
		switch k {
		case "short_message":
			items[k] = v
		default:
			items["_"+k] = v
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
