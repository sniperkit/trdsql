package trdsql

import (
	"bufio"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	runewidth "github.com/mattn/go-runewidth"
)

// VfOut is Vertical Format output
type VfOut struct {
	writer    *bufio.Writer
	termWidth int
	hsize     int
	header    []string
	count     int
}

func (trdsql *TRDSQL) vfOutNew() Output {
	var err error
	vf := &VfOut{}
	vf.writer = bufio.NewWriter(trdsql.OutStream)
	vf.termWidth, _, err = terminal.GetSize(0)
	if err != nil {
		vf.termWidth = 40
	}
	return vf
}

// First is preparation
func (vf *VfOut) First(columns []string) error {
	vf.header = make([]string, len(columns))
	vf.hsize = 0
	for i, col := range columns {
		if vf.hsize < runewidth.StringWidth(col) {
			vf.hsize = runewidth.StringWidth(col)
		}
		vf.header[i] = col
	}
	return nil
}

// RowWrite is Actual output
func (vf *VfOut) RowWrite(values []interface{}, columns []string) error {
	vf.count++
	fmt.Fprintf(vf.writer,
		"---[ %d]%s\n", vf.count, strings.Repeat("-", (vf.termWidth-16)))
	for i, col := range vf.header {
		v := vf.hsize - runewidth.StringWidth(col)
		fmt.Fprintf(vf.writer,
			"%s%s | %-s\n",
			strings.Repeat(" ", v+2),
			col,
			valString(values[i]))
	}
	return nil
}

// Last is flush
func (vf *VfOut) Last() error {
	return vf.writer.Flush()
}
