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

package setting

import (
	"fmt"
//	"net/url"
	"os"
//	"path/filepath"
//	"strings"
	"time"

	"github.com/Unknwon/goconfig"
//	"github.com/howeyc/fsnotify"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
//	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils/captcha"
//	"github.com/beego/compress"
//	"github.com/beego/i18n"
//	"github.com/beego/social-auth"
//	"github.com/beego/social-auth/apps"
)


var (
	AppName             string
	AppVer              string
	AppHost             string
	AppUrl              string
	AppLogo             string
	EnforceRedirect     bool
	AvatarURL           string
	LoginRememberDays   int
	LoginMaxRetries     int
	LoginFailedBlocks   int
    TimeZone            string
    DateFormat          string
	DateTimeFormat      string
	DateTimeShortFormat string
	CookieRememberName  string
	CookieUserName      string
	)

var (
	Cfg     *goconfig.ConfigFile
	Cache   cache.Cache
	Captcha *captcha.Captcha
)

var (
	GlobalConfPath   = "conf/global/app.ini"
//	CompressConfPath = "conf/compress.json"
)

// LoadConfig loads configuration file.
func LoadConfig() *goconfig.ConfigFile {
	var err error

	//if fh, _ := os.OpenFile(AppConfPath, os.O_RDONLY|os.O_CREATE, 0600); fh != nil {
//		fh.Close()
	//}

	// Load configuration, set app version and log level.
	Cfg, err = goconfig.LoadConfigFile(GlobalConfPath)

	if Cfg == nil {
		//Cfg, err = goconfig.LoadConfigFile(AppConfPath)
		if err != nil {
			fmt.Println("Fail to load configuration file: " + err.Error())
			os.Exit(2)
		}
	} 

	Cfg.BlockMode = false

	// set time zone of wetalk system
	TimeZone = Cfg.MustValue("app", "time_zone", "UTC")
	if _, err := time.LoadLocation(TimeZone); err == nil {
		os.Setenv("TZ", TimeZone)
	} else {
		fmt.Println("Wrong time_zone: " + TimeZone + " " + err.Error())
		os.Exit(2)
	}

	beego.RunMode = Cfg.MustValue("app", "run_mode")
	beego.HttpPort = Cfg.MustInt("app", "http_port")

//	IsProMode = beego.RunMode == "pro"
//	if IsProMode {
//		beego.SetLevel(beego.LevelInfo)
//	}

	// cache system
//	Cache, err = cache.NewCache("memory", `{"interval":360}`)

//	Captcha = captcha.NewCaptcha("/captcha/", Cache)
	//Captcha.FieldIdName = "CaptchaId"
//	Captcha.FieldCaptchaName = "Captcha"

	// session settings
	beego.SessionOn = true
	beego.SessionProvider = Cfg.MustValue("session", "session_provider", "file")
	beego.SessionSavePath = Cfg.MustValue("session", "session_path", "sessions")
	beego.SessionName = Cfg.MustValue("session", "session_name", "Cloud_In_a_Box")
	beego.SessionCookieLifeTime = Cfg.MustInt("session", "session_life_time", 0)
	beego.SessionGCMaxLifetime = Cfg.MustInt64("session", "session_gc_time", 86400)

	beego.EnableXSRF = true
	// xsrf token expire time
	beego.XSRFExpire = 86400 * 365

//	driverName := Cfg.MustValue("orm", "driver_name", "mysql")
//	dataSource := Cfg.MustValue("orm", "data_source", "root:root@/wetalk?charset=utf8&loc=UTC")
//	maxIdle := Cfg.MustInt("orm", "max_idle_conn", 30)
//	maxOpen := Cfg.MustInt("orm", "max_open_conn", 50)

	// set default database
//	err = orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
//	if err != nil {
//		beego.Error(err)
//	}
//	orm.RunCommand()

//	err = orm.RunSyncdb("default", false, false)
//	if err != nil {
//		beego.Error(err)
//	}

	reloadConfig()

//	if SphinxEnabled {
		// for search config
//		SphinxHost = Cfg.MustValue("search", "sphinx_host", "127.0.0.1:9306")
//		SphinxMaxConn = Cfg.MustInt("search", "sphinx_max_conn", 5)
//		orm.RegisterDriver("sphinx", orm.DR_MySQL)
//	}

//	social.DefaultAppUrl = AppUrl

	//settingLocales()
	//settingCompress()

	//configWatcher()

	return Cfg
}

func reloadConfig() {
	AppName = Cfg.MustValue("app", "app_name", "Cloud In a Box")
	beego.AppName = AppName

	AppHost = Cfg.MustValue("app", "app_host", "127.0.0.1:8085")
	AppUrl = Cfg.MustValue("app", "app_url", "http://127.0.0.1:8085/")
	AppLogo = Cfg.MustValue("app", "app_logo", "https://s3-ap-southeast-1.amazonaws.com/megampub/images/logo-megam160x43w.png")
	AvatarURL = Cfg.MustValue("app", "avatar_url")

	EnforceRedirect = Cfg.MustBool("app", "enforce_redirect")

	DateFormat = Cfg.MustValue("app", "date_format")
	DateTimeFormat = Cfg.MustValue("app", "datetime_format")
	DateTimeShortFormat = Cfg.MustValue("app", "datetime_short_format")

	LoginRememberDays = Cfg.MustInt("app", "login_remember_days", 7)
	LoginMaxRetries = Cfg.MustInt("app", "login_max_retries", 5)
	LoginFailedBlocks = Cfg.MustInt("app", "login_failed_blocks", 10)

	CookieRememberName = Cfg.MustValue("app", "cookie_remember_name")
	CookieUserName = Cfg.MustValue("app", "cookie_user_name")
	
}
