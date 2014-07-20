package main

import (
	"fmt"
	"os"

	"github.com/armon/consul-api"
	"github.com/codegangsta/cli"
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

func kvkeys(c *cli.Context) {
	prefix := c.Args().First()
	keys, _, _ := kv().Keys(prefix, "", nil)

	for _, key := range keys {
		fmt.Printf("%s\n", key)
	}
}

func kvlist(c *cli.Context) {
	prefix := c.Args().First()
	pairs, _, _ := kv().List(prefix, nil)

	for _, pair := range pairs {
		fmt.Printf("key: %s, value: %s\n", pair.Key, string(pair.Value))
	}
}

func kvget(c *cli.Context) {
	key := c.Args().First()

	pair, _, err := kv().Get(key, nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	if pair == nil {
		fmt.Printf("couldn't find key '%s'", key)
		return
	}

	fmt.Printf("key: %s, value: %s\n", pair.Key, string(pair.Value))
}

func main() {

	app := cli.NewApp()
	app.Name = "consulcli"
	app.Usage = "consul api cli client"

	app.Commands = []cli.Command{
		{
			Name:   "kvget",
			Usage:  "get an item from the kv store",
			Action: kvget,
		},
		{
			Name:   "kvkeys",
			Usage:  "list keys in the kv store",
			Action: kvkeys,
		},
		{
			Name:   "kvlist",
			Usage:  "list items in the kv store",
			Action: kvlist,
		},
	}

	app.Run(os.Args)
}
