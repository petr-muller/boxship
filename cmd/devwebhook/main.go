package main

import (
	"flag"
	"fmt"
	"os"

	"sigs.k8s.io/prow/pkg/phony"
)

const devHMAC = "devhmac"

func main() {
	var address, eventType, payload string
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.StringVar(&address, "address", "http://localhost:8888/hook", "Webhook endpoint address")
	fs.StringVar(&eventType, "event", "", "GitHub event type (e.g., pull_request, issue_comment)")
	fs.StringVar(&payload, "payload", "", "Path to JSON payload file")
	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if eventType == "" || payload == "" {
		fmt.Fprintln(os.Stderr, "both --event and --payload are required")
		os.Exit(1)
	}

	data, err := os.ReadFile(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading payload file: %v\n", err)
		os.Exit(1)
	}

	if err := phony.SendHook(address, eventType, data, []byte(devHMAC)); err != nil {
		fmt.Fprintf(os.Stderr, "error sending webhook: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("webhook sent successfully")
}
