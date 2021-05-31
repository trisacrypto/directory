package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "gdsutil"
	app.Version = pkg.Version()
	app.Usage = "utilities for operating the GDS service and database"
	app.Commands = []cli.Command{
		{
			Name:     "ldb:keys",
			Usage:    "list the keys currently in the leveldb store",
			Category: "leveldb",
			Action:   ldbKeys,
			Before:   openLevelDB,
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
			Name:      "ldb:get",
			Usage:     "get the value for the specified key",
			ArgsUsage: "key [key ...]",
			Category:  "leveldb",
			Action:    ldbGet,
			Before:    openLevelDB,
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
				cli.StringFlag{
					Name:  "o, out",
					Usage: "write the fetched key to directory if specified, otherwise printed",
				},
			},
		},
		{
			Name:      "ldb:put",
			Usage:     "insert key/value pair into database, loading from disk if necessary",
			ArgsUsage: "key/path [value]\n\n   Examples:\n   gdsutil ldb:put foo bar\n   gdsutil ldb:put foo bar.json\n   gdsutil ldb:put foo.json",
			Category:  "leveldb",
			Action:    ldbPut,
			Before:    openLevelDB,
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
					Name:  "f, format",
					Usage: "format of the data (raw, json, pb)",
					Value: "json",
				},
			},
		},
		{
			Name:      "ldb:delete",
			Usage:     "delete the leveldb record for the specified key(s)",
			ArgsUsage: "key [key ...]",
			Category:  "leveldb",
			Action:    ldbDelete,
			Before:    openLevelDB,
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
	}

	app.Run(os.Args)
}

//===========================================================================
// LevelDB Actions
//===========================================================================

var (
	ldb            *leveldb.DB
	vaspPrefix     = []byte("vasps::")
	certreqsPrefix = []byte("certreqs::")
	indexPrefix    = []byte("index::")
	sequencePrefix = []byte("sequence::")
	indexNames     = []byte("index::names")
	indexCountries = []byte("index::countries")
	sequencePK     = []byte("sequence::pks")
)

