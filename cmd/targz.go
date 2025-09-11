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

func targz(from, to *string) error {
	uuid := syscall.Getuid()

	tgzFromDir := strings.TrimSuffix(*from, "/")
	tgzToDir := strings.TrimSuffix(*to, "/")

	switch uuid {
	case 0:
		return comprx(&tgzFromDir, &tgzToDir)
	default:
		return fmt.Errorf("need to grant sudo privileges")
	}

	return nil
}

func comprx(from, to *string) error {
	logger.Info("logfx - creating archive - start task", "from dir", fromDir, "to dir", toDir)

	archive := archiveName(from)
	fmt.Println(archive)

	out, err := os.Create(*to + "/" + archive)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := createArchive(out, from); err != nil {
		return err
	}

	logger.Info("logfx - archive created - end task", "from dir", fromDir, "to dir", toDir)
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
		fname := *filepath + "/" + file.Name()
		logger.Info("logfx - adding file", "file", fname)

		f, err := os.Open(fname)
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
