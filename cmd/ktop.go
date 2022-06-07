package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yannickepstein/ktop/pkg/kapi"
	"github.com/yannickepstein/ktop/pkg/pods"
	"golang.org/x/net/context"
)

var ErrNamespaceMissing = errors.New("namespace missing")

func main() {
	namespace := flag.String("n", "default", "Namespace")
	flag.Parse()
	if namespace == nil {
		fmt.Fprintln(os.Stderr, ErrNamespaceMissing)
		os.Exit(1)
	}
	client, err := kapi.NewClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	statsTable, err := pods.Stats(client)(ctx, *namespace)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(statsTable)
}
