package temp_method

import "testing"

func TestHTTPDownloader(t *testing.T) {
	downloader := NewHTTPDownloader()
	downloader.Download("https://qq.com/downlaod")
}

func TestFTPDownloader(t *testing.T) {
	downloader := NewFTPDownloader()
	downloader.Download("https://qq.com/downlaod")
}
