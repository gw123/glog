package glog

import (
	"fmt"
	"github.com/pkg/profile"
	"log"
	"os"
	"strings"
	"time"
)

const (
	PModeCpu = "cpu"
	PModeMem = "mem"
	PModeMutex = "mutex"
	PModeThread = "thread"
	PModeTrace = "trace"
)

var gProfile *Profile

type Profile struct {
	Profile  interface{ Stop() }
	mode     string
	period   time.Duration
	stopFlag bool
	rootPath string
}

func init() {
	gProfile = &Profile{}
	SetProfile(ProfileMode(""), ProfilePath(""), ProfilePeriod(0))
}

func ProfileMode(mode string) func(*Profile) {
	return func(p *Profile) {
		if mode == "" {
			mode = "mem"
		}
		p.mode = mode
	}
}

func ProfilePath(rootPath string) func(*Profile) {
	return func(p *Profile) {
		if rootPath == "" {
			arr := strings.Split(os.Args[0], string(os.PathSeparator))
			rootPath = "/tmp/pprof/" + arr[len(arr)-1]
		}
		p.rootPath = rootPath
	}
}

func ProfilePeriod(t time.Duration) func(*Profile) {
	return func(p *Profile) {
		if t < time.Minute {
			t = time.Hour
		}
		p.period = t
	}
}

func SetProfile(options ...func(*Profile)) {
	for _, option := range options {
		option(gProfile)
	}
	return
}

func StartProfile() {
	if gProfile != nil {
		gProfile.start()
	}
}

func StopProfile() {
	if gProfile != nil {
		gProfile.stop()
	}
}

func (p *Profile) stopRecodeProfile() {
	if p.Profile != nil {
		p.Profile.Stop()
		p.Profile = nil
	}
}

func (p *Profile) stop() {
	p.stopFlag = true
	p.stopRecodeProfile()
	fmt.Println("停止记录")
}

func (p *Profile) start() *Profile {
	p.stopFlag = false

	go func() {
		for !p.stopFlag {
			p.stopRecodeProfile()
			//pprof 会把有:的当做url地址 所以文件路径不能有:
			//path := fmt.Sprintf("./logs/pprof/%d_%d:%d", time.Now().Day(), time.Now().Hour(), time.Now().Minute())

			path := fmt.Sprintf(p.rootPath+"/%d_%d_%d", time.Now().Day(), time.Now().Hour(), time.Now().Minute())
			fmt.Println(path)
			err := os.MkdirAll(path, 0755)
			if err != nil {
				log.Print(err)
				return
			}

			switch p.mode {
			case "cpu":
				p.Profile = profile.Start(profile.CPUProfile,
					profile.ProfilePath(path),
					profile.NoShutdownHook,
				)
			case "mem":
				p.Profile = profile.Start(
					profile.MemProfile,
					profile.ProfilePath(path),
					profile.NoShutdownHook,
				)
			case "mutex":
				p.Profile = profile.Start(profile.MutexProfile,
					profile.ProfilePath(path),
					profile.NoShutdownHook,
				)
			case "thread":
				p.Profile = profile.Start(profile.ThreadcreationProfile,
					profile.ProfilePath(path),
					profile.NoShutdownHook, )
			case "trace":
				p.Profile = profile.Start(profile.TraceProfile,
					profile.ProfilePath(path),
					profile.NoShutdownHook, )
			}
			time.Sleep(p.period)
		}
	}()
	return p
}
