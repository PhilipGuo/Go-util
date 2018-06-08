package bufmanager

type ibufmanager interface {
	// write data
	WriteData([]byte) bool
	// read data
	ReadData(uint16) ([]byte, bool)
	// read data frame
	ReadDataPrep(uint16) ([]byte, bool)
	// clear data
	ClearData(uint16) bool
	// get data count
	DataCount() uint16
}
