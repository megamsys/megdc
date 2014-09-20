/*
** Copyright [2012-2014] [Megam Systems]
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

package orm

import (
    "database/sql"
    
    "github.com/coopernurse/gorp"
    "log"
)


func OpenDB() *sql.DB {
	
	// connect to db using standard Go database/sql API
    // use whatever database/sql driver you wish
    db, err := sql.Open("sqlite3", "./cloudinabox.db")
    CheckErr(err, "sql.Open failed")
    return db
}

func GetDBMap(db *sql.DB) *gorp.DbMap {
	// construct a gorp DbMap
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
    return dbmap
}

func InitDB(dbmap *gorp.DbMap) error {
	
    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(Users{}, "users").SetKeys(true, "Id")
    dbmap.AddTableWithName(Servers{}, "servers").SetKeys(true, "Id")
    // create the table. in a production system you'd generally
    // use a migration tool, or create the tables via scripts
    err := dbmap.CreateTablesIfNotExists()
    CheckErr(err, "Create tables failed")
    return err
}

func ConnectToTable(dbmap *gorp.DbMap, tablename string, field interface{}) error {
	// add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(field, tablename).SetKeys(true, "Id")
    return nil
}

func DeleteRowFromServerName(dbmap *gorp.DbMap, serverName string) error {
	// delete row manually via Exec
    _, err := dbmap.Exec("delete from servers where Name=?", serverName)
    CheckErr(err, "Exec failed")
    return err
}

func CheckErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
