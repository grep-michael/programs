package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	pcielookup "github.com/grep-michael/goPCIe/PCIeLookup"
	"github.com/grep-michael/goPCIe/PCIeTable"
)

func main() {
	vendor := flag.String("vendor", "", "Vendor to look up")
	device := flag.String("device", "", "device lookup")
	class := flag.String("class", "", "class lookup")
	sources := flag.String("sources", "/usr/share/misc/pci.ids,/usr/share/hwdata/pci.ids", "comma separated list of source files")
	flag.Parse()

	table := &pcietable.PCITable{}
	srclist := strings.Split(*sources, ",")
	for _, source := range srclist {
		err := table.LoadSource(source)
		if err != nil {
			fmt.Printf("Failed to process source %s\n", source)
		}
	}

	var ven *pcietable.Vendor
	var dev *pcietable.Device

	status := ""

	if *class != "" {
		lookupClass(*class, srclist)
	}

	if *vendor == "" || (*device == "" && *vendor != "") {
		flag.Usage()
		os.Exit(1)
	}
	if *vendor != "" {
		var found bool
		ven, found = table.FindVendor(*vendor)
		if !found {
			fmt.Printf("Failed to find vendor with id \"%s\"\n", *vendor)
			os.Exit(1)
		}
		status += fmt.Sprintf("Vendor: %s (ID:%s)\n", ven.Name, ven.VendorID)
	}

	if *device != "" {
		var found bool
		dev, found = ven.FindDevice(*device)
		if !found {
			fmt.Printf("Failed to find device with id \"%s\" in vendor %s(ID: %s)\n", *device, ven.Name, ven.VendorID)
			os.Exit(1)
		}
		status += fmt.Sprintf("\tDevice: %s(ID:%s)\n", dev.Name, dev.DeviceID)

	}
	status += fmt.Sprintf("Loaded from: %s", strings.Join(table.Sources, ","))
	fmt.Println(status)
}

func lookupClass(class string, sources []string) {
	for _, src := range sources {
		result, err := pcielookup.PCIeClassLookupFromSource(class, src)
		if err == nil {
			fmt.Printf("%s\n", result.Class)
			os.Exit(0)
		}
	}
	fmt.Printf("Failed to find %s class")
	os.Exit(1)
}
