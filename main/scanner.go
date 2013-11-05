package main

import (
	"fmt"
	"github.com/jsli/cp_release/constant"
	"github.com/jsli/cp_release/policy"
	"github.com/jsli/cp_release/release"
	"github.com/jsli/gtbox/pathutil"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	SCANNER_LOG      = constant.LOGS_ROOT + "scanner.log"
	FLAG_SCAN_DETAIL = true
	SCAN_INTERNAL    = 100
)

var (
	dir_list = []string{
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
	logOutput, err := os.OpenFile(SCANNER_LOG, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logOutput)
	log.SetPrefix("[SCANNER]  ")
}

func main() {
	dal, err := release.NewDal()
	if err != nil {
		log.Printf("New DAL failed: %s\n", err)
		panic(err)
	}
	defer dal.Close()

	counter := 0
	for {
		counter = counter + 1
		log.Printf("scan %d times ----------------------------begin\n", counter)
		for _, path := range dir_list {
			if exist, err := pathutil.IsExist(path); !exist && err == nil {
				//			pathutil.MkDir(path)
				continue
			}

			mode := constant.PATH_TO_MODE[path]
			sim := constant.MODE_TO_SIM[mode]

			err := ScanDir(dal, path, mode, sim)
			if err != nil {
				log.Printf("Scan PANIC [%s] failed : %s\n", path, err)
				panic(err)
			}
		}
		log.Printf("scan %d times ----------------------------end\n", counter)
		time.Sleep(SCAN_INTERNAL * time.Second)
	}
}

func ScanDir(dal *release.Dal, path string, mode string, sim string) error {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("read dir failed : %s\n", err)
		return err
	}

	err = CheckRecord(dal, mode)
	if err != nil {
		log.Printf("Check CP release failed : %s\n", err)
	}

	for _, info := range fileInfos {
		if info.Mode().IsDir() {
			err := ProcessDir(info, dal, path, mode, sim, false)
			if err != nil {
				log.Printf("ProcessDir failed: %s", err)
			}
		}
	}
	return nil
}

func CheckRecord(dal *release.Dal, mode string) error {
	query := fmt.Sprintf("SELECT * FROM %s where mode='%s' and flag=%d", constant.TABLE_CP, mode, constant.AVAILABLE_FLAG)
	cp_list, err := release.FindCpReleaseList(dal, query)
	if err != nil {
		return err
	}

	for _, cp := range cp_list {
		full_path := fmt.Sprintf("%s%s", constant.CP_SERVER_MIRROR_ROOT, cp.RelPath)
		exist, err := pathutil.IsExist(full_path)
		if err != nil {
			continue
		}

		if !exist {
			log.Printf("CheckDir db and fs unmatched, delete db record: %s", cp)
			cp.Delete(dal)
		} else {
			ProcessDetail(cp, dal)
		}
	}
	return nil
}

func ProcessDir(info os.FileInfo, dal *release.Dal, path string, mode string, sim string, force bool) error {
	version := policy.ExtractVersion(info.Name())
	if version == "" { //illegal version fmt, ignore
		return fmt.Errorf("Illegal version format : %s", info.Name())
	}
	rel_path := fmt.Sprintf("%s/%s", path, info.Name())[constant.PATH_PREFIX_LEN:]
	cp, err := release.FindCpReleaseByPath(dal, rel_path)
	if err != nil {
		return err
	}
	if cp != nil {
		if !force {
			return fmt.Errorf("Existed CP release : %s", cp)
		} else {
			log.Printf("Existed CP release, delete arbi&grbi for force updating : %s", cp)
			release.DeleteArbiByCpId(dal, cp.Id)
			release.DeleteGrbiByCpId(dal, cp.Id)
			release.DeleteRficByCpId(dal, cp.Id)
		}
	} else {
		cp = &release.CpRelease{}
		cp.Mode = mode
		cp.Sim = sim
		cp.Version = version
		cp.VersionScalar = policy.QuantitateVersion(version)
		cp.LastModifyTs = time.Now().Unix()
		cp.Flag = constant.AVAILABLE_FLAG
		cp.RelPath = rel_path

		slice := strings.Split(rel_path, "/")
		cp.Prefix = strings.TrimSuffix(slice[2], version)

		log.Printf("Find new CP release : %s\n", cp)
		id, err := cp.Save(dal)
		if err != nil {
			cp.Id = -1
			log.Printf("Save CP release failed: %s\n", err)
		} else {
			cp.Id = id
			log.Printf("Save CP release success: %d | %s\n", id, cp)
		}
	}

	//find detail information
	if cp.Id > 0 && FLAG_SCAN_DETAIL {
		ProcessDetail(cp, dal)
	} else {
		return fmt.Errorf("Neither finding or saving CP release success! in [%s]", info.Name())
	}

	return nil
}

