package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/jsli/cp_release/constant"
	"github.com/jsli/cp_release/release"
	"github.com/jsli/gtbox/pathutil"
	"log"
	"os"
)

const (
	MONITOR_LOG = constant.LOGS_ROOT + "monitor.log"
)

var (
	watcher_map = make(map[string]*fsnotify.Watcher)
	dir_list    = []string{
		constant.HLWB_ROOT,
		constant.HLWB_DSDS_ROOT,
		constant.HLTD_ROOT,
		constant.HLTD_DSDS_ROOT,
		constant.LTG_ROOT,
		constant.LWG_ROOT,
	}

	logOutput *os.File
)

func init() {
	logOutput, err := os.OpenFile(MONITOR_LOG, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logOutput)
	log.SetPrefix("[MONITOR]  ")
}

func main() {
	done := make(chan bool)
	for _, path := range dir_list {
		if exist, err := pathutil.IsExist(path); !exist && err == nil {
			//			pathutil.MkDir(path)
			continue
		}

		err := MonitorDir(path)
		if err != nil {
			log.Printf("Monitor [%s] failed : %s\n", path, err)
			panic(err)
		}
	}
	<-done
}

func MonitorDir(path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	watcher_map[path] = watcher

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				PreprocessEvent(ev)
			case err := <-watcher.Error:
				panic(err)
			}
		}
	}()

	log.Printf("Begin monitor [%s]\n", path)
	err = watcher.Watch(path)
	if err != nil {
		return err
	}

	return nil
}

func PreprocessEvent(event *fsnotify.FileEvent) {
	path := event.Name
	//ignore 5 root dir
	for _, p := range dir_list {
		if p == path {
			return
		}
	}

	if event.IsDelete() {
		ProcessDeleteEvent(event)
	} else {
		if isDir, err := pathutil.IsDir(path); isDir && err == nil {
			if event.IsCreate() {
				ProcessCreateEvent(event)
			} else if event.IsModify() {
				ProcessModifyEvent(event)
			} else if event.IsRename() {
				ProcessRenameEvent(event)
			}
		}
	}
}

func ProcessDeleteEvent(event *fsnotify.FileEvent) {
	log.Printf("Delete [%s]\n", event.Name)

	dal, err := release.NewDal()
	if err != nil {
		log.Printf("New DAL failed: %s\n", err)
		return
	}
	defer dal.Close()

	cp := getCpByRelPath(event.Name, dal)
	if cp != nil {
		_, err := cp.Delete(dal)
		if err != nil {
			log.Printf("Delete failed: %s\n", err)
		} else {
			log.Printf("Delete success: %s\n", cp)
			log.Printf("CP release deleted, delete arbi&grbi together : %s", cp)
			release.DeleteArbiByCpId(dal, cp.Id)
			release.DeleteGrbiByCpId(dal, cp.Id)
		}
	}
}

func ProcessModifyEvent(event *fsnotify.FileEvent) {
	log.Printf("Modify [%s]\n", event.Name)

	dal, err := release.NewDal()
	if err != nil {
		log.Printf("New DAL failed: %s\n", err)
		return
	}
	defer dal.Close()

	cp := getCpByRelPath(event.Name, dal)
	if cp != nil {
		if cp.Flag == constant.AVAILABLE_FLAG {
			log.Printf("CP release modified, delete arbi&grbi for updating in scanner : %s", cp)
			release.DeleteArbiByCpId(dal, cp.Id)
			release.DeleteGrbiByCpId(dal, cp.Id)
		}
	}
}

func ProcessCreateEvent(event *fsnotify.FileEvent) {
	log.Printf("Create [%s]\n", event.Name)
}

func ProcessRenameEvent(event *fsnotify.FileEvent) {
	fmt.Printf("rename : %s\n", event.Name)
}

func getCpByRelPath(full_path string, dal *release.Dal) *release.CpRelease {
	parent_path := pathutil.ParentPath(full_path)
	mode := pathutil.BaseName(parent_path)
	rel_path := full_path[constant.PATH_PREFIX_LEN:]
	cp, err := release.FindCpReleaseByPath(dal, rel_path)
	if err != nil {
		log.Printf("Find cp failed by [%s]: %s\n", rel_path, err)
		return nil
	}
	return cp
}
