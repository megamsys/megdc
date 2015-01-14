package app

import (
	"bytes"
	"fmt"
	"github.com/megamsys/cloudinabox/models/orm"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
	"log"
	"strings"
	"crypto/rand"
    "math/big"
//	"time"
)

const layout = "Jan 2, 2006 at 3:04pm (MST)"
const (
  /* 
	opennebulapreinstall  = "bash conf/trusty/opennebula/one_preinstall.sh"
	opennebulaverify      = "bash conf/trusty/opennebula/one_verify.sh"
	opennebulapostinstall = "bash conf/trusty/opennebula/one_postinstall.sh"
	opennebulainstall     = "bash conf/trusty/opennebula/one_install.sh"
	opennebulascpssh      = "bash conf/trusty/opennebula/scp_ssh.sh"
	opennebulahostverify  = "bash conf/trusty/opennebulahost/host_verify.sh"
	opennebulahostinstall = "bash conf/trusty/opennebulahost/host_install.sh"
	megam                 = "bash conf/trusty/megam/megam.sh"
	cobbler               = "bash conf/trusty/cobblerd/cobbler.sh install"
	ceph                  = "bash conf/trusty/ceph/ceph_install.sh install osd1="/storage1" osd2="/storage2" osd3="/storage3""
	cephone               = "bash conf/trusty/ceph/ceph_one_install.sh"
	*/
	opennebulapreinstall  = "bash conf/trusty/opennebula/one_preinstall_test.sh"
	opennebulaverify      = "bash conf/trusty/opennebula/one_verify_test.sh"
	opennebulapostinstall = "bash conf/trusty/opennebula/one_postinstall_test.sh"
	opennebulainstall     = "bash conf/trusty/opennebula/one_install_test.sh"
	opennebulascpssh      = "bash conf/trusty/opennebula/scp_ssh_test.sh"
	opennebulahostverify  = "bash conf/trusty/opennebulahost/host_verify_test.sh"
	opennebulahostinstall = "bash conf/trusty/opennebulahost/host_install_test.sh"
	megam                 = "bash conf/trusty/megam/megam_test.sh"
	cobbler               = "bash conf/trusty/cobblerd/cobbler_test.sh"
	ceph                  = "bash conf/trusty/ceph/ceph_install_test.sh"
	cephone               = "bash conf/trusty/ceph/ceph_one_install_test.sh"
   
)

func CIBExecutor(cib *CIB) (action.Result, error) {
	var e exec.OsExecutor
	var b bytes.Buffer
	var commandWords []string
	commandWords = strings.Fields(cib.Command)
	fmt.Println(commandWords, len(commandWords))

	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, &b, &b); err != nil {
			log.Printf("stderr:%s", b)
			return nil, err
		}
	}

	log.Printf("%s", b)
	return &b, nil
}

/*
* Step 1: Install Megam packages. This invokes the script for the platform (trusty) - megam.sh
 */
var megamInstall = action.Action{
	Name: "megamInstall",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		cib.Command = megam
		exec, err := CIBExecutor(&cib)
		if err != nil {
			fmt.Println("server insert error")
			return &cib, err
		}
		// write server details in database
		// insert rows - auto increment PKs will be set properly after the insert
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)
		newserver := orm.NewServer("MEGAM")
		orm.ConnectToTable(dbmap, "servers", newserver)
		err = dbmap.Insert(&newserver)
		defer db.Close()
		if err != nil {
			fmt.Println("server insert error")
			return &cib, err
		}
		return exec, err
	},
	Backward: func(ctx action.BWContext) {
		//app := ctx.FWResult.(*App)
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)
		err := orm.DeleteRowFromServerName(dbmap, "MEGAM")
		if err != nil {
			fmt.Println("Server delete error")
			//return &cib, err
		}
		defer db.Close()
		log.Printf(" Nothing to recover")
	},
	MinParams: 1,
}

