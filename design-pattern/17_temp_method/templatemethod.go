package temp_method

// 模板方法 使用继承机制，把通用步骤和通用方法放到父类中，把具体实现延迟到子类中实现。
import "log"

type Downloader interface {
	Download(uri string)
}

type implement interface {
	download()
	save()
}

type template struct {
	implement
	uri string
}

func newTemplate(impl implement) *template {
	return &template{
		implement: impl,
	}
}

func (t *template) Download(uri string) {
	t.uri = uri
	log.Println("prepare downloading")
	t.implement.download()
	t.implement.save()
	log.Println("finish downloading")
}

func (t *template) save() {
	log.Println("default save")
}

// http downloader

type HTTPDownloader struct {
	*template
}

func NewHTTPDownloader() Downloader {
	downloader := &HTTPDownloader{}
	t := newTemplate(downloader)
	downloader.template = t
	return downloader
}

func (d *HTTPDownloader) download() {
	log.Printf("downlaod %s via http", d.uri)
}

func (d *HTTPDownloader) save() {
	log.Println("http save")
}

// ftp downloader

type FTPDownloader struct {
	*template
}

func NewFTPDownloader() Downloader {
	downloader := &FTPDownloader{}
	t := newTemplate(downloader)
	downloader.template = t
	return downloader
}

func (d *FTPDownloader) download() {
	log.Printf("download %s via ftp", d.uri)
}
