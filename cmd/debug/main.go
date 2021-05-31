package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	rvasp "github.com/trisacrypto/testnet/pkg/rvasp/pb/v1"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto/aesgcm"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/handler"
	"github.com/trisacrypto/trisa/pkg/trisa/mtls"
	"github.com/trisacrypto/trisa/pkg/trust"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "debug"
	app.Version = "beta"
	app.Usage = "debugging utilities for the TRISA directory service"
	app.Commands = []cli.Command{
		{
			Name:     "store:keys",
			Usage:    "list the keys currently in the leveldb store",
			Category: "store",
			Action:   storeKeys,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				cli.BoolFlag{
					Name:  "s, stringify",
					Usage: "stringify keys otherwise they are base64 encoded",
				},
				cli.StringFlag{
					Name:  "p, prefix",
					Usage: "specify a prefix to filter keys on",
				},
			},
		},
		{
			Name:      "store:get",
			Usage:     "get the value for the specified key",
			Category:  "store",
			Action:    storeGet,
			ArgsUsage: "key [key ...]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				cli.BoolFlag{
					Name:  "b, b64decode",
					Usage: "specify the keys as base64 encoded values which must be decoded",
				},
			},
		},
		{
			Name:     "store:put",
			Usage:    "put the value for the specified key",
			Category: "store",
			Action:   storePut,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				cli.BoolFlag{
					Name:  "b, b64decode",
					Usage: "specify the key and value as base64 encoded strings which must be decoded",
				},
				cli.StringFlag{
					Name:  "k, key",
					Usage: "the key to put the value to",
				},
				cli.StringFlag{
					Name:  "v, value",
					Usage: "the value to put to the database (or specify json document)",
				},
				cli.StringFlag{
					Name:  "p, path",
					Usage: "path to a JSON document containing the value",
				},
			},
		},
		{
			Name:      "store:delete",
			Usage:     "delete the leveldb record for the specified key(s)",
			Category:  "store",
			Action:    storeDelete,
			ArgsUsage: "key [key ...]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "d, db",
					Usage:  "dsn to connect to trisa directory storage",
					EnvVar: "GDS_DATABASE_URL",
				},
				cli.BoolFlag{
					Name:  "b, b64decode",
					Usage: "specify the keys as base64 encoded values which must be decoded",
				},
			},
		},
		{
			Name:      "decrypt",
			Usage:     "decrypt base64 encoded ciphertext with an HMAC signature",
			ArgsUsage: "ciphertext hmac",
			Category:  "cipher",
			Action:    cipherDecrypt,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "k, key",
					Usage:  "secret key to decrypt the cipher text",
					EnvVar: "GDS_SECRET_KEY",
				},
			},
		},
		{
			Name:     "transfer",
			Usage:    "send a generated unary transfer request to VASP",
			Category: "trisa",
			Action:   transfer,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "e, endpoint",
					Usage:  "endpoint to send the transfer request to",
					EnvVar: "TRISA_ENDPOINT",
					Value:  "localhost:4435",
				},
				cli.StringFlag{
					Name:   "c, certs",
					Usage:  "path to client certificates",
					EnvVar: "RVASP_CERT_PATH",
					Value:  "fixtures/certs/bob.gz",
				},
				cli.StringFlag{
					Name:   "t, trust-chain",
					Usage:  "path to trust chain pool",
					EnvVar: "RVASP_TRUST_CHAIN_PATH",
					Value:  "fixtures/certs/trisa.zip",
				},
				cli.StringFlag{
					Name:   "k, pubkey",
					Usage:  "path to public key of the beneficiary",
					EnvVar: "BENEFICIARY_PUBKEY",
				},
				cli.StringFlag{
					Name:  "d, transaction",
					Usage: "path to JSON data to load transaction data from",
				},
				cli.StringFlag{
					Name:  "i, identity",
					Usage: "path to JSON data to load identity data from",
				},
			},
		},
		{
			Name:     "stream",
			Usage:    "open a transfer stream and send transfer requests to VASP",
			Category: "trisa",
			Action:   transferStream,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "e, endpoint",
					Usage:  "endpoint to send the transfer request to",
					EnvVar: "TRISA_ENDPOINT",
					Value:  "localhost:4435",
				},
				cli.StringFlag{
					Name:   "c, certs",
					Usage:  "path to client certificates",
					EnvVar: "RVASP_CERT_PATH",
					Value:  "fixtures/certs/bob.gz",
				},
				cli.StringFlag{
					Name:   "t, trust-chain",
					Usage:  "path to trust chain pool",
					EnvVar: "RVASP_TRUST_CHAIN_PATH",
					Value:  "fixtures/certs/trisa.zip",
				},
			},
		},
		{
			Name:     "address",
			Usage:    "send a confirm adress request to VASP",
			Category: "trisa",
			Action:   confirmAddress,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "e, endpoint",
					Usage:  "endpoint to send the transfer request to",
					EnvVar: "TRISA_ENDPOINT",
					Value:  "localhost:4435",
				},
				cli.StringFlag{
					Name:   "c, certs",
					Usage:  "path to client certificates",
					EnvVar: "RVASP_CERT_PATH",
					Value:  "fixtures/certs/bob.gz",
				},
				cli.StringFlag{
					Name:   "t, trust-chain",
					Usage:  "path to trust chain pool",
					EnvVar: "RVASP_TRUST_CHAIN_PATH",
					Value:  "fixtures/certs/trisa.zip",
				},
			},
		},
		{
			Name:     "keys",
			Usage:    "Send a key exchange request to the beneficiary VASP",
			Category: "trisa",
			Action:   keyExchange,
			Before:   initClient,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "e, endpoint",
					Usage:  "endpoint to send the transfer request to",
					EnvVar: "TRISA_ENDPOINT",
					Value:  "localhost:4435",
				},
				cli.StringFlag{
					Name:   "c, certs",
					Usage:  "path to client certificates",
					EnvVar: "RVASP_CERT_PATH",
					Value:  "fixtures/certs/bob.gz",
				},
				cli.StringFlag{
					Name:   "t, trust-chain",
					Usage:  "path to trust chain pool",
					EnvVar: "RVASP_TRUST_CHAIN_PATH",
					Value:  "fixtures/certs/trisa.zip",
				},
				cli.StringFlag{
					Name:  "k, key",
					Usage: "path to signing key to send (required)",
				},
			},
		},
		{
			Name:     "status",
			Usage:    "send a health check request to the VASP",
			Category: "trisa",
			Action:   healthCheck,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "e, endpoint",
					Usage:  "endpoint to send the transfer request to",
					EnvVar: "TRISA_ENDPOINT",
					Value:  "localhost:4435",
				},
				cli.StringFlag{
					Name:   "c, certs",
					Usage:  "path to client certificates",
					EnvVar: "RVASP_CERT_PATH",
					Value:  "fixtures/certs/bob.gz",
				},
				cli.StringFlag{
					Name:   "t, trust-chain",
					Usage:  "path to trust chain pool",
					EnvVar: "RVASP_TRUST_CHAIN_PATH",
					Value:  "fixtures/certs/trisa.zip",
				},
				cli.UintFlag{
					Name:  "a, attempts",
					Usage: "set the number of previous attempts",
				},
				cli.DurationFlag{
					Name:  "l, last-checked",
					Usage: "set the last checked field as this long ago",
				},
			},
		},
	}

	app.Run(os.Args)
}

