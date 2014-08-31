package main

import (
	"github.com/megamsys/libgo/cmd"
	"launchpad.net/gnuflag"
)

/////// Server node

type CIBNodeStart struct {
	fs *gnuflag.FlagSet
}

func (g *CIBNodeStart) Info() *cmd.Info {
	desc := `starts the cib node daemon.


`
	return &cmd.Info{
		Name:    "startnode",
		Usage:   `startnode`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *CIBNodeStart) Run(context *cmd.Context, client *cmd.Client) error {
	RunNode()
	return nil
}

func (c *CIBNodeStart) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("cib", gnuflag.ExitOnError)
	}
	return c.fs
}

/*

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

func (c *CIBUpdate) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("cib", gnuflag.ExitOnError)
		c.fs.StringVar(&c.name, "name", "", "name: app/service host name to update (eg: mobcom.megam.co)")
		c.fs.StringVar(&c.name, "n", "", "n: app/service host name to update (eg: mobcom.megam.co)")
		c.fs.StringVar(&c.status, "status", "", "status: app/server status to update (supported: running, notrunning)")
		c.fs.StringVar(&c.status, "s", "", "s: app/server status to update (supported: running, notrunning)")
	}
	return c.fs
}

func (c *CIBUpdate) Run(ctx *cmd.Context, client *cmd.Client) error {
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
*/
