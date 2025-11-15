package mimetype

import (
	"sync"

	"github.com/gabriel-vasile/mimetype/internal/magic"
)

// mimetype stores the list of MIME types in a tree structure with
// "application/octet-stream" at the root of the hierarchy. The hierarchy
// approach minimizes the number of checks that need to be done on the input
// and allows for more precise results once the base type of file has been
// identified.
//
// root is a detector which passes for any slice of bytes.
// When a detector passes the check, the children detectors
// are tried in order to find a more accurate MIME type.
var root = newMIME("application/octet-stream", "",
	func([]byte, uint32) bool { return true },
	xpm, sevenZ, zip, pdf, fdf, ole, ps, psd, p7s, ogg, png, jpg, jxl, jp2, jpx,
	jpm, jxs, gif, webp, exe, elf, ar, tar, xar, bz2, fits, tiff, bmp, lotus, ico,
	mp3, flac, midi, ape, musePack, amr, wav, aiff, au, mpeg, quickTime, mp4, webM,
	avi, flv, mkv, asf, aac, voc, m3u, rmvb, gzip, class, swf, crx, ttf, woff,
	woff2, otf, ttc, eot, wasm, shx, dbf, dcm, rar, djvu, mobi, lit, bpg, cbor,
	sqlite3, dwg, nes, lnk, macho, qcp, icns, hdr, mrc, mdb, accdb, zstd, cab,
	rpm, xz, lzip, torrent, cpio, tzif, xcf, pat, gbr, glb, cabIS, jxr, parquet,
	oneNote, chm,
	// Keep text last because it is the slowest check.
	text,
)

// errMIME is returned from Detect functions when err is not nil.
// Detect could return root for erroneous cases, but it needs to lock mu in order to do so.
// errMIME is same as root but it does not require locking.
var errMIME = newMIME("application/octet-stream", "", func([]byte, uint32) bool { return false })

// mu guards access to the root MIME tree. Access to root must be synchronized with this lock.
var mu = &sync.RWMutex{}

