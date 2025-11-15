package parser

// Option represents parser's option.
type Option func(p *parser)

// AllowDuplicateMapKey allow the use of keys with the same name in the same map,
// but by default, this is not permitted.
func AllowDuplicateMapKey() Option {
	return func(p *parser) {
		p.allowDuplicateMapKey = true
	}
}
