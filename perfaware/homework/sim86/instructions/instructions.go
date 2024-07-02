package instructions

import (
	"bufio"
	"fmt"
)

type BitType int

const (
	BitsLiteral BitType = iota
	BitsD
	BitsS
	BitsW
	BitsRM
	BitsMOD
	BitsREG
	BitsDisp
	BitsDataLo
	BitsDataHi
	BitsAddr
)

type Bit struct {
	Type  BitType
	Size  uint8
	Value byte
}

func B(b string) Bit {
	var value byte
	for i := 0; i < len(b); i++ {
		value = value << 1
		if b[i] == '1' {
			value |= 1
		}
	}
	return Bit{Type: BitsLiteral, Size: uint8(len(b)), Value: value}
}

var (
	D         = Bit{Type: BitsD, Size: 1}
	S         = Bit{Type: BitsS, Size: 1}
	W         = Bit{Type: BitsW, Size: 1}
	RM        = Bit{Type: BitsRM, Size: 3}
	MOD       = Bit{Type: BitsMOD, Size: 2}
	REG       = Bit{Type: BitsREG, Size: 3}
	DISP      = Bit{Type: BitsDisp, Size: 8}
	ADDR      = Bit{Type: BitsAddr, Size: 16}
	DATA      = Bit{Type: BitsDataLo, Size: 8}
	DATA_IF_W = Bit{Type: BitsDataHi, Size: 8}
)

func ImplD(value byte) Bit {
	return Bit{Type: BitsD, Size: 0, Value: value}
}

func ImplREG(value byte) Bit {
	return Bit{Type: BitsREG, Size: 0, Value: value}
}

func ImplMOD(value byte) Bit {
	return Bit{Type: BitsMOD, Size: 0, Value: value}
}

func ImplRM(value byte) Bit {
	return Bit{Type: BitsRM, Size: 0, Value: value}
}

type OpType uint8

const (
	OpNone OpType = iota
	OpMov
	OpAdd
	OpSub
	OpCmp
	OpJmp
)

var opToString = map[OpType]string{
	OpNone: "",
	OpMov:  "mov",
	OpAdd:  "add",
	OpSub:  "sub",
	OpCmp:  "cmp",
	OpJmp:  "jmp",
}

type DecodeScheme struct {
	Mnemonic OpType
	Bits     []Bit
}
type Operation struct {
	Op      OpType
	Literal []Bit

	D,
	S,
	W,
	RM,
	REG,
	MOD,
	DISP_LO,
	DISP_HI,
	ADDR_LO,
	ADDR_HI,
	DATA_LO,
	DATA_HI *Bit
}

var OpNotFound = Operation{
	Op: OpNone,
}

var noOperation = DecodeScheme{OpNone, nil}

var movRegMemReg = DecodeScheme{OpMov, []Bit{B("100010"), D, W, MOD, REG, RM}}
var movImmedRegMem = DecodeScheme{OpMov, []Bit{B("1100011"), W, MOD, B("000"), RM, DATA, DATA_IF_W, ImplD(0)}}
var movImmedReg = DecodeScheme{OpMov, []Bit{B("1011"), W, REG, DATA, DATA_IF_W, ImplD(1)}}
var movMem2Acc = DecodeScheme{OpMov, []Bit{B("1010000"), W, ADDR, ImplD(1), ImplMOD(0), ImplREG(0)}}
var movAcc2Mem = DecodeScheme{OpMov, []Bit{B("1010001"), W, ADDR, ImplD(0), ImplMOD(0), ImplREG(0)}}

