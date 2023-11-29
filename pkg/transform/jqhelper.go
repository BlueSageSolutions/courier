package transform

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/itchyny/gojq"
)

// FilterByQuery ...
func FilterByQuery(inputJSON string, inputQuery string) ([]byte, error) {
	query, err := gojq.Parse(inputQuery)
	if err != nil {
		return nil, err
	}
	return Run([]byte(inputJSON), query)
}

func Iter(iter gojq.Iter) (result []byte, err error) {
	var iteratedValues []interface{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}
		iteratedValues = append(iteratedValues, v)
		result, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}
	if len(iteratedValues) > 0 {
		return getConsolidatedResult(iteratedValues)
	}
	return
}

func getConsolidatedResult(iterValues []interface{}) ([]byte, error) {
	var result []byte
	comma := ","
	for _, val := range iterValues {
		v, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		v = []byte(string(v) + comma)
		result = append(result, v...)
	}
	result = []byte(strings.TrimSuffix(string(result), comma))
	return result, nil
}

func Run(inputJSON []byte, query *gojq.Query) ([]byte, error) {
	var inputVar interface{}
	err := json.Unmarshal(inputJSON, &inputVar)
	if err != nil {
		return nil, err
	}

	iter := query.Run(inputVar)
	return Iter(iter)
}

// FilterByQueryWithStruct ...
func FilterByQueryWithStruct(inputVar interface{}, inputQuery string) (interface{}, error) {
	query, err := gojq.Parse(inputQuery)
	if err != nil {
		return nil, err
	}
	iter := query.Run(inputVar) // or query.RunWithContext
	var previousResult interface{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}
		previousResult = v
	}

	return previousResult, nil
}

// GroupArrayOfObjectsByKey ...
func GroupArrayOfObjectsByKey(inputJSON string, inputkey string) ([]byte, error) {
	inputQuery := `[group_by(.` + inputkey + `?)[] | {(.[0].` + inputkey + `?): .} | to_entries[]]  | reduce .[] as $i ({}; .[$i.key?] = $i.value?)`
	return FilterByQuery(inputJSON, inputQuery)
}

// GetStringValue Excute the jq query on given payload to get the specific value in JSON.
func GetStringValue(payload string, query string) (string, error) {
	jqQuery, err := gojq.Parse(query)
	if err != nil {
		return "", err
	}
	processedPayload, err := Run([]byte(payload), jqQuery)

	if err != nil {
		return "", err
	}
	if len(processedPayload) == 0 {
		return payload, errors.New("error while processing the payload, Processed payload's length is zero")
	}
	return string(processedPayload), nil
}