// The list of nodes appended to the root node.
var (
	xz   = newMIME("application/x-xz", ".xz", magic.Xz)
	gzip = newMIME("application/gzip", ".gz", magic.Gzip).alias(
		"application/x-gzip", "application/x-gunzip", "application/gzipped",
		"application/gzip-compressed", "application/x-gzip-compressed",
		"gzip/document")
	sevenZ = newMIME("application/x-7z-compressed", ".7z", magic.SevenZ)
	// APK must be checked before JAR because APK is a subset of JAR.
	// This means APK should be a child of JAR detector, but in practice,
	// the decisive signature for JAR might be located at the end of the file
	// and not reachable because of library readLimit.
	zip = newMIME("application/zip", ".zip", magic.Zip, docx, pptx, xlsx, epub, apk, jar, odt, ods, odp, odg, odf, odc, sxc, kmz, visio).
		alias("application/x-zip", "application/x-zip-compressed")
	tar = newMIME("application/x-tar", ".tar", magic.Tar)
	xar = newMIME("application/x-xar", ".xar", magic.Xar)
	bz2 = newMIME("application/x-bzip2", ".bz2", magic.Bz2)
	pdf = newMIME("application/pdf", ".pdf", magic.PDF).
		alias("application/x-pdf")
	fdf   = newMIME("application/vnd.fdf", ".fdf", magic.Fdf)
	xlsx  = newMIME("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", ".xlsx", magic.Xlsx)
	docx  = newMIME("application/vnd.openxmlformats-officedocument.wordprocessingml.document", ".docx", magic.Docx)
	pptx  = newMIME("application/vnd.openxmlformats-officedocument.presentationml.presentation", ".pptx", magic.Pptx)
	visio = newMIME("application/vnd.ms-visio.drawing.main+xml", ".vsdx", magic.Visio)
	epub  = newMIME("application/epub+zip", ".epub", magic.Epub)
	jar   = newMIME("application/java-archive", ".jar", magic.Jar).
		alias("application/jar", "application/jar-archive", "application/x-java-archive")
	apk = newMIME("application/vnd.android.package-archive", ".apk", magic.APK)
	ole = newMIME("application/x-ole-storage", "", magic.Ole, msi, aaf, msg, xls, pub, ppt, doc)
	msi = newMIME("application/x-ms-installer", ".msi", magic.Msi).
		alias("application/x-windows-installer", "application/x-msi")
	aaf = newMIME("application/octet-stream", ".aaf", magic.Aaf)
	doc = newMIME("application/msword", ".doc", magic.Doc).
		alias("application/vnd.ms-word")
	ppt = newMIME("application/vnd.ms-powerpoint", ".ppt", magic.Ppt).
		alias("application/mspowerpoint")
	pub = newMIME("application/vnd.ms-publisher", ".pub", magic.Pub)
	xls = newMIME("application/vnd.ms-excel", ".xls", magic.Xls).
		alias("application/msexcel")
	msg  = newMIME("application/vnd.ms-outlook", ".msg", magic.Msg)
	ps   = newMIME("application/postscript", ".ps", magic.Ps)
	fits = newMIME("application/fits", ".fits", magic.Fits).alias("image/fits")
	ogg  = newMIME("application/ogg", ".ogg", magic.Ogg, oggAudio, oggVideo).
		alias("application/x-ogg")
	oggAudio = newMIME("audio/ogg", ".oga", magic.OggAudio)
	oggVideo = newMIME("video/ogg", ".ogv", magic.OggVideo)
	text     = newMIME("text/plain", ".txt", magic.Text, svg, html, xml, php, js, lua, perl, python, ruby, json, ndJSON, rtf, srt, tcl, csv, tsv, vCard, iCalendar, warc, vtt, shell, netpbm, netpgm, netppm, netpam)
	xml      = newMIME("text/xml", ".xml", magic.XML, rss, atom, x3d, kml, xliff, collada, gml, gpx, tcx, amf, threemf, xfdf, owl2, xhtml).
			alias("application/xml")
	xhtml   = newMIME("application/xhtml+xml", ".html", magic.XHTML)
	json    = newMIME("application/json", ".json", magic.JSON, geoJSON, har, gltf)
	har     = newMIME("application/json", ".har", magic.HAR)
	csv     = newMIME("text/csv", ".csv", magic.CSV)
	tsv     = newMIME("text/tab-separated-values", ".tsv", magic.TSV)
	geoJSON = newMIME("application/geo+json", ".geojson", magic.GeoJSON)
	ndJSON  = newMIME("application/x-ndjson", ".ndjson", magic.NdJSON)
	html    = newMIME("text/html", ".html", magic.HTML)
	php     = newMIME("text/x-php", ".php", magic.Php)
	rtf     = newMIME("text/rtf", ".rtf", magic.Rtf).alias("application/rtf")
	js      = newMIME("text/javascript", ".js", magic.Js).
		alias("application/x-javascript", "application/javascript")
	srt = newMIME("application/x-subrip", ".srt", magic.Srt).
		alias("application/x-srt", "text/x-srt")
	vtt    = newMIME("text/vtt", ".vtt", magic.Vtt)
	lua    = newMIME("text/x-lua", ".lua", magic.Lua)
	perl   = newMIME("text/x-perl", ".pl", magic.Perl)
	python = newMIME("text/x-python", ".py", magic.Python).
		alias("text/x-script.python", "application/x-python")
	ruby = newMIME("text/x-ruby", ".rb", magic.Ruby).
		alias("application/x-ruby")
	shell = newMIME("text/x-shellscript", ".sh", magic.Shell).
		alias("text/x-sh", "application/x-shellscript", "application/x-sh")
	tcl = newMIME("text/x-tcl", ".tcl", magic.Tcl).
		alias("application/x-tcl")
	vCard     = newMIME("text/vcard", ".vcf", magic.VCard)
	iCalendar = newMIME("text/calendar", ".ics", magic.ICalendar)
	svg       = newMIME("image/svg+xml", ".svg", magic.Svg)
	rss       = newMIME("application/rss+xml", ".rss", magic.Rss).
			alias("text/rss")
	owl2    = newMIME("application/owl+xml", ".owl", magic.Owl2)
	atom    = newMIME("application/atom+xml", ".atom", magic.Atom)
	x3d     = newMIME("model/x3d+xml", ".x3d", magic.X3d)
	kml     = newMIME("application/vnd.google-earth.kml+xml", ".kml", magic.Kml)
	kmz     = newMIME("application/vnd.google-earth.kmz", ".kmz", magic.KMZ)
	xliff   = newMIME("application/x-xliff+xml", ".xlf", magic.Xliff)
	collada = newMIME("model/vnd.collada+xml", ".dae", magic.Collada)
	gml     = newMIME("application/gml+xml", ".gml", magic.Gml)
	gpx     = newMIME("application/gpx+xml", ".gpx", magic.Gpx)
	tcx     = newMIME("application/vnd.garmin.tcx+xml", ".tcx", magic.Tcx)
	amf     = newMIME("application/x-amf", ".amf", magic.Amf)
	threemf = newMIME("application/vnd.ms-package.3dmanufacturing-3dmodel+xml", ".3mf", magic.Threemf)
	png     = newMIME("image/png", ".png", magic.Png, apng)
	apng    = newMIME("image/vnd.mozilla.apng", ".png", magic.Apng)
	jpg     = newMIME("image/jpeg", ".jpg", magic.Jpg)
	jxl     = newMIME("image/jxl", ".jxl", magic.Jxl)
	jp2     = newMIME("image/jp2", ".jp2", magic.Jp2)
	jpx     = newMIME("image/jpx", ".jpf", magic.Jpx)
	jpm     = newMIME("image/jpm", ".jpm", magic.Jpm).
		alias("video/jpm")
	jxs  = newMIME("image/jxs", ".jxs", magic.Jxs)
	xpm  = newMIME("image/x-xpixmap", ".xpm", magic.Xpm)
	bpg  = newMIME("image/bpg", ".bpg", magic.Bpg)
	gif  = newMIME("image/gif", ".gif", magic.Gif)
	webp = newMIME("image/webp", ".webp", magic.Webp)
	tiff = newMIME("image/tiff", ".tiff", magic.Tiff)
	bmp  = newMIME("image/bmp", ".bmp", magic.Bmp).
		alias("image/x-bmp", "image/x-ms-bmp")
	// lotus check must be done before ico because some ico detection is a bit
	// relaxed and some lotus files are wrongfully identified as ico otherwise.
	lotus = newMIME("application/vnd.lotus-1-2-3", ".123", magic.Lotus123)
	ico   = newMIME("image/x-icon", ".ico", magic.Ico)
	icns  = newMIME("image/x-icns", ".icns", magic.Icns)
	psd   = newMIME("image/vnd.adobe.photoshop", ".psd", magic.Psd).
		alias("image/x-psd", "application/photoshop")
	heic    = newMIME("image/heic", ".heic", magic.Heic)
	heicSeq = newMIME("image/heic-sequence", ".heic", magic.HeicSequence)
	heif    = newMIME("image/heif", ".heif", magic.Heif)
	heifSeq = newMIME("image/heif-sequence", ".heif", magic.HeifSequence)
	hdr     = newMIME("image/vnd.radiance", ".hdr", magic.Hdr)
	avif    = newMIME("image/avif", ".avif", magic.AVIF)
	mp3     = newMIME("audio/mpeg", ".mp3", magic.Mp3).
		alias("audio/x-mpeg", "audio/mp3")
	flac = newMIME("audio/flac", ".flac", magic.Flac)
	midi = newMIME("audio/midi", ".midi", magic.Midi).
		alias("audio/mid", "audio/sp-midi", "audio/x-mid", "audio/x-midi")
	ape      = newMIME("audio/ape", ".ape", magic.Ape)
	musePack = newMIME("audio/musepack", ".mpc", magic.MusePack)
	wav      = newMIME("audio/wav", ".wav", magic.Wav).
			alias("audio/x-wav", "audio/vnd.wave", "audio/wave")
	aiff = newMIME("audio/aiff", ".aiff", magic.Aiff).alias("audio/x-aiff")
	au   = newMIME("audio/basic", ".au", magic.Au)
	amr  = newMIME("audio/amr", ".amr", magic.Amr).
		alias("audio/amr-nb")
	aac  = newMIME("audio/aac", ".aac", magic.AAC)
	voc  = newMIME("audio/x-unknown", ".voc", magic.Voc)
	aMp4 = newMIME("audio/mp4", ".mp4", magic.AMp4).
		alias("audio/x-mp4a")
	m4a = newMIME("audio/x-m4a", ".m4a", magic.M4a)
	m3u = newMIME("application/vnd.apple.mpegurl", ".m3u", magic.M3u).
		alias("audio/mpegurl")
	m4v  = newMIME("video/x-m4v", ".m4v", magic.M4v)
	mj2  = newMIME("video/mj2", ".mj2", magic.Mj2)
	dvb  = newMIME("video/vnd.dvb.file", ".dvb", magic.Dvb)
	mp4  = newMIME("video/mp4", ".mp4", magic.Mp4, avif, threeGP, threeG2, aMp4, mqv, m4a, m4v, heic, heicSeq, heif, heifSeq, mj2, dvb)
	webM = newMIME("video/webm", ".webm", magic.WebM).
		alias("audio/webm")
	mpeg      = newMIME("video/mpeg", ".mpeg", magic.Mpeg)
	quickTime = newMIME("video/quicktime", ".mov", magic.QuickTime)
	mqv       = newMIME("video/quicktime", ".mqv", magic.Mqv)
	threeGP   = newMIME("video/3gpp", ".3gp", magic.ThreeGP).
			alias("video/3gp", "audio/3gpp")
	threeG2 = newMIME("video/3gpp2", ".3g2", magic.ThreeG2).
		alias("video/3g2", "audio/3gpp2")
	avi = newMIME("video/x-msvideo", ".avi", magic.Avi).
		alias("video/avi", "video/msvideo")
	flv = newMIME("video/x-flv", ".flv", magic.Flv)
	mkv = newMIME("video/x-matroska", ".mkv", magic.Mkv)
	asf = newMIME("video/x-ms-asf", ".asf", magic.Asf).
		alias("video/asf", "video/x-ms-wmv")
	rmvb  = newMIME("application/vnd.rn-realmedia-vbr", ".rmvb", magic.Rmvb)
	class = newMIME("application/x-java-applet", ".class", magic.Class)
	swf   = newMIME("application/x-shockwave-flash", ".swf", magic.SWF)
	crx   = newMIME("application/x-chrome-extension", ".crx", magic.CRX)
	ttf   = newMIME("font/ttf", ".ttf", magic.Ttf).
		alias("font/sfnt", "application/x-font-ttf", "application/font-sfnt")
	woff    = newMIME("font/woff", ".woff", magic.Woff)
	woff2   = newMIME("font/woff2", ".woff2", magic.Woff2)
	otf     = newMIME("font/otf", ".otf", magic.Otf)
	ttc     = newMIME("font/collection", ".ttc", magic.Ttc)
	eot     = newMIME("application/vnd.ms-fontobject", ".eot", magic.Eot)
	wasm    = newMIME("application/wasm", ".wasm", magic.Wasm)
	shp     = newMIME("application/vnd.shp", ".shp", magic.Shp)
	shx     = newMIME("application/vnd.shx", ".shx", magic.Shx, shp)
	dbf     = newMIME("application/x-dbf", ".dbf", magic.Dbf)
	exe     = newMIME("application/vnd.microsoft.portable-executable", ".exe", magic.Exe)
	elf     = newMIME("application/x-elf", "", magic.Elf, elfObj, elfExe, elfLib, elfDump)
	elfObj  = newMIME("application/x-object", "", magic.ElfObj)
	elfExe  = newMIME("application/x-executable", "", magic.ElfExe)
	elfLib  = newMIME("application/x-sharedlib", ".so", magic.ElfLib)
	elfDump = newMIME("application/x-coredump", "", magic.ElfDump)
	ar      = newMIME("application/x-archive", ".a", magic.Ar, deb).
		alias("application/x-unix-archive")
	deb = newMIME("application/vnd.debian.binary-package", ".deb", magic.Deb)
	rpm = newMIME("application/x-rpm", ".rpm", magic.RPM)
	dcm = newMIME("application/dicom", ".dcm", magic.Dcm)
	odt = newMIME("application/vnd.oasis.opendocument.text", ".odt", magic.Odt, ott).
		alias("application/x-vnd.oasis.opendocument.text")
	ott = newMIME("application/vnd.oasis.opendocument.text-template", ".ott", magic.Ott).
		alias("application/x-vnd.oasis.opendocument.text-template")
	ods = newMIME("application/vnd.oasis.opendocument.spreadsheet", ".ods", magic.Ods, ots).
		alias("application/x-vnd.oasis.opendocument.spreadsheet")
	ots = newMIME("application/vnd.oasis.opendocument.spreadsheet-template", ".ots", magic.Ots).
		alias("application/x-vnd.oasis.opendocument.spreadsheet-template")
	odp = newMIME("application/vnd.oasis.opendocument.presentation", ".odp", magic.Odp, otp).
		alias("application/x-vnd.oasis.opendocument.presentation")
	otp = newMIME("application/vnd.oasis.opendocument.presentation-template", ".otp", magic.Otp).
		alias("application/x-vnd.oasis.opendocument.presentation-template")
	odg = newMIME("application/vnd.oasis.opendocument.graphics", ".odg", magic.Odg, otg).
		alias("application/x-vnd.oasis.opendocument.graphics")
	otg = newMIME("application/vnd.oasis.opendocument.graphics-template", ".otg", magic.Otg).
		alias("application/x-vnd.oasis.opendocument.graphics-template")
	odf = newMIME("application/vnd.oasis.opendocument.formula", ".odf", magic.Odf).
		alias("application/x-vnd.oasis.opendocument.formula")
	odc = newMIME("application/vnd.oasis.opendocument.chart", ".odc", magic.Odc).
		alias("application/x-vnd.oasis.opendocument.chart")
	sxc = newMIME("application/vnd.sun.xml.calc", ".sxc", magic.Sxc)
	rar = newMIME("application/x-rar-compressed", ".rar", magic.RAR).
		alias("application/x-rar")
	djvu    = newMIME("image/vnd.djvu", ".djvu", magic.DjVu)
	mobi    = newMIME("application/x-mobipocket-ebook", ".mobi", magic.Mobi)
	lit     = newMIME("application/x-ms-reader", ".lit", magic.Lit)
	sqlite3 = newMIME("application/vnd.sqlite3", ".sqlite", magic.Sqlite).
		alias("application/x-sqlite3")
	dwg = newMIME("image/vnd.dwg", ".dwg", magic.Dwg).
		alias("image/x-dwg", "application/acad", "application/x-acad",
			"application/autocad_dwg", "application/dwg", "application/x-dwg",
			"application/x-autocad", "drawing/dwg")
	warc    = newMIME("application/warc", ".warc", magic.Warc)
	nes     = newMIME("application/vnd.nintendo.snes.rom", ".nes", magic.Nes)
	lnk     = newMIME("application/x-ms-shortcut", ".lnk", magic.Lnk)
	macho   = newMIME("application/x-mach-binary", ".macho", magic.MachO)
	qcp     = newMIME("audio/qcelp", ".qcp", magic.Qcp)
	mrc     = newMIME("application/marc", ".mrc", magic.Marc)
	mdb     = newMIME("application/x-msaccess", ".mdb", magic.MsAccessMdb)
	accdb   = newMIME("application/x-msaccess", ".accdb", magic.MsAccessAce)
	zstd    = newMIME("application/zstd", ".zst", magic.Zstd)
	cab     = newMIME("application/vnd.ms-cab-compressed", ".cab", magic.Cab)
	cabIS   = newMIME("application/x-installshield", ".cab", magic.InstallShieldCab)
	lzip    = newMIME("application/lzip", ".lz", magic.Lzip).alias("application/x-lzip")
	torrent = newMIME("application/x-bittorrent", ".torrent", magic.Torrent)
	cpio    = newMIME("application/x-cpio", ".cpio", magic.Cpio)
	tzif    = newMIME("application/tzif", "", magic.TzIf)
	p7s     = newMIME("application/pkcs7-signature", ".p7s", magic.P7s)
	xcf     = newMIME("image/x-xcf", ".xcf", magic.Xcf)
	pat     = newMIME("image/x-gimp-pat", ".pat", magic.Pat)
	gbr     = newMIME("image/x-gimp-gbr", ".gbr", magic.Gbr)
	xfdf    = newMIME("application/vnd.adobe.xfdf", ".xfdf", magic.Xfdf)
	glb     = newMIME("model/gltf-binary", ".glb", magic.GLB)
	gltf    = newMIME("model/gltf+json", ".gltf", magic.GLTF)
	jxr     = newMIME("image/jxr", ".jxr", magic.Jxr).alias("image/vnd.ms-photo")
	parquet = newMIME("application/vnd.apache.parquet", ".parquet", magic.Par1).
		alias("application/x-parquet")
	netpbm  = newMIME("image/x-portable-bitmap", ".pbm", magic.NetPBM)
	netpgm  = newMIME("image/x-portable-graymap", ".pgm", magic.NetPGM)
	netppm  = newMIME("image/x-portable-pixmap", ".ppm", magic.NetPPM)
	netpam  = newMIME("image/x-portable-arbitrarymap", ".pam", magic.NetPAM)
	cbor    = newMIME("application/cbor", ".cbor", magic.CBOR)
	oneNote = newMIME("application/onenote", ".one", magic.One)
	chm     = newMIME("application/vnd.ms-htmlhelp", ".chm", magic.CHM)
)
