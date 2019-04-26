package main

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/morya/utils/log"
)

type M map[string]interface{}

type Drone struct {
	wg   sync.WaitGroup
	lock sync.Mutex

	ctx      context.Context
	cancel   context.CancelFunc
	echo     *echo.Echo
	renderer echo.Renderer

	chanEvent chan *Event
}

func newDrone() *Drone {
	ctx, cancel := context.WithCancel(context.Background())
	d := &Drone{
		ctx:       ctx,
		cancel:    cancel,
		chanEvent: make(chan *Event, 10),
	}
	d.wg.Add(1)
	return d
}

func (p *Drone) onHook(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (p *Drone) onProjectList(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (p *Drone) onProjectLog(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (p *Drone) onProjectStart(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (p *Drone) build(evt *Event) {
	return
}

func (p *Drone) builder(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return

		case evt := <-p.chanEvent:
			p.build(evt)
		}
	}
}

func (p *Drone) checkRepo(prj *Project) {
	repo := newRepo(prj.LocalPath)
	if !repo.HasUpdate() {
		log.Debug("repo has no update")
		return
	}

	repo.Pull()

	serverCfg := config.Servers[prj.Server]
	sshConn, err := newSshConn(serverCfg)
	if err != nil {
		log.InfoErrorf(err, "connect server failed")
		return
	}

	for _, task := range prj.PreTasks {
		var output = &bytes.Buffer{}
		sshConn.Exec(output, task.Command, task.Env)
		log.Info(output.String())
	}
}

func (d *Drone) warden(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			for _, prj := range config.Projects {
				d.checkRepo(prj)
			}
		}
	}
}

func (p *Drone) Run(listenAddr string) {
	defer p.wg.Done()

	p.echo = echo.New()

	p.echo.GET("/webhook", p.onHook)
	p.echo.POST("/webhook", p.onHook)
	apiGroup := p.echo.Group("/api")
	{
		apiGroup.GET("/project/list", p.onProjectList)
		apiGroup.GET("/p/:prj/log", p.onProjectLog)
		apiGroup.POST("/p/:prj/start", p.onProjectStart)
	}
	p.echo.Static("/static", "static")
	p.echo.File("/", "static/index.html")

	go p.warden(p.ctx, p.cancel)
	go p.builder(p.ctx, p.cancel)

	defer p.cancel()

	err := p.echo.Start(listenAddr)
	if err != nil {
		log.InfoError(err)
	}
}

func (p *Drone) Stop() {
	p.echo.Shutdown(context.Background())
	p.wg.Wait()

}
