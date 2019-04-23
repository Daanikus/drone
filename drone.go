package main

import (
	"context"
	"github.com/morya/utils/log"
	"sync"

	"github.com/labstack/echo"
)

type Drone struct {
	wg   sync.WaitGroup
	lock sync.Mutex

	echo *echo.Echo
}

func newDrone() *Drone {
	d := &Drone{}
	d.wg.Add(1)
	return d
}

func (p *Drone) onHook(c echo.Context) error {
	return nil
}

func (p *Drone) Run(listenAddr string) {
	defer p.wg.Done()

	p.echo = echo.New()
	p.echo.GET("/webhook", p.onHook)
	p.echo.POST("/webhook", p.onHook)

	err := p.echo.Start(listenAddr)
	if err != nil {
		log.InfoError(err)
	}
}

func (p *Drone) Stop() {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.echo == nil {
		return
	}
	p.echo.Shutdown(context.Background())
	p.wg.Wait()
	p.echo = nil
}
