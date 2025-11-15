package magic

import (
	"bytes"
)

/*
 NOTE:

 In May 2003, two Internet RFCs were published relating to the format.
 The Ogg bitstream was defined in RFC 3533 (which is classified as
 'informative') and its Internet content type (application/ogg) in RFC
 3534 (which is, as of 2006, a proposed standard protocol). In
 September 2008, RFC 3534 was obsoleted by RFC 5334, which added
 content types video/ogg, audio/ogg and filename extensions .ogx, .ogv,
 .oga, .spx.

 See:
 https://tools.ietf.org/html/rfc3533
 https://developer.mozilla.org/en-US/docs/Web/HTTP/Configuring_servers_for_Ogg_media#Serve_media_with_the_correct_MIME_type
 https://github.com/file/file/blob/master/magic/Magdir/vorbis
*/

// Ogg matches an Ogg file.
func Ogg(raw []byte, limit uint32) bool {
	return bytes.HasPrefix(raw, []byte("\x4F\x67\x67\x53\x00"))
}

// OggAudio matches an audio ogg file.
func OggAudio(raw []byte, limit uint32) bool {
	return len(raw) >= 37 && (bytes.HasPrefix(raw[28:], []byte("\x7fFLAC")) ||
		bytes.HasPrefix(raw[28:], []byte("\x01vorbis")) ||
		bytes.HasPrefix(raw[28:], []byte("OpusHead")) ||
		bytes.HasPrefix(raw[28:], []byte("Speex\x20\x20\x20")))
}

// OggVideo matches a video ogg file.
func OggVideo(raw []byte, limit uint32) bool {
	return len(raw) >= 37 && (bytes.HasPrefix(raw[28:], []byte("\x80theora")) ||
		bytes.HasPrefix(raw[28:], []byte("fishead\x00")) ||
		bytes.HasPrefix(raw[28:], []byte("\x01video\x00\x00\x00"))) // OGM video
}
