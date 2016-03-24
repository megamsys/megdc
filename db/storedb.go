package db

import (
//"fmt"
)
const (
	SETTINGS = "setting"
)

func StoreDB(data interface{},dbip string) error{
		t := TableInfo{
			Name: SETTINGS,
			Pks: []string{"Name"},
			Ccms: []string{},
			Db: dbip,
		}
 err := Write(t, data)
 return err
}

