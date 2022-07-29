package gonuts

import (
	"math"
	"strconv"
)

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func BytesToNiceString(size int64) (formattedString string) {
	var newSize float64
	sizeString := ""
	if size > 1024*1024*1024*1024 {
		newSize = BytesToTB(size, 2)
		sizeString = strconv.FormatFloat(newSize, 'f', -1, 64) + "TB"
	} else if size > 1024*1024*1024 {
		newSize = BytesToGB(size, 2)
		sizeString = strconv.FormatFloat(newSize, 'f', -1, 64) + "GB"
	} else if size > 1024*1024 {
		newSize = BytesToMB(size, 2)
		sizeString = strconv.FormatFloat(newSize, 'f', -1, 64) + "MB"
	} else if size > 1024 {
		newSize = BytesToKB(size, 2)
		sizeString = strconv.FormatFloat(newSize, 'f', -1, 64) + "KB"
	} else if size > 99 {
		newSize = BytesToKB(size, 2)
		sizeString = strconv.FormatFloat(newSize, 'f', -1, 64) + "KB"
	} else {
		sizeString = strconv.FormatInt(size, 10) + "B"
	}
	// L.Debugf("[BytesToNiceString] from [%d] to [%s]", size, sizeString)
	return sizeString
}

func BytesToKB(size int64, digits int) (kilobytes float64) {
	newSize := Round(float64(size)/1024, .5, digits)
	return newSize
}
func BytesToMB(size int64, digits int) (megabytes float64) {
	newSize := Round(float64(size)/1024/1024, .5, digits)
	return newSize
}
func BytesToGB(size int64, digits int) (terrabytes float64) {
	newSize := Round(float64(size)/1024/1024/1024, .5, digits)
	return newSize
}
func BytesToTB(size int64, digits int) (terrabytes float64) {
	newSize := Round(float64(size)/1024/1024/1024/1024, .5, digits)
	return newSize
}
