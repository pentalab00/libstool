package main

import (
	"bytes"
	"context"
	"net"
	"strconv"
	"time"

	"github.com/pentalab00/winrm"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

func RunWinRM(
	ctx context.Context,

	address string, // xx.xx.xx.xx:22
	user string,
	password string,
	cert []byte,
	key []byte,
	https bool,
	timeout int,
	cmdline string,
) ([]byte, error) {

	host, p, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	port, _ := strconv.Atoi(p)

	var client *winrm.Client

	if len(key) > 0 {
		params := winrm.NewParameters("PT60S", "en-US", 153600)
		params.TransportDecorator = func() winrm.Transporter {
			return &winrm.ClientAuthRequest{}
		}
		client, err = winrm.NewClientWithParameters(
			winrm.NewEndpoint(host, port, https, https, nil, cert, key, time.Duration(timeout)*time.Second),
			user, password, params)
		if err != nil {
			return nil, err
		}
	} else {
		client, err = winrm.NewClient(
			winrm.NewEndpoint(host, port, https, https, nil, nil, nil, time.Duration(timeout)*time.Second),
			user, password)
		if err != nil {
			return nil, err
		}
	}

	done := make(chan error)
	var bout, berr []byte
	go func() {
		stdout, stderr, _, err := client.RunPSWithContextWithString(ctx, cmdline, "")
		bout = convEUCKRtoUTF8([]byte(stdout))
		berr = convEUCKRtoUTF8([]byte(stderr))
		// bout, berr = []byte(stdout), []byte(stderr)
		bout = bytes.TrimSpace(bout)
		berr = bytes.TrimSpace(berr)

		if err != nil {
			if done != nil {
				done <- err
			}
			return
		}

		if done != nil {
			done <- nil
		}
	}()

	select {
	case err := <-done:
		if err != nil {
			return berr, err
		}
		return bout, nil

	case <-ctx.Done():
		done = nil
		return bout, ctx.Err()
	}
}

func convEUCKRtoUTF8(in []byte) []byte {
	var bufs bytes.Buffer

	wr := transform.NewWriter(&bufs, korean.EUCKR.NewDecoder())

	wr.Write(in)
	wr.Close()
	return bufs.Bytes()
}
