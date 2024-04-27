package service

import (
	filemonitor "fileguard/internal/file"
	"github.com/kardianos/service"
	"log"
	"os"
)

type program struct {
	FolderPath string
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	filemonitor.WatchFolder(p.FolderPath)
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func CreateWindowsService(folderPath string) {
	svcConfig := &service.Config{
		Name:        "FileGuard",
		DisplayName: "FileGuard",
		Description: "Description of my service.",
	}

	prg := &program{
		FolderPath: folderPath,
	}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the program is run in service mode (i.e., installed as a service)
	// If not, install and start the service.
	if len(os.Args) > 1 && os.Args[1] == "install" {
		if err := svc.Install(); err != nil {
			log.Fatal("Failed to install service:", err)
		}
		//		log.Println("Service installed successfully.")
		//		return
		//	}
		//
		//	// If run in service mode, directly call svc.Run()
		if err := svc.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
