package gonuts

import "strings"

// THIS IS NOT REALLY SAFE! JUST A ROUGH WAY TO MAKE IT A LITTLE HARDER

var SANITIZE_SQLSAFER = []string{"`", "´", "'", " or ", " OR ", "=", ";", ":", "(", ")"}

func SanitizeString(badStringsList []string, stringToClean string) (cleanString string) {
	cleanString = stringToClean
	for _, badThing := range badStringsList {
		cleanString = strings.ReplaceAll(cleanString, badThing, "")
	}
	return cleanString
}
