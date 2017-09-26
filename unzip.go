package main

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"

	pb "github.com/nhite/pb-nhite"
)

func unzip(body *pb.Body) error {
	workdir, err := ioutil.TempDir("", ".terraformgrpc")
	if err != nil {
		return err
	}
	err = os.Chdir(workdir)
	if err != nil {
		return err
	}

	// We have all the file
	// Now let's extract the zipfile
	r, err := zip.NewReader(bytes.NewReader(body.Zipfile), int64(len(body.Zipfile)))
	if err != nil {
		return err
	}
	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(f.Name, 0700)
			if err != nil {
				return err
			}
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		localFile, err := os.Create(f.Name)
		if err != nil {
			return err
		}
		_, err = io.Copy(localFile, rc)
		if err != nil {
			return err
		}

		rc.Close()
	}
	return nil
}
