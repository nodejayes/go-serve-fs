package goservefs

import (
	"embed"
	"fmt"
	"net/http"
	"path"
	"strings"
)

type Config struct {
	UseClientSideRouter bool
	MainPath string
}

func getMimeType(ext string) string {
	switch ext {
	case ".aac":
		return "audio/aac"
	case ".abw":
		return "application/x-abiword"
	case ".apng":
		return "image/apng"
	case ".arc":
		return "application/x-freearc"
	case ".avif":
		return "image/avif"
	case ".avi":
		return "video/x-msvideo"
	case ".azw":
		return "application/vnd.amazon.ebook"
	case ".bmp":
		return "image/bmp"
	case ".bz":
		return "application/x-bzip"
	case ".bz2":
		return "application/x-bzip2"
	case ".csh":
		return "application/x-csh"
	case ".csv":
		return "text/csv"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".eot":
		return "application/vnd.ms-fontobject"
	case ".epub":
		return "application/epub+zip"
	case ".gz":
		return "application/gzip"
	case ".gif":
		return "image/gif"
	case ".ics":
		return "text/calendar"
	case ".jar":
		return "application/java-archive"
	case ".mid", ".midi":
		return "audio/x-midi"
	case ".mp3":
		return "audio/mpeg"
	case ".mp4":
		return "video/mp4"
	case ".mpeg":
		return "video/mpeg"
	case ".mpkg":
		return "application/vnd.apple.installer+xml"
	case ".odp":
		return "application/vnd.oasis.opendocument.presentation"
	case ".ods":
		return "application/vnd.oasis.opendocument.spreadsheet"
	case ".odt":
		return "application/vnd.oasis.opendocument.text"
	case ".oga":
		return "audio/ogg"
	case ".ogv":
		return "video/ogg"
	case ".ogx":
		return "application/ogg"
	case ".opus":
		return "audio/opus"
	case ".otf":
		return "font/otf"
	case ".php":
		return "application/x-httpd-php"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".rar":
		return "application/vnd.rar"
	case ".rtf":
		return "application/rtf"
	case ".sh":
		return "application/x-sh"
	case ".tar":
		return "application/x-tar"
	case ".tif", ".tiff":
		return "image/tiff"
	case ".ts":
		return "video/mp2t"
	case ".ttf":
		return "font/ttf"
	case ".vsd":
		return "application/vnd.visio"
	case ".wav":
		return "audio/wav"
	case ".weba":
		return "audio/webm"
	case ".webm":
		return "video/webm"
	case ".webp":
		return "image/webp"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".xhtml":
		return "application/xhtml+xml"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".xml":
		return "application/xml"
	case ".xul":
		return "application/vnd.mozilla.xul+xml"
	case ".zip":
		return "application/zip"
	case ".3gp":
		return "video/3gpp"
	case ".3gp2":
		return "video/3gpp2"
	case ".7zip":
		return "application/x-7z-compressed"
	case ".jsonld":
		return "application/ld+json"
	case ".json":
		return "application/json"
	case ".js", ".mjs":
		return "text/javascript"
	case ".css":
		return "text/css"
	case ".htm", ".html":
		return "text/html"
	case ".pdf":
		return "application/pdf"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".ico":
		return "image/vnd.microsoft.icon"
	case ".svg":
		return "	image/svg+xml"
	case ".bin":
		return "application/octet-stream"
	}
	return "text/plain"
}

func detectMainFolderName(fileSystem embed.FS) string {
	entries, err := fileSystem.ReadDir(".")
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		return strings.Split(entry.Name(), "/")[0]
	}
	return ""
}

func ConnectFS(fileSystem embed.FS, config *Config) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ext := path.Ext(req.URL.Path)
		requestPath := req.URL.Path
		if strings.HasSuffix(requestPath, "/") {
			ext = ".html"
			requestPath = fmt.Sprintf("%sindex.html", requestPath)
		}
		filePath := fmt.Sprintf("%s%s", detectMainFolderName(fileSystem), requestPath)
		content, err := fileSystem.ReadFile(filePath)
		if err != nil {
			if config.UseClientSideRouter {
				ext = ".html"
				filePath := fmt.Sprintf("%s%s", detectMainFolderName(fileSystem), config.MainPath)
				content, err := fileSystem.ReadFile(filePath)
				if err != nil {
					res.WriteHeader(http.StatusNotFound)
					res.Write([]byte{})
					return
				}
				res.Header().Set("Content-Type", getMimeType(ext))
				res.WriteHeader(http.StatusOK)
				res.Write(content)
				return
			}
			res.WriteHeader(http.StatusNotFound)
			res.Write([]byte{})
			return
		}
		res.Header().Set("Content-Type", getMimeType(ext))
		res.WriteHeader(http.StatusOK)
		res.Write(content)
	}
}
