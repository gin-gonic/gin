package binding

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Validate(obj interface{}) error {

	typ := reflect.TypeOf(obj)
	value := reflect.ValueOf(obj)

	// Check to ensure we are getting a valid
	// pointer for manipulation.
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		value = value.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {

		field := typ.Field(i)
		fieldValue := value.Field(i).Interface()
		zero := reflect.Zero(field.Type).Interface()

		// Validate nested and embedded structs (if pointer, only do so if not nil)
		if field.Type.Kind() == reflect.Struct ||
			(field.Type.Kind() == reflect.Ptr && !reflect.DeepEqual(zero, fieldValue)) {
			if err := Validate(fieldValue); err != nil {
				return err
			}
		}

		if field.Tag.Get("validate") != "" || field.Tag.Get("binding") != "" {
			// Break validate field into array
			array := strings.Split(field.Tag.Get("validate"), "|")

			// Legacy Support for binding.
			if array[0] == "" {
				array = strings.Split(field.Tag.Get("binding"), "|")
			}

			// Do the hard work of checking all assertions
			for setting := range array {

				match := array[setting]

				switch {
				case "required" == match:
					if err := required(field, fieldValue, zero); err != nil {
						return err
					}
				case "email" == match:
					if err := regex(`regex:^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`, fieldValue); err != nil {
						return err
					}
				case strings.HasPrefix(match, "min:"):
					if err := min(match, fieldValue); err != nil {
						return err
					}
				case strings.HasPrefix(match, "max:"):
					if err := max(match, fieldValue); err != nil {
						return err
					}
				case strings.HasPrefix(match, "in:"):
					if err := in(match, fieldValue); err != nil {
						return err
					}
				case strings.HasPrefix(match, "regex:"):
					if err := regex(match, fieldValue); err != nil {
						return err
					}
				default:
					panic("The field " + match + " is not a valid validation check.")
				}
			}
		}
	}

	return nil
}

// Check that the following function features
// the required field. May need to check for
// more special cases like since passing in null
// is the same as 0 for int type checking.
func required(field reflect.StructField, value, zero interface{}) error {

	if reflect.DeepEqual(zero, value) {
		if _, ok := value.(int); !ok {
			return errors.New("The required field " + field.Name + " was not submitted.")
		}
	}

	return nil
}

// Check that the passed in field is a valid email
// Need to improve error logging for this method
// Currently only supports strings, ints
func in(field string, value interface{}) error {

	if data, ok := value.(string); ok {
		if len(data) == 0 {
			return nil
		}

		valid := strings.Split(field[3:], ",")

		for option := range valid {
			if valid[option] == data {
				return nil
			}
		}

	} else {
		return errors.New("The value passed in for IN could not be converted to a string.")
	}

	return errors.New("In did not match any of the expected values.")
}

func min(field string, value interface{}) error {

	if data, ok := value.(int); ok {

		min := field[strings.Index(field, ":")+1:]

		if minNum, ok := strconv.ParseInt(min, 0, 64); ok == nil {

			if int64(data) >= minNum {
				return nil
			} else {
				return errors.New("The data you passed in was smaller then the allowed minimum.")
			}

		}
	}

	return errors.New("The value passed in for MIN could not be converted to an int.")
}

func max(field string, value interface{}) error {

	if data, ok := value.(int); ok {

		max := field[strings.Index(field, ":")+1:]

		if maxNum, ok := strconv.ParseInt(max, 0, 64); ok == nil {
			if int64(data) <= maxNum {
				return nil
			} else {
				return errors.New("The data you passed in was larger than the maximum.")
			}

		}
	}

	return errors.New("The value passed in for MAX could not be converted to an int.")
}

// Regex handles the general regex call and also handles
// the regex email.
func regex(field string, value interface{}) error {

	reg := field[strings.Index(field, ":")+1:]

	if data, ok := value.(string); ok {
		if len(data) == 0 {
			return nil
		} else if err := match_regex(reg, []byte(data)); err != nil {
			return err
		}
	} else if data, ok := value.(int); ok {
		if err := match_regex(reg, []byte(strconv.Itoa(data))); err != nil {
			return err
		}
	} else {
		return errors.New("The value passed in for REGEX could not be converted to a string or int.")
	}

	return nil
}

// Helper function for regex.
func match_regex(reg string, data []byte) error {

	if match, err := regexp.Match(reg, []byte(data)); err == nil && match {
		return nil
	} else {
		return errors.New("Your regex did not match or was not valid.")
	}
}
