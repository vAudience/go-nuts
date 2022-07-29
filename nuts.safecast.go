package gonuts

import "encoding/json"

func AssignToCastStringOr(val any, fallbackVal string) string {
	castVal, ok := val.(string)
	if ok {
		return castVal
	}
	return fallbackVal
}

func AssignToCastBoolOr(val any, fallbackVal bool) bool {
	castVal, ok := val.(bool)
	if ok {
		return castVal
	}
	return fallbackVal
}

func AssignToCastInt64Or(val any, fallbackVal int64) int64 {
	// L.Debugf("trying to cast to int64 : %d", val)
	switch expType := val.(type) {
	case float64:
		return int64(expType)
	case json.Number:
		castVal, err := expType.Int64()
		if err == nil {
			return castVal
		}
	}
	return fallbackVal
}
