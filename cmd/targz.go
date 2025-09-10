package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"time"
)

func targz(filepath *string) error {
	uuid := syscall.Getuid()
	switch uuid {
	case 0:
		fmt.Println("Is a root")
		return comprx(filepath)
	default:
		return comprx(filepath)
		//return fmt.Errorf("not root")
	}

	return nil
}

func comprx(filepath *string) error {
	archive := archiveName(filepath)
	fmt.Println(archive)

	out, err := os.Create(archive)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := createArchive(out, filepath); err != nil {
		return err
	}

	return nil
}

func createArchive(buf io.Writer, filepath *string) error {
	gzw := gzip.NewWriter(buf)
	defer gzw.Close()

	trw := tar.NewWriter(gzw)
	defer trw.Close()

	files, err := os.ReadDir(*filepath)
	if err != nil {
		return err
	}

	if err := addFiles(trw, filepath, &files); err != nil {
		return err
	}

	return nil
}

func addFiles(trw *tar.Writer, filepath *string, files *[]os.DirEntry) error {
	for _, file := range *files {
		f, err := os.Open(*filepath + "/" + file.Name())
		if err != nil {
			return err
		}
		defer f.Close()

		fInfo, err := f.Stat()
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fInfo, fInfo.Name())
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(file.Name(), *filepath)
		if err := trw.WriteHeader(header); err != nil {
			return err
		}

		if _, err := io.Copy(trw, f); err != nil {
			return err
		}
	}

	return nil
}

func archiveName(filepath *string) string {
	path := strings.Split(*filepath, "/")
	date := time.Now().Format(time.DateOnly)
	times := time.Now().Format(time.TimeOnly)

	return fmt.Sprintf("%s_%s-%s-%s.tar.gz",
		"logfx",
		path[len(path)-1],
		strings.Replace(date, "-", "", -1),
		strings.Replace(times, ":", "", -1),
	)
}
