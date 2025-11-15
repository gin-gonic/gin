package magic

import (
	"bytes"
	"encoding/binary"
)

// Shp matches a shape format file.
// https://www.esri.com/library/whitepapers/pdfs/shapefile.pdf
func Shp(raw []byte, limit uint32) bool {
	if len(raw) < 112 {
		return false
	}

	if binary.BigEndian.Uint32(raw[0:4]) != 9994 ||
		binary.BigEndian.Uint32(raw[4:8]) != 0 ||
		binary.BigEndian.Uint32(raw[8:12]) != 0 ||
		binary.BigEndian.Uint32(raw[12:16]) != 0 ||
		binary.BigEndian.Uint32(raw[16:20]) != 0 ||
		binary.BigEndian.Uint32(raw[20:24]) != 0 ||
		binary.LittleEndian.Uint32(raw[28:32]) != 1000 {
		return false
	}

	shapeTypes := []int{
		0,  // Null shape
		1,  // Point
		3,  // Polyline
		5,  // Polygon
		8,  // MultiPoint
		11, // PointZ
		13, // PolylineZ
		15, // PolygonZ
		18, // MultiPointZ
		21, // PointM
		23, // PolylineM
		25, // PolygonM
		28, // MultiPointM
		31, // MultiPatch
	}

	for _, st := range shapeTypes {
		if st == int(binary.LittleEndian.Uint32(raw[108:112])) {
			return true
		}
	}

	return false
}

// Shx matches a shape index format file.
// https://www.esri.com/library/whitepapers/pdfs/shapefile.pdf
func Shx(raw []byte, limit uint32) bool {
	return bytes.HasPrefix(raw, []byte{0x00, 0x00, 0x27, 0x0A})
}
