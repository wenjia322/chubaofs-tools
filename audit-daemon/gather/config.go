package gather

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

const exclusionDir = "logs"

func parseConfig(configPath string) {
	all, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("read %s has err:[%s]", configPath, err.Error()))
	}

	reg := regexp.MustCompile(`\s+`)
	for _, line := range strings.Split(string(all), "\n") {
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		if strings.TrimSpace(line)[0] == '#' {
			continue
		}

		split := reg.Split(line, -1)
		ipAddr := split[0]
		srcDir := split[1]
		pattern := split[2]
		dstDir := split[3]

		if _, found := workers[ipAddr]; !found {
			workers[ipAddr] = &Worker{
				addr:    ipAddr,
				srcDir:  srcDir,
				dstDir:  dstDir,
				pattern: pattern,
				jobs:    make(map[string]*Job),
			}
		}
		wk := workers[ipAddr]

		var subDirs []string
		if subDirs, err = remoteDirs(ipAddr, srcDir, exclusionDir); err != nil { // exclude dir "logs"
			panic(fmt.Sprintf("list remote dir err: addr[%s], dir[%s], err[%s]", ipAddr, srcDir, err.Error()))
		}

		for _, subDir := range subDirs {
			srcSubDir := path.Join(srcDir, subDir)
			dstSubDir := path.Join(dstDir, subDir)

			wk.createJob(srcSubDir, dstSubDir)
		}

	}
}
