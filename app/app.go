package app

import (
	"encoding/json"
	"log"
	"github.com/megamsys/cloudinabox/fs"	
	"github.com/megamsys/cloudinabox/action"	
	"github.com/megamsys/cloudinabox/app/bind"
	"github.com/megamsys/cloudinabox/db"
	"regexp"
)

var (
	cnameRegexp = regexp.MustCompile(`^[a-zA-Z0-9][\w-.]+$`)
	fsystem fs.Fs
)

// Appreq is the main type in megam. An app represents a real world application.
// This struct holds information about the app: its name, address, list of
// teams that have access to it, used platform, etc.
type App struct {
	Env                map[string]bind.EnvVar
    Email              string `json:"email"`
	ApiKey             string `json:"api_key"`
	InstallPackage     string `json:"install_package"`
	NeedMegam          string `json:"need_megam"`	
	ClusterName        string `json:"cluster_name"`
	NodeIp             string `json:"node_ip"`
	NodeName           string `json:"node_name"`
	Action             string `json:"action"`
	Command            string
   // AppReqs *AppRequests
}


// MarshalJSON marshals the app in json format. It returns a JSON object with
//the following keys: name, framework, teams, units, repository and ip.
func (app *App) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["name"] = app.Name
	result["platform"] = app.Platform
	//result["repository"] = repository.ReadWriteURL(app.Name)
	result["ip"] = app.Ip
	result["cname"] = app.CName
	result["launched"] = app.State == "launched"
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

func filesystem() fs.Fs {
	if fsystem == nil {
		fsystem = fs.OsFs{}
	}
	return fsystem
}

// Get queries the database and fills the App object with data retrieved from
// the database. It uses the name of the app as filter in the query, so you can
// provide this field:
//
//     app := App{Name: "myapp"}
//     err := app.Get()
//     // do something with the app
/*func (app *App) Get(reqId string) error {
log.Printf("Get message %v", reqId)
	if app.Type != "addon" {
	conn, err := db.Conn("appreqs")
	if err != nil {	
		return err
	}	
	appout := &AppRequests{}
	conn.FetchStruct(reqId, appout)	
	app.AppReqs = appout
	defer conn.Close()
	} else {
	  conn, err := db.Conn("addonconfigs")
	if err != nil {	
		return err
	}	
	appout := &AppConfigurations{}
	conn.FetchStruct(reqId, appout)	
	app.AppConf = appout
	log.Printf("Get message from riak  %v", appout)
	defer conn.Close()
	}
	//fetch it from riak.
	// conn.Fetch(app.id)
	// store stuff back in the appreq object.
	return nil
}*/

// StartsApp creates a new app.
//
// Starts the app :
func Remove(app *App) error {
	actions := []*action.Action{&remove}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
	}
	return nil
}

//
// verify needed packages
//
func install(app *App) error {
	actions := []*action.Action{&install}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(app)
	if err != nil {
		return &AppLifecycleError{app: app.Name, Err: err}
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

// Env returns app.Env
func (app *App) Envs() map[string]bind.EnvVar {
	return app.Env
}

func (app *App) GetNodeIp() string {
	return app.NodeIp
}

func (app *App) GetNodeName() string {
    return app.NodeName
}    

func (app *App) GetAction() string {
    return app.Action
}    


