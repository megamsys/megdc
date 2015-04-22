package app

import (
	"encoding/json"
	"github.com/megamsys/libgo/action"
	"fmt"
)

type App struct {
    Email              string `json:"email"`
	ApiKey             string `json:"api_key"`
	InstallPackage     string `json:"install_package"`
	NeedMegam          string `json:"need_megam"`
	ClusterName        string `json:"cluster_name"`
	NodeIp             string `json:"node_ip"`
	NodeName           string `json:"node_name"`
	Action             string `json:"action"`
	Command            string
}

// MarshalJSON marshals the app in json format. It returns a JSON object with
//the following keys: name, framework, teams, units, repository and ip.
func (app *App) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["email"] = app.Email
	result["api_key"] = app.ApiKey
	result["install_package"] = app.InstallPackage
	result["need_megam"] = app.NeedMegam
	result["cluster_name"] = app.ClusterName
	result["node_ip"] = app.NodeIp
	result["node_name"] = app.NodeName
	result["action"] = app.Action
	return json.Marshal(&result)
}

//UnmarshalJSON parse app configuration json using App struct.
func (a *App) UnmarshalJSON(b []byte) error {

	var f interface{}
	json.Unmarshal(b, &f)

	m := f.(map[string]interface{})
    a.Email = m["email"].(string)
    a.ApiKey  = m["api_key"].(string)
    a.InstallPackage  = m["install_package"].(string)
    a.NeedMegam = m["need_megam"].(string)
    a.ClusterName  =  m["cluster_name"].(string)
    a.NodeIp   =  m["node_ip"].(string)
    a.NodeName  = m["node_name"].(string)
    a.Action  =  m["action"].(string)

    return nil
}


type CIB struct {
	Command    	  	string
	LocalIP       	string   `json:"localip"`
	LocalHost     	string   `json:"localhost"`
	LocalDisk     	string   `json:"localdisk"`
	RemoteIP       	string   `json:"remoteip"`
	RemoteHost     	string   `json:"remotehost"`  
	RemoteDisk     	string   `json:"remotedisk"`
	Master          bool     `json:"master"`
}


//
// this executes all actions for megam install
//
func MegamInstall() error {
	fmt.Println("app entry")
	actions := []*action.Action{&megamInstall}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(&CIB{})
	if err != nil {
		return err
	}
	return nil
}

//
// this executes all actions for cobbler install
//
func CobblerInstall() error {
	actions := []*action.Action{&cobblerInstall}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(&CIB{})
	if err != nil {
		return err
	}
	return nil
}

//
// this executes all actions for opennebula install
//
func NebulaInstall() error {
	actions := []*action.Action{&opennebulaVerify, &opennebulaPreInstall, &opennebulaInstall, &opennebulaPostInstall}
    //actions := []*action.Action{&nebulaInstall}
	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(&CIB{})
	if err != nil {
		return err
	}
	return nil
}

//
// this executes all actions for opennebula install
//
func OpenNebulaHostMasterInstall() error {
	actions := []*action.Action{&opennebulaSCPSSH, &opennebulaHostMasterVerify, &opennebulaHostMasterInstall, &cephInstall, &cephOneInstall}
	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(&CIB{})
	if err != nil {
		return err
	}
	return nil
}

//
// this executes all actions for opennebula install
//
//MEGAM CHANGE
//Add a new script for installing only ceph without configuration inside slave system
func OpenNebulaHostNodeInstall() error {
	actions := []*action.Action{&opennebulaSCPSSH, &opennebulaHostMasterVerify, &opennebulaHostNodeInstall, &cephInstall, &cephOneInstall}
	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(&CIB{})
	if err != nil {
		return err
	}
	return nil
}

//
// this executes all actions for opennebula install
//
func HANodeInstall() error {
	actions := []*action.Action{&megamInstall, &cobblerInstall, &opennebulaVerify, &opennebulaPreInstall, &opennebulaInstall, &opennebulaPostInstall, &opennebulaHostMasterVerify, &opennebulaHostMasterInstall, &cephInstall, &cephOneInstall}
	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(&CIB{})
	if err != nil {
		return err
	}
	return nil
}

//
// this executes all actions for opennebula install
//
func SCPSSHInstall() error {
	actions := []*action.Action{&opennebulaSCPSSH}
	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(&CIB{})
	if err != nil {
		return err
	}
	return nil
}

func HAProxyInstall(cib *CIB, Stype string) error {
	if Stype == "MASTER" {
		actions := []*action.Action{&haHooksInstall, &haProxyInstall, &haMegamInstall}

		pipeline := action.NewPipeline(actions...)
		err := pipeline.Execute(cib)
		if err != nil {
			return err
		}
	} else {
		actions := []*action.Action{&haMegamInstall}

		pipeline := action.NewPipeline(actions...)
		err := pipeline.Execute(cib)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetEmail returns the email of the app.
func (app *App) GetEmail() string {
	return app.Email
}

// GetApiKey returns the api_key of the app.
func (app *App) GetApiKey() string {
	return app.ApiKey
}

// GetInstallPackage returns the package.
func (app *App) GetInstallPackage() string {
	return app.InstallPackage
}

// GetNeedMegam returns the need megam of the app.
func (app *App) GetNeedMegam() string {
	return app.NeedMegam
}

func (app *App) GetClusterName() string {
	return app.ClusterName
}

func (app *App) GetNodeIp() string {
	return app.NodeIp
}

func (app *App) GetNodeName() string {
    return app.NodeName
}

