package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/armon/consul-api"
)

type runStatus int

const (
	runOK runStatus = iota
	runErr
	runExit
)

func help() {
	fmt.Println("kvlist, kvget, kvkeys")
}

func client() *consulapi.Client {
	client, _ := consulapi.NewClient(consulapi.DefaultConfig())
	return client
}

func kv() *consulapi.KV {
	return client().KV()
}

func kvkeys(prefix string) {
	keys, _, _ := kv().Keys(prefix, "", nil)

	for _, key := range keys {
		fmt.Printf("%s\n", key)
	}
}

func kvlist(prefix string) {
	pairs, _, _ := kv().List(prefix, nil)

	for _, pair := range pairs {
		fmt.Printf("key: %s, value: %s\n", pair.Key, string(pair.Value))
	}
}

func kvget(key string) {
	pair, _, err := kv().Get(key, nil)
	if err != nil {
		fmt.Errorf("err: %v\n", err)
		return
	}

	if pair == nil {
		fmt.Errorf("couldn't find key '%s'", key)
		return
	}

	fmt.Printf("key: %s, value: %s\n", pair.Key, string(pair.Value))
}

func runner(cmd string) runStatus {
	split := strings.Split(strings.TrimSpace(cmd), " ")

	switch split[0] {
	case "kvkeys":
		if len(split) > 1 {
			kvkeys(split[1])
			return runOK
		} else {
			kvkeys("")
			return runOK
		}
	case "kvlist":
		if len(split) > 1 {
			kvlist(split[1])
			return runOK
		} else {
			kvlist("")
			return runOK
		}
	case "kvget":
		if len(split) > 1 {
			kvget(split[1])
			return runOK
		}
	case "help":
		help()
		return runOK
	default:
		return runErr
	}

	return runErr
}

func main() {

	done := false

	for !done {
		// cmd := []string{}
		fmt.Print("> ")
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			cmd := scanner.Text()
			if err := scanner.Err(); err != nil {
				fmt.Println("err: %s", err)
			} else {
				status := runner(cmd)

				switch status {
				case runErr:
					fmt.Printf("unknown command %s\n", cmd)
				}

			}

			fmt.Print("> ")
		}
	}

}
