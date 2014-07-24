package main

import (
	"github.com/megamsys/libgo/amqp"
//	"github.com/megamsys/cloudinabox/app"
	"log"
	"os"
	"os/signal"
	//"strings"
	"regexp"
	"sync"
	"syscall"
	"time"
)

const (
	// queue actions
	runningApp = "running"
	startApp   = "start"
	stopApp    = "stop"
	buildApp   = "build"
	restartApp = "restart"
	addonApp   = "addon"
	queueName  = "gulpd-app"
)

var (
	qfactory      amqp.QFactory
	_queue        amqp.Q
	_handler      amqp.Handler
	o             sync.Once
	signalChannel chan<- os.Signal
	nameRegexp    = regexp.MustCompile(`^[a-z][a-z0-9-]{0,62}$`)
)

func RunServer(dry bool) {
	log.Printf("Gulpd starting at %s", time.Now())
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)
	handler().Start()
	log.Printf("Gulpd at your service.")
	<-signalChannel
	log.Println("Gulpd killed |_|.")
}

func StopServer(bark bool) {
	log.Printf("Gulpd stopping at %s", time.Now())
	handler().Stop()
	close(signalChannel)
	log.Println("Gulpd finished |-|.")
}

func setQueue() {
	var err error
	qfactory, err = amqp.Factory()

	if err != nil {
		log.Fatalf("Failed to get the queue instance: %s", err)
	}
	_handler, err = qfactory.Handler(handle, queueName)
	if err != nil {
		log.Fatalf("Failed to create the queue handler: %s", err)
	}

	_queue, err = qfactory.Get(queueName)

	if err != nil {
		log.Fatalf("Failed to get the queue instance: %s", err)
	}
}

func aqueue() amqp.Q {
	o.Do(setQueue)
	return _queue
}

func handler() amqp.Handler {
	o.Do(setQueue)
	return _handler
}


// handle is the function called by the queue handler on each message.
//This is getting bulky. We need to move it out to apps and others.
func handle(msg *amqp.Message) {
	log.Printf("Handling message %v", msg)

	/*switch strings.ToLower(msg.Action) {
	case "ganeti install":
		if len(msg.Args) < 1 {
			log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
		}
		//for now comment it out.
	  ap := app.App{Email: msg.Email, ApiKey: msg.ApiKey, InstallPackage: msg.InstallPackage, NeedMegam: msg.NeedMegam, ClusterName: msg.ClusterName, NodeIp: msg.NodeIp, NodeName: msg.NodeName, Action: msg.Action}


		if err := app.GanetiInstall(&ap); err != nil {
			log.Printf("Error handling %q: Installation process failed:\n%s.", msg.Action, err)
			return
		}

		msg.Delete()
		break
	case "opennebula install":
		if len(msg.Args) < 1 {
			log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
		}
		//stick the id from msg.
		ap := app.App{Email: msg.Email, ApiKey: msg.ApiKey, InstallPackage: msg.InstallPackage, NeedMegam: msg.NeedMegam, ClusterName: msg.ClusterName, NodeIp: msg.NodeIp, NodeName: msg.NodeName, Action: msg.Action}

		if err := app.OpennebulaInstall(&ap); err != nil {
			log.Printf("Error handling %q: Installation process failed:\n%s.", msg.Action, err)
			return
		}

		msg.Delete()
		break
	case "remove":
		if len(msg.Args) < 1 {
			log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
		}
		//stick the id from msg.
		ap := app.App{Email: "myapp", ApiKey: "RIPAB"}

		if err := app.Remove(&ap); err != nil {
			log.Printf("Error handling %q: Remove process failed:\n%s.", msg.Action, err)
			return
		}
		msg.Delete()
		break

	default:
		log.Printf("Error handling %q: invalid action.", msg.Action)
		msg.Delete()
	}
	*/
}