var cobblerInstall = action.Action{
	Name: "cobblerInstall",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		cib.Command = cobbler
		exec, err1 := CIBExecutor(&cib)
		if err1 != nil {
			fmt.Println("server insert error")
			return &cib, err1
		}
		// write server details in database
		// insert rows - auto increment PKs will be set properly after the insert
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)
		newserver := orm.NewServer("COBBLER")
		orm.ConnectToTable(dbmap, "servers", newserver)
		err := dbmap.Insert(&newserver)
		if err != nil {
			fmt.Println("server insert error======>")
			return &cib, err
		}
		return exec, err1

	},
	Backward: func(ctx action.BWContext) {
		//app := ctx.FWResult.(*App)
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)
		err := orm.DeleteRowFromServerName(dbmap, "COBBLER")
		if err != nil {
			log.Printf("Server delete error")
			///return &cib, err
		}
		log.Printf(" Nothing to recover")
	},
	MinParams: 1,
}

var opennebulaVerify = action.Action{
	Name: "opennebulaVerify",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		cib.Command = opennebulaverify
		exec, err1 := CIBExecutor(&cib)

		return exec, err1
	},
	Backward: func(ctx action.BWContext) {
		log.Printf("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var opennebulaInstall = action.Action{
	Name: "opennebulaInstall",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		cib.Command = opennebulainstall
		exec, err1 := CIBExecutor(&cib)

		return exec, err1
	},
	Backward: func(ctx action.BWContext) {
		log.Printf("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var opennebulaPreInstall = action.Action{
	Name: "opennebulaPreInstall",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		cib.Command = opennebulapreinstall
		exec, err1 := CIBExecutor(&cib)

		return exec, err1
	},
	Backward: func(ctx action.BWContext) {
		log.Printf("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var opennebulaPostInstall = action.Action{
	Name: "opennebulaPostInstall",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		cib.Command = opennebulapostinstall
		exec, err1 := CIBExecutor(&cib)
		if err1 != nil {
			fmt.Println("server insert error")
			return &cib, err1
		}
		// write server details in database

		// insert rows - auto increment PKs will be set properly after the insert
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)
		newserver := orm.NewServer("OPENNEBULA")
		orm.ConnectToTable(dbmap, "servers", newserver)
		err := dbmap.Insert(&newserver)

		if err != nil {
			fmt.Println("server insert error======>")
			return &cib, err
		}
		return exec, err1
	},
	Backward: func(ctx action.BWContext) {
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)
		err := orm.DeleteRowFromServerName(dbmap, "OPENNEBULA")
		if err != nil {
			log.Printf("Server delete error")
			///return &cib, err
		}
		log.Printf(" Nothing to recover")
	},
	MinParams: 1,
}

var opennebulaSCPSSH = action.Action{
	Name: "opennebulaSCPSSH",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		var server orm.Servers
		db := orm.OpenDB()
	    dbmap := orm.GetDBMap(db)
	    err := dbmap.SelectOne(&server, "select * from servers where Name=?", "OPENNEBULAHOST")
	    fmt.Println(err)
		cib.Command = opennebulascpssh + " " + server.IP
		exec, err1 := CIBExecutor(&cib)

		return exec, err1
	},
	Backward: func(ctx action.BWContext) {
		log.Printf("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var opennebulaHostMasterVerify = action.Action{
	Name: "opennebulaHostVerify",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		cib.Command = opennebulahostverify
		exec, err1 := CIBExecutor(&cib)

		return exec, err1
	},
	Backward: func(ctx action.BWContext) {
		log.Printf("[%s] Nothing to recover")
	},
	MinParams: 1,
}

var opennebulaHostMasterInstall = action.Action{
	Name: "opennebulaHostInstall",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		cib.Command = opennebulahostinstall
		exec, err1 := CIBExecutor(&cib)
		if err1 != nil {
			fmt.Println("server insert error")
			return &cib, err1
		}
		// write server details in database

		// insert rows - auto increment PKs will be set properly after the insert
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)
		newserver := orm.NewServer("OPENNEBULAHOST")
		orm.ConnectToTable(dbmap, "servers", newserver)
		err := dbmap.Insert(&newserver)

		if err != nil {
			fmt.Println("server insert error======>")
			return &cib, err
		}
		return exec, err1
	},
	Backward: func(ctx action.BWContext) {
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)
		err := orm.DeleteRowFromServerName(dbmap, "OPENNEBULAHOST")
		if err != nil {
			log.Printf("Server delete error")
			///return &cib, err
		}
		log.Printf(" Nothing to recover")
	},
	MinParams: 1,
}

