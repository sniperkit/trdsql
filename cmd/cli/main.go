package main

import (
	"os"

	trdsql "github.com/sniperkit/trdsql/pkg/trdsql"
)

func main() {
	tr := &trdsql.TRDSQL{OutStream: os.Stdout, ErrStream: os.Stderr}
	os.Exit(tr.RunCLI(os.Args))
}
