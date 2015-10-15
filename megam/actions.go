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
	"bufio"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
	"io"
	"os"
	"path"
	"strings"
)

const layout = "Jan 2, 2006 at 3:04pm (MST)"
const (
	CMD_PREFIX      = "bash conf/trusty/megam/"
	CMD_SUFFIX      = ".sh"
	CMD_SUFFIX_TEST = "_test.sh"

	LOGPATH = "/var/log/megam/logs"
)

type runMachineActionsArgs struct {
	command   string
	operation string
	writer    io.Writer
}

/*
* Step 1: Install Megam packages. This invokes the script for the platform (trusty) - megam.sh
 */
func megamInstall(pack string) *action.Action {
	return &action.Action{
		Name: pack,
		Forward: func(ctx action.FWContext) (action.Result, error) {
			args := ctx.Params[0].(*runMachineActionsArgs)
			args.command = CMD_PREFIX + pack + "_" + args.operation + CMD_SUFFIX_TEST
			//args.command = CMD_PREFIX + pack + "_" + args.operation + CMD_SUFFIX

			exec, err := Execute(args)
			if err != nil {
				fmt.Println("server insert error")
				return &args, err
			}
			return exec, err
		},
		Backward: func(ctx action.BWContext) {
			log.Printf(" Nothing to recover")
		},
		MinParams: 1,
	}
}

func Execute(args *runMachineActionsArgs) (action.Result, error) {

	var e exec.OsExecutor
	var commandWords []string
	commandWords = strings.Fields(args.command)

	basePath := LOGPATH
	dir := path.Join(basePath, "megdc")

	fileOutPath := path.Join(dir, "megdc_out")
	fileErrPath := path.Join(dir, "megdc_err")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Info("Creating directory: %s\n", dir)
		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return nil, errm
		}
	}
	// open output file
	fout, outerr := os.OpenFile(fileOutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if outerr != nil {
		return nil, outerr
	}
	defer fout.Close()
	// open Error file
	ferr, errerr := os.OpenFile(fileErrPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errerr != nil {
		return nil, errerr
	}
	defer ferr.Close()

	foutwriter := bufio.NewWriter(fout)
	ferrwriter := bufio.NewWriter(ferr)

	defer ferrwriter.Flush()
	defer foutwriter.Flush()

	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, args.writer, args.writer); err != nil {
			//if err := e.Execute(commandWords[0], commandWords[1:], nil, foutwriter, ferrwriter); err != nil {
			return nil, err
		}
	}

	return &args, nil

}
