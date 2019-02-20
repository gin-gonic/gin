## Struct level validations

Validations can also be registered at the `struct` level when field level validations
don't make much sense. This can also be used to solve cross-field validation elegantly.
Additionally, it can be combined with tag validations. Struct Level validations run after
the structs tag validations.

### Example requests

```shell
# Validation errors are generated for struct tags as well as at the struct level
$ curl -s -X POST http://localhost:8085/user \
	-H 'content-type: application/json' \
	-d '{}' | jq
{
  "error": "Key: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'User.FirstName' Error:Field validation for 'FirstName' failed on the 'fnameorlname' tag\nKey: 'User.LastName' Error:Field validation for 'LastName' failed on the 'fnameorlname' tag",
  "message": "User validation failed!"
}

# Validation fails at the struct level because neither first name nor last name are present
$ curl -s -X POST http://localhost:8085/user \
    -H 'content-type: application/json' \
	-d '{"email": "george@vandaley.com"}' | jq
{
  "error": "Key: 'User.FirstName' Error:Field validation for 'FirstName' failed on the 'fnameorlname' tag\nKey: 'User.LastName' Error:Field validation for 'LastName' failed on the 'fnameorlname' tag",
  "message": "User validation failed!"
}

# No validation errors when either first name or last name is present
$ curl -X POST http://localhost:8085/user \
    -H 'content-type: application/json' \
	-d '{"fname": "George", "email": "george@vandaley.com"}'
{"message":"User validation successful."}

$ curl -X POST http://localhost:8085/user \
    -H 'content-type: application/json' \
	-d '{"lname": "Contanza", "email": "george@vandaley.com"}'
{"message":"User validation successful."}

$ curl -X POST http://localhost:8085/user \
    -H 'content-type: application/json' \
	-d '{"fname": "George", "lname": "Costanza", "email": "george@vandaley.com"}'
{"message":"User validation successful."}
```

### Useful links

- Validator docs - https://godoc.org/gopkg.in/go-playground/validator.v8#Validate.RegisterStructValidation
- Struct level example - https://github.com/go-playground/validator/blob/v8.18.2/examples/struct-level/struct_level.go
- Validator release notes - https://github.com/go-playground/validator/releases/tag/v8.7
