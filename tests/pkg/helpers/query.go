package helpers

import (
	"log"
	"strconv"
)

// CheckPageParams validates page number and size
// returns page number and size else an error
func CheckPageParams(pageNum, pageSize string) (int, int, error) {
	pageNumInt, err := strconv.ParseUint(pageNum, 10, 0)
	if err != nil {
		log.Println("page num not parsed as int: ", err)
		return 0, 0, err
	}

	pageSizeInt, err := strconv.ParseUint(pageSize, 10, 0)
	if err != nil {
		log.Println("page size not parsed as int: ", err)
		return 0, 0, err
	}

	return int(pageNumInt), int(pageSizeInt), nil
}

// CalculateOffsetAndLimit returns offset, limit
func CalculateOffsetAndLimit(pageNum, pageSize int) (int, int) {
	offset := (pageNum - 1) * pageSize
	limit := pageSize

	return offset, limit
}
