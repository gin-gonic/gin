package urn

type Option func(Machine)

func WithParsingMode(mode ParsingMode) Option {
	return func(m Machine) {
		m.WithParsingMode(mode)
	}
}
