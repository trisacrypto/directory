package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ExtractGzip extracts a gzipped archive to the specified directory and
// returns the path to the root level of the extracted archive.
// Caller must ensure that the destination directory exists.
func ExtractGzip(file, destDir string, skipHidden bool) (root string, err error) {
	var (
		f  *os.File
		gr *gzip.Reader
	)

	// Read the gzip file.
	if f, err = os.Open(file); err != nil {
		return "", err
	}
	defer f.Close()
	if gr, err = gzip.NewReader(f); err != nil {
		return "", err
	}
	defer gr.Close()

	// Write the contents to the temporary directory.
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(filepath.Join(destDir, hdr.Name), os.FileMode(hdr.Mode)); err != nil {
				return "", err
			}
			if root == "" {
				root = filepath.Join(destDir, hdr.Name)
			}
		case tar.TypeReg:
			var reg *os.File
			if skipHidden && hdr.Name[0] == '.' {
				// Skip hidden files if requested.
				continue
			}
			if reg, err = os.Create(filepath.Join(destDir, hdr.Name)); err != nil {
				return "", err
			}
			if _, err = io.Copy(reg, tr); err != nil {
				reg.Close()
				return "", err
			}
			reg.Close()
		default:
			return "", fmt.Errorf("extracting %s: unknown type flag: %c", hdr.Name, hdr.Typeflag)
		}
	}
	return root, nil
}

// WriteGzip compresses a directory to a gzipped archive.
func WriteGzip(dir, file string) (err error) {
	var (
		f *os.File
	)
	// Create a gzip file.
	if f, err = os.Create(file); err != nil {
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	defer gw.Close()

	// Create a tar file.
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Write the DB to the tar file.
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		var hdr *tar.Header
		if hdr, err = tar.FileInfoHeader(info, ""); err != nil {
			return err
		}
		hdr.Name = path[len(dir):]
		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		var tmp *os.File
		if tmp, err = os.Open(path); err != nil {
			return err
		}
		defer tmp.Close()
		if _, err = io.Copy(tw, tmp); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
