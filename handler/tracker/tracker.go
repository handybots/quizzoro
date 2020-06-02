package tracker

func IsSpam(id int64, data string) bool {
	defer Data.Set(id, data)

	prevData, ok := Data.Get(id)
	if !ok {
		prevData = data
		return false
	}

	return prevData == data
}
