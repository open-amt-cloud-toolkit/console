package sqldb

func B2i(b bool) int {
	if b {
		return 1
	}

	return 0
}

func I2b(i int) bool {
	return i == 1
}
