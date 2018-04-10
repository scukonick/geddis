package main

import (
	"log"
	"os"
	"sort"

	"encoding/json"

	"github.com/scukonick/geddis/cli/client"
	"github.com/urfave/cli"
)

func cliSetString(c *cli.Context) error {
	url := c.String("url")
	client := geddiclient.NewClient(url)

	key := c.String("key")
	value := c.String("value")
	ttl := c.Int("ttl")

	client.SetString(key, value, ttl)

	return nil
}

func cliGetString(c *cli.Context) error {
	url := c.String("url")
	client := geddiclient.NewClient(url)

	key := c.String("key")

	client.GetString(key)

	return nil
}

func cliSetArr(c *cli.Context) error {
	url := c.String("url")
	client := geddiclient.NewClient(url)

	key := c.String("key")
	value := c.StringSlice("value")
	ttl := c.Int("ttl")

	client.SetArr(key, value, ttl)

	return nil
}

func cliGetArr(c *cli.Context) error {
	url := c.String("url")
	client := geddiclient.NewClient(url)

	key := c.String("key")

	client.GetArr(key)

	return nil
}

func cliGetArrIndex(c *cli.Context) error {
	url := c.String("url")
	client := geddiclient.NewClient(url)

	key := c.String("key")
	index := c.Int("index")

	client.GetArrIndex(key, int32(index))

	return nil
}

func cliSetMap(c *cli.Context) error {
	url := c.String("url")
	client := geddiclient.NewClient(url)

	key := c.String("key")
	value := c.String("map")
	ttl := c.Int("ttl")

	m := make(map[string]string, 10)

	err := json.Unmarshal([]byte(value), &m)
	if err != nil {
		return err
	}

	client.SetMap(key, m, ttl)

	return nil
}

func cliGetMap(c *cli.Context) error {
	url := c.String("url")
	client := geddiclient.NewClient(url)

	key := c.String("key")

	client.GetMap(key)

	return nil
}

func cliGetMapSubKey(c *cli.Context) error {
	url := c.String("url")
	client := geddiclient.NewClient(url)

	key := c.String("key")
	subkey := c.String("subkey")

	client.GetMapSubKey(key, subkey)

	return nil
}

func cliDelete(c *cli.Context) error {
	url := c.String("url")
	client := geddiclient.NewClient(url)

	log.Println("deleting")
	key := c.String("key")

	client.Delete(key)

	return nil
}

func cliKeys(c *cli.Context) error {
	url := c.String("url")
	client := geddiclient.NewClient(url)

	key := c.String("key")

	client.Keys(key)

	return nil
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url, u",
			Value: "http://127.0.0.1:8080",
			Usage: "Base URL of Geddis",
		},
	}

	keyFlag := cli.StringFlag{
		Name:  "key",
		Usage: "Key of the element",
	}

	subKeyFlag := cli.StringFlag{
		Name:  "subkey",
		Usage: "Key of the element in map",
	}

	ttlFlag := cli.IntFlag{
		Name:  "ttl",
		Value: 0,
		Usage: "ttl for the element",
	}

	keyValueFlags := []cli.Flag{
		keyFlag,
		ttlFlag,
		cli.StringFlag{
			Name:  "value",
			Usage: "value for the element",
		},
	}

	arrKVFlag := []cli.Flag{
		keyFlag,
		ttlFlag,
		cli.StringSliceFlag{
			Name:  "value",
			Usage: "value for the element",
		},
	}

	indexFlag := cli.IntFlag{
		Name:  "index",
		Usage: "index of the element in the array",
	}

	mapFlag := cli.StringFlag{
		Name:  "map",
		Usage: "json representation of the map",
	}

	app.Commands = []cli.Command{
		{
			Name:  "strings",
			Usage: "commands for work with string values",
			Subcommands: []cli.Command{
				{
					Name:   "set",
					Usage:  "Set value for string identified by key",
					Flags:  append(app.Flags, keyValueFlags...),
					Action: cliSetString,
				},
				{
					Name:   "get",
					Usage:  "Get value of the string identified by key",
					Flags:  append(app.Flags, keyFlag),
					Action: cliGetString,
				},
			},
		},
		{
			Name:  "arrays",
			Usage: "commands for work with array values",
			Subcommands: []cli.Command{
				{
					Name:   "set",
					Usage:  "Set value for array identified by key",
					Flags:  append(app.Flags, arrKVFlag...),
					Action: cliSetArr,
				},
				{
					Name:   "get",
					Usage:  "Get value of the array identified by key",
					Flags:  append(app.Flags, keyFlag),
					Action: cliGetArr,
				},
				{
					Name:   "getIndex",
					Usage:  "Get value of element 'index' from the array identified by key",
					Flags:  append(app.Flags, keyFlag, indexFlag),
					Action: cliGetArrIndex,
				},
			},
		},
		{
			Name:  "maps",
			Usage: "commands for work with map values",
			Subcommands: []cli.Command{
				{
					Name:   "set",
					Usage:  "Set value for map identified by key",
					Flags:  append(app.Flags, keyFlag, ttlFlag, mapFlag),
					Action: cliSetMap,
				},
				{
					Name:   "get",
					Usage:  "Get value of the map identified by key",
					Flags:  append(app.Flags, keyFlag),
					Action: cliGetMap,
				},
				{
					Name:   "getSubKey",
					Usage:  "Get value of element 'subkey' from the map identified by key",
					Flags:  append(app.Flags, keyFlag, subKeyFlag),
					Action: cliGetMapSubKey,
				},
			},
		},
		{
			Name:  "common",
			Usage: "other commands",
			Subcommands: []cli.Command{
				{
					Name:   "delete",
					Usage:  "Delete value by key",
					Flags:  append(app.Flags, keyFlag),
					Action: cliDelete,
				},
				{
					Name:   "keys",
					Usage:  "Return keys corresponding to provided submask",
					Flags:  append(app.Flags, keyFlag),
					Action: cliKeys,
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
