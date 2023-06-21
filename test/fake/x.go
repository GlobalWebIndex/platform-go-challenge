package fake

func PtrStr(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

func PtrInt32(i int32) *int32 {
	if i == 0 {
		return nil
	}
	return &i
}

func PtrUInt32(i uint32) *uint32 {
	if i == 0 {
		return nil
	}
	return &i
}

func PtrBool(ni, b bool) *bool {
	if ni {
		return nil
	}
	return &b
}
