package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type PNM [5]byte

func (p PNM) String() string {
	return string(p[:])
}

func (marshable PNM) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshable.String())
}

type PSN [4]byte

func (p PSN) String() string {
	return fmt.Sprintf("0x%0x", binary.BigEndian.Uint32(p[:]))
}

func (marshable PSN) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshable.String())
}

type MDT [2]byte

func (m MDT) String() string {
	date_int := binary.BigEndian.Uint16(m[:])
	date_int = date_int & 0b0000111111111111 //we just removed the reserved bits, fuck em
	year := 2000 + (date_int >> 4)
	month := date_int & 0b1111
	return fmt.Sprintf("%02d/%d", month, year)
}

func (marshable MDT) MarshalJSON() ([]byte, error) {
	return json.Marshal(marshable.String())
}

type CID struct {
	ManID       byte   `json:"ManufacturerID"`
	OEMID       uint16 `json:"EOMID"`
	ProdName    PNM    `json:"ProductName"`
	ProdRev     byte   `json:"ProductRevision"`
	ProdSerial  PSN    `json:"ProductSerial"`
	ManufacDate MDT    `json:"ManufacturerDate"`
	CRC         byte   `json:"CRC"`
}

func parseCID(bytes []byte) CID {
	cid := CID{}
	cid.ManID = bytes[0]
	cid.OEMID = binary.BigEndian.Uint16(bytes[1:3])
	cid.ProdName = [5]byte(bytes[3:8])
	cid.ProdRev = bytes[8]
	cid.ProdSerial = [4]byte(bytes[9:13])
	cid.ManufacDate = [2]byte(bytes[13:15])
	cid.CRC = bytes[15] >> 1
	return cid
}

func main() {

	var cidInput string

	if len(os.Args) == 2 {
		cidInput = os.Args[1]
	} else {
		cidInput = readSTDin()
	}

	cidInput = strings.TrimSpace(cidInput)
	if cidInput == "" {
		fmt.Printf("No cid Passed\n\tUsage: cid <hex string>\n\tOr:    cat cidfile | readcid\n\n")
		os.Exit(1)
	}

	bytes, err := hex.DecodeString(cidInput)
	if err != nil {
		fmt.Printf("Error converting cid to bytes:\n\t%+v\n", err)
		os.Exit(1)
	}

	if len(bytes) != 16 {
		fmt.Printf("CID length incorrect\n\tExpected: 16\n\tActual: %d\n", len(bytes))
		os.Exit(1)
	}

	cid := parseCID(bytes)
	js, err := json.MarshalIndent(cid, "", "    ")
	if err != nil {
		fmt.Printf("Error Converting cid to json:\n\t%+v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(js))
}

func readSTDin() string {
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error checking stdin:\n\t%+v\n", err)
		os.Exit(1)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return ""
	}

	pipedData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading stdin:\n\t%+v\n", err)
		os.Exit(1)
	}

	return string(pipedData)
}
