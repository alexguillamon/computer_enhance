package registers

import "fmt"

const (
	Full = iota
	Low
	High
)

type Register struct {
	value *[2]uint8
	rtype uint8
	Name  string
}

func (r *Register) Get() uint16 {
	if r.rtype == Low {
		return uint16(r.value[0])
	} else if r.rtype == High {
		return uint16(r.value[1])
	} else {
		return uint16(r.value[1])<<8 | uint16(r.value[0])
	}

}
func (r *Register) Put(val uint16) {
	if r.rtype == Low {
		r.value[0] = uint8(val)
	} else if r.rtype == High {
		r.value[1] = uint8(val)
	} else {
		r.value[0] = uint8(val)
		r.value[1] = uint8(val >> 8)
	}
}
func (r *Register) String() string {
	val := r.Get()
	if r.rtype == Low {
		return fmt.Sprintf("%s %b (lo)", r.Name, val)
	} else if r.rtype == High {
		return fmt.Sprintf("%s %b (hi)", r.Name, val)
	} else {
		return fmt.Sprintf("%s %b", r.Name, val)
	}
}

var (
	AX = Register{value: &[2]uint8{}, Name: "ax"}
	AL = Register{value: AX.value, rtype: Low, Name: "al"}
	AH = Register{value: AX.value, rtype: High, Name: "ah"}

	BX = Register{value: &[2]uint8{}, Name: "bx"}
	BL = Register{value: BX.value, rtype: Low, Name: "bl"}
	BH = Register{value: BX.value, rtype: High, Name: "bh"}

	CX = Register{value: &[2]uint8{}, Name: "cx"}
	CL = Register{value: CX.value, rtype: Low, Name: "cl"}
	CH = Register{value: CX.value, rtype: High, Name: "ch"}

	DX = Register{value: &[2]uint8{}, Name: "dx"}
	DL = Register{value: DX.value, rtype: Low, Name: "dl"}
	DH = Register{value: DX.value, rtype: High, Name: "dh"}

	SP = Register{value: &[2]uint8{}, Name: "sp"}
	BP = Register{value: &[2]uint8{}, Name: "bp"}
	SI = Register{value: &[2]uint8{}, Name: "si"}
	DI = Register{value: &[2]uint8{}, Name: "di"}

	// Array of all registers
	RegistersArray = []*Register{
		&AL,
		&CL,
		&DL,
		&BL,
		&AH,
		&CH,
		&DH,
		&BH,
		&AX,
		&CX,
		&DX,
		&BX,
		&SP,
		&BP,
		&SI,
		&DI,
	}
)

func Get(bw, idx uint8) Register {
	if idx > 15 {
		panic("idx is bigger than then number of registers")
	}
	if bw == 0b1 {
		idx = 1<<3 | idx
	}
	return *RegistersArray[idx]
}

func Print() {
	for i := 0; i < len(RegistersArray); i++ {
		fmt.Printf("%s\n", RegistersArray[i])
	}
}
