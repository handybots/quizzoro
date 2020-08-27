package tracker

var exceptions = []string{"privacy"}

func IsSpam(id int64, data string) bool {
	for _, e := range exceptions {
		if e == data {
			return false
		}
	}

	defer Data.Set(id, data)

	prevData, ok := Data.Get(id)
	if !ok {
		prevData = data
		return false
	}

	return prevData == data
}
