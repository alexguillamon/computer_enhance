package registers

const (
	Full = iota
	Low
	High
)

type Register struct {
	value *[2]uint8
	rtype uint8
}

func (r *Register) Get() uint16 {
	if r.rtype == Low {
		return uint16(r.value[0])
	} else if r.rtype == High {
		return uint16(r.value[1])
	} else {
		return uint16(r.value[0])<<8 | uint16(r.value[1])
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

var (
	AX = Register{value: &[2]uint8{}}
	AL = Register{value: AX.value, rtype: Low}
	AH = Register{value: AX.value, rtype: High}

	BX = Register{value: &[2]uint8{}}
	BL = Register{value: BX.value, rtype: Low}
	BH = Register{value: BX.value, rtype: High}

	CX = Register{value: &[2]uint8{}}
	CL = Register{value: CX.value, rtype: Low}
	CH = Register{value: CX.value, rtype: High}

	DX = Register{value: &[2]uint8{}}
	DL = Register{value: DX.value, rtype: Low}
	DH = Register{value: DX.value, rtype: High}

	SP = Register{value: &[2]uint8{}}
	DP = Register{value: &[2]uint8{}}
	SI = Register{value: &[2]uint8{}}
	DI = Register{value: &[2]uint8{}}
)
