package main

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/mattn/go-colorable"

	"github.com/labstack/echo"
	glog "github.com/labstack/gommon/log"
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

func (p *Drone) checkRepo(name string, prj *ProjectConfig) {
	repo, err := newRepo(prj.GitURL, prj.LocalPath, prj.GitKey)
	if err != nil {
		log.InfoError(err, "check repository failed")
		return
	}

	hasUpdate, err := repo.HasUpdate()
	if err != nil {
		log.InfoErrorf(err, "project [%v] pull new changes failed", name)
		return
	}
	if !hasUpdate {
		log.Infof("project [%v] checked, no update.", name)
		return
	}

	serverCfg := config.Servers[prj.Server]
	sshConn, err := newSshConn(serverCfg)
	if err != nil {
		log.InfoErrorf(err, "[%v] connect server [%v] failed", name, serverCfg.IP)
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

	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			for prjName, prjCfg := range config.Projects {
				d.checkRepo(prjName, prjCfg)
			}
		}
	}
}

func (p *Drone) Run(listenAddr string) {
	defer p.wg.Done()

	p.echo = echo.New()
	p.echo.HideBanner = true

	p.echo.Logger = glog.New("echo")
	p.echo.Logger.SetOutput(colorable.NewColorableStderr())

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
