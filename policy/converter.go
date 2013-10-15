package policy

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	VersionPrefixLength   = len("_CP_")
	VersionExtractPattern = regexp.MustCompile(`_CP_\d+\.\d+\.\w{3}$`)
)

func ExtractVersion(name string) string {
	var version string = ""
	version = VersionExtractPattern.FindString(name)
	if version != "" {
		version = version[VersionPrefixLength:]
	}
	return version
}

func QuantitateVersion(version string) int64 {
	v_list := strings.Split(version, ".")
	if len(v_list) != 3 {
		return -1
	}

	/*****slice 0*****/
	v_list_0 := v_list[0]
	v_list_0_int, err := strconv.ParseInt(v_list_0, 10, 64)
	if err != nil {
		return -1
	}
	v_list_0_int = v_list_0_int * 1000000000
	//	fmt.Printf("0 str = %s, 0 int = %d\n", v_list_0, v_list_0_int)

	/*****slice 1*****/
	v_list_1 := v_list[1]
	v_list_1_int, err := strconv.ParseInt(v_list_1, 10, 64)
	if err != nil {
		return -1
	}
	v_list_1_int = v_list_1_int * 1000000
	//	fmt.Printf("1 str = %s, 1 int = %d\n", v_list_1, v_list_1_int)

	/*****slice 2*****/
	v_list_2 := strings.ToUpper(v_list[2])
	//	fmt.Printf("original = %s\n", v_list[2])
	//	fmt.Printf("converted = %s\n", v_list_2)
	v_list_2_bytes := []byte(v_list_2)
	if len(v_list_2_bytes) != 3 {
		return -1
	}
	v_list_2_str := fmt.Sprintf("%d%d%d", v_list_2_bytes[0], v_list_2_bytes[1], v_list_2_bytes[2])
	v_list_2_int, err := strconv.ParseInt(v_list_2_str, 10, 64)
	if err != nil {
		return -1
	}
	//	fmt.Printf(" 2 ori_str = %s, str = %s, int = %d\n", v_list_2, v_list_2_str, v_list_2_int)

	v_int := v_list_0_int + v_list_1_int + v_list_2_int
	//	fmt.Printf("final int version = %d\n", v_int)
	return v_int
}
