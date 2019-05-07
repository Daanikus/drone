package main

import (
	"io"

	"github.com/morya/utils/log"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	// "compress/gzip"
)

type TarGzip struct {
	Filename string
}

func copyFolder(mem, src billy.Filesystem) {
	root := src.Root()
	log.Infof("root %v", root)
	dirs, _ := src.ReadDir(root)
	log.Infof("root dir %v, files = ", dirs)
	for _, d := range dirs {
		if d.IsDir() {
			target, err := src.Chroot(d.Name())
			if err != nil {
				log.InfoError(err)
				continue
			}
			mem2, _ := mem.Chroot(d.Name())
			log.Infof("copy folder %s", d.Name())
			copyFolder(mem2, target)
		} else {
			s, _ := src.Open(d.Name())
			n, _ := mem.Create(d.Name())
			io.Copy(n, s)
			s.Close()
			n.Close()

			log.Infof("copy file %s", d.Name())
		}
	}
}

func (p *TarGzip) Clone(src billy.Filesystem) {
	mem := memfs.New()
	log.Info("do copy")
	copyFolder(mem, src)
}
