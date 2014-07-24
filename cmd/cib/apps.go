package main

import (
	"fmt"
	"github.com/megamsys/libgo/cmd"
	//	"io/ioutil"
	"launchpad.net/gnuflag"
	//"log"
	//	"net/http"
)

type AppStart struct{}

func (AppStart) Run(context *cmd.Context,client *cmd.Client) error {
	appName := context.Args[0]
	/*platform := context.Args[1]
	b := bytes.NewBufferString(fmt.Sprintf(`{"name":"%s","platform":"%s"}`, appName, platform))
	log.Printf("This is just a crappy print %s", b)
	 url, err := cmd.GetURL("/apps")
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	out := make(map[string]string)
	err = json.Unmarshal(result, &out)
	if err != nil {
		return err
	}
	*/
	fmt.Fprintf(context.Stdout, "App %q is being started!\n", appName)
	fmt.Fprintln(context.Stdout, "Use appreqs list to check the status of the app.")
	return nil
}

func (AppStart) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "startapp",
		Usage:   "startapp <appname> <lifecycle_when>",
		Desc:    "starts the installed app.",
		MinArgs: 1,
	}
}

type AppStop struct {
	yes bool
	fs  *gnuflag.FlagSet
}

func (c *AppStop) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "stopapp",
		Usage: "stopapp appname [--lifecycle-when-yes]",
		Desc: `stops the installed app.

If you don't provide the app name, megam will try to guess it.`,
		MinArgs: 1,
	}
}

func (c *AppStop) Run(context *cmd.Context, client *cmd.Client) error {
	appName := context.Args[0]

	var answer string
	if !c.yes {
		fmt.Fprintf(context.Stdout, `Are you sure you want to remove app "%s"? (y/n) `, appName)
		fmt.Fscanf(context.Stdin, "%s", &answer)
		if answer != "y" {
			fmt.Fprintln(context.Stdout, "Abort.")
			return nil
		}
	}
	/*url, err := cmd.GetURL(fmt.Sprintf("/apps/%s", appName))
	if err != nil {
		return err
	}
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	fmt.Fprintf(context.Stdout, `App "%s" successfully removed!`+"\n", appName)
	*/
	return nil
}

func (c *AppStop) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("appremove", gnuflag.ExitOnError)
		c.fs.BoolVar(&c.yes, "assume-yes", false, "Don't ask for confirmation, just remove the app.")
		c.fs.BoolVar(&c.yes, "y", false, "Don't ask for confirmation, just remove the app.")
	}
	return c.fs
}
