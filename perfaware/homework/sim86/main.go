package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sim86/registers"
)

type Instruction func(bytes []byte, byteReader *bufio.Reader) string

var instructionLookup = map[uint8]Instruction{
	// Mov
	0b10001000: regMemReg("mov"),
	0b10001001: regMemReg("mov"),
	0b10001010: regMemReg("mov"),
	0b10001011: regMemReg("mov"),
	0b10110000: immedToReg("mov"),
	0b10110001: immedToReg("mov"),
	0b10110010: immedToReg("mov"),
	0b10110011: immedToReg("mov"),
	0b10110100: immedToReg("mov"),
	0b10110101: immedToReg("mov"),
	0b10110110: immedToReg("mov"),
	0b10110111: immedToReg("mov"),
	0b10111000: immedToReg("mov"),
	0b10111001: immedToReg("mov"),
	0b10111010: immedToReg("mov"),
	0b10111011: immedToReg("mov"),
	0b10111100: immedToReg("mov"),
	0b10111101: immedToReg("mov"),
	0b10111110: immedToReg("mov"),
	0b10111111: immedToReg("mov"),
	0b11000110: immedtoRegMem("mov"),
	0b11000111: immedtoRegMem("mov"),
	0b10100000: memToAccAccToMem,
	0b10100001: memToAccAccToMem,
	0b10100010: memToAccAccToMem,
	0b10100011: memToAccAccToMem,
	// Add/Sub/Cmp
	0b10000000: immedtoRegMem(""), // add,sub,cmp
	0b10000001: immedtoRegMem(""), // add,sub,cmp
	0b10000010: immedtoRegMem(""), // add,sub,cmp
	0b10000011: immedtoRegMem(""), // add,sub,cmp
	// Add
	0b00000000: regMemReg("add"),
	0b00000001: regMemReg("add"),
	0b00000010: regMemReg("add"),
	0b00000011: regMemReg("add"),
	0b00000100: immedToReg("add"),
	0b00000101: immedToReg("add"),
	// Sub
	0b00101000: regMemReg("sub"),
	0b00101001: regMemReg("sub"),
	0b00101010: regMemReg("sub"),
	0b00101011: regMemReg("sub"),
	0b00101100: immedToReg("sub"),
	0b00101101: immedToReg("sub"),
	// Cmp
	0b00111000: regMemReg("cmp"),
	0b00111001: regMemReg("cmp"),
	0b00111010: regMemReg("cmp"),
	0b00111011: regMemReg("cmp"),
	0b00111100: immedToReg("cmp"),
	0b00111101: immedToReg("cmp"),
	// Jump
	0b01110100: jump("je"),
	0b01111100: jump("jl"),
	0b01111110: jump("jle"),
	0b01110010: jump("jb"),
	0b01110110: jump("jbe"),
	0b01111010: jump("jp"),
	0b01110000: jump("jo"),
	0b01111000: jump("js"),
	0b01110101: jump("jne"),
	0b01111101: jump("jnl"),
	0b01111111: jump("jnle"),
	0b01110011: jump("jnb"),
	0b01110111: jump("jnbe"),
	0b01111011: jump("jnp"),
	0b01110001: jump("jno"),
	0b01111001: jump("jns"),
	0b11100010: jump("loop"),
	0b11100001: jump("loopz"),
	0b11100000: jump("loopnz"),
	0b11100011: jump("jcxz"),
}

func notImplemented(_ byte, _ *bufio.Reader) string {
	return "\nnot implemented"
}

func regMemReg(mnemonic string) Instruction {
	return func(bytes []byte, br *bufio.Reader) string {
		var left string
		var right string
		b1 := bytes[0]
		b2 := bytes[1]

		byteOrWord := getBits(b1, 0, 0)
		destination := getBits(b1, 1, 1)

		reg := registers.Get(byteOrWord, getBits(b2, 3, 5))
		left = reg.Name
		right = handleModRegRm(byteOrWord, b2, br)

		if destination == 0b0 {
			temp := left
			left = right
			right = temp
		}
		return formatInst(mnemonic, left, right)
	}
}

func immedToReg(mnemonic string) Instruction {
	return func(bytes []byte, br *bufio.Reader) string {
		b1 := bytes[0]
		b2 := bytes[1]
		byteOrWord := getBits(b1, 3, 3)
		regIdx := getBits(b1, 0, 2)
		if mnemonic != "mov" {
			byteOrWord = getBits(b1, 0, 0)
			regIdx = 0b0
		}
		reg := registers.Get(byteOrWord, regIdx)

		num := int16(int8(b2))
		if byteOrWord == 0b1 {
			b3, err := br.ReadByte()
			if err != nil {
				return fmt.Sprint(err)
			}
			num = int16(uint16(b3)<<8 | uint16(b2))
		}
		if mnemonic == "mov" {
			reg.Put(uint16(num))
		} else if mnemonic == "add" {
			reg.Put(reg.Get() + uint16(num))
		} else {
			reg.Put(reg.Get() - uint16(num))
		}

		return formatInst(mnemonic, reg.Name, fmt.Sprint(num))
	}
}

func jump(mnemonic string) Instruction {
	return func(bytes []byte, br *bufio.Reader) string {
		b2 := bytes[1]
		num := int8(b2)
		return fmt.Sprintf("\n%s %s", mnemonic, fmt.Sprint(num))
	}
}

