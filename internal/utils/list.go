package utils

import "gwi_api/internal/domain"

func Contains[T domain.Comparable](arr []T, value interface{}) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
