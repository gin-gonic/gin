package binding

import (
	"errors"
	"fmt"
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
					fmt.Println("email...")
					if !reflect.DeepEqual(zero, fieldValue) {
						if err := email(fieldValue); err != nil {
							return err
						}
					}
				case strings.Contains(match, "digit:"):
					if !reflect.DeepEqual(zero, fieldValue) {
						if err := digit(match, fieldValue); err != nil {
							return err
						}
					}
				case strings.Contains(match, "digits_between:"):
					if !reflect.DeepEqual(zero, fieldValue) {
						if err := digits_between(match, fieldValue); err != nil {
							return err
						}
					}
				case strings.Contains(match, "min:"):
					if !reflect.DeepEqual(zero, fieldValue) {
						if err := min(match, fieldValue); err != nil {
							return err
						}
					}
				case strings.Contains(match, "max:"):
					if !reflect.DeepEqual(zero, fieldValue) {
						if err := max(match, fieldValue); err != nil {
							return err
						}
					}
				case strings.Contains(match, "in:"):
					if !reflect.DeepEqual(zero, fieldValue) {
						if err := in(match, fieldValue); err != nil {
							return err
						}
					}
				case strings.Contains(match, "regex:"):
					if !reflect.DeepEqual(zero, fieldValue) {
						if err := regex(match, fieldValue); err != nil {
							return err
						}
					}
				default:
					// Temp logging to check for errors
					errors.New(array[setting] + " is not a valid validation type.")
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
func email(value interface{}) error {
	// Email Regex Checker
	var emailRegex string = `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`

	if data, ok := value.(string); ok {
		if match, _ := regexp.Match(emailRegex, []byte(data)); match {
			return nil
		} else {
			return errors.New("A valid email address was not entered.")
		}
	} else {
		return errors.New("Email was not able to convert the passed in data to a []byte.")
	}
}

// Check that the passed in field is a valid email
// Need to improve error logging for this method
// Currently only supports strings, ints
func in(field string, value interface{}) error {

	if data, ok := value.(string); ok {

		valid := strings.Split(field[3:], ",")

		for option := range valid {
			if valid[option] == data {
				return nil
			}
		}

		return errors.New("In did not match any of the expected values.")

	} else if data, ok := value.(int); ok {
		// This will run with passed in data is an int
		valid := strings.Split(field[3:], ",")

		for option := range valid {
			// Check for convertion to valid int
			if valint, err := strconv.ParseInt(valid[option], 0, 64); err == nil {
				if valint == int64(data) {
					return nil
				}
			}
		}

		return errors.New("In did not match any of the expected values.")

	} else {
		return errors.New("in, was not able to convert the data passed in to a string.")
	}

}

// Check that the passed in field is exactly X digits
func digit(field string, value interface{}) error {

	if data, ok := value.(int); ok {
		// Unpack number of digits it should be.
		digit := field[6:]

		if digits, check := strconv.ParseInt(digit, 0, 64); check == nil {

			if int64(len(strconv.FormatInt(int64(data), 10))) == digits {
				return nil
			} else {
				return errors.New("The data you passed in was not the right number of digits.")
			}

		} else {
			return errors.New("Digit must check for a number.")
		}
	}

	return errors.New("The number passed into digit was not an int.")
}

func digits_between(field string, value interface{}) error {

	if data, ok := value.(int); ok {

		digit := strings.Split(field[15:], ",")

		if digitSmall, ok := strconv.ParseInt(digit[0], 0, 64); ok == nil {

			if digitLarge, okk := strconv.ParseInt(digit[1], 0, 64); okk == nil {

				num := int64(len(strconv.FormatInt(int64(data), 10)))

				if num >= digitSmall && num <= digitLarge {
					return nil
				} else {
					return errors.New("The data you passed in was not the right number of digits.")
				}
			}
		}
	}

	return errors.New("The value passed into digits_between could not be converted to an int.")
}

func min(field string, value interface{}) error {

	if data, ok := value.(int); ok {

		min := field[4:]

		if minNum, ok := strconv.ParseInt(min, 0, 64); ok == nil {

			if int64(data) >= minNum {
				return nil
			} else {
				return errors.New("The data you passed in was smaller then the allowed minimum.")
			}

		}
	}

	return errors.New("The value passed into min could not be converted to an int.")
}

func max(field string, value interface{}) error {

	if data, ok := value.(int); ok {

		max := field[4:]

		if maxNum, ok := strconv.ParseInt(max, 0, 64); ok == nil {
			if int64(data) <= maxNum {
				return nil
			} else {
				return errors.New("The data you passed in was larger than the maximum.")
			}

		}
	}

	return errors.New("The value passed into max could not be converted to an int.")
}

func regex(field string, value interface{}) error {
	// Email Regex Checker

	reg := field[6:]

	if data, ok := value.(string); ok {
		if match, err := regexp.Match(reg, []byte(data)); err == nil && match {
			return nil
		} else {
			return errors.New("Your regex did not match or was not valid.")
		}
	} else {
		return errors.New("Regex was not able to convert the passed in data to a string.")
	}
}
