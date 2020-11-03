// Package sysinfo system information
package sysinfo

import (
	"log"

	elastic "github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
	info "github.com/zcalusic/sysinfo"
)

// XSysInfoS data
type XSysInfoS struct {
	Info   infoX
	Memory memoryX
	types.GoInfo
}

type memoryX struct {
	Host    *types.HostMemoryInfo
	Process types.MemoryInfo
}
type infoX struct {
	Info       info.SysInfo
	Host       types.HostInfo
	Process    types.ProcessInfo
	Processess []types.ProcessInfo
}

// SysInfo info
func SysInfo() (*XSysInfoS, error) {
	s := new(XSysInfoS)
	var si info.SysInfo
	si.GetSysInfo()

	host, err := elastic.Host()
	if err != nil {
		return s, err
	}
	process, err := elastic.Self()
	if err != nil {
		return s, err
	}

	mem, err := process.Memory()
	if err != nil {
		return s, err
	}
	s.Memory.Process = mem

	memH, err := host.Memory()
	if err != nil {
		return s, err
	}
	s.Memory.Host = memH

	s.Info.Info = si
	s.Info.Host = host.Info()
	s.Info.Process, err = process.Info()
	if err != nil {
		log.Println(err)
	}

	s.GoInfo = elastic.Go()
	cc, err := elastic.Processes()
	if err != nil {
		return s, err
	}
	for _, p := range cc {
		dd, err1 := p.Info()
		u, err := p.User()
		if err1 != nil || err != nil {
			log.Println(err)
			continue
		}
		log.Println(u)
		s.Info.Processess = append(s.Info.Processess, dd)
	}
	return s, nil
}
