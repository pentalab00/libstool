package sample

import (
	"context"
	"os"
	"plugin"
	"syscall"
	"testing"
	"time"

	"golang.org/x/term"
)

func Test_winrm(t *testing.T) {
	if !term.IsTerminal(syscall.Stdout) {
		return
	}

	t.Log("start")

	addr := os.Getenv("ADDRESS")
	id := os.Getenv("ID")
	password := os.Getenv("PASSWORD")

	p, err := plugin.Open("../libstool.so")
	if err != nil {
		t.Fatal(err)
	}
	f, err := p.Lookup("RunWinRM")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	out, err := f.(func(
		ctx context.Context,

		address string, // xx.xx.xx.xx:1234
		user string,
		password string,
		cert []byte,
		key []byte,
		https bool,
		timeout int,
		cmdline string,
	) ([]byte, error))(ctx, addr, id, password, []byte{}, []byte{}, false, 10, "dir")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(out)
	t.Log(string(out))
}
