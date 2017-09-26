package main

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/grpc"

	pbBackend "github.com/nhite/pb-backend"
	pb "github.com/nhite/pb-nhite"
)

// getLocalFiles check whether a directory identified by the ID exists in the workingdir
// if not, it creates it and fetch the data from the backend

func (g *grpcCommands) getLocalFiles(id string) error {
	cwd := filepath.Join(g.workingDir, id)
	// Get fileinfos
	var fi os.FileInfo
	var err error
	fi, err = os.Stat(cwd)
	switch {
	case err == nil:
		if fi.IsDir() {
			log.Println("IsDirectory")
			err = os.Chdir(cwd)
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("Not a directory")
	case os.IsNotExist(err):
		break
	default:
		return err
	}
	// TODO check if the directory is accessible
	// Creating the temporary directory
	err = os.Mkdir(cwd, 0700)
	if err != nil {
		return err
	}
	err = os.Chdir(cwd)
	if err != nil {
		return err
	}
	// Fetch from the backend
	fetchClient, err := (*g.backend).Fetch(context.Background(), &pbBackend.ElementID{id}, grpc.MaxCallRecvMsgSize(65536))
	log.Println("DEBUG, fetchClient")
	if err != nil {
		return err
	}
	element, err := fetchClient.Recv()
	if err != nil {
		return err
	}
	return unzip(&pb.Body{
		Zipfile: element.Body,
	})
}

func unzip(body *pb.Body) error {

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
