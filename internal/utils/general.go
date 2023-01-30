package utils

func MakeRange(min, max uint64) []uint64 {
	if max <= min {
		return nil
	}
	answer := make([]uint64, max-min)
	for i := 0; uint64(i) < (max - min); i++ {
		answer[i] = min + uint64(i)
	}
	return answer
}
