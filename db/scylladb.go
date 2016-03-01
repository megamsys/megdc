package db

import (
	log "github.com/Sirupsen/logrus"
//	"github.com/megamsys/gocassa"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/db"
)


type TableInfo struct {
	Name    string
	Pks     []string
	Ccms    []string
  Db      string
  Query   map[string]string
}

//A global function which helps to avoid passing config of riak everywhere.
func newDBConn(noips string) (*db.ScyllaDB, error) {
	 r,err := db.NewScyllaDB(db.ScyllaDBOpts{
  		KeySpaceName: "testing",
  		NodeIps:      []string{noips},
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
  t, err := newDBConn(tinfo.Db)
	if err != nil {
		return nil
	}
	log.Debugf("%s (%s, %s)", cmd.Colorfy("  > [scylla] fetch", "blue", "", "bold"),tinfo.Name)
	tbl := t.Table(tinfo.Name, tinfo.Pks, tinfo.Ccms, data)
  /*errors := tbl.T.(gocassa.TableChanger).CreateIfNotExist()
  if errors != nil {
    fmt.Println(errors)
    return nil
  }*/
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
