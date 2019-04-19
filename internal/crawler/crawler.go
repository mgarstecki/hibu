package crawler

import (
	"crypto/md5"
	"io"
	"os"

	"github.com/spf13/afero"
)

type CrawlResult struct {
	Err error

	File string
	Hash []byte
}

type Crawler struct {
	filesystem afero.Fs
}

func NewCrawler(filesystem afero.Fs) *Crawler {
	return &Crawler{
		filesystem: filesystem,
	}
}

func (this *Crawler) Crawl(baseDir string, resultsChannel chan CrawlResult) error {
	err := afero.Walk(this.filesystem, baseDir, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			file, err := this.filesystem.Open(path)
			if err != nil {
				resultsChannel <- CrawlResult{
					Err:  err,
					File: path,
				}

				return nil
			}

			defer file.Close()

			md5 := md5.New()
			_, err = io.Copy(md5, file)
			if err != nil {
				resultsChannel <- CrawlResult{
					Err:  err,
					File: path,
				}

				return nil
			}

			sum := md5.Sum(nil)
			resultsChannel <- CrawlResult{
				File: path,
				Hash: sum,
			}
		}

		return nil
	})

	close(resultsChannel)

	if err != nil {
		return err
	}

	return nil
}
