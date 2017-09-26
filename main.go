package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hashicorp/terraform/command"
	"github.com/kelseyhightower/envconfig"
	pbBackend "github.com/nhite/pb-backend"
	pb "github.com/nhite/pb-nhite"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type configuration struct {
	ListenAddress   string `envconfig:"LISTEN_ADDRESS" required:"true" default:"127.0.0.1:1234"`
	MaxMessageSize  int    `envconfig:"MAX_RECV_MSG_SIZE" required:"true" default:"16500545"`
	BackendAddress  string `envconfig:"BACKEND_ADDRESS" required:"true"`
	BackendCertFile string `envconfig:"BACKEND_CERT_FILE" required:"true"`
	CertFile        string `envconfig:"CERT_FILE" required:"true"`
	KeyFile         string `envconfig:"KEY_FILE" required:"true"`
}

const envPrefix = "nhite"

var (
	config configuration
	// Build date
	Build string
	// Version number
	Version     string
	versionFlag bool
)

func initBackend(config *configuration) (*pbBackend.BackendClient, error) {

	creds, err := credentials.NewClientTLSFromFile(config.BackendCertFile, "")
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(config.BackendAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}
	//defer conn.Close()
	client := pbBackend.NewBackendClient(conn)
	return &client, nil
}

func main() {
	flag.BoolVar(&versionFlag, "v", false, "Display version then exit")
	flag.Parse()
	if versionFlag {
		if Version == "" {
			Version = "dev"
		}
		fmt.Printf("%v version %v, build %v\n", os.Args[0], Version, Build)
		os.Exit(0)
	}
	if len(os.Args) > 1 {
		envconfig.Usage(envPrefix, &config)
		os.Exit(1)

	}
	err := envconfig.Process(envPrefix, &config)
	if err != nil {
		envconfig.Usage(envPrefix, &config)
		fmt.Println(err)
		os.Exit(1)
	}

	// Try to access the backend
	backend, err := initBackend(&config)
	if err != nil {
		log.Fatalf("Cannot reach backend: %v", err)
	}

	log.Println("Listening on " + config.ListenAddress)
	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
	if err != nil {
		log.Fatal("could not load TLS keys: ", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.MaxRecvMsgSize(config.MaxMessageSize))
	// PluginOverrides are paths that override discovered plugins, set from
	// the config file.
	var PluginOverrides command.PluginOverrides

	meta := command.Meta{
		Color:            false,
		GlobalPluginDirs: globalPluginDirs(),
		PluginOverrides:  &PluginOverrides,
		Ui:               &grpcUI{},
	}

	pb.RegisterTerraformServer(grpcServer, &grpcCommands{
		meta:       meta,
		backend:    backend,
		workingDir: os.TempDir(),
	})
	log.Fatal(grpcServer.Serve(listener))
}
