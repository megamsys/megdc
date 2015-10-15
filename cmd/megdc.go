/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package cmd

import (
	//"fmt"
	//	"strings"
	"github.com/megamsys/megdc/megam"
	"github.com/megamsys/megdc/one"
	"github.com/megamsys/megdc/ceph"
	"github.com/spf13/cobra"
	"os"
)

var megdcCmd = &cobra.Command{Use: "megdc"}

//Execute adds all child commands to the root command HugoCmd and sets flags appropriately.
func Execute() {

	AddCommands()
	if err := megdcCmd.Execute(); err != nil {
		// the err is already logged by Cobra
		os.Exit(-1)
	}
}

var s = `
Usage:{{if .Runnable}}
  {{.UseLine}}{{if .HasFlags}} [flags]{{end}}{{end}}{{if .HasSubCommands}}
  {{ .CommandPath}} [command]{{end}}{{if gt .Aliases 0}}
  
Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}

Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasLocalFlags}}
  
Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasSubCommands }}
  
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}

`

//AddCommands adds child commands to the root command HugoCmd.
func AddCommands() {
	megam.Register(megdcCmd)
	one.Register(megdcCmd)
	ceph.Register(megdcCmd)
}
