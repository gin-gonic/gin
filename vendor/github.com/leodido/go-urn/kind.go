package urn

type Kind int

const (
	NONE Kind = iota
	RFC2141
	RFC7643
	RFC8141
)
