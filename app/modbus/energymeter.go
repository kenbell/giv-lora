package modbus

type regDataType int

const (
	RegTypeUnknown regDataType = iota
	RegTypeFloat
)

type RegType struct {
	Name     string
	Unit     string
	DataType regDataType
}

var EnergyMeterRegisters = map[uint16]RegType{
	0x0010: {"Voltage", "V", RegTypeFloat},
	0x004E: {"Frequency", "Hz", RegTypeFloat},
	0x0052: {"Current", "A", RegTypeFloat},
	0x0092: {"Active Power", "W", RegTypeFloat},
	0x00D2: {"Apparent Power", "VA", RegTypeFloat},
	0x0112: {"Reactive Power", "VAr", RegTypeFloat},
	0x0152: {"Power Factor", "", RegTypeFloat},
	0x0160: {"Import Active Energy", "kWh", RegTypeFloat},
	0x0166: {"Export Active Energy", "kWh", RegTypeFloat},
}
