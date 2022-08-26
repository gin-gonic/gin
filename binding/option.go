package binding

type options struct {
	params map[string]interface{}
	query  map[string]interface{}
}

type Option func(opt *options) error

func WithParams(params map[string]interface{}) Option {
	return func(opt *options) error {
		opt.params = params
		return nil
	}
}