func storeKeys(c *cli.Context) (err error) {
	if c.String("db") == "" {
		return cli.NewExitError("specify path to leveldb database", 1)
	}

	var db *leveldb.DB
	if db, err = leveldb.OpenFile(c.String("db"), nil); err != nil {
		return cli.NewExitError(err, 1)
	}
	defer db.Close()

	var prefix *util.Range
	if prefixs := c.String("prefix"); prefixs != "" {
		prefix = util.BytesPrefix([]byte(prefixs))
	}

	iter := db.NewIterator(prefix, nil)
	defer iter.Release()

	stringify := c.Bool("stringify")
	for iter.Next() {
		if stringify {
			fmt.Printf("- %s\n", string(iter.Key()))
		} else {
			fmt.Printf("- %s\n", base64.RawStdEncoding.EncodeToString(iter.Key()))
		}
	}

	if err = iter.Error(); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func storeGet(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		return cli.NewExitError("specify at least one key to fetch", 1)
	}
	if c.String("db") == "" {
		return cli.NewExitError("specify path to leveldb database", 1)
	}

	var db *leveldb.DB
	if db, err = leveldb.OpenFile(c.String("db"), nil); err != nil {
		return cli.NewExitError(err, 1)
	}
	defer db.Close()

	b64decode := c.Bool("b64decode")
	for _, keys := range c.Args() {
		var key []byte
		if b64decode {
			if key, err = base64.RawStdEncoding.DecodeString(keys); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			key = []byte(keys)
		}

		var data []byte
		if data, err = db.Get(key, nil); err != nil {
			return cli.NewExitError(err, 1)
		}

		// Unmarshall the thing
		var value interface{}

		// Determine how to unmarshall the data
		if bytes.HasPrefix(key, []byte("vasps")) {
			vasp := new(pb.VASP)
			if err = proto.Unmarshal(data, vasp); err != nil {
				return cli.NewExitError(err, 1)
			}
			value = vasp
		} else if bytes.HasPrefix(key, []byte("certreqs")) {
			careq := new(models.CertificateRequest)
			if err = proto.Unmarshal(data, careq); err != nil {
				return cli.NewExitError(err, 1)
			}
			value = careq
		} else if bytes.Equal(key, []byte("index::names")) {
			value = make(map[string]string)
			if err = json.Unmarshal(data, &value); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else if bytes.Equal(key, []byte("index::countries")) {
			value = make(map[string][]string)
			if err = json.Unmarshal(data, &value); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else if bytes.Equal(key, []byte("sequence::pks")) {
			pk, n := binary.Uvarint(data)
			if n <= 0 {
				return cli.NewExitError("could not parse sequence", 1)
			}
			value = pk
		} else {
			return cli.NewExitError("could not determine unmarshall type", 1)
		}

		// Marshall the JSON representation
		var out []byte
		if out, err = json.MarshalIndent(value, "", "  "); err != nil {
			return cli.NewExitError(err, 1)
		}
		fmt.Println(string(out))
	}

	return nil
}

func storePut(c *cli.Context) (err error) {
	if c.String("key") == "" {
		return cli.NewExitError("must specify a key to put to", 1)
	}
	if c.String("value") != "" && c.String("path") != "" {
		return cli.NewExitError("specify either value or path, not both", 1)
	}
	if c.String("db") == "" {
		return cli.NewExitError("specify path to leveldb database", 1)
	}

	var db *leveldb.DB
	if db, err = leveldb.OpenFile(c.String("db"), nil); err != nil {
		return cli.NewExitError(err, 1)
	}
	defer db.Close()

	var key, data, value []byte
	keys := c.String("key")
	b64decode := c.Bool("b64decode")

	if b64decode {
		if key, err = base64.RawStdEncoding.DecodeString(keys); err != nil {
			return cli.NewExitError(err, 1)
		}
	} else {
		key = []byte(keys)
	}

	if c.String("value") != "" {
		if b64decode {
			// If value is b64 encoded then we just assume it's data to put directly
			if value, err = base64.RawStdEncoding.DecodeString(c.String("value")); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			data = []byte(keys)
		}
	}

	if c.String("path") != "" {
		if data, err = ioutil.ReadFile(c.String("path")); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	// Quick spot check
	if len(data) == 0 && len(value) == 0 {
		return cli.NewExitError("no value to put to database", 1)
	}

	if len(data) > 0 && len(value) > 0 {
		return cli.NewExitError("both data and value specified?", 1)
	}

	if len(data) > 0 {
		jsonpb := &protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}

		// Unmarshall the thing from JSON then
		// Marshall the database representation
		if bytes.HasPrefix(key, []byte("vasps")) {
			vasp := new(pb.VASP)
			if err = jsonpb.Unmarshal(data, vasp); err != nil {
				return cli.NewExitError(err, 1)
			}
			if value, err = proto.Marshal(vasp); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else if bytes.HasPrefix(key, []byte("certreqs")) {
			careq := new(models.CertificateRequest)
			if err = jsonpb.Unmarshal(data, careq); err != nil {
				return cli.NewExitError(err, 1)
			}
			if value, err = proto.Marshal(careq); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else if bytes.Equal(key, []byte("index::names")) {
			var names map[string]string
			if err = json.Unmarshal(data, &names); err != nil {
				return cli.NewExitError(err, 1)
			}
			if value, err = json.Marshal(names); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else if bytes.Equal(key, []byte("index::countries")) {
			var countries map[string][]string
			if err = json.Unmarshal(data, &countries); err != nil {
				return cli.NewExitError(err, 1)
			}
			if value, err = json.Marshal(countries); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else if bytes.Equal(key, []byte("sequence::pks")) {
			var pk uint64
			if err = json.Unmarshal(data, &pk); err != nil {
				return cli.NewExitError(err, 1)
			}
			value = make([]byte, binary.MaxVarintLen64)
			binary.PutUvarint(value, pk)
		} else {
			return cli.NewExitError("could not determine unmarshal type", 1)
		}
	}

	// Final spot check
	if len(value) == 0 {
		return cli.NewExitError("no value marshalled", 1)
	}

	// Put the key/value to the database
	if err = db.Put(key, value, nil); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func storeDelete(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		return cli.NewExitError("specify at least one key to fetch", 1)
	}
	if c.String("db") == "" {
		return cli.NewExitError("specify path to leveldb database", 1)
	}

	var db *leveldb.DB
	if db, err = leveldb.OpenFile(c.String("db"), nil); err != nil {
		return cli.NewExitError(err, 1)
	}
	defer db.Close()

	b64decode := c.Bool("b64decode")
	for _, keys := range c.Args() {
		var key []byte
		if b64decode {
			if key, err = base64.RawStdEncoding.DecodeString(keys); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			key = []byte(keys)
		}

		if err = db.Delete(key, nil); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	return nil
}

// TODO: package this all up somewhere!

const nonceSize = 12

func cipherDecrypt(c *cli.Context) (err error) {
	if c.NArg() != 2 {
		return cli.NewExitError("must specify ciphertext and hmac arguments", 1)
	}

	var secret string
	if secret = c.String("key"); secret == "" {
		return cli.NewExitError("cipher key required", 1)
	}

	var ciphertext, signature []byte
	if ciphertext, err = base64.RawStdEncoding.DecodeString(c.Args()[0]); err != nil {
		return cli.NewExitError(fmt.Errorf("could not decode ciphertext: %s", err), 1)
	}
	if signature, err = base64.RawStdEncoding.DecodeString(c.Args()[1]); err != nil {
		return cli.NewExitError(fmt.Errorf("could not decode signature: %s", err), 1)
	}

	if len(ciphertext) == 0 {
		return cli.NewExitError("empty cipher text", 1)
	}

	// Create a 32 byte signature of the key
	hash := sha256.New()
	hash.Write([]byte(secret))
	key := hash.Sum(nil)

	// Separate the data from the nonce
	data := ciphertext[:len(ciphertext)-nonceSize]
	nonce := ciphertext[len(ciphertext)-nonceSize:]

	// Validate HMAC signature
	if err = validateHMAC(key, data, signature); err != nil {
		return cli.NewExitError(err, 1)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	plainbytes, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Println(string(plainbytes))
	return nil
}

func createHMAC(key, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("cannot sign empty data")
	}
	hm := hmac.New(sha256.New, key)
	hm.Write(data)
	return hm.Sum(nil), nil
}

func validateHMAC(key, data, sig []byte) error {
	hmac, err := createHMAC(key, data)
	if err != nil {
		return err
	}

	if !bytes.Equal(sig, hmac) {
		return errors.New("HMAC mismatch")
	}
	return nil
}

var (
	client api.TRISANetworkClient
)

func initClient(c *cli.Context) (err error) {
	var (
		cc    *grpc.ClientConn
		sz    *trust.Serializer
		certs *trust.Provider
		pool  trust.ProviderPool
	)

	// TODO: can we reuse the serializer for two different file types?
	if sz, err = trust.NewSerializer(false); err != nil {
		return cli.NewExitError(err, 1)
	}

	if certs, err = sz.ReadFile(c.String("certs")); err != nil {
		return cli.NewExitError(err, 1)
	}

	if pool, err = sz.ReadPoolFile(c.String("trust-chain")); err != nil {
		return cli.NewExitError(err, 1)
	}

	opts := make([]grpc.DialOption, 0, 1)
	var opt grpc.DialOption
	if opt, err = mtls.ClientCreds(c.String("endpoint"), certs, pool); err != nil {
		return cli.NewExitError(err, 1)
	}
	opts = append(opts, opt)

	if cc, err = grpc.Dial(c.String("endpoint"), opts...); err != nil {
		return cli.NewExitError(err, 1)
	}

	client = api.NewTRISANetworkClient(cc)
	return nil
}

type transferData struct {
	Identity    *ivms101.IdentityPayload
	Transaction *rvasp.Transaction
}

func transfer(c *cli.Context) (err error) {
	var key *rsa.PublicKey
	var data *transferData

	if keyPath := c.String("pubkey"); keyPath != "" {
		var (
			sz  *trust.Serializer
			p   *trust.Provider
			crt *x509.Certificate
			ok  bool
		)
		if sz, err = trust.NewSerializer(false); err != nil {
			return cli.NewExitError(err, 1)
		}
		if p, err = sz.ReadFile(keyPath); err != nil {
			return cli.NewExitError(err, 1)
		}
		if crt, err = p.GetLeafCertificate(); err != nil {
			return cli.NewExitError(err, 1)
		}
		if key, ok = crt.PublicKey.(*rsa.PublicKey); !ok {
			return cli.NewExitError("RSA public key required in cert", 1)
		}

	} else {
		return cli.NewExitError("specify public key path of beneficiary", 1)
	}

	data = new(transferData)
	if transactionPath := c.String("transaction"); transactionPath != "" {
		var f []byte
		if f, err = ioutil.ReadFile(transactionPath); err != nil {
			return cli.NewExitError(err, 1)
		}

		var t rvasp.Transaction
		if err = protojson.Unmarshal(f, &t); err != nil {
			return cli.NewExitError(err, 1)
		}

		data.Transaction = &t

	} else {
		return cli.NewExitError("specify path to json transaction data", 1)
	}

	if identityPath := c.String("identity"); identityPath != "" {
		var f []byte
		if f, err = ioutil.ReadFile(identityPath); err != nil {
			return cli.NewExitError(err, 1)
		}

		var t ivms101.IdentityPayload
		if err = protojson.Unmarshal(f, &t); err != nil {
			return cli.NewExitError(err, 1)
		}

		data.Identity = &t

	} else {
		return cli.NewExitError("specify path to json identity data", 1)
	}

	// Create the transaction data
	payload := &api.Payload{}
	if payload.Identity, err = anypb.New(data.Identity); err != nil {
		return cli.NewExitError(err, 1)
	}

	if payload.Transaction, err = anypb.New(data.Transaction); err != nil {
		return cli.NewExitError(err, 1)
	}

	// Encrypt the transaction data
	var cipher *aesgcm.AESGCM
	if cipher, err = aesgcm.New(nil, nil); err != nil {
		return cli.NewExitError(err, 1)
	}
	var req *api.SecureEnvelope
	if req, err = handler.Seal(uuid.New().String(), payload, cipher, key); err != nil {
		return cli.NewExitError(err, 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	printJSON(req)
	var rep *api.SecureEnvelope
	if rep, err = client.Transfer(ctx, req); err != nil {
		return cli.NewExitError(err, 1)
	}

	printJSON(rep)

	// TODO: load private key of originator and print response
	var env *handler.Envelope
	if env, err = handler.Open(rep, nil); err != nil {
		return cli.NewExitError(err, 1)
	}
	return printJSON(env)
}

func transferStream(c *cli.Context) (err error) {
	return cli.NewExitError("not implemented yet", 1)
}

func confirmAddress(c *cli.Context) (err error) {
	return cli.NewExitError("not implemented yet", 1)
}

func keyExchange(c *cli.Context) (err error) {
	var key *x509.Certificate
	if keyPath := c.String("key"); keyPath != "" {
		var (
			sz *trust.Serializer
			p  *trust.Provider
		)

		if sz, err = trust.NewSerializer(false); err != nil {
			return cli.NewExitError(err, 1)
		}

		if p, err = sz.ReadFile(keyPath); err != nil {
			return cli.NewExitError(err, 1)
		}
		if key, err = p.GetLeafCertificate(); err != nil {
			return cli.NewExitError(err, 1)
		}
	} else {
		return cli.NewExitError("specify key path to load keys from", 1)
	}

	req := &api.SigningKey{
		Version:            int64(key.Version),
		Signature:          key.Signature,
		SignatureAlgorithm: key.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: key.PublicKeyAlgorithm.String(),
		NotBefore:          key.NotBefore.Format(time.RFC3339),
		NotAfter:           key.NotAfter.Format(time.RFC3339),
	}

	if req.Data, err = x509.MarshalPKIXPublicKey(key.PublicKey); err != nil {
		return cli.NewExitError(err, 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *api.SigningKey
	if rep, err = client.KeyExchange(ctx, req); err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

func healthCheck(c *cli.Context) (err error) {
	var (
		cc     *grpc.ClientConn
		client api.TRISAHealthClient
		sz     *trust.Serializer
		certs  *trust.Provider
		pool   trust.ProviderPool
	)

	// TODO: can we reuse the serializer for two different file types?
	if sz, err = trust.NewSerializer(false); err != nil {
		return cli.NewExitError(err, 1)
	}

	if certs, err = sz.ReadFile(c.String("certs")); err != nil {
		return cli.NewExitError(err, 1)
	}

	if pool, err = sz.ReadPoolFile(c.String("trust-chain")); err != nil {
		return cli.NewExitError(err, 1)
	}

	opts := make([]grpc.DialOption, 0, 1)
	var opt grpc.DialOption
	if opt, err = mtls.ClientCreds(c.String("endpoint"), certs, pool); err != nil {
		return cli.NewExitError(err, 1)
	}
	opts = append(opts, opt)

	if cc, err = grpc.Dial(c.String("endpoint"), opts...); err != nil {
		return cli.NewExitError(err, 1)
	}

	client = api.NewTRISAHealthClient(cc)
	req := &api.HealthCheck{
		Attempts: uint32(c.Uint("attempts")),
	}

	if lastCheckedAgo := c.Duration("last-checked"); lastCheckedAgo != 0 {
		req.LastCheckedAt = time.Now().Add(-1 * lastCheckedAgo).Format(time.RFC3339)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *api.ServiceState
	if rep, err = client.Status(ctx, req); err != nil {
		return cli.NewExitError(err, 1)
	}

	return printJSON(rep)
}

// helper function to print JSON response and exit
func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Println(string(data))
	return nil
}
