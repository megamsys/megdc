// Copyright © 2013-2015 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//Package commands defines and implements command-line commands and flags used by Hugo. Commands and flags are implemented using
//cobra.
package commands

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"
	
)

var All, Nilavu, MegamGateway, Megamd, MegamCommon, MegamMonitor bool

//root megam command
var cmdMegam = &cobra.Command{Use: "megam"}  

var megamCmdI, megamCmdU *cobra.Command

//Install - subcommand for megam 
var cmdMegamInstall = &cobra.Command{
        Use:   "install [Install megam packages]",
        Short: "Install any megam package",
        Long:  `Description: Install any megam package.
        [ Nilavu, MegamGateway, Megamd..]
        `,
        Run: func(cmd *cobra.Command, args []string) {
        	Initialize()
            fmt.Println("Megam Install: " + strings.Join(args, " "))
        },
    }

//Uninstall - subcommand for megam
var cmdMegamUninstall = &cobra.Command{
        Use:   "uninstall [Uninstall megam packages]",
        Short: "Uninstall any megam package",
        Long:  `Description: Uninstall any megam package.
        [ Nilavu, MegamGateway, Megamd..]
        `,
        Run: func(cmd *cobra.Command, args []string) {
        	Initialize()
            fmt.Println("Megam Uninstall: " + strings.Join(args, " "))
        },
    }


 func megamFlags() {
 	cmdMegamInstall.PersistentFlags().BoolVarP(&All, "all", "a", false, "install all megam packages")
 	cmdMegamInstall.PersistentFlags().BoolVarP(&Nilavu, "nilavu", "n", false, "install nilavu package")
 	cmdMegamInstall.PersistentFlags().BoolVarP(&MegamGateway, "megamgateway", "g", false, "install megam gateway package")
 	cmdMegamInstall.PersistentFlags().BoolVarP(&Megamd, "megamd", "d", false, "install megamd package")
 	megamCmdI = cmdMegamInstall
 	megamCmdU = cmdMegamUninstall
 }
 
 // InitializeConfig initializes a config file with sensible default configuration flags.
func Initialize() {

	if megamCmdI.PersistentFlags().Lookup("all").Changed {
		fmt.Println("=========install==========all")
		return
	}
	
	if megamCmdI.PersistentFlags().Lookup("nilavu").Changed {
		fmt.Println("==========install=========nilavu")
		return
	}
	
	if megamCmdI.PersistentFlags().Lookup("megamgateway").Changed {
		fmt.Println("==========install=========gateway")
		return
	}
	
	if megamCmdI.PersistentFlags().Lookup("megamd").Changed {
		fmt.Println("==========install=========megamd")
		return
	}
	
	if megamCmdU.PersistentFlags().Lookup("all").Changed {
		fmt.Println("==========uninstall=========all")
		return
	}
	
	if megamCmdU.PersistentFlags().Lookup("nilavu").Changed {
		fmt.Println("===========uninstall========nilavu")
		return
	}
	
	if megamCmdU.PersistentFlags().Lookup("megamgateway").Changed {
		fmt.Println("===========uninstall========gateway")
		return
	}
	
	if megamCmdU.PersistentFlags().Lookup("megamd").Changed {
		fmt.Println("===========uninstall========megamd")
		return
	}
}
