package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("./listing_0038_many_register_mov")
	if err != nil {
		return
	}
	defer file.Close()

	breader := bufio.NewReader(file)
	newfile := "bits 16\n\n"
	for {
		b1, err1 := breader.ReadByte()
		b2, err2 := breader.ReadByte()
		if err1 != nil || err2 != nil {
			break
		}

		var bits = make([]string, 0, 6)
		first := fmt.Sprintf("%b", b1)
		second := fmt.Sprintf("%b", b2)
		bits = append(bits, first[0:6], string(first[6]), string(first[7]))
		bits = append(bits, second[0:2], second[2:5], second[5:8])
		asm, err := buildASM(bits)
		if err != nil {
			fmt.Println(err)
		}
		newfile += asm + "\n"

	}
	fmt.Println(newfile)
	err = os.WriteFile("newmany.asm", ([]byte)(newfile), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

var bitsToInst = map[string]string{
	"100010": "mov",
}
var regField = map[string]map[string]string{
	"0": {
		"000": "al",
		"001": "cl",
		"010": "dl",
		"011": "bl",
		"100": "ah",
		"101": "ch",
		"110": "dh",
		"111": "bh"},
	"1": {
		"000": "ax",
		"001": "cx",
		"010": "dx",
		"011": "bx",
		"100": "sp",
		"101": "bp",
		"110": "si",
		"111": "di",
	},
}

var regFieldW1 = map[string]string{}

func buildASM(fields []string) (string, error) {
	destination := fields[1]
	opOnWord := fields[2]
	inst, ok := bitsToInst[fields[0]]
	if !ok {
		return "", fmt.Errorf("no mapping instruction")
	}

	var reg string
	var regm string
	rname, ok := regField[opOnWord][fields[4]]
	if !ok {
		return "", fmt.Errorf("no mapping for register: %v,  %v", opOnWord, fields[4])
	}
	reg = rname
	rname, ok = regField[opOnWord][fields[5]]
	if !ok {
		return "", fmt.Errorf("no mapping for register: %v,  %v", opOnWord, fields[5])
	}
	regm = rname

	if destination == "1" {
		return fmt.Sprintf("%s %s, %s", inst, reg, regm), nil
	} else {
		return fmt.Sprintf("%s %s, %s", inst, regm, reg), nil
	}
}
