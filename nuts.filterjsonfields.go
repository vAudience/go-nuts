package gonuts

import "encoding/json"

// from here: https://stackoverflow.com/questions/17306358/removing-fields-from-struct-or-hiding-them-in-json-response by https://stackoverflow.com/users/7496198/chhaileng
func FilterJsonFields(obj any, ignoreFields ...string) (string, error) {
	toJson, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	if len(ignoreFields) == 0 {
		return string(toJson), nil
	}
	toMap := map[string]any{}
	json.Unmarshal([]byte(string(toJson)), &toMap)
	for _, field := range ignoreFields {
		delete(toMap, field)
	}
	toJson, err = json.Marshal(toMap)
	if err != nil {
		return "", err
	}
	return string(toJson), nil
}
