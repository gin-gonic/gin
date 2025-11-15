package magic

var (
	// Sqlite matches an SQLite database file.
	Sqlite = prefix([]byte{
		0x53, 0x51, 0x4c, 0x69, 0x74, 0x65, 0x20, 0x66,
		0x6f, 0x72, 0x6d, 0x61, 0x74, 0x20, 0x33, 0x00,
	})
	// MsAccessAce matches Microsoft Access dababase file.
	MsAccessAce = offset([]byte("Standard ACE DB"), 4)
	// MsAccessMdb matches legacy Microsoft Access database file (JET, 2003 and earlier).
	MsAccessMdb = offset([]byte("Standard Jet DB"), 4)
)