func ldbKeys(c *cli.Context) (err error) {
	defer ldb.Close()

	var prefix *util.Range
	if prefixs := c.String("prefix"); prefixs != "" {
		prefix = util.BytesPrefix([]byte(prefixs))
	}

	iter := ldb.NewIterator(prefix, nil)
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

func ldbGet(c *cli.Context) (err error) {
	defer ldb.Close()

	var out string
	if out = c.String("out"); out != "" {
		// Check that out is a directory
		var info fs.FileInfo
		if info, err = os.Stat(out); err != nil {
			return cli.NewExitError(err, 1)
		}
		if !info.IsDir() {
			return cli.NewExitError("specify a directory to write files out to", 1)
		}
	}

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
		if data, err = ldb.Get(key, nil); err != nil {
			return cli.NewExitError(err, 1)
		}

		// Unmarshal the thing
		var (
			jsonValue interface{}
			pbValue   proto.Message
		)

		// Determine how to unmarshal the data
		if bytes.HasPrefix(key, vaspPrefix) {
			vasp := new(pb.VASP)
			if err = proto.Unmarshal(data, vasp); err != nil {
				return cli.NewExitError(err, 1)
			}
			pbValue = vasp
		} else if bytes.HasPrefix(key, certreqsPrefix) {
			careq := new(models.CertificateRequest)
			if err = proto.Unmarshal(data, careq); err != nil {
				return cli.NewExitError(err, 1)
			}
			pbValue = careq
		} else if bytes.Equal(key, indexNames) {
			jsonValue = make(map[string]string)
			if err = json.Unmarshal(data, &jsonValue); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else if bytes.Equal(key, indexCountries) {
			jsonValue = make(map[string][]string)
			if err = json.Unmarshal(data, &jsonValue); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else if bytes.HasPrefix(key, sequencePrefix) {
			pk, n := binary.Uvarint(data)
			if n <= 0 {
				return cli.NewExitError("could not parse sequence", 1)
			}
			jsonValue = pk
		} else {
			return cli.NewExitError("could not determine unmarshal type", 1)
		}

		// Marshal JSON representation
		var outdata []byte
		if jsonValue != nil {
			if outdata, err = json.MarshalIndent(jsonValue, "", "  "); err != nil {
				return cli.NewExitError(err, 1)
			}
		}
		if pbValue != nil {
			jsonpb := protojson.MarshalOptions{
				Multiline:       true,
				Indent:          "  ",
				AllowPartial:    true,
				UseProtoNames:   true,
				UseEnumNumbers:  false,
				EmitUnpopulated: true,
			}
			if outdata, err = jsonpb.Marshal(pbValue); err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if out != "" {
			path := filepath.Join(out, string(key)+".json")
			if err = ioutil.WriteFile(path, outdata, 0644); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			fmt.Println(string(outdata) + "\n")
		}
	}

	return nil
}

func ldbPut(c *cli.Context) (err error) {
	defer ldb.Close()

	if c.NArg() == 0 || c.NArg() > 2 {
		return cli.NewExitError("specify path, key and path, or key and value as arguments", 1)
	}

	// Determine the key and value as follows:
	// If only one argument, assume it is a path; key is basename, value is data
	// If two arguments, check if second value is a path, otherwise second value is a key
	args := c.Args()
	format := strings.ToLower(c.String("format"))
	var key, data, value []byte
	if c.NArg() == 1 {
		path := args.Get(0)
		name := filepath.Base(path)
		ext := filepath.Ext(name)
		if strings.TrimLeft(ext, ".") != format {
			return cli.NewExitError(fmt.Errorf("mismatch file extension %q and data format %q: specify --format", ext, format), 1)
		}

		key = []byte(strings.TrimSuffix(name, ext))
		if data, err = ioutil.ReadFile(path); err != nil {
			return cli.NewExitError(err, 1)
		}
	} else {
		key = []byte(args.Get(0))

		// Determine if second argument is a path
		varg := args.Get(1)
		if isFile(varg) {
			ext := filepath.Ext(varg)
			if strings.TrimLeft(ext, ".") != format {
				return cli.NewExitError(fmt.Errorf("mismatch file extension %q and data format %q: specify --format", ext, format), 1)
			}
			if data, err = ioutil.ReadFile(varg); err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			data = []byte(varg)
		}

	}

	// Perform base64 decoding if necessary
	b64decode := c.Bool("b64decode")
	if b64decode {
		if key, err = base64.RawStdEncoding.DecodeString(string(key)); err != nil {
			return cli.NewExitError(err, 1)
		}
		if data, err = base64.RawStdEncoding.DecodeString(string(data)); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	// Quick spot check
	if len(data) == 0 || len(key) == 0 {
		return cli.NewExitError("no key or value found", 1)
	}

	// Unmarshal the thing from data then
	// Marshal the database representation
	if format != "raw" && format != "bytes" {
		jsonpb := &protojson.UnmarshalOptions{
			AllowPartial:   true,
			DiscardUnknown: true,
		}

		if bytes.HasPrefix(key, vaspPrefix) {
			vasp := new(pb.VASP)
			switch format {
			case "json":
				if err = jsonpb.Unmarshal(data, vasp); err != nil {
					return cli.NewExitError(err, 1)
				}
				if value, err = proto.Marshal(vasp); err != nil {
					return cli.NewExitError(err, 1)
				}
			case "pb", "proto", "protobuf":
				if err = proto.Unmarshal(data, vasp); err != nil {
					return cli.NewExitError(err, 1)
				}
				value = data
			default:
				return cli.NewExitError(fmt.Errorf("cannot unmarshal VASP format %q", format), 1)
			}

		} else if bytes.HasPrefix(key, certreqsPrefix) {
			careq := new(models.CertificateRequest)
			switch format {
			case "json":
				if err = jsonpb.Unmarshal(data, careq); err != nil {
					return cli.NewExitError(err, 1)
				}
				if value, err = proto.Marshal(careq); err != nil {
					return cli.NewExitError(err, 1)
				}
			case "pb", "proto", "protobuf":
				if err = proto.Unmarshal(data, careq); err != nil {
					return cli.NewExitError(err, 1)
				}
				value = data
			default:
				return cli.NewExitError(fmt.Errorf("cannot unmarshal Certificate Request format %q", format), 1)
			}
		} else if bytes.HasPrefix(key, indexPrefix) {
			switch format {
			case "json":
				value = data
			default:
				return cli.NewExitError(fmt.Errorf("cannot unmarshal index format %q", format), 1)
			}
		} else if bytes.Equal(key, sequencePK) {
			switch format {
			case "json":
				var pk uint64
				if err = json.Unmarshal(data, &pk); err != nil {
					return cli.NewExitError(err, 1)
				}
				value = make([]byte, binary.MaxVarintLen64)
				binary.PutUvarint(value, pk)
			default:
				return cli.NewExitError(fmt.Errorf("cannot unmarshal sequence format %q", format), 1)
			}

		} else {
			return cli.NewExitError("could not determine unmarshal type from key", 1)
		}
	} else {
		// Raw or bytes data is just the data
		value = data
	}

	// Final spot check
	if len(value) == 0 {
		return cli.NewExitError("no value marshaled", 1)
	}

	// Put the key/value to the database
	if err = ldb.Put(key, value, nil); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func ldbDelete(c *cli.Context) (err error) {
	defer ldb.Close()
	if c.NArg() == 0 {
		return cli.NewExitError("specify at least one key to delete", 1)
	}

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

		if err = ldb.Delete(key, nil); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	return nil
}

//===========================================================================
// LevelDB Helper Functions
//===========================================================================

func openLevelDB(c *cli.Context) (err error) {
	var uri string
	if uri = c.String("db"); uri == "" {
		return cli.NewExitError("specify path to leveldb database", 1)
	}

	var dsn *store.DSN
	if dsn, err = store.ParseDSN(uri); err != nil {
		return cli.NewExitError(err, 1)
	}

	if dsn.Scheme != "leveldb" && dsn.Scheme != "ldb" {
		return cli.NewExitError("this action requires a leveldb DSN", 1)
	}

	if ldb, err = leveldb.OpenFile(dsn.Path, nil); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func isFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return fi.Mode().IsRegular()
}

//===========================================================================
// Cipher Actions
//===========================================================================

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

//===========================================================================
// Cipher Helper Functions
//===========================================================================

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
