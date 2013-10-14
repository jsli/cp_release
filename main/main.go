package main

import (
	"fmt"
	"github.com/jsli/cp_release/policy"
	"github.com/jsli/cp_release/release"
	"github.com/jsli/cp_release/constant"
	"github.com/jsli/gtbox/pathutil"
)

func main() {
	full_path := "/home/manson/OTA/CP_RELEASE/HL/HLWB/HLWB_CP_1.63.0001"
	parent_path := pathutil.ParentPath(full_path)
	mode := pathutil.BaseName(parent_path)
	fmt.Println(parent_path)
	fmt.Println(mode)
	rel_path := full_path[constant.MODE_TO_PREFIX_LEN[mode]:]
	fmt.Println(rel_path)
}

func testExtractVersion() {
	policy.QuantitateVersion(policy.ExtractVersion("HLTD_DSDS_CP_3.26.P10_Test"))
	policy.QuantitateVersion(policy.ExtractVersion("HLWB_CP_1.33.0040"))
	policy.QuantitateVersion(policy.ExtractVersion("HLWB_CP_1.57.002"))
	policy.QuantitateVersion(policy.ExtractVersion("HLWB_CP_1.50.M12"))
	policy.QuantitateVersion(policy.ExtractVersion("HLWB_CP_1.50.L72"))
	policy.QuantitateVersion(policy.ExtractVersion("HLWB_CP_1.50.l72"))
}

func testDetail() {
	arbi_list, err := policy.FindArbi("HLTD/HLTD_CP_2.52.000")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(arbi_list)

	grbi_list, err := policy.FindGrbi("HLTD/HLTD_CP_2.52.000")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(grbi_list)

	fmt.Println("------------------------------------------------")

	dal, err := release.NewDal()
	if err != nil {
		fmt.Printf("New DAL failed: %s\n", err)
		return
	}
	defer dal.Close()

	arbi, err := release.FindArbiByPath(dal, "HLWB/HLWB_CP_1.65.000/Seagull_SS_DIALOG_MNH/HL_WB_CP_SS_DIALOG_MYNAH_WP.bin")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(arbi)

	grbi, err := release.FindGrbiByPath(dal, "HLWB/HLWB_CP_1.65.000/HLWB_MSA_1.65.000/MNH/HELAN_A0_M16_AI_Flash.bin")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(grbi)
}
