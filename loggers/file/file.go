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

package file

import (
	//	"errors"
	"fmt"
	//	"io"
	//	"net/http"
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/loggers"
	"os"
	"path"
	"text/tabwriter"
	"time"
)

func init() {
	loggers.Register("file", fileManager{})
}

type fileManager struct{}

func (m fileManager) Notify(boxName string, messages []loggers.Boxlog) error {

	basePath := "/var/log/megam/logs"
	dir := path.Join(basePath, boxName)

	filePath := path.Join(dir, boxName+"_log")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Debugf("Creating directory: %s\n", dir)
		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return errm
		}
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Errorf("Error on logs notify: %s", err.Error())
		return err
	}

	defer f.Close()

	for _, msg := range messages {
		if _, err = f.WriteString(String(msg.Date, msg.Message)); err != nil {
			log.Errorf("Error on logs notify: %s", err.Error())
			return err
		}
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
