package gather

import (
	"os"
	"path"
	"sync"
	"time"

	"github.com/chubaofs/chubaofs-tools/audit-daemon/util"
)

var ipSyncMap sync.Map // key: The path to store the synchronization file, value: ip of machine which did the file come from

func getNodeIP(dir string) string {
	addr, _ := ipSyncMap.Load(dir)
	return addr.(string)
}

var workers = make(map[string]*Worker)

type Worker struct {
	addr    string
	srcDir  string
	dstDir  string
	pattern string
	jobs    map[string]*Job
}

type Job struct {
	src     string
	dist    string
	pattern string
}

func toWork(w *Worker) {

	for {
		if util.Stop {
			break
		}
		w.updateJobs()
		for _, job := range w.jobs {
			w.toJob(job)
		}

		time.Sleep(30 * time.Second)
	}
}

// Applicable to each machine only one directory needs to be synchronized
func (w *Worker) updateJobs() {
	srcDir := w.srcDir
	dstDir := w.dstDir

	var subDirs []string
	var err error
	if subDirs, err = remoteDirs(w.addr, srcDir, exclusionDir); err != nil { // exclude dir "logs"
		util.LOG.Errorf("list remote dir err: addr[%s], dir[%s], err[%s]", w.addr, srcDir, err.Error())
		return
	}

	for _, subDir := range subDirs {
		srcSubDir := path.Join(srcDir, subDir)
		dstSubDir := path.Join(dstDir, subDir)
		if _, exist := w.jobs[srcSubDir]; exist {
			continue
		}
		w.createJob(srcSubDir, dstSubDir)
	}
}

func (w *Worker) createJob(srcSubDir, dstSubDir string) {
	util.LOG.Debugf("create new job: addr[%v], src[%v], dst[%v], pattern[%v]", w.addr, srcSubDir, dstSubDir, w.pattern)
	if fi, err := os.Stat(dstSubDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dstSubDir, os.ModePerm); err != nil {
				util.LOG.Errorf("make dir err: dir[%s], err[%s]", dstSubDir, err.Error())
				return
			}
			if err := os.MkdirAll(path.Join(dstSubDir, "archive"), os.ModePerm); err != nil {
				util.LOG.Errorf("make dir err: dir[%s], err[%s]", path.Join(dstSubDir, "archive"), err.Error())
				return
			}
		} else {
			util.LOG.Errorf("stat dir err: dir[%s], err[%s]", dstSubDir, err.Error())
			return
		}
	} else if !fi.IsDir() {
		util.LOG.Errorf("make dir err: [%s] is not dir", dstSubDir)
		return
	}

	w.jobs[srcSubDir] = &Job{
		src:     srcSubDir,
		pattern: w.pattern,
		dist:    dstSubDir,
	}
	ipSyncMap.Store(dstSubDir, w.addr)
}
