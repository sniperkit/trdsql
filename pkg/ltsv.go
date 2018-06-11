package trdsql

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// LTSVIn provides methods of the Input interface
type LTSVIn struct {
	reader    *bufio.Reader
	firstRow  map[string]string
	delimiter string
	name      []string
}

// LTSVOut provides methods of the Output interface
type LTSVOut struct {
	writer    *bufio.Writer
	delimiter string
	results   map[string]string
}

func (trdsql *TRDSQL) ltsvInputNew(r io.Reader) (Input, error) {
	lr := &LTSVIn{}
	lr.reader = bufio.NewReader(r)
	lr.delimiter = "\t"
	return lr, nil
}

// FirstRead is read input to determine column of table
func (lr *LTSVIn) FirstRead() ([]string, error) {
	var err error
	lr.firstRow, lr.name, err = lr.read()
	if err != nil {
		return nil, err
	}
	debug.Printf("Column Name: [%v]", strings.Join(lr.name, ","))
	return lr.name, nil
}

// FirstRowRead is read the first row
func (lr *LTSVIn) FirstRowRead(list []interface{}) []interface{} {
	for i := range lr.name {
		list[i] = lr.firstRow[lr.name[i]]
	}
	return list
}

// RowRead is read 2row or later
func (lr *LTSVIn) RowRead(list []interface{}) ([]interface{}, error) {
	record, _, err := lr.read()
	if err != nil {
		return list, err
	}
	for i := range lr.name {
		list[i] = record[lr.name[i]]
	}
	return list, nil
}

func (lr *LTSVIn) read() (map[string]string, []string, error) {
	line, err := lr.readline()
	if err != nil {
		return nil, nil, err
	}
	columns := strings.Split(line, lr.delimiter)
	lvs := make(map[string]string)
	keys := make([]string, 0, len(columns))
	for _, column := range columns {
		kv := strings.SplitN(column, ":", 2)
		if len(kv) != 2 {
			return nil, nil, errors.New("LTSV format error")
		}
		lvs[kv[0]] = kv[1]
		keys = append(keys, kv[0])
	}
	return lvs, keys, nil
}

func (lr *LTSVIn) readline() (string, error) {
	for {
		line, _, err := lr.reader.ReadLine()
		if err != nil {
			return "", err
		}
		tline := strings.TrimSpace(string(line))
		if len(tline) != 0 {
			return tline, nil
		}
	}
}

func (trdsql *TRDSQL) ltsvOutNew() Output {
	lw := &LTSVOut{}
	lw.delimiter = "\t"
	lw.writer = bufio.NewWriter(trdsql.OutStream)
	return lw
}

// First is preparation
func (lw *LTSVOut) First(columns []string) error {
	lw.results = make(map[string]string, len(columns))
	return nil
}

// RowWrite is Actual output
func (lw *LTSVOut) RowWrite(values []interface{}, columns []string) error {
	results := make([]string, len(values))
	for i, col := range values {
		results[i] = columns[i] + ":" + valString(col)
	}
	str := strings.Join(results, lw.delimiter) + "\n"
	_, err := lw.writer.Write([]byte(str))
	return err
}

// Last is flush
func (lw *LTSVOut) Last() error {
	return lw.writer.Flush()
}
