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

func agent() *consulapi.Agent {
	return client().Agent()
}

func kvkeys(c *cli.Context) {
	prefix := c.Args().First()
	keys, _, _ := kv().Keys(prefix, "", nil)

	for _, key := range keys {
		fmt.Printf("%s\n", key)
	}
}

func kvDelTree(c *cli.Context) {
	prefix := c.Args().First()
	_, err := kv().DeleteTree(prefix, nil)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	fmt.Printf("deleted %s\n", prefix)
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

func findNode(nodeName string) bool {
	members, err := agent().Members(false)
	if err != nil {
		fmt.Printf("err: can't list members\n")
		return false
	}

	for _, member := range members {
		if member.Name == nodeName {
			return true
		}
	}

	return false
}

func nodeEject(c *cli.Context) {
	nodeName := c.Args().First()

	if findNode(nodeName) {
		err := agent().ForceLeave(nodeName)
		if err != nil {
			fmt.Printf("err: %s\n", err)
			return
		}
	} else {
		fmt.Printf("err: cound't find node %s\n", nodeName)
	}

	fmt.Printf("removed %s\n", nodeName)
}

func nodeList(c *cli.Context) {
	members, err := agent().Members(false)
	if err != nil {
		fmt.Printf("err: can't list members\n")
		return
	}

	for _, member := range members {
		fmt.Println(member.Name)
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "consulcli"
	app.Usage = "consul api cli client"
	app.Version = "0.3.0"

	app.Commands = []cli.Command{
		{
			Name:   "kv-get",
			Usage:  "get an item from the kv store",
			Action: kvget,
		},
		{
			Name:   "kv-keys",
			Usage:  "list keys in the kv store",
			Action: kvkeys,
		},
		{
			Name:   "kv-list",
			Usage:  "list items in the kv store",
			Action: kvlist,
		},
		{
			Name:   "kv-deltree",
			Usage:  "delete trees in the kv store",
			Action: kvDelTree,
		},
		{
			Name:   "node-eject",
			Usage:  "eject node",
			Action: nodeEject,
		},
		{
			Name:   "node-list",
			Usage:  "list nodes",
			Action: nodeList,
		},
	}

	app.Run(os.Args)
}
