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
package megam

import (
	"fmt"
	"strings"
//	"reflect"
	"time"
	"bytes"
	"io"
	"github.com/spf13/cobra"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/megdc/logwriter"
	"github.com/megamsys/megdc/loggers"
	_"github.com/megamsys/megdc/loggers/file"
	_"github.com/megamsys/megdc/loggers/terminal"
)

const (
	ALL = "all"
	NILAVU = "nilavu"
	MEGAMGATEWAY = "megamgateway"
	MEGAMD = "megamd"
	INSTALL = "install"
	UNINSTALL = "uninstall"
)

type MegamLog struct {}

var Logger loggers.Logger

var All, Nilavu, MegamGateway, Megamd, MegamCommon, MegamMonitor bool

//root megam command
var cmdMegam = &cobra.Command{Use: "megam"}  

var megamCmdI, megamCmdU *cobra.Command

func Register(megdcCmd *cobra.Command) {
	megdcCmd.AddCommand(cmdMegam)
	
	cmdMegam.AddCommand(cmdMegamInstall)
	cmdMegam.AddCommand(cmdMegamUninstall)
	megamFlags()
	//cmdMegam.SetHelpTemplate(s)
}

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
 	cmdMegamInstall.PersistentFlags().BoolVarP(&All, ALL, "a", false, "install all megam packages")
 	cmdMegamInstall.PersistentFlags().BoolVarP(&Nilavu, NILAVU, "n", false, "install nilavu package")
 	cmdMegamInstall.PersistentFlags().BoolVarP(&MegamGateway, MEGAMGATEWAY, "g", false, "install megam gateway package")
 	cmdMegamInstall.PersistentFlags().BoolVarP(&Megamd, MEGAMD, "d", false, "install megamd package")
 	
 	cmdMegamUninstall.PersistentFlags().BoolVarP(&All, ALL, "a", false, "uninstall all megam packages")
 	cmdMegamUninstall.PersistentFlags().BoolVarP(&Nilavu, NILAVU, "n", false, "uninstall nilavu package")
 	cmdMegamUninstall.PersistentFlags().BoolVarP(&MegamGateway, MEGAMGATEWAY, "g", false, "uninstall megam gateway package")
 	cmdMegamUninstall.PersistentFlags().BoolVarP(&Megamd, MEGAMD, "d", false, "uninstall megamd package")
 	megamCmdI = cmdMegamInstall
 	megamCmdU = cmdMegamUninstall
 }
 
 // InitializeConfig initializes a config file with sensible default configuration flags.
func Initialize() {
	var outBuffer bytes.Buffer
	m := MegamLog{}
	logWriter := logwriter.LogWriter{Box: &m}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)

	if megamCmdI.PersistentFlags().Lookup(ALL).Changed {
		createPipelineAll(writer, INSTALL)
		return
	}
	
	if megamCmdI.PersistentFlags().Lookup(NILAVU).Changed {
		createPipeline(NILAVU, writer, INSTALL)
		return
	}
	
	if megamCmdI.PersistentFlags().Lookup(MEGAMGATEWAY).Changed {
		createPipeline(MEGAMGATEWAY, writer, INSTALL)
		return
	}
	
	if megamCmdI.PersistentFlags().Lookup(MEGAMD).Changed {
		createPipeline(MEGAMD, writer, INSTALL)
		return
	}
	
	if megamCmdU.PersistentFlags().Lookup(ALL).Changed {
		createPipelineAll(writer, UNINSTALL)
		return
	}
	
	if megamCmdU.PersistentFlags().Lookup(NILAVU).Changed {
		createPipeline(NILAVU, writer, UNINSTALL)
		return
	}
	
	if megamCmdU.PersistentFlags().Lookup(MEGAMGATEWAY).Changed {
		createPipeline(MEGAMGATEWAY, writer, UNINSTALL)
		return
	}
	
	if megamCmdU.PersistentFlags().Lookup(MEGAMD).Changed {
		createPipeline(MEGAMD, writer, UNINSTALL)
		return
	}
}

func createPipelineAll(w io.Writer, operation string) error {
	actions := []*action.Action{
		megamInstall(NILAVU),
		megamInstall(MEGAMGATEWAY),
		megamInstall(MEGAMD),
	}
	pipeline := action.NewPipeline(actions...)
	
	args := runMachineActionsArgs{
		operation:     operation,
		writer:        w,
		}

	err := pipeline.Execute(&args)
	if err != nil {
		log.Errorf("error on execute create pipeline for install all packages", err)
		return err
	}
	return nil
}

func createPipeline(pack string, w io.Writer, operation string) error {
	actions := []*action.Action{
		megamInstall(pack),
	}
	pipeline := action.NewPipeline(actions...)
	
	args := runMachineActionsArgs{
		operation:     operation,
		writer:        w,
		}

	err := pipeline.Execute(&args)
	if err != nil {
		log.Errorf("error on execute create pipeline for install - %s", pack, err)
		return err
	}
	return nil
}

// Log adds a log message to the app. Specifying a good source is good so the
// user can filter where the message come from.
func (box *MegamLog) Log(message, source string) error {
	messages := strings.Split(message, "\n")
	loggers_list := []string{ "file", "terminal"}
	logs := make([]loggers.Boxlog, 0, len(messages))
	for _, msg := range messages {
		if msg != "" {
			bl := loggers.Boxlog{
				Date:    time.Now().In(time.UTC),
				Message: msg,
				Source:  source,
				Name:    source,
			}
			logs = append(logs, bl)
		}
	}
	if len(logs) > 0 {
		for _, logger := range loggers_list {
			a, err := loggers.Get(logger)

			if err != nil {
				log.Errorf("fatal error, couldn't located the Logger %s", logger)
				return err
			}

			Logger = a

			if initializableLogger, ok := Logger.(loggers.InitializableLogger); ok {
				log.Debugf("Notify to [%s] Logger ", logger)
				err = initializableLogger.Notify(source, logs)
				if err != nil {
					log.Errorf("fatal error, couldn't initialize the Logger %s", logger)
					return err
				} 
			}
			//_ = notify(box.Name+"."+box.DomainName, logs)
		}
	}
	return nil
}

