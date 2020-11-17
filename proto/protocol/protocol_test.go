package protocol

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadWriter(t *testing.T) {
	var msg = new(Message)
	msg.Body = bytes.Repeat([]byte{12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12}, 1024)
	msg.Operate = OpHeartBeat
	msg.Version = 1

	buf := bytes.NewBuffer(nil)
	rw := bufio.NewWriter(buf)
	n, err := Write(rw, msg)
	require.Nil(t, err)
	t.Logf("length: %d", n)

	r := bufio.NewReader(buf)
	n2, _, err := Read(r)
	require.Nil(t, err)

	require.Equal(t, n, n2)
}
