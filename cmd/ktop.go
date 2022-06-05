package main

import (
	"fmt"
	"os"
	"time"

	"github.com/yannickepstein/ktop/pkg/kapi"
	"github.com/yannickepstein/ktop/pkg/pods"
	"golang.org/x/net/context"
)

func main() {
	client, err := kapi.NewClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	pods.Display(ctx, client)
}
