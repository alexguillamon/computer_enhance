package main

import (
	"bufio"
	"fmt"
	"os"
)

var bitsToInst = map[uint8]string{
	0b100010:  "mov",
	0b1100011: "mov",
	0b1011:    "mov",
	0b101000:  "mov",
}

var regField = map[uint8]map[uint8]string{
	0b0: {
		0b000: "al",
		0b001: "cl",
		0b010: "dl",
		0b011: "bl",
		0b100: "ah",
		0b101: "ch",
		0b110: "dh",
		0b111: "bh"},
	0b1: {
		0b000: "ax",
		0b001: "cx",
		0b010: "dx",
		0b011: "bx",
		0b100: "sp",
		0b101: "bp",
		0b110: "si",
		0b111: "di",
	},
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
	// Calculate the width of the bit range
	width := end - start + 1

	// Create a mask for the width
	mask := uint8((1 << width) - 1)

	// Shift right to align the desired bits to the least significant position and apply the mask
	return (b >> start) & mask
}

func formatInst(inst, left, right string) string {
	return fmt.Sprintf(
		"\n%s %s, %s",
		inst,
		left,
		right,
	)
}

func bracket(s string) string {
	return fmt.Sprintf("[%s]", s)
}

func printB(b byte) {
	fmt.Printf("%b\n", b)
}

func main() {
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	breader := bufio.NewReader(file)
	newfile := "bits 16\n"
Loop:
	for {
		var opCode uint8
		var inst string
		var left string
		var right string

		b1, err1 := breader.ReadByte()
		if err1 != nil {
			break
		}

		for i := 0; i < 5; i++ {
			opInst, ok := bitsToInst[b1>>i]
			if ok {
				opCode = b1 >> i
				inst = opInst
				break
			}
		}
		if inst == "" {
			break Loop
		}

		byteOrWord := getBits(b1, 0, 0)
		b2, err2 := breader.ReadByte()
		if err2 != nil {
			break Loop
		}
		mod := getBits(b2, 6, 7)
		reg := getBits(b2, 3, 5)
		regm := getBits(b2, 0, 2)
		regs, okreg := regField[byteOrWord][reg]
		regms, okregm := regmField[regm]
		if !okreg || !okregm {
			fmt.Printf("no mapping  inst %s, reg %b, regm %b", inst, reg, regm)
			break Loop
		}

		switch opCode {
		case 0b100010:
			destination := getBits(b1, 1, 1)
			switch mod {
			case 0b11:
				regms, okregm := regField[byteOrWord][regm]
				if !okregm {
					fmt.Printf("no mapping  inst %s, reg %b, regm %b", inst, reg, regm)
					break Loop
				}

				left = regs
				right = regms
			case 0b00:
				if regm == 0b110 {
					b3, err3 := breader.ReadByte()
					b4, err4 := breader.ReadByte()
					if err3 != nil || err4 != nil {
						fmt.Println(err3)
						fmt.Println(err4)
						break Loop
					}
					dirAdd := fmt.Sprintf("%d", uint16(b4)<<8|uint16(b3))

					left = regs
					right = bracket(dirAdd)
					break
				}

				left = regs
				right = bracket(regms)
			case 0b01:
				b3, err := breader.ReadByte()
				if err != nil {
					fmt.Println(err)
					break Loop
				}
				displ := fmt.Sprint(int8(b3))

				left = regs
				right = bracket(regms + " + " + displ)
			case 0b10:
				b3, err3 := breader.ReadByte()
				b4, err4 := breader.ReadByte()
				if err3 != nil || err4 != nil {
					fmt.Println(err3)
					fmt.Println(err4)
					break Loop
				}
				displ := fmt.Sprintf("%d", int16(uint16(b4)<<8|uint16(b3)))

				left = regs
				right = bracket(regms + " + " + displ)
			}
			if destination == 0b0 {
				temp := left
				left = right
				right = temp
			}
		case 0b1011:
			byteOrWord := getBits(b1, 3, 3)
			reg := getBits(b1, 0, 2)
			regs, ok := regField[byteOrWord][reg]
			if !ok {
				break Loop
			}

			num := uint16(b2)
			if byteOrWord == 0b1 {
				b3, err := breader.ReadByte()
				if err != nil {
					fmt.Println(err)
					break Loop
				}
				num = uint16(b3)<<8 | uint16(b2)
			}
			left = regs
			right = fmt.Sprint(num)

		case 0b1100011:
			switch mod {
			case 0b00:
				if regm == 0b110 {
					b3, err3 := breader.ReadByte()
					b4, err4 := breader.ReadByte()
					if err3 != nil || err4 != nil {
						fmt.Println(err3)
						fmt.Println(err4)
						break Loop
					}
					dirAdd := fmt.Sprintf("%d", uint16(b4)<<8|uint16(b3))

					left = bracket(dirAdd)
					break
				}
				left = bracket(regms)
			case 0b01:
				b3, err := breader.ReadByte()
				if err != nil {
					fmt.Println(err)
					break Loop
				}
				displ := fmt.Sprint(int8(b3))

				left = bracket(regms + " + " + displ)
			case 0b10:
				b3, err3 := breader.ReadByte()
				b4, err4 := breader.ReadByte()
				if err3 != nil || err4 != nil {
					fmt.Println(err3)
					fmt.Println(err4)
					break Loop
				}
				displ := fmt.Sprintf("%d", int16(uint16(b4)<<8|uint16(b3)))

				left = bracket(regms + " + " + displ)
			}
			var num uint16
			length := "byte"
			bfirst, err := breader.ReadByte()
			if err != nil {
				fmt.Println(err)
				break Loop
			}
			num = uint16(bfirst)
			if byteOrWord == 0b1 {
				bsecond, err := breader.ReadByte()
				if err != nil {
					fmt.Println(err)
					break Loop
				}
				num = uint16(bsecond)<<8 | num
				length = "word"
			}
			right = fmt.Sprintf("%s %d", length, num)

		case 0b101000:
			memToAcc := getBits(b1, 1, 1) == 0b0
			regs := regField[byteOrWord][0b000]

			b3, err := breader.ReadByte()
			if err != nil {
				break Loop
			}

			dirAdd := uint16(b3)<<8 | uint16(b2)
			left = regs
			right = bracket(fmt.Sprint(dirAdd))
			if !memToAcc {
				tmp := left
				left = right
				right = tmp
			}
		}

		line := formatInst(inst, left, right)
		fmt.Print(line)
		newfile += line
	}

	err = os.WriteFile("new"+fileName+".asm", ([]byte)(newfile), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