var Table = map[uint8]DecodeScheme{
	// Mov
	0b10001000: movRegMemReg,
	0b10001001: movRegMemReg,
	0b10001010: movRegMemReg,
	0b10001011: movRegMemReg,
	0b10110000: movImmedRegMem,
	0b10110001: movImmedRegMem,
	0b10110010: movImmedRegMem,
	0b10110011: movImmedRegMem,
	0b10110100: movImmedRegMem,
	0b10110101: movImmedRegMem,
	0b10110110: movImmedRegMem,
	0b10110111: movImmedRegMem,
	0b10111000: movImmedRegMem,
	0b10111001: movImmedRegMem,
	0b10111010: movImmedRegMem,
	0b10111011: movImmedRegMem,
	0b10111100: movImmedRegMem,
	0b10111101: movImmedRegMem,
	0b10111110: movImmedRegMem,
	0b10111111: movImmedRegMem,
	0b11000110: movImmedReg,
	0b11000111: movImmedReg,
	0b10100000: movMem2Acc,
	0b10100001: movMem2Acc,
	0b10100010: movAcc2Mem,
	0b10100011: movAcc2Mem,
	// Add/Sub/Cmp
	0b10000000: noOperation, // immedtoRegMem(""), // add,sub,cmp
	0b10000001: noOperation, // immedtoRegMem(""), // add,sub,cmp
	0b10000010: noOperation, // immedtoRegMem(""), // add,sub,cmp
	0b10000011: noOperation, // immedtoRegMem(""), // add,sub,cmp
	// Add
	0b00000000: noOperation, // regMemReg("add"),
	0b00000001: noOperation, // regMemReg("add"),
	0b00000010: noOperation, // regMemReg("add"),
	0b00000011: noOperation, // regMemReg("add"),
	0b00000100: noOperation, // immedToReg("add"),
	0b00000101: noOperation, // immedToReg("add"),
	// Sub
	0b00101000: noOperation, // regMemReg("sub"),
	0b00101001: noOperation, // regMemReg("sub"),
	0b00101010: noOperation, // regMemReg("sub"),
	0b00101011: noOperation, // regMemReg("sub"),
	0b00101100: noOperation, // immedToReg("sub"),
	0b00101101: noOperation, // immedToReg("sub"),
	// Cmp
	0b00111000: noOperation, // regMemReg("cmp"),
	0b00111001: noOperation, // regMemReg("cmp"),
	0b00111010: noOperation, // regMemReg("cmp"),
	0b00111011: noOperation, // regMemReg("cmp"),
	0b00111100: noOperation, // immedToReg("cmp"),
	0b00111101: noOperation, // immedToReg("cmp"),
	// Jump
	0b01110100: noOperation, // jump("je"),
	0b01111100: noOperation, // jump("jl"),
	0b01111110: noOperation, // jump("jle"),
	0b01110010: noOperation, // jump("jb"),
	0b01110110: noOperation, // jump("jbe"),
	0b01111010: noOperation, // jump("jp"),
	0b01110000: noOperation, // jump("jo"),
	0b01111000: noOperation, // jump("js"),
	0b01110101: noOperation, // jump("jne"),
	0b01111101: noOperation, // jump("jnl"),
	0b01111111: noOperation, // jump("jnle"),
	0b01110011: noOperation, // jump("jnb"),
	0b01110111: noOperation, // jump("jnbe"),
	0b01111011: noOperation, // jump("jnp"),
	0b01110001: noOperation, // jump("jno"),
	0b01111001: noOperation, // jump("jns"),
	0b11100010: noOperation, // jump("loop"),
	0b11100001: noOperation, // jump("loopz"),
	0b11100000: noOperation, // jump("loopnz"),
	0b11100011: noOperation, // jump("jcxz"),
}

func Parse(b1 byte, br *bufio.Reader) (op Operation) {
	decodeScheme, ok := Table[b1]
	if !ok {
		panic("op code not found can't continue decoding")
	}

	op.Op = decodeScheme.Mnemonic

	var currentByte = b1
	var currentBit uint8 = 0
	var bitsRead uint8 = 0
	for _, bits := range decodeScheme.Bits {
		if bitsRead == 8 {
			b, err := br.ReadByte()
			if err != nil {
				panic("Error reading byte")
			}
			currentByte = b
		}

		if currentBit == 8 {
			currentBit = 0
		}

		bits.Value = readBits(currentByte, currentBit, bits.Size-1)
		switch bits.Type {
		case BitsLiteral:
			op.Literal = append(op.Literal, bits)
		case BitsD:
			op.D = &bits
		case BitsS:
			op.S = &bits
		case BitsW:
			op.W = &bits
		case BitsMOD:
			op.MOD = &bits
		case BitsRM:
			op.RM = &bits
			switch op.MOD.Value {
			case 0b00:
				if bits.Value == 0b110 {
					bLo, errLo := br.ReadByte()
					bHi, errHi := br.ReadByte()
					if errLo != nil || errHi != nil {
						panic(fmt.Sprintf("%v, %v", errLo, errHi))
					}
					op.ADDR_LO.Value = bLo
					op.ADDR_HI.Value = bHi
				}
			case 0b01:
				b, err := br.ReadByte()
				if err != nil {
					panic(fmt.Sprintf("%v", err))
				}
				op.DISP_LO.Value = b
			case 0b10:
				bLo, errLo := br.ReadByte()
				bHi, errHi := br.ReadByte()
				if errLo != nil || errHi != nil {
					panic(fmt.Sprintf("%v, %v", errLo, errHi))
				}
				op.DISP_LO.Value = bLo
				op.DISP_HI.Value = bHi
			}
		case BitsREG:
			op.REG = &bits
		case BitsDataLo:
			b, err := br.ReadByte()
			if err != nil {
				panic(fmt.Sprintf("%v", err))
			}
			op.DATA_LO.Value = b
		case BitsDataHi:
			if op.W.Value == 0b1 {
				b, err := br.ReadByte()
				if err != nil {
					panic(fmt.Sprintf("%v", err))
				}
				op.DATA_HI.Value = b
			}
		case BitsAddr:
			bLo, errLo := br.ReadByte()
			bHi, errHi := br.ReadByte()
			if errLo != nil || errHi != nil {
				panic(fmt.Sprintf("%v, %v", errLo, errHi))
			}
			op.ADDR_LO.Value = bLo
			op.ADDR_HI.Value = bHi
		}

		currentBit += bits.Size
		bitsRead += bits.Size
	}

	return op
}

func readBits(b byte, start, end uint8) uint8 {
	if start > end || end >= 8 {
		panic("Invalid start or end position")
	}
	width := end - start + 1

	mask := uint8((1 << width) - 1)

	return (b >> start) & mask
}
