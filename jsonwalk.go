package jsonwalk

func Object(jsonObject []byte, fn func(key, value []byte) error) error {
	_, err := parseObject(jsonObject, fn)
	return err
}

func Array(fn func(value []byte) error) error {
	return nil
}
