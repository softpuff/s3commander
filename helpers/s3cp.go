package helpers

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync/atomic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type progressWriter struct {
	written int64
	writer  io.WriterAt
	size    int64
}

func (pw *progressWriter) WriteAt(p []byte, off int64) (int, error) {
	atomic.AddInt64(&pw.written, int64(len(p)))

	percDownloaded := float32(pw.written*100) / float32(pw.size)

	fmt.Printf("File size:%d downloaded:%d percentage:%.2f%%\n", pw.size, pw.written, percDownloaded)

	return pw.writer.WriteAt(p, off)
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func GetFileSize(svc *s3.S3, bucket string, prefix string) (filesize int64, err error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(prefix),
	}

	resp, err := svc.HeadObject(params)
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}

func parseFileName(key string) (filename string) {
	ss := strings.Split(key, "/")
	s := ss[len(ss)-1]
	return s
}

func (c *AWSConfig) CpS3file(key string, bucket string, dest string, debug bool) error {
	filename := parseFileName(key)
	if debug {
		fmt.Printf("Filename: %s\n", filename)
	}
	s3Client := s3.New(c.Session)
	downloader := s3manager.NewDownloader(c.Session)
	size, err := GetFileSize(s3Client, bucket, key)
	if err != nil {
		return err
	}
	if debug {
		fmt.Printf("%s size: %v\n", key, size)
	}
	println("Starting download, size:", byteCountDecimal(size))
	if dest == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		dest = cwd
	}
	if debug {
		fmt.Printf("Destination path: %s\n", path.Join(dest, filename))
	}

	temp, err := ioutil.TempFile(dest, "dlS3File-tmp-")
	if err != nil {
		return err
	}
	tempfile := temp.Name()
	if debug {
		fmt.Printf("Destination temp file: %s\n", tempfile)
	}

	writer := &progressWriter{
		writer:  temp,
		size:    size,
		written: 0,
	}
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	if _, err := downloader.Download(writer, params); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed! Deleting tempfile %s", tempfile)
		os.Remove(tempfile)
	}

	if err := temp.Close(); err != nil {
		return err
	}

	if err := os.Rename(temp.Name(), path.Join(dest, filename)); err != nil {
		return err
	}
	fmt.Println("File downloaded at:", filename)
	fmt.Println("")
	return nil
}
