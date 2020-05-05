package utils

func ConvertBytesToHash(data []byte) [32]byte {
	var hashArr [32]byte

	if len(data) >= 32 {
		copy(hashArr[:], data[:32])
	}

	return hashArr
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
