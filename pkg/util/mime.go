package util

type Mime string

const (
	// MimeTextPlain is the MIME type for plain text
	MimeTextPlain Mime = "text/plain"
	// MimeTextHTML is the MIME type for HTML
	MimeTextHTML Mime = "text/html"
	// MimeTextMarkdown is the MIME type for markdown
	MimeTextMarkdown Mime = "text/markdown"
	// MimeImageJPEG is the MIME type for JPEG images
	MimeImageJPEG Mime = "image/jpeg"
	// MimeImagePNG is the MIME type for PNG images
	MimeImagePNG Mime = "image/png"
	// MimeImageGIF is the MIME type for GIF images
	MimeImageGIF Mime = "image/gif"
	// MimeImageBMP is the MIME type for BMP images
	MimeImageBMP Mime = "image/bmp"
	// MimeImageSVG is the MIME type for SVG images
	MimeImageSVG Mime = "image/svg+xml"
	// MimeAudioMPEG is the MIME type for MPEG audio
	MimeAudioMPEG Mime = "audio/mpeg"
	// MimeAudioOGG is the MIME type for OGG audio
	MimeAudioOGG Mime = "audio/ogg"
	MimeAudioOGA Mime = "audio/oga"
	// MimeAudioWAV is the MIME type for WAV audio
	MimeAudioWAV Mime = "audio/wav"
	// MimeVideoMPEG is the MIME type for MPEG video
	MimeVideoMPEG Mime = "video/mpeg"
	// MimeVideoMP4 is the MIME type for MP4 video
	MimeVideoMP4 Mime = "video/mp4"
	// MimeVideoOGG is the MIME type for OGG video
	MimeVideoOGG Mime = "video/ogg"
	// MimeVideoQuickTime is the MIME type for QuickTime video
	MimeVideoQuickTime Mime = "video/quicktime"
	// MimeApplicationJSON is the MIME type for JSON
	MimeApplicationJSON Mime = "application/json"
	// MimeApplicationXML is the MIME type for XML
	MimeApplicationXML Mime = "application/xml"
	// MimeApplicationPDF is the MIME type for PDF
	MimeApplicationPDF Mime = "application/pdf"
	// MimeApplicationZIP is the MIME type for ZIP archives
	MimeApplicationZIP Mime = "application/zip"
	// MimeApplicationGZIP is the MIME type for GZIP archives
	MimeApplicationGZIP Mime = "application/gzip"
	// MimeApplicationTar is the MIME type for TAR archives
	MimeApplicationTar Mime = "application/tar"
	// MimeApplicationXZ is the MIME type for XZ archives
	MimeApplicationXZ Mime = "application/x-xz"
	MimeOctetStream   Mime = "application/octet-stream"
)

var MimeArray = []Mime{
	MimeTextPlain,
	MimeTextHTML,
	MimeTextMarkdown,
	MimeImageJPEG,
	MimeImagePNG,
	MimeImageGIF,
	MimeImageBMP,
	MimeImageSVG,
	MimeAudioMPEG,
	MimeAudioOGG,
	MimeAudioOGA,
	MimeAudioWAV,
	MimeVideoMPEG,
	MimeVideoMP4,
	MimeVideoOGG,
	MimeVideoQuickTime,
	MimeApplicationJSON,
	MimeApplicationXML,
	MimeApplicationPDF,
	MimeApplicationZIP,
	MimeApplicationGZIP,
	MimeApplicationTar,
	MimeApplicationXZ,
	MimeOctetStream,
}

type MimeExtension string

var MimeExtensionMap = map[Mime]MimeExtension{
	MimeTextPlain:       "txt",
	MimeTextHTML:        "html",
	MimeTextMarkdown:    "md",
	MimeImageJPEG:       "jpg",
	MimeImagePNG:        "png",
	MimeImageGIF:        "gif",
	MimeImageBMP:        "bmp",
	MimeImageSVG:        "svg",
	MimeAudioMPEG:       "mp3",
	MimeAudioOGG:        "ogg",
	MimeAudioOGA:        "oga",
	MimeAudioWAV:        "wav",
	MimeVideoMPEG:       "mpeg",
	MimeVideoMP4:        "mp4",
	MimeVideoOGG:        "ogg",
	MimeVideoQuickTime:  "mov",
	MimeApplicationJSON: "json",
	MimeApplicationXML:  "xml",
	MimeApplicationPDF:  "pdf",
	MimeApplicationZIP:  "zip",
	MimeApplicationGZIP: "gzip",
	MimeApplicationTar:  "tar",
	MimeApplicationXZ:   "xz",
	MimeOctetStream:     "",
}
