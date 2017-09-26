package main

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"

	"github.com/hashicorp/terraform/command"
	pbBackend "github.com/nhite/pb-backend"
	pb "github.com/nhite/pb-nhite"
)

type grpcCommands struct {
	meta    command.Meta
	backend *pbBackend.BackendClient
}

func (g *grpcCommands) Push(stream pb.Terraform_PushServer) error {
	workdir, err := ioutil.TempDir("", ".terraformgrpc")
	if err != nil {
		return err
	}
	err = os.Chdir(workdir)
	if err != nil {
		return err
	}

	body, err := stream.Recv()
	if err == io.EOF || err == nil {
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
	}
	if err != nil {
		return err
	}
	return stream.SendAndClose(&pb.Id{
		Tmpdir: workdir,
	})
}

func (g *grpcCommands) Init(ctx context.Context, in *pb.Arg) (*pb.Output, error) {
	err := os.Chdir(in.WorkingDir)
	if err != nil {
		return &pb.Output{int32(0), nil, nil}, err
	}

	tfCommand := &command.InitCommand{
		Meta: g.meta,
	}
	var stdout []byte
	var stderr []byte
	myUI := &grpcUI{
		stdout: stdout,
		stderr: stderr,
	}
	tfCommand.Meta.Ui = myUI
	ret := int32(tfCommand.Run(in.Args))
	return &pb.Output{ret, myUI.stdout, myUI.stderr}, err
}

func (g *grpcCommands) Apply(ctx context.Context, in *pb.Arg) (*pb.Output, error) {
	err := os.Chdir(in.WorkingDir)
	if err != nil {
		return &pb.Output{int32(0), nil, nil}, err
	}

	tfCommand := &command.ApplyCommand{
		Meta:       g.meta,
		ShutdownCh: ctx.Done(),
	}
	var stdout []byte
	var stderr []byte
	myUI := &grpcUI{
		stdout: stdout,
		stderr: stderr,
	}
	tfCommand.Meta.Ui = myUI
	ret := int32(tfCommand.Run(in.Args))
	return &pb.Output{ret, myUI.stdout, myUI.stderr}, err
}

func (g *grpcCommands) Plan(ctx context.Context, in *pb.Arg) (*pb.Output, error) {
	err := os.Chdir(in.WorkingDir)
	if err != nil {
		return &pb.Output{int32(0), nil, nil}, err
	}

	tfCommand := &command.PlanCommand{
		Meta: g.meta,
	}
	var stdout []byte
	var stderr []byte
	myUI := &grpcUI{
		stdout: stdout,
		stderr: stderr,
	}
	tfCommand.Meta.Ui = myUI
	ret := int32(tfCommand.Run(in.Args))
	return &pb.Output{ret, myUI.stdout, myUI.stderr}, err
}
