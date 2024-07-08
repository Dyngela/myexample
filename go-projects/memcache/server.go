package memcache

import (
	"errors"
	"fmt"
	"github.com/tidwall/evio"
	"log"
	"sync"
)

type conn struct {
	is   evio.InputStream
	addr string
}

type Command struct {
	Operation byte
	Key       string
	Value     string
}

const (
	MaxKeyLen   = 255
	MaxValueLen = 1024 * 1024 // 1 MB

)

var (
	IncompleteCommand = errors.New("incomplete command")
	InvalidOperation  = errors.New("invalid operation")
	Operators         = map[byte]string{
		'*': "SET",
		'&': "GET",
		'#': "DELETE",
		'!': "FLUSH",
		'@': "CLEANUP",
	}
)

type Server struct {
	mu   sync.Mutex
	keys map[string]string
}

func NewServer() *Server {
	return &Server{
		keys: make(map[string]string),
		mu:   sync.Mutex{},
	}
}

func (s *Server) Start(config *Config) {
	var events evio.Events

	events.NumLoops = -1
	events.LoadBalance = evio.RoundRobin
	events.Serving = func(srv evio.Server) (action evio.Action) {
		Logger.Debug().Msgf("echo server started on port %s (loops: %d)", config.ServerPort, srv.NumLoops)
		return
	}
	events.Opened = func(ec evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		ec.SetContext(&conn{})
		return
	}
	events.Closed = func(ec evio.Conn, err error) (action evio.Action) {
		// fmt.Printf("closed: %v\n", ec.RemoteAddr())
		return
	}
	events.Data = s.handleData

	switch config.ServerMode {
	case TCP:
		log.Fatal(evio.Serve(events, fmt.Sprintf("%s://:%s", config.ServerMode, config.ServerPort)))
	case UDP:
		log.Fatal(evio.Serve(events, fmt.Sprintf("%s://:%s", config.ServerMode, config.ServerPort)))
	case BOTH:
		go log.Fatal(evio.Serve(events, fmt.Sprintf("%s://:%s?reuseport=true", "udp", config.ServerPort)))
		go log.Fatal(evio.Serve(events, fmt.Sprintf("%s://:%s?reuseport=true", "tcp", config.ServerPort)))
	}
}

func (s *Server) handleData(ec evio.Conn, in []byte) (out []byte, action evio.Action) {
	Logger.Info().Msgf("received: %s", string(in))
	Logger.Debug().Msgf("echo server received %d bytes: %s", len(in), in)
	if in == nil {
		Logger.Debug().Msgf("wake from %s\n", ec.RemoteAddr())
		return nil, evio.Close
	}
	ctx := ec.Context().(*conn)
	buffer := ctx.is.Begin(in)

	for {
		command, remaining, err := s.readCommand(buffer)
		if err != nil {
			if errors.Is(err, IncompleteCommand) {
				break
			}
			out = append(out, fmt.Sprintf("ERR %s\n", err.Error())...)
			action = evio.Close
			return
		}

		out = append(out, s.processCommand(command)...)
		buffer = remaining
		if len(buffer) == 0 {
			break
		}
	}

	return
}

func (s *Server) readCommand(data []byte) (*Command, []byte, error) {
	if len(data) < 256 {
		return nil, data, IncompleteCommand
	}
	operation := data[0]
	_, ok := Operators[operation]
	if !ok {
		return nil, data, InvalidOperation
	}

	key := string(data[1:256])
	if len(key) > MaxKeyLen {
		return nil, nil, fmt.Errorf("key length exceeds maximum of %d bytes", MaxKeyLen)
	}

	endIdx := -1
	for i := 256; i < len(data); i++ {
		if data[i] == '\x00' {
			endIdx = i
			break
		}
	}
	if endIdx == -1 {
		return nil, data, IncompleteCommand
	}

	value := string(data[256:endIdx])
	if len(value) > MaxValueLen {
		return nil, nil, fmt.Errorf("value length exceeds maximum of %d bytes", MaxValueLen)
	}

	remaining := data[endIdx+1:]

	return &Command{
		Operation: operation,
		Key:       key,
		Value:     value,
	}, remaining, nil
}

func (s *Server) processCommand(command *Command) []byte {
	var out []byte

	switch command.Operation {
	case '*':
		fmt.Printf("Set operation - Key: %s, Value: %s\n", command.Key, command.Value)
		s.mu.Lock()
		s.keys[command.Key] = command.Value
		s.mu.Unlock()
		out = append(out, "OK\n"...)
	case '&':
		fmt.Printf("Get operation - Key: %s\n", command.Key)
		s.mu.Lock()
		value, ok := s.keys[command.Key]
		s.mu.Unlock()
		if !ok {
			out = append(out, "NULL\n"...)
		} else {
			out = append(out, value+"\n"...)
		}
	default:
		fmt.Printf("Unknown operation: %c\n", command.Operation)
		out = append(out, fmt.Sprintf("ERR unknown command '%c'\n", command.Operation)...)
	}

	return out
}
