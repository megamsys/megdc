package main

import (
	"github.com/megamsys/cloudinabox/amqp"
	"github.com/megamsys/cloudinabox/app"
	"log"
	"os"
	"os/signal"
	"strings"
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

	switch strings.ToLower(msg.Action) {
	case restartApp:
		if len(msg.Args) < 1 {
			log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
		}
		//stick the id from msg.
		ap := app.App{Name: "myapp", Id: "RIPAB"}

		if err := ap.Get(msg.Id); err != nil {
			log.Printf("Error handling %q: Riak didn't cooperate:\n%s.", msg.Action, err)
			return
		}
				
		log.Printf("Handling message %#v", ap.GetAppReqs())
		err := app.StopApp(&ap)
		if err != nil {
			log.Printf("Error handling %q. App failed to stop:\n%s.", msg.Action, err)
			return
		}
		
		err = app.StartApp(&ap)
		if err != nil {
			log.Printf("Error handling %q. App failed to start:\n%s.", msg.Action, err)
			return
		}

		msg.Delete()
		break
	case startApp:
		if len(msg.Args) < 1 {
			log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
		}
		//stick the id from msg.
		ap := app.App{Name: "myapp", Id: "RIPAB"}

		if err := ap.Get(msg.Id); err != nil {
			log.Printf("Error handling %q: Riak didn't cooperate:\n%s.", msg.Action, err)
			return
		}
		log.Printf("Handling message %#v", ap.GetAppReqs())
		err := app.StartApp(&ap)
		if err != nil {
			log.Printf("Error handling %q. App failed to start:\n%s.", msg.Action, err)
			return
		}
		msg.Delete()
		break
	case stopApp:
		if len(msg.Args) < 1 {
			log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
		}
		//stick the id from msg.
		ap := app.App{Name: "myapp", Id: "RIPAB"}

		if err := ap.Get(msg.Id); err != nil {
			log.Printf("Error handling %q: Riak didn't cooperate:\n%s.", msg.Action, err)
			return
		}
		log.Printf("Handling message %#v", ap.GetAppReqs())
		err := app.StopApp(&ap)
		if err != nil {
			log.Printf("Error handling %q. App failed to stop:\n%s.", msg.Action, err)
			return
		}

		msg.Delete()
		break	
	case buildApp:
		if len(msg.Args) < 1 {
			log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
		}
		//stick the id from msg.
		ap := app.App{Name: "myapp", Id: "RIPAB"}

		if err := ap.Get(msg.Id); err != nil {
			log.Printf("Error handling %q: Riak didn't cooperate:\n%s.", msg.Action, err)
			return
		}
		log.Printf("Handling message %#v", ap.GetAppReqs())
		err := app.BuildApp(&ap)
		if err != nil {
			log.Printf("Error handling %q. App failed to build:\n%s.", msg.Action, err)
			return
		}

		msg.Delete()
		break
	case runningApp:
		if len(msg.Args) < 1 {
			log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
		}
		//stick the i
		ap := app.App{Name: msg.Id, Id: "RIPAB"}

		log.Printf("Handling message %#v", ap.Name)
		err := app.LaunchedApp(&ap)
		if err != nil {
			log.Printf("Error handling %q. App failed to launch:\n%s.", msg.Action, err)
			return
		}

		msg.Delete()
		break	
    case addonApp: 
       if len(msg.Args) < 1 {
			log.Printf("Error handling %q: this action requires at least 1 argument.", msg.Action)
		}
		ap := app.App{Name: "myapp", Id: "RIPAB", Type: "addon"}
		
		if err := ap.Get(msg.Id); err != nil {
			log.Printf("Error handling %q: Riak didn't cooperate:\n%s.", msg.Action, err)
			return
		}
		
		log.Printf("Handling message %#v", ap.GetAppConf())
		err := app.AddonApp(&ap)
		if err != nil {
			log.Printf("Error handling %q. Addon failed to App:\n%s.", msg.Action, err)
			return
		}
		msg.Delete()
		break	
		
	default:
		log.Printf("Error handling %q: invalid action.", msg.Action)
		msg.Delete()
	}
}
