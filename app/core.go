package core

import (
	"log"
	"os"
	"os/signal"
	"prj/app/config"
	"prj/app/renderer"
	"syscall"
)

func App() {
	subscribeOSSignals()
	conf := config.GetConfig()

	renderer.Init()

	view := renderer.View{
		Projects: conf.Projects,
	}

	c := make(chan int)
	go func() {
		select {
		case selectedProject := <-renderer.GetSelectedProject():
			log.Printf("Project selected: %s\n", selectedProject)
			c <- 1
		}
	}()

	renderer.RenderView(view)

	<-c
	os.Exit(0)
}

func subscribeOSSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		os.Exit(0)
	}()
}
