package main

import (
	"os"

	"github.com/aws-controllers-k8s/ack-scanner/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
