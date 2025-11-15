package urn

type ParsingMode int

const (
	Default ParsingMode = iota
	RFC2141Only
	RFC7643Only
	RFC8141Only
)

const DefaultParsingMode = RFC2141Only
