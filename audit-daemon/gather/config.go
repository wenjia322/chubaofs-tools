package gather

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

const exclusionDir = "log"

func parseConfig(configPath string) {
	ipSyncMap = make(map[string]string)
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
				addr: ipAddr,
			}
		}
		wk := workers[ipAddr]

		var subDirs []string
		if subDirs, err = remoteDirs(ipAddr, srcDir, exclusionDir); err != nil { // exclude dir "logs"
			panic(fmt.Sprintf("list remote dir err: addr[%s], dir[%s], err[%s]", ipAddr, srcDir, err.Error()))
		}

		for _, subDir := range subDirs {
			dstSubDir := path.Join(dstDir, subDir)
			wk.jobs = append(wk.jobs, &Job{
				src:     path.Join(srcDir, subDir),
				pattern: pattern,
				dist:    dstSubDir,
			})

			if fi, err := os.Stat(dstSubDir); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(dstSubDir, os.ModePerm); err != nil {
						panic(err)
					}
					if err := os.MkdirAll(path.Join(dstSubDir, "archive"), os.ModePerm); err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
			} else if !fi.IsDir() {
				panic(fmt.Sprintf("%s is not dir", dstSubDir))
			}

			ipSyncMap[dstSubDir] = ipAddr
		}

	}
}
