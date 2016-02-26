package db

import (
  "fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gocassa"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/db"
)


type S struct{
  sy *db.ScyllaDB
}

type TableInfo struct {
	Name    string
	Pks     []string
	Ccms    []string
  Query   map[string]string
}

var noips = []string{"103.56.92.24"}
//A global function which helps to avoid passing config of riak everywhere.
func newDBConn() (*db.ScyllaDB, error) {
	 r,err := db.NewScyllaDB(db.ScyllaDBOpts{
  		KeySpaceName: "testing",
  		NodeIps:      noips,
  		Username:     "",
  		Password:     "",
  		Debug:        true,
  	})
    if err != nil {
      return nil, err
    }
	return r, nil
}


func newScyllaTable(tinfo TableInfo, data interface{}) (*db.ScyllaTable) {
  t, err := newDBConn()
	if err != nil {
		return nil
	}
	log.Debugf("%s (%s, %s)", cmd.Colorfy("  > [scylla] fetch", "blue", "", "bold"),tinfo.Name)
	tbl := t.Table(tinfo.Name, tinfo.Pks, tinfo.Ccms, data)
  errors := tbl.T.(gocassa.TableChanger).CreateIfNotExist()
  if errors != nil {
    fmt.Println(errors)
    return nil
  }
	return tbl
}

func ReadWhere(tinfo TableInfo, data interface{}) error {
   d := newScyllaTable(tinfo, data)
   if d != nil {
     err := d.ReadWhere(db.Alt(tinfo.Query), data)
     if err != nil {
       return err
     }
   }
   return nil
}

func Write(tinfo TableInfo, data interface{}) error {
  t := newScyllaTable(tinfo, data)
  err := t.Upsert(data)
  if err !=nil {
    return err
  }
  return nil
}
