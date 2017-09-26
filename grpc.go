package main

import (
	"crypto/sha256"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/hashicorp/terraform/command"
	pbBackend "github.com/nhite/pb-backend"
	pb "github.com/nhite/pb-nhite"
)

type grpcCommands struct {
	meta       command.Meta
	workingDir string
	backend    *pbBackend.BackendClient
}

func (g *grpcCommands) Push(stream pb.Terraform_PushServer) error {
	var chksum [32]byte
	body, err := stream.Recv()
	if err == io.EOF || err == nil {
		chksum = sha256.Sum256(body.Zipfile)
		// TODO context
		// TODO msgSize
		pushClient, err := (*g.backend).Store(context.Background(), grpc.MaxCallRecvMsgSize(65536))
		if err != nil {
			return err
		}

		err = pushClient.Send(&pbBackend.Element{
			ID: &pbBackend.ElementID{
				ID: fmt.Sprintf("%x", chksum),
			},
			Body: body.Zipfile,
		})
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return stream.SendAndClose(&pb.Id{
		Sha256: fmt.Sprintf("%x", chksum),
	})
}

func (g *grpcCommands) Init(ctx context.Context, in *pb.Arg) (*pb.Output, error) {
	err := g.getLocalFiles(in.WorkingDir)
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
	err := g.getLocalFiles(in.WorkingDir)
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
	err := g.getLocalFiles(in.WorkingDir)
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
