package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type Adr_Detail struct {
	country string
	state   string
	city    string
}

type Rules struct {
	Name            string
	MainDistributor *Rules
	include         []string
	exclude         []string
}

var (
	RuleMap                             map[string]*Rules
	Distributor_Name, Location, Src_Key string
	Src_det, Dest_det                   Adr_Detail
)

func init() {
	setRules_fn()          // Rules Assign
	l_arg_input := os.Args // Argument Parsing
	if len(l_arg_input) > 2 {
		Distributor_Name = os.Args[1]
		Location = os.Args[2]
		fmt.Println("Name: ", Distributor_Name)
		fmt.Println("Location: ", Location)
	}

	file, err := os.Open("cities.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	lineCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		key := strings.ToUpper(record[3]) + "," + strings.ToUpper(record[4]) + "," + strings.ToUpper(record[5])
		l_source := FormSrcKey_fn(key) //Removing space between name to form Source key
		if l_source == Location {
			Src_Key = l_source
			fmt.Println("Location is Available in the City List")
			break
		}
		lineCount += 1
	}
}

func FormSrcKey_fn(l_key string) (r_key string) {
	l_contains := strings.ContainsAny(l_key, " ")
	if l_contains {
		r_key = strings.Replace(l_key, " ", "", -1)
	} else {
		return l_key
	}
	return
}

func setRules_fn() {
	Rule1 := &Rules{}
	Rule1.Name = "DISTRIBUTOR1"
	Rule1.include = []string{"INDIA", "UNITEDSTATES"}
	Rule1.exclude = []string{"CHENNAI,TAMILNADU,INDIA", "KARNATAKA,INDIA"}

	Rule2 := &Rules{}
	Rule2.Name = "DISTRIBUTOR2"
	Rule2.MainDistributor = Rule1
	Rule2.include = []string{"INDIA"}
	Rule2.exclude = []string{"TAMILNADU,INDIA"}

	Rule3 := &Rules{}
	Rule3.Name = "DISTRIBUTOR3"
	Rule3.MainDistributor = Rule2
	Rule3.include = []string{"HUBLI,KARNATAKA,INDIA"}

	RuleMap = make(map[string]*Rules)
	RuleMap["DISTRIBUTOR1"] = Rule1
	RuleMap["DISTRIBUTOR2"] = Rule2
	RuleMap["DISTRIBUTOR3"] = Rule3
	return
}

func get_adrdetail_fn(dest_key string) (r_dest Adr_Detail) { // Converting Location into list of names
	loc := strings.Split(dest_key, ",")
	if len(loc) >= 3 {
		r_dest.city = loc[0]
		r_dest.state = loc[1]
		r_dest.country = loc[2]
	} else if len(loc) >= 2 {
		r_dest.state = loc[0]
		r_dest.country = loc[1]
	} else if len(loc) >= 1 {
		r_dest.country = loc[0]
	}
	return
}

func CheckPermission_fn(l_DistributorRule *Rules) (r_non_access, r_access bool) {
	if l_DistributorRule.MainDistributor != nil {
		l_main_dist := l_DistributorRule.MainDistributor
		l_non_access, l_access := CheckPermission_fn(l_main_dist)
		if l_non_access == true {
			return l_non_access, l_access
		}
	}

	if l_DistributorRule.exclude != nil {
		fmt.Println("***** Checking Rules With Non Access Region *****")
		l_permission_exclude := CheckList_fn(l_DistributorRule.exclude) // Check permission for Exclude list
		if l_permission_exclude == true {
			fmt.Println(l_DistributorRule.Name, " has no permission")
			r_non_access = true
		}
	}
	if r_non_access == false {
		fmt.Println("Location is not in Restricted Region !!!! ")
		fmt.Println("\n***** Checking Rules With Access Region *****")
		l_permission_include := CheckList_fn(l_DistributorRule.include) // Check permission for Include List
		if l_permission_include {
			fmt.Println(l_DistributorRule.Name, "has Permission")
			r_access = true
		} else {
			fmt.Println(l_DistributorRule.Name, " has no permission")
			r_access = false
		}
	}
	return
}

func CheckList_fn(l_list []string) (result bool) {
	for index := range l_list {
		l_list_det := get_adrdetail_fn(l_list[index])
		if l_list_det.country != "" {
			if l_list_det.country == Dest_det.country {
				result = true
			} else {
				result = false
			}
		}
		if l_list_det.state != "" {
			if l_list_det.state == Dest_det.state {
				result = true
			} else {
				result = false
			}
		}
		if l_list_det.city != "" {
			if l_list_det.city == Dest_det.city {

				result = true
			} else {
				result = false
			}
		}
		if result == true {
			break
		}
	}
	fmt.Println("Location Detail is matched with list")
	return
}

func main() {
	if Src_Key == "" {
		fmt.Println("OOPS !! Location not in List")
		return
	}
	Dest_det = get_adrdetail_fn(Location)
	l_DistributorRules := RuleMap[Distributor_Name]
	l_non_access, l_access := CheckPermission_fn(l_DistributorRules)
	if l_non_access == false && l_access == true {
		fmt.Println("\n******* ", Distributor_Name, "having access in ", Location, " *******")
	} else {
		fmt.Println("\n******* ", Distributor_Name, "having no access in ", Location, " *******")
	}
	return
}
