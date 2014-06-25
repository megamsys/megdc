package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/megamsys/cloudinabox/cmd"
	"launchpad.net/gnuflag"
	"strconv"
	"net/http"
)

type GulpStart struct {
	manager *cmd.Manager
	fs      *gnuflag.FlagSet
	dry     bool
}

func (g *GulpStart) Info() *cmd.Info {
	desc := `starts the gulpd daemon, and connects to queue.

If you use the '--dry' flag gulpd will do a dry run(parse conf/jsons) and exit.

`
	return &cmd.Info{
		Name:    "start",
		Usage:   `start [--dry] [--config]`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *GulpStart) Run(context *cmd.Context, client *cmd.Client) error {
	// The struc will also have the c.manager
	// c.manager
	// Now using this value start the queue.
	RunServer(c.dry)        
	return nil
}

func (c *GulpStart) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("gulpd", gnuflag.ExitOnError)
		c.fs.BoolVar(&c.dry, "config", false, "config: the configuration file to use")
		c.fs.BoolVar(&c.dry, "c", false, "dry-run: does not start the gulpd (for testing purpose)")
		c.fs.BoolVar(&c.dry, "dry", false, "dry-run: does not start the gulpd (for testing purpose)")
		c.fs.BoolVar(&c.dry, "d", false, "dry-run: does not start the gulpd (for testing purpose)")
	}
	return c.fs
}

type GulpStop struct {
	fs   *gnuflag.FlagSet
	bark bool
}

func (g *GulpStop) Info() *cmd.Info {
	desc := `stops the gulpd daemon, and shutsdown the queue.

If you use the '--bark' flag gulpd will notify daemon status.

`
	return &cmd.Info{
		Name:    "stop",
		Usage:   `stop [--bark]`,
		Desc:    desc,
		MinArgs: 0,
	}
}

//The stop has a design issue.
func (c *GulpStop) Run(context *cmd.Context, client *cmd.Client) error {
	// Now using this value stop the queue.
	StopServer(c.bark)
	return nil
}

func (c *GulpStop) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("gulpd", gnuflag.ExitOnError)
		c.fs.BoolVar(&c.bark, "bark", false, "bark: does a notify of the daemon status (to rmq)")
		c.fs.BoolVar(&c.bark, "b", false, "bark: does a notify of the daemon status (to rmq)")
	}
	return c.fs
}

type GulpUpdate struct {
	fs     *gnuflag.FlagSet
	name   string
	status string
}

func (c *GulpUpdate) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "update",
		Usage:   "update",
		Desc:    "Update service data, using [email/api_key] from the configuration file.",
		MinArgs: 0,
	}
}

func (c *GulpUpdate) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("gulpd", gnuflag.ExitOnError)
		c.fs.StringVar(&c.name, "name", "", "name: app/service host name to update (eg: mobcom.megam.co)")
		c.fs.StringVar(&c.name, "n", "", "n: app/service host name to update (eg: mobcom.megam.co)")
		c.fs.StringVar(&c.status, "status", "", "status: app/server status to update (supported: running, notrunning)")
		c.fs.StringVar(&c.status, "s", "", "s: app/server status to update (supported: running, notrunning)")
	}
	return c.fs
}

func (c *GulpUpdate) Run(ctx *cmd.Context, client *cmd.Client) error {
	if len(c.status) <= 0 || len(c.name) <= 0 {
		fmt.Println("Nothing to update.")
		return nil
	}

//we need to move into a struct
	tmpinp := map[string]string{
		"node_name":     c.name,
		"accounts_id":   "",
		"status":        c.status,
		"appdefnsid":    "",
		"boltdefnsid":   "",
		"new_node_name": "",
	}

//and this as well. 
	jsonMsg, err := json.Marshal(tmpinp)

	if err != nil {
		return err
	}
	
	authly, err := cmd.NewAuthly("/nodes/update", jsonMsg)
	if err != nil {
		return err
	}

	url, err := cmd.GetURL("/nodes/update")
	if err != nil {
		return err
	}
	
	fmt.Println("==> " + url)
	authly.JSONBody = jsonMsg
	
	err = authly.AuthHeader()
	if err != nil {
		return err
	}
	client.Authly = authly

	request, err := http.NewRequest("POST", url, bytes.NewReader(jsonMsg))
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
    fmt.Println(strconv.Itoa(resp.StatusCode) + " ....code")
	if resp.StatusCode == http.StatusNoContent {
		fmt.Fprintln(ctx.Stdout, "Service successfully updated.")
	}
	return nil
}

