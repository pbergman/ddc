package ddc

import (
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type debug struct {
	inner io.ReadWriteCloser
	outer io.Writer
}

func (d *debug) Read(p []byte) (int, error) {

	n, err := d.inner.Read(p)

	if err != nil {
		return n, err
	}

	_, _ = fmt.Fprintf(d.outer, "> %s\n", strings.Replace(strings.TrimSpace(hex.Dump(p[:n])), "\n", "\n> ", -1))

	return n, err
}

func (d *debug) Write(p []byte) (int, error) {
	n, err := d.inner.Write(p)

	if err != nil {
		return n, err
	}

	_, _ = fmt.Fprintf(d.outer, "< %s\n", strings.Replace(strings.TrimSpace(hex.Dump(p[:n])), "\n", "\n< ", -1))

	return n, err
}

func (d *debug) Close() error {
	_ = d.inner.Close()

	if v, o := d.outer.(io.Closer); o {
		_ = v.Close()
	}

	return nil
}
