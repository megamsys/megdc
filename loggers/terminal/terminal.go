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

package terminal

import (
	"fmt"
//	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megdc/loggers"
	//"strings"
	"time"
	"text/tabwriter"
	"github.com/megamsys/libgo/cmd"	
	"bytes"
)

func init() {
	loggers.Register("terminal", terminalManager{})
}

type terminalManager struct{}


func (m terminalManager) Notify(boxName string, messages []loggers.Boxlog) error {

	for _, msg := range messages {
		fmt.Println(String(msg.Date, msg.Message))
	}
	return nil
}

func String(date time.Time, msg string) string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	t := date
	b.Write([]byte(cmd.Colorfy("["+t.Format("Mon Jan _2 15:04:05 2006")+"] ", "white", "", "bold") + "\t" +
		cmd.Colorfy(msg, "white", "", "")))
	fmt.Fprintln(w)
	w.Flush()
	return b.String()
}