func ProcessDetail(cp *release.CpRelease, dal *release.Dal) error {
	err := ProcessArbi(cp, dal)
	if err != nil {
		log.Printf("ProcessArbi failed: %s", err)
	}

	err = ProcessGrbi(cp, dal)
	if err != nil {
		log.Printf("ProcessGrbi failed: %s", err)
	}

	err = ProcessRfic(cp, dal)
	if err != nil {
		log.Printf("ProcessRfic failed: %s", err)
	}

	return nil
}

func ProcessArbi(cp *release.CpRelease, dal *release.Dal) error {
	arbi_list, err := policy.FindArbi(cp.RelPath, cp.Mode)
	if err != nil {
		return err
	} else {
		for _, arbi_path := range arbi_list {
			arbi_rel_path := arbi_path[constant.PATH_PREFIX_LEN:]
			arbi, err := release.FindArbiByPath(dal, arbi_rel_path)
			if err == nil && arbi != nil {
				log.Printf("Existed arbi : %s\n", arbi)
				//id unmatched, delete itself
				if arbi.CpId != cp.Id {
					log.Printf("Id unmatched cp[%d] <--> arbi[%d] : delete", cp.Id, arbi.CpId)
					arbi.Delete(dal)
				} else {
					continue
				}
			}
			arbi = &release.Arbi{}
			arbi.CpId = cp.Id
			arbi.Flag = constant.AVAILABLE_FLAG
			arbi.LastModifyTs = time.Now().Unix()
			arbi.RelPath = arbi_rel_path
			log.Printf("Found arbi in [%s] : %s\n", cp.RelPath, arbi)
			id, err := arbi.Save(dal)
			if err != nil {
				log.Printf("Save ARBI failed: %s\n", err)
			} else {
				log.Printf("Save ARBI success: %d | %s\n", id, arbi)
			}
		}
	}
	return nil
}

func ProcessRfic(cp *release.CpRelease, dal *release.Dal) error {
	rfic_list, err := policy.FindRfic(cp.RelPath, cp.Mode)
	if err != nil {
		return err
	} else {
		for _, rfic_path := range rfic_list {
			rfic_rel_path := rfic_path[constant.PATH_PREFIX_LEN:]
			rfic, err := release.FindRficByPath(dal, rfic_rel_path)
			if err == nil && rfic != nil {
				log.Printf("Existed rfic : %s\n", rfic)
				//id unmatched, delete itself
				if rfic.CpId != cp.Id {
					log.Printf("Id unmatched cp[%d] <--> rfic[%d] : delete", cp.Id, rfic.CpId)
					rfic.Delete(dal)
				} else {
					continue
				}
			}
			rfic = &release.Rfic{}
			rfic.CpId = cp.Id
			rfic.Flag = constant.AVAILABLE_FLAG
			rfic.LastModifyTs = time.Now().Unix()
			rfic.RelPath = rfic_rel_path
			log.Printf("Found rfic in [%s] : %s\n", cp.RelPath, rfic)
			id, err := rfic.Save(dal)
			if err != nil {
				log.Printf("Save RFIC failed: %s\n", err)
			} else {
				log.Printf("Save RFIC success: %d | %s\n", id, rfic)
			}
		}
	}
	return nil
}

func ProcessGrbi(cp *release.CpRelease, dal *release.Dal) error {
	grbi_list, err := policy.FindGrbi(cp.RelPath, cp.Mode)
	if err != nil {
		return err
	} else {
		for _, grbi_path := range grbi_list {
			grbi_rel_path := grbi_path[constant.PATH_PREFIX_LEN:]
			grbi, err := release.FindGrbiByPath(dal, grbi_rel_path)
			if err == nil && grbi != nil {
				log.Printf("Existed grbi : %s\n", grbi)
				//id unmatched, delete itself
				if grbi.CpId != cp.Id {
					log.Printf("Id unmatched cp[%d] <--> grbi[%d] : delete", cp.Id, grbi.CpId)
					grbi.Delete(dal)
				} else {
					continue
				}
			}
			grbi = &release.Grbi{}
			grbi.CpId = cp.Id
			grbi.Flag = constant.AVAILABLE_FLAG
			grbi.LastModifyTs = time.Now().Unix()
			grbi.RelPath = grbi_rel_path
			log.Printf("Found grbi in [%s] : %s\n", cp.RelPath, grbi)
			id, err := grbi.Save(dal)
			if err != nil {
				log.Printf("Save GRBI failed: %s\n", err)
			} else {
				log.Printf("Save GRBI success: %d | %s\n", id, grbi)
			}
		}
	}
	return nil
}
