package policy

import (
	"fmt"
	"github.com/jsli/cp_release/constant"
	"io/ioutil"
	"strings"
)

func FindArbi(rel_path string, mode string) ([]string, error) {
	//	arbi_list := make([]string, 0, 5)
	full_path := fmt.Sprintf("%s%s", constant.CP_SERVER_MIRROR_ROOT, rel_path)
	arbi_list, err := doFindArbi(full_path)
	if err != nil {
		return nil, err
	}
	return arbi_list, nil
}

func doFindArbi(path string) ([]string, error) {
	arbi_list := make([]string, 0, 5)
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, info := range fileInfos {
		if info.Mode().IsDir() {
			list, err := doFindArbi(fmt.Sprintf("%s/%s", path, info.Name()))
			if err != nil {
				return nil, err
			}
			arbi_list = append(arbi_list, list...)
		} else if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".bin") {
			if !strings.HasSuffix(info.Name(), "Flash.bin") && !strings.Contains(path, "/RFIC/") {
				bin_path := fmt.Sprintf("%s/%s", path, info.Name())
				//				fmt.Printf("append : %s\n", bin_path)
				arbi_list = append(arbi_list, bin_path)
			}
		}
	}
	return arbi_list, nil
}

func FindRfic(rel_path string, mode string) ([]string, error) {
	//	arbi_list := make([]string, 0, 5)
	full_path := fmt.Sprintf("%s%s", constant.CP_SERVER_MIRROR_ROOT, rel_path)
	rfic_list, err := doFindRfic(full_path)
	if err != nil {
		return nil, err
	}
	return rfic_list, nil
}

func doFindRfic(path string) ([]string, error) {
	rfic_list := make([]string, 0, 5)
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, info := range fileInfos {
		if info.Mode().IsDir() {
			list, err := doFindRfic(fmt.Sprintf("%s/%s", path, info.Name()))
			if err != nil {
				return nil, err
			}
			rfic_list = append(rfic_list, list...)
		} else if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".bin") {
			if strings.Contains(path, "/RFIC/") {
				bin_path := fmt.Sprintf("%s/%s", path, info.Name())
				//				fmt.Printf("append : %s\n", bin_path)
				rfic_list = append(rfic_list, bin_path)
			}
		}
	}
	return rfic_list, nil
}

func FindGrbi(rel_path string, mode string) ([]string, error) {
	//	arbi_list := make([]string, 0, 5)
	full_path := fmt.Sprintf("%s%s", constant.CP_SERVER_MIRROR_ROOT, rel_path)
	arbi_list, err := doFindGrbi(full_path)
	if err != nil {
		return nil, err
	}
	return arbi_list, nil
}

func doFindGrbi(path string) ([]string, error) {
	arbi_list := make([]string, 0, 5)
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, info := range fileInfos {
		if info.Mode().IsDir() {
			list, err := doFindGrbi(fmt.Sprintf("%s/%s", path, info.Name()))
			if err != nil {
				return nil, err
			}
			arbi_list = append(arbi_list, list...)
		} else if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), "Flash.bin") {
			bin_path := fmt.Sprintf("%s/%s", path, info.Name())
			//			fmt.Printf("append : %s\n", bin_path)
			arbi_list = append(arbi_list, bin_path)
		}
	}
	return arbi_list, nil
}
