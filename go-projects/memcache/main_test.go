package memcache

import (
	"bufio"
	"github.com/joho/godotenv"
	"net"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	teardown()

	os.Exit(exitCode)
}

func setup() {
	err := godotenv.Load("../.env.test.local")
	if err != nil {
	}
}

func teardown() {
}

func makeTcpRequest(data []byte) (string, error) {
	c, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		return "", err
	}
	defer c.Close()
	_, err = c.Write(data)
	if err != nil {
		return "", err
	}

	resp, err := readResponse(c)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func readResponse(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return response, nil
}
