package generatebuildinfo

import (
	"os/exec"
	"solopg/app/cmdversion/types/buildinfo"
	"solopg/config"
	"strings"

	plumprimitive "github.com/cyril-fl/plum/primitive"
)

var (
	Version = ""
	Commit  = ""
	Date    = ""
)

type run struct {
	plumprimitive.Process
	cache
}

type cache struct {
	version string
	commit  string
	date    string
}

type shellcmd = string

const (
	GET_COMMIT  = "git rev-parse --short HEAD 2>/dev/null || echo none"
	GET_DATE    = "date -u +%Y-%m-%dT%H:%M:%SZ"
	GET_VERSION = "git describe --tags --always --dirty 2>/dev/null || echo dev"
)

func Process() *run {
	return &run{}
}
func (p *run) Run() {
	executeByEnv(executablenv{
		devfunc:  p.makeDevEnv,
		prodfunc: p.makeProdEnv,
	})
}

// Getters & Setters
func (p *run) GetResult() buildinfo.BuildInfo {
	infos := buildinfo.New(buildinfo.InfoTemplate{
		Version: p.cache.version,
		Commit:  p.cache.commit,
		Date:    p.cache.date,
	})

	return infos
}

// Methods
func (p *run) makeDevEnv() {
	p.makeVersion()
	p.makeCommit()
	p.makeDate()
}

func (p *run) makeProdEnv() {
	p.cache.version = Version
	p.cache.commit = Commit
	p.cache.date = Date
}

func (p *run) makeCommit() {
	p.cache.commit = p.executeCmd(GET_COMMIT)
}

func (p *run) makeDate() {
	p.cache.date = p.executeCmd(GET_DATE)
}

func (p *run) makeVersion() {
	p.cache.version = p.executeCmd(GET_VERSION)
}

func (p *run) executeCmd(cmd shellcmd) string {
	c := exec.Command("sh", "-c", cmd)
	out, err := c.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// Helpers
type executablenv struct {
	devfunc  func()
	prodfunc func()
}

func executeByEnv(e executablenv) {
	if config.Current.IsDev() {
		e.devfunc()
	} else {
		e.prodfunc()
	}
}
