package filter

import (
	"os"
	"path/filepath"
	"strings"
)

// binaryExts is a set of file extensions that are known to be binary.
// Checking this first avoids the cost of opening the file and reading 512 bytes.
var binaryExts = map[string]bool{
	// Images
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".bmp": true,
	".ico": true, ".webp": true, ".tiff": true, ".tif": true, ".svgz": true,
	// Audio/Video
	".mp3": true, ".mp4": true, ".wav": true, ".flac": true, ".avi": true,
	".mkv": true, ".mov": true, ".wmv": true, ".ogg": true, ".m4a": true,
	// Archives
	".zip": true, ".tar": true, ".gz": true, ".bz2": true, ".xz": true,
	".7z": true, ".rar": true, ".tgz": true, ".zst": true,
	// Executables / Binaries
	".exe": true, ".dll": true, ".so": true, ".dylib": true, ".o": true,
	".a": true, ".class": true, ".pyc": true, ".pyd": true,
	// Documents (binary formats)
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	".ppt": true, ".pptx": true, ".odt": true,
	// Fonts
	".ttf": true, ".otf": true, ".woff": true, ".woff2": true, ".eot": true,
	// Data / Other
	".sqlite": true, ".db": true, ".iso": true, ".dmg": true, ".pkg": true,
	".deb": true, ".rpm": true, ".msi": true, ".lock": true,
	".wasm": true, ".proto": true,
}

// HasBinaryExt returns true if the file has a known binary extension.
// This is a cheap check (no I/O) and should be called before IsBinary.
func HasBinaryExt(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return binaryExts[ext]
}

// IsBinary opens the file and reads the first 512 bytes to detect binary content.
// This is more expensive than HasBinaryExt and should only be called as a fallback.
func IsBinary(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && n == 0 {
		return false
	}
	return containsNullByte(buf[:n])
}

func containsNullByte(buf []byte) bool {
	for _, b := range buf {
		if b == 0 {
			return true
		}
	}
	return false
}
