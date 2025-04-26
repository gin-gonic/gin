package query

import (
	"errors"
	"fmt"
	"strings"
)

// MaxNestedMapDepth is the maximum depth of nesting of single map key in query params.
const MaxNestedMapDepth = 100

type queryKeyType int

const (
	filteredUnsupported queryKeyType = iota
	filteredMap
	filteredRejected
	mapType
	arrayType
	emptyKeyValue
	valueType
)

// GetMap returns a map, which satisfies conditions.
func GetMap(query map[string][]string, filteredKey string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	getAll := filteredKey == ""
	var allErrors = make([]error, 0)
	for key, value := range query {
		kType := getType(key, filteredKey, getAll)
		switch kType {
		case filteredUnsupported:
			allErrors = append(allErrors, fmt.Errorf("invalid access to map %s", key))
			continue
		case filteredMap:
			fallthrough
		case mapType:
			path, err := parsePath(key)
			if err != nil {
				allErrors = append(allErrors, err)
				continue
			}
			if !getAll {
				path = path[1:]
			}
			err = setValueOnPath(result, path, value)
			if err != nil {
				allErrors = append(allErrors, err)
				continue
			}
		case arrayType:
			err := setValueOnPath(result, []string{keyWithoutArraySymbol(key), ""}, value)
			if err != nil {
				allErrors = append(allErrors, err)
				continue
			}
		case filteredRejected:
			continue
		case emptyKeyValue:
			result[key] = value[0]
		case valueType:
			fallthrough
		default:
			err := setValueOnPath(result, []string{key}, value)
			if err != nil {
				allErrors = append(allErrors, err)
				continue
			}
		}
	}
	if len(allErrors) > 0 {
		return nil, errors.Join(allErrors...)
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

// getType is an internal function to get the type of query key.
func getType(key string, filteredKey string, getAll bool) queryKeyType {
	if getAll {
		if isMap(key) {
			return mapType
		}
		if isArray(key) {
			return arrayType
		}
		if key == "" {
			return emptyKeyValue
		}
		return valueType
	}
	if isFilteredKey(key, filteredKey) {
		if isMap(key) {
			return filteredMap
		}
		return filteredUnsupported
	}
	return filteredRejected
}

// isFilteredKey is an internal function to check if k is accepted when searching for map with given key.
func isFilteredKey(k string, filteredKey string) bool {
	return k == filteredKey || strings.HasPrefix(k, filteredKey+"[")
}

// isMap is an internal function to check if k is a map query key.
func isMap(k string) bool {
	i := strings.IndexByte(k, '[')
	j := strings.IndexByte(k, ']')
	return j-i > 1 || (i >= 0 && j == -1)
}

// isArray is an internal function to check if k is an array query key.
func isArray(k string) bool {
	i := strings.IndexByte(k, '[')
	j := strings.IndexByte(k, ']')
	return j-i == 1
}

// keyWithoutArraySymbol is an internal function to remove array symbol from query key.
func keyWithoutArraySymbol(key string) string {
	return key[:len(key)-2]
}

// parsePath is an internal function to parse key access path.
// For example, key[foo][bar] will be parsed to ["foo", "bar"].
func parsePath(k string) ([]string, error) {
	firstKeyEnd := strings.IndexByte(k, '[')
	if firstKeyEnd == -1 {
		return nil, fmt.Errorf("invalid access to map key %s", k)
	}
	first, rawPath := k[:firstKeyEnd], k[firstKeyEnd:]

	split := strings.Split(rawPath, "]")

	// Bear in mind that split of the valid map will always have "" as the last element.
	if split[len(split)-1] != "" {
		return nil, fmt.Errorf("invalid access to map key %s", k)
	}
	if len(split)-1 > MaxNestedMapDepth {
		return nil, fmt.Errorf("maximum depth [%d] of nesting in map exceeded [%d]", MaxNestedMapDepth, len(split)-1)
	}

	// -2 because after split the last element should be empty string.
	last := len(split) - 2

	paths := []string{first}
	for i := 0; i <= last; i++ {
		p := split[i]

		// this way we can handle both error cases: foo] and [foo[bar
		if strings.LastIndex(p, "[") == 0 {
			p = p[1:]
		} else {
			return nil, fmt.Errorf("invalid access to map key %s", p)
		}
		if p == "" && i != last {
			return nil, fmt.Errorf("unsupported array-like access to map key %s", k)
		}
		paths = append(paths, p)
	}
	return paths, nil
}

// setValueOnPath is an internal function to set value a path on dicts.
func setValueOnPath(dicts map[string]interface{}, paths []string, value []string) error {
	nesting := len(paths)
	previousLevel := dicts
	currentLevel := dicts
	for i, p := range paths {
		if isLast(i, nesting) {
			// handle setting value
			if isArrayOnPath(p) {
				previousLevel[paths[i-1]] = value
			} else {
				// if there was already a value set, then it is an error to set a different value.
				if _, ok := currentLevel[p]; ok {
					return fmt.Errorf("trying to set different types at the same key [%s]", p)
				}
				currentLevel[p] = value[0]
			}
		} else {
			// handle subpath of the map
			switch currentLevel[p].(type) {
			case map[string]any:
				// if there was map, and we have to set array, then it is an error
				if isArrayOnPath(paths[i+1]) {
					return fmt.Errorf("trying to set different types at the same key [%s]", p)
				}
			case nil:
				// initialize map if it is not set here yet
				currentLevel[p] = make(map[string]any)
			default:
				// if there was different value then a map, then it is an error to set a map here.
				return fmt.Errorf("trying to set different types at the same key [%s]", p)
			}
			previousLevel = currentLevel
			currentLevel = currentLevel[p].(map[string]any)
		}
	}
	return nil
}

// isArrayOnPath is an internal function to check if the current parsed map path is an array.
func isArrayOnPath(p string) bool {
	return p == ""
}

// isLast is an internal function to check if the current level is the last level.
func isLast(i int, nesting int) bool {
	return i == nesting-1
}
