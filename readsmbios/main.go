package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	smbiosdata "github.com/grep-michael/SMBIOS_parser/SMBiosLib/SMBiosData"
)

func main() {

	tableFlag := flag.String("table", "/sys/firmware/dmi/tables/DMI", "Path to load dmi table bytes from")
	epsFlag := flag.String("eps", "/sys/firmware/dmi/tables/smbios_entry_point", "Path to load SMBios Entry point header bytes from")

	var help bool

	flag.BoolVar(&help, "h", false, "Print Usage")
	flag.BoolVar(&help, "help", false, "Print Usage")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(1)
	}

	table, err := LoadRaw(*tableFlag)
	handleError(err)
	eps, err := LoadRaw(*epsFlag)
	handleError(err)

	smbios := smbiosdata.NewSMBiosData(eps, table)
	err = smbios.LoadDMITable()
	handleError(err)

	js, err := json.MarshalIndent(smbios.DMITable, "", "    ")
	handleError(err)
	fmt.Println(string(js))

}
func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func LoadRaw(path string) ([]byte, error) {
	ret, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error Reading %s:\n\t%s", path, err.Error())
	}
	return ret, err
}