var opennebulaHostNodeInstall = action.Action{
	Name: "opennebulaHostNodeInstall",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		cib.Command = opennebulahostinstall
		exec, err1 := CIBExecutor(&cib)
		if err1 != nil {
			fmt.Println("node insert error")
			return &cib, err1
		}
		
		return exec, err1
	},
	Backward: func(ctx action.BWContext) {
		log.Printf(" Nothing to recover")
	},
	MinParams: 1,
}

/*var opennebulaHostInstall = action.Action{
	Name: "opennebulaHostInstall",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		var server orm.Servers
		cib.Command = opennebulahostinstall
		exec, err1 := CIBExecutor(&cib)
		if err1 != nil {
			fmt.Println("server insert error")
			return &cib, err1
		}
		// write server details in database

		// insert rows - auto increment PKs will be set properly after the insert
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)
		nodename := "OPENNEBULAHOST"

		server = orm.NewServer("OPENNEBULAHOST")
		err := dbmap.SelectOne(&server, "select * from servers where Name=?", nodename)
		if err != nil {
			fmt.Println("server select error======>")
			return &cib, err
		}
		err3 := orm.DeleteRowFromServerName(dbmap, nodename)
		if err3 != nil {
			log.Printf("Server delete error")
			return &cib, err3
		}
		time := time.Now()
		update_server := orm.Servers{Id: server.Id, Name: server.Name, Install: true, IP: server.IP, InstallDate: server.InstallDate, UpdateDate: time.Format(layout)}
		orm.ConnectToTable(dbmap, "servers", update_server)
		err2 := dbmap.Insert(&update_server)
		//server.Install = true
		//_, err2 := dbmap.Update(&server)
		if err2 != nil {
			fmt.Println("server insert error======>")
			return &cib, err2
		}

		return exec, err1
	},
	Backward: func(ctx action.BWContext) {
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)
		err := orm.DeleteRowFromServerName(dbmap, "OPENNEBULAHOST")
		if err != nil {
			log.Printf("Server delete error")
			///return &cib, err
		}
		log.Printf(" Nothing to recover")
	},
	MinParams: 1,
}*/


/*
* Step 4: Install Ceph Storage
 */
var cephInstall = action.Action{
	Name: "cephInstall",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		cib.Command = ceph
		exec, err := CIBExecutor(&cib)
		if err != nil {
			fmt.Println("server insert error")
			return &cib, err
		}
		return exec, err
	},
	Backward: func(ctx action.BWContext) {
		log.Printf(" Nothing to recover")
	},
	MinParams: 1,
}

var cephOneInstall = action.Action{
	Name: "cephOneInstall",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var cib CIB
		cib.Command = cephone
		exec, err := CIBExecutor(&cib)
		if err != nil {
			fmt.Println("server insert error")
			return &cib, err
		}
		return exec, err
	},
	Backward: func(ctx action.BWContext) {
		log.Printf(" Nothing to recover")
	},
	MinParams: 1,
}

func randString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    symbols := big.NewInt(int64(len(alphanum)))
    states := big.NewInt(0)
    states.Exp(symbols, big.NewInt(int64(n)), nil)
    r, err := rand.Int(rand.Reader, states)
    if err != nil {
        panic(err)
    }
    var bytes = make([]byte, n)
    r2 := big.NewInt(0)
    symbol := big.NewInt(0)
    for i := range bytes {
        r2.DivMod(r, symbols, symbol)
        r, r2 = r2, r
        bytes[i] = alphanum[symbol.Int64()]
    }
    return string(bytes)
}
