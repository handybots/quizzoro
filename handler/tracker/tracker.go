package tracker

func IsSpam(id int64, data string) bool {
	prevData, ok := Data.Get(id)
	if !ok {
		prevData = data
		Data.Set(id, data)
	}

	if prevData == data {
		return true
	}

	Data.Set(id, data)
	return false
}
