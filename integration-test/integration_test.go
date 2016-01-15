package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats"
	"github.com/stretchr/testify/require"
)

func TestRemote(t *testing.T) {
	host := os.Getenv("DH")
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:4222", host))
	require.Nil(t, err)
	sub, err := nc.SubscribeSync("test")
	require.Nil(t, err)
	resp, err := http.Post(fmt.Sprintf("http://%s:8080/test", host), "text/clear", strings.NewReader("hallo"))
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	msg, err := sub.NextMsg(2 * time.Second)
	require.Nil(t, err)
	require.Equal(t, "hallo", string(msg.Data))
}