func immedtoRegMem(mnemonic string) Instruction {
	return func(bytes []byte, br *bufio.Reader) string {
		var left string
		var right string
		b1 := bytes[0]
		b2 := bytes[1]

		byteOrWord := getBits(b1, 0, 0)
		signedExt := getBits(b1, 1, 1)
		_, reg, _ := parseModRegRM(b2)

		left = handleModRegRm(byteOrWord, b2, br)

		var num int16
		length := "byte"
		bfirst, err := br.ReadByte()
		if err != nil {
			return fmt.Sprint(err)
		}

		num = int16(bfirst)
		if byteOrWord == 0b1 && signedExt == 0b0 {
			bsecond, err := br.ReadByte()
			if err != nil {
				return fmt.Sprint(err)
			}

			num = int16(bsecond)<<8 | num
			length = "word"
		}
		_ = signedExt
		right = fmt.Sprintf("%s %d", length, num)

		if reg == 0b0 {
			if mnemonic != "mov" {
				mnemonic = "add"
			}
		} else if reg == 0b101 {
			mnemonic = "sub"
		} else if reg == 0b111 {
			mnemonic = "cmp"
		}

		return formatInst(mnemonic, left, right)
	}
}

func memToAccAccToMem(bytes []byte, br *bufio.Reader) string {
	b1 := bytes[0]
	b2 := bytes[1]
	memToAcc := getBits(b1, 1, 1) == 0b0
	byteOrWord := getBits(b1, 0, 0)
	regs := registers.Get(byteOrWord, 0b000).Name

	b3, err := br.ReadByte()
	if err != nil {
		return fmt.Sprint(err)
	}

	dirAdd := uint16(b3)<<8 | uint16(b2)
	left := regs
	right := bracket(fmt.Sprint(dirAdd))
	if !memToAcc {
		return formatInst("mov", right, left)
	}
	return formatInst("mov", left, right)

}

func main() {
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	br := bufio.NewReader(file)
	newfile := "bits 16\n"

	for {
		b1, err1 := br.ReadByte()
		b2, err2 := br.ReadByte()
		if err1 != nil || err2 != nil {
			if err1 == io.EOF {
				break
			}
			newfile += fmt.Sprintf("error parsing bytes 1&2: %v, %v", err1, err2)
			break
		}

		opInst, ok := instructionLookup[b1]
		if ok {
			inst := opInst([]byte{b1, b2}, br)
			fmt.Print(inst)
			newfile += inst
		}
	}
	fmt.Println()
	registers.Print()
	dir := filepath.Dir(fileName)
	newFileName := filepath.Join(dir, "new"+filepath.Base(fileName)+".asm")
	err = os.WriteFile(newFileName, ([]byte)(newfile), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func parseModRegRM(b2 byte) (uint8, uint8, uint8) {
	mod := getBits(b2, 6, 7)
	reg := getBits(b2, 3, 5)
	regm := getBits(b2, 0, 2)
	return mod, reg, regm
}
func handleModRegRm(byteOrWord, b2 byte, br *bufio.Reader) string {
	var right string
	mod, _, regm := parseModRegRM(b2)
	regms, okregm := regmField[regm]
	if !okregm {
		regms = "not found"
	}
	switch mod {
	case 0b11:
		right = registers.Get(byteOrWord, regm).Name
	case 0b00:
		if regm == 0b110 {
			b3, err3 := br.ReadByte()
			b4, err4 := br.ReadByte()
			if err3 != nil || err4 != nil {
				return fmt.Sprintf("%v, %v", err3, err4)
			}
			dirAdd := fmt.Sprintf("%d", uint16(b4)<<8|uint16(b3))

			right = bracket(dirAdd)
			break
		}

		right = bracket(regms)
	case 0b01:
		b3, err := br.ReadByte()
		if err != nil {
			return fmt.Sprint(err)
		}
		displ := fmt.Sprint(int8(b3))

		right = bracket(regms + " + " + displ)
	case 0b10:
		b3, err3 := br.ReadByte()
		b4, err4 := br.ReadByte()
		if err3 != nil || err4 != nil {
			return fmt.Sprintf("%v, %v", err3, err4)
		}
		displ := fmt.Sprintf("%d", int16(uint16(b4)<<8|uint16(b3)))

		right = bracket(regms + " + " + displ)
	}
	return right
}

var regmField = map[uint8]string{
	0b000: "bx + si",
	0b001: "bx + di",
	0b010: "bp + si",
	0b011: "bp + di",
	0b100: "si",
	0b101: "di",
	0b110: "bp",
	0b111: "bx",
}

func getBits(b byte, start, end uint8) uint8 {
	if start > end || end >= 8 {
		panic("Invalid start or end position")
	}
	width := end - start + 1

	mask := uint8((1 << width) - 1)

	return (b >> start) & mask
}

func bracket(s string) string {
	return fmt.Sprintf("[%s]", s)
}

func printB(b byte) {
	fmt.Printf("%b\n", b)
}

func formatInst(inst, left, right string) string {
	return fmt.Sprintf(
		"\n%s %s, %s",
		inst,
		left,
		right,
	)
}
