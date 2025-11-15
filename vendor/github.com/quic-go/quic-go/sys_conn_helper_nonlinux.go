//go:build !linux

package quic

func forceSetReceiveBuffer(c any, bytes int) error { return nil }
func forceSetSendBuffer(c any, bytes int) error    { return nil }

func appendUDPSegmentSizeMsg([]byte, uint16) []byte { return nil }
func isGSOError(error) bool                         { return false }
func isPermissionError(err error) bool              { return false }
