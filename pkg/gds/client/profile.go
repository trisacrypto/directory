package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/trisacrypto/directory/pkg/gds/admin/v2"
	"github.com/trisacrypto/directory/pkg/gds/store"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Profile contains the client-side configuration to connect to a specifc GDS instance.
// Profiles are loaded first from the YAML configuration file and then can be overrided
// by the CLI context if the user specifies a value via an environment variable or flag.
type Profile struct {
	Directory   *DirectoryProfile `yaml:"directory"`              // directory configuration
	Admin       *AdminProfile     `yaml:"admin"`                  // admin api configuration
	Replica     *ReplicaProfile   `yaml:"replica"`                // replica configuration
	DatabaseURL string            `yaml:"database_url,omitempty"` // localhost only: the dsn to the leveldb database, usually $GDS_DATABASE_URL
	Timeout     time.Duration     `yaml:"timeout,omitempty"`      // default timeout to create contexts for API connections, if not specified defaults to 30 seconds
	Active      bool              `yaml:"active,omitempty"`       // if this is the active profile, it is treated as the default profile
}

type DirectoryProfile struct {
	Endpoint string `yaml:"endpoint"`           // the GDS endpoint to connect to the gRPC directory service, also $TRISA_DIRECTORY_URL
	Insecure bool   `yaml:"insecure,omitempty"` // do not connect to the directory endpoint with TLS
}

type AdminProfile struct {
	Endpoint  string            `yaml:"endpoint"`             // the Admin URL to connect to the Admin API, also $TRISA_DIRECTORY_ADMIN_URL
	Audience  string            `yaml:"audience,omitempty"`   // the Audience for local token generation auth, usually $GDS_ADMIN_AUDIENCE
	TokenKeys map[string]string `yaml:"token_keys,omitempty"` // the token keys identifier and paths for local token generation auth, usually $GDS_ADMIN_TOKEN_KEYS
}

type ReplicaProfile struct {
	Endpoint string `yaml:"endpoint"`           // the replica endpoint to connect to the anti-entropy service
	Insecure bool   `yaml:"insecure,omitempty"` // do not connect to the replica endpoint with TLS
}

// Update the specified profile with the CLI context.
func (p *Profile) Update(c *cli.Context) error {
	if endpoint := c.String("endpoint"); endpoint != "" {
		p.Directory.Endpoint = endpoint
	}

	if endpoint := c.String("admin-endpoint"); endpoint != "" {
		p.Admin.Endpoint = endpoint
	}

	if insecure := c.Bool("no-secure"); insecure {
		p.Directory.Insecure = insecure
	}

	if dburl := c.String("db"); dburl != "" {
		p.DatabaseURL = dburl
	}

	return nil
}

// Context returns a default context with the timeout specified or 30 seconds by default.
func (p *Profile) Context() (context.Context, context.CancelFunc) {
	if p.Timeout > 0 {
		return context.WithTimeout(context.Background(), p.Timeout)
	}
	return context.WithTimeout(context.Background(), 30*time.Second)
}

// OpenLevelDB opens a leveldb database using the DSN supplied for gdsutil commands.
func (p *Profile) OpenLevelDB() (ldb *leveldb.DB, err error) {
	if p.DatabaseURL == "" {
		return nil, errors.New("please specify a leveldb DSN to connect to the database")
	}

	var dsn *store.DSN
	if dsn, err = store.ParseDSN(p.DatabaseURL); err != nil {
		return nil, err
	}

	if dsn.Scheme != "leveldb" && dsn.Scheme != "ldb" {
		return nil, fmt.Errorf("cannot open leveldb database with %q scheme", dsn.Scheme)
	}

	return leveldb.OpenFile(dsn.Path, nil)
}

// Connect to the TRISA Directory Service and return a gRPC client
func (p *DirectoryProfile) Connect() (_ api.TRISADirectoryClient, err error) {
	var opts []grpc.DialOption
	if p.Insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		config := &tls.Config{}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	}

	// Connect the directory client
	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(p.Endpoint, opts...); err != nil {
		return nil, err
	}
	return api.NewTRISADirectoryClient(cc), nil
}

// Connect to the GDS Admin API and return an admin client
func (p *AdminProfile) Connect() (client admin.DirectoryAdministrationClient, err error) {
	// Connect the admin client
	if client, err = admin.New(p.Endpoint); err != nil {
		return nil, err
	}

	// Attempt to login the admin client
	// TODO: use audience and token keys stored on profile
	if err = client.Login(context.TODO()); err != nil {
		return nil, err
	}
	return client, nil
}
