// Code generated by go-bindata. DO NOT EDIT.
// sources:
// templates/deliver_certs.html (1.577kB)
// templates/deliver_certs.txt (1.267kB)
// templates/reject_registration.html (587B)
// templates/reject_registration.txt (469B)
// templates/review_request.html (1.225kB)
// templates/review_request.txt (978B)
// templates/verify_contact.html (652B)
// templates/verify_contact.txt (528B)

package emails

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _deliver_certsHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x54\x4d\x6f\x1c\x37\x0c\x3d\x77\x7e\x05\x6b\xf4\xd0\x02\xf6\x0e\x92\xa3\x31\x1d\x34\x8d\x93\x76\xd1\xc0\x0d\xbc\xdb\x02\x3d\x72\x25\xee\x8e\x6a\x49\x54\x29\xce\x6e\x27\x81\xff\x7b\xa0\x99\xfd\x86\x6f\xfa\x7a\xe4\xe3\xd3\x23\x9b\xd4\xfe\x4e\xde\x33\x7c\xfd\x0a\xb3\x47\x0c\x04\x2f\x2f\xb7\x4d\x9d\xda\xaa\x6a\x52\xfb\x0f\xf7\x02\xcb\xa7\xf9\xe2\x1d\x44\xd2\x1d\xcb\x33\x08\x6d\x5c\x56\x41\x75\x1c\xa1\xc3\x0c\x2b\xa2\x08\x98\x92\xf0\x96\xec\xf7\x30\x42\x58\x36\x18\xdd\x97\xab\x47\x1b\xc1\xa8\x64\x2b\x67\x29\xaa\xd3\x01\x0c\x89\xba\xb5\x33\xa8\x94\x61\x8b\xde\x59\x54\x17\x37\x30\x94\x18\x81\xc2\x8a\x24\x77\x2e\x81\x32\xb0\x76\x74\xa0\xb2\xbf\x81\xad\x43\xd0\x8e\xa6\xd3\xea\x37\xcf\x2b\xf4\xf0\xe0\x84\x8c\xb2\x0c\x33\x78\xa7\x8a\xa6\x23\x5b\xf0\xda\xb9\x0c\x14\xd0\x79\x40\x21\xf8\xfc\xc7\xfb\xc5\x9b\xb7\x40\xd1\xc8\x90\x94\xec\x25\x95\x5c\xde\xa3\x16\x1e\x95\xc1\x08\x2e\x24\x4f\x81\xa2\x9e\xd2\x41\x12\x56\x36\xec\xa1\xcf\x85\x72\x58\x7e\x5a\xc0\xce\x69\xb7\x67\x7a\x90\xeb\xc0\x55\x19\xe8\x7f\xd3\x61\xdc\x50\xb5\x14\xdc\x92\x87\xa7\xde\x13\x18\x0e\xc9\x3b\x8c\x86\xc0\xc5\x35\x4b\x18\x35\x9b\x1d\x7f\x60\xd9\x11\x24\x71\x01\x65\x00\x4b\x8a\xce\x67\xe0\xf5\xa4\x90\x3d\x94\x0a\x14\x55\x86\xb1\x30\xcc\xb0\x66\xef\x79\x97\xef\xf7\x31\x7a\xdf\x56\xdf\x35\xde\xb5\x4d\x56\xe1\xb8\x69\xe7\x0f\xf7\x4d\xbd\x5f\x8f\xdf\xfe\xf7\xfc\x01\x5e\x5e\x9a\xda\xbb\xb6\x02\x38\x7f\xfa\x34\x7e\x36\x09\xd9\x93\xae\x57\xe0\xd3\x93\xe3\x8b\x63\xb0\x8b\xb4\xef\x39\x04\x8e\x50\x3c\x76\x15\x62\xba\xd9\x9b\xef\x15\xe4\x82\xc4\xa1\x87\xc7\xbe\x48\x79\x85\x9d\xee\xa6\xab\xd7\xd1\x1f\xa2\x4d\xec\xa2\x5e\x01\x0f\xc7\x47\x50\x53\x17\xa5\x46\xcd\x19\x2c\x8d\xc6\x98\x74\x3e\xf7\xc6\x2d\x34\x14\xda\x81\x7b\xd8\x39\xef\x21\x52\xb1\x57\x77\x34\x54\xc2\x9c\x77\x2c\xb6\xa9\x29\xb4\x27\x13\x09\x19\x72\x5b\xb2\xb0\xeb\x28\x96\x13\x58\x3b\xc9\x0a\xb9\x5f\x05\xa7\xc5\x7f\x63\xa2\xf3\xde\x9a\xc1\x19\x8d\x92\xa1\x8f\x5f\x5c\x4a\xd7\x56\xe5\x58\x2e\x2b\xc3\x21\x60\xb4\xe0\x5d\xa4\xdb\x31\x41\xf1\x6d\x9f\x09\x1a\xc3\x96\x5a\x4e\x14\x73\xf6\x4d\x3d\xee\xce\x5c\x02\x3f\x1e\x6b\x59\x15\xa7\x71\x18\xdb\xa1\xd8\x35\x2a\xc9\xc4\xeb\x50\xd5\x4f\x07\x4f\x25\xa1\xf6\x07\xd8\x07\x85\xf4\x6c\xf2\x9b\xb7\x70\xe7\x22\xcc\x1f\x3f\xce\x3f\x7d\x98\xa5\xb2\xe5\x5e\xe1\xcf\xbf\x96\xe3\x81\x11\x85\xbb\xc8\x96\x72\x53\x17\xf0\xa8\xf3\x47\x16\x08\x2c\x17\xce\x2f\x05\xb9\xa8\xb4\x91\x69\x14\x8c\x2d\x75\xea\xba\x7d\x5b\xdd\x42\xf2\x84\x99\x20\x13\x01\xf7\x52\x59\x36\x7d\x69\xd0\x29\x06\x2a\x34\x08\x9d\xd0\xfa\xe7\x9b\x4e\x35\xe5\xfb\xba\x56\x71\x19\x67\x96\xb6\xf5\x4d\x7b\x5c\x37\x35\xb6\x33\x98\x8f\xfd\x04\x1d\x6e\x09\x30\x0e\xd5\x7f\x3d\xe5\x12\x27\x4f\x4a\x06\x1c\xc0\x70\x54\x34\x0a\x7d\xbe\x08\x5e\xe6\x89\xf2\x3d\xda\xe0\xe2\x2f\x53\x54\xc7\x37\xed\xe5\xbe\x24\x01\x96\xea\x5f\x76\xe5\x4b\x4a\x85\x45\xd4\xd7\x19\xde\x95\xf2\x72\x42\x43\xb3\xec\xd1\x3c\xcf\x0c\x87\x9b\x76\x51\x96\x50\xc6\x47\x24\x3f\x91\xfe\x3c\x09\x60\x19\x22\x2b\x08\x25\x3f\xec\x07\x82\x1f\x2e\x07\xde\x69\x98\xfc\x4a\x59\xe1\x89\x36\x28\x36\xdf\x36\x2b\x81\xba\xad\x26\x5d\xaf\x87\x27\x2c\x48\xb6\xce\x10\x2c\x09\xc3\x88\xff\x16\x00\x00\xff\xff\x32\xc6\x9e\x0e\x29\x06\x00\x00")

func deliver_certsHtmlBytes() ([]byte, error) {
	return bindataRead(
		_deliver_certsHtml,
		"deliver_certs.html",
	)
}

func deliver_certsHtml() (*asset, error) {
	bytes, err := deliver_certsHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "deliver_certs.html", size: 1577, mode: os.FileMode(0644), modTime: time.Unix(1632262825, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x36, 0x7d, 0x10, 0x39, 0xf8, 0xa, 0x55, 0xa4, 0x1d, 0x9c, 0x6f, 0xb1, 0xd0, 0x2c, 0x26, 0xec, 0x7, 0x47, 0x7b, 0x10, 0xf1, 0xcc, 0x91, 0x48, 0x87, 0x4, 0xf4, 0xcb, 0xf1, 0x3e, 0x22, 0x60}}
	return a, nil
}

var _deliver_certsTxt = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x54\xc1\x6e\x1b\x37\x10\xbd\xef\x57\xbc\xde\x5a\xc0\x5e\x23\x39\xfa\xd4\x34\x4e\x5a\xa1\x81\x1b\x58\x6a\x81\x1e\x47\xe4\x48\x3b\x35\xc9\x61\xc9\x59\xa9\x9b\xc0\xff\x5e\x70\x57\xb2\xe5\x20\xc7\x9d\x79\x1c\xbe\x79\xef\x71\x7f\xe3\x10\x14\x5f\xbf\xa2\xbf\xa7\xc8\x78\x7a\xba\xea\xba\xbf\x75\x2c\xd8\x3c\xac\xd6\xef\x90\xd8\x8e\x5a\x1e\x51\x78\x2f\xd5\x0a\x99\x68\xc2\x40\x15\x5b\xe6\x04\xca\xb9\xe8\x81\xfd\x0f\x98\x8f\x68\xd9\x53\x92\x2f\xdf\x80\xf6\x85\x92\xb1\x87\x78\x4e\x26\x36\xc1\x71\x31\xd9\x89\x23\xe3\x8a\x03\x05\xf1\x64\x92\xf6\x98\xda\x8c\xc8\x71\xcb\xa5\x0e\x92\x61\x0a\xb5\x81\xcf\x54\x4e\x1d\x1c\x84\x60\x03\x9f\xaa\xbf\x06\xdd\x52\xc0\x9d\x14\x76\xa6\x65\xea\xf1\xce\x8c\xdc\xc0\xbe\x9d\xb7\x41\x2a\x38\x92\x04\x50\x61\x7c\xfe\xfd\xfd\xfa\xcd\x5b\x70\x72\x65\xca\x8d\xd3\x2b\x2a\xb5\xe1\xc9\x1a\x0f\x38\x4a\x90\x98\x03\x47\x4e\x76\x71\x5d\x2e\x6a\xea\x34\x60\xac\x8d\x72\xdc\x7c\x5a\xe3\x28\x36\x9c\x98\x9e\xe5\x3a\x73\x35\x05\xff\xe7\x06\x4a\x7b\xc6\xa6\xd0\x81\x03\x1e\xc6\xc0\x70\x1a\x73\x10\x4a\x8e\x21\x69\xa7\x25\xce\x9a\xf5\x5d\xb7\x19\x18\xb9\x48\xa4\x32\xc1\xb3\x91\x84\x0a\xdd\x2d\xd2\xf8\xf3\x8e\xe0\x64\x65\x9a\x37\xa2\x8a\x9d\x86\xa0\xc7\x7a\xdb\x75\xab\xbb\xdb\xd9\xc9\xbf\x56\x77\x78\x7a\xea\x1e\x66\xcf\xb8\xb0\x7f\x91\x67\x01\xbc\x74\x9e\x1b\xed\xc0\x7b\x8d\x51\x13\x5a\x10\x16\xdc\x52\x38\x05\xa3\x5b\x73\x11\x0a\xb8\x1f\xdb\x6e\x0b\x60\x29\x2d\x95\x06\xf9\x90\x7c\x56\x49\xb6\x74\xcf\x5f\xad\xd3\x6d\x14\x9e\x67\xdd\x97\x6d\x2e\xa5\xbf\x9a\x35\x3f\x4a\x08\x48\xdc\x9c\x1b\x9e\xbd\xca\x54\xeb\x51\x8b\x7f\xb1\xa6\xb0\x63\x39\xb0\xc7\x71\xe0\x34\x57\x76\x52\xaa\xa1\x8e\xdb\x28\xd6\x5c\x9d\xe7\x5f\x26\xb6\xc7\xc5\xed\x6d\xf8\x98\xbe\x48\xce\xdf\x06\x40\xd3\xdc\x74\x1a\x23\x25\x8f\x20\x89\xaf\x9e\xd3\x30\x56\x86\x66\x4e\xb5\x86\x0b\xd5\xf1\xe3\x33\xf3\x6d\x73\x4e\xe3\x9c\xab\xe6\x7b\x32\x2e\x0b\x95\xf3\x0e\x3f\xdd\x76\xdd\x79\x44\x7e\x74\xf5\xcd\x5b\x5c\x4b\xc2\xea\xfe\xe3\xea\xd3\x87\x3e\xb7\x4f\x1d\x0d\x7f\xfc\xb9\x99\x0b\xae\x18\xae\x93\x7a\xae\x5d\xf7\x51\x0b\xa2\x96\x57\x71\x69\x7c\x25\x19\xef\xcb\xf2\x7e\xe6\x1c\xbe\x44\xf5\x94\xc5\x2b\xe4\xc0\x54\x19\x07\xa9\x62\x98\x83\xa4\x6e\x6c\xb9\x5e\xa6\x90\x61\x30\xcb\xf5\xf6\xe6\xc6\x8a\x54\xea\x3d\x1f\x6e\x7a\xac\xe6\xd4\x61\xa0\x03\x83\xd2\x84\x7f\x47\xae\x0d\x7f\x32\x2b\xd2\x04\xa7\xc9\xc8\x19\xc6\xda\x86\x90\x8f\x92\x7e\x5e\x46\x88\x42\x0b\xfe\x51\x69\xba\x35\x9e\xed\xda\x75\x20\xf7\x88\xf6\x1a\x12\x07\xcc\xc0\xeb\xc6\xb0\x66\x72\xdc\xd7\xd6\xed\x9d\xc6\x1e\x9f\x17\xc2\x5e\x91\xd4\x50\x38\x87\xe9\x14\xfe\x30\xbd\x7e\xd5\x7d\xd7\xfd\xc2\xd5\xf0\xc0\x7b\x2a\xbe\x5e\x75\xdf\xff\x29\x60\xcd\xe5\x20\x8e\xb1\x61\x8a\xdd\xff\x01\x00\x00\xff\xff\xd6\xa5\x72\xd2\xf3\x04\x00\x00")

func deliver_certsTxtBytes() ([]byte, error) {
	return bindataRead(
		_deliver_certsTxt,
		"deliver_certs.txt",
	)
}

func deliver_certsTxt() (*asset, error) {
	bytes, err := deliver_certsTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "deliver_certs.txt", size: 1267, mode: os.FileMode(0644), modTime: time.Unix(1632262825, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xee, 0x62, 0x6e, 0xe5, 0x26, 0x48, 0xcc, 0x1, 0x1c, 0x9f, 0x28, 0x99, 0x7, 0xc, 0x25, 0x3b, 0x86, 0xba, 0x6f, 0xbc, 0x7b, 0x6c, 0xba, 0x6b, 0x40, 0x58, 0x92, 0x5, 0xfc, 0xb4, 0xf0, 0xf6}}
	return a, nil
}

var _reject_registrationHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x91\x41\x8f\xd3\x40\x0c\x85\xef\xf9\x15\x4f\x7b\xe1\x52\x25\xf7\x6a\x88\x60\x55\x09\x7a\x41\xa8\x5b\xb8\x3b\x1d\xb7\x19\x34\x19\x07\x8f\xd3\x12\xad\xfa\xdf\xd1\x84\x76\x45\x2f\x7b\x4b\x3c\xf6\xf3\xfb\xfc\xdc\xd8\x7e\xe5\x18\x05\xaf\xaf\xa8\xbf\xd1\xc0\xb8\x5e\x57\xae\x19\xdb\xaa\x72\x63\xfb\x23\x1d\x45\x6d\x4a\x64\x1c\xe7\x15\xac\x67\xec\x77\xdb\x97\xcf\xd8\xf1\x39\xf0\x05\xcf\x42\xea\xd1\x53\x86\xf2\x2f\x3e\x18\x7b\xcc\x32\xe9\xad\xe9\x4b\x94\x8e\x22\x36\x41\xf9\x60\xa2\x33\x5e\x58\xcf\xe1\xc0\x95\xf2\x29\x64\x53\xb2\x20\x09\xca\xbf\x27\xce\x56\x63\xdf\x33\x94\x29\x4b\x5a\x16\xdd\xea\xb8\xfc\xaf\x1e\x32\x28\xe3\x28\x31\xca\x25\xaf\xdd\x62\x73\x8a\x6d\x05\xb8\x18\x5a\x97\x4d\x25\x9d\xda\xed\x66\xed\x9a\xdb\xf7\x02\xf6\x73\xbb\xc1\xf5\xea\x9a\x18\xee\xad\xa5\xba\xfb\xb7\xec\xfe\xe0\x9a\xa2\x54\xb0\xbf\x47\xa6\xcc\x48\x62\x0c\xeb\xc9\x10\x12\x44\x3d\x2b\x4c\x70\x0c\x7f\x16\x7f\xa3\x4a\x17\x79\xc8\x18\x38\x15\x10\xf6\xa0\x4e\xce\xbc\x2a\x27\xf8\x10\x23\x7a\x3a\x73\x19\xc8\x53\x37\x04\xab\x08\x0f\xd8\x34\x8e\x4c\xb1\xc6\x6d\x57\xe6\xe4\x41\xf7\x32\x44\xb1\xd0\x07\x49\x19\x26\x95\x23\xf4\xca\xc7\x8f\x4f\x03\x85\x68\xb2\x26\x3f\x84\xf4\xc9\x34\x64\xaa\x83\x3c\xb5\x8f\xff\xae\xa1\xf6\x4d\xd9\x4b\x01\x81\xf2\x18\x67\xf8\x25\x8b\x38\x17\x5f\xd6\x87\x0c\x2e\x82\xf5\x5b\xe0\xcf\xe5\xe2\x3b\x3e\x91\xfa\xbc\x72\x9d\xa2\x69\xab\xf7\xc3\xc4\x9e\x69\x58\xe6\xff\x06\x00\x00\xff\xff\x78\xd8\xb7\xea\x4b\x02\x00\x00")

func reject_registrationHtmlBytes() ([]byte, error) {
	return bindataRead(
		_reject_registrationHtml,
		"reject_registration.html",
	)
}

func reject_registrationHtml() (*asset, error) {
	bytes, err := reject_registrationHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "reject_registration.html", size: 587, mode: os.FileMode(0644), modTime: time.Unix(1622582528, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xfd, 0xdb, 0xf6, 0x6, 0xa7, 0xb8, 0x68, 0x11, 0x89, 0xea, 0x97, 0x3a, 0xf1, 0x23, 0xe7, 0xc2, 0xb7, 0x2b, 0x89, 0x53, 0xbe, 0xcf, 0xe3, 0x46, 0xe1, 0x35, 0xed, 0xec, 0x37, 0xd2, 0x40, 0x7d}}
	return a, nil
}

var _reject_registrationTxt = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x90\xcf\x4e\xf3\x30\x10\xc4\xef\x7e\x8a\xb9\x7d\x97\x28\x0f\xd0\xd3\x47\x55\x09\x7a\x41\xa8\x2d\xdc\x37\xf5\xb6\x31\x72\xbc\x61\xbd\x49\x89\xaa\xbc\x3b\x72\x28\x08\x4e\xdc\xec\xfd\x33\xf3\xdb\x79\xe0\x18\x05\xd7\x2b\xea\x47\xea\x18\xf3\x5c\x39\xf7\x9c\x4e\xa2\x36\x24\x32\x8e\x53\x05\x6b\x19\x87\xdd\x76\x7f\x87\x1d\x8f\x81\x2f\x58\x0b\xa9\x47\x4b\x19\xca\xaf\x7c\x34\xf6\x98\x64\xd0\xdb\xd0\x7d\x94\x86\x22\x36\x41\xf9\x68\xa2\x13\xf6\xac\x63\x38\x32\x94\xcf\x21\x9b\x92\x05\x49\x50\x7e\x1b\x38\x5b\x8d\x43\x5b\x3a\x94\x25\x2d\x46\xb7\x3a\x2e\x3f\xd5\x43\x06\x65\x9c\x24\x46\xb9\xe4\x95\x73\xdb\xcd\x6a\x41\x7e\xd9\x6e\x30\xcf\xce\x95\xf7\xee\x53\xa3\x7c\x9f\x22\x53\x66\x24\x31\x86\xb5\x64\x08\x09\xa2\x9e\x15\x26\x38\x85\xf7\xc5\xa8\x57\x69\x22\x77\x19\x1d\xa7\x42\xc4\x1e\xd4\xc8\xc8\x55\xb9\xe5\x5f\x8c\x68\x69\xe4\xb2\x90\x87\xa6\x0b\x06\xfa\xcd\x4f\x7d\xcf\x14\x6b\xdc\xbc\x32\x27\x0f\xfa\x2a\x43\x14\xcb\x19\x41\x52\x2e\x1a\xe4\xbb\x90\xfe\x9b\x86\x4c\x75\x90\xef\x2d\x2f\x05\x12\xca\x7d\x9c\xe0\x97\xc0\xe2\x54\xe6\xad\x0d\x19\xdc\x51\x88\xb5\x73\xeb\x92\x87\xf2\x99\xd4\xe7\xca\xfd\x11\xf2\x81\xa9\x73\x1f\x01\x00\x00\xff\xff\x2a\x11\x51\x8f\xd5\x01\x00\x00")

func reject_registrationTxtBytes() ([]byte, error) {
	return bindataRead(
		_reject_registrationTxt,
		"reject_registration.txt",
	)
}

func reject_registrationTxt() (*asset, error) {
	bytes, err := reject_registrationTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "reject_registration.txt", size: 469, mode: os.FileMode(0644), modTime: time.Unix(1622582528, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x63, 0xbf, 0x13, 0xc7, 0xbb, 0x61, 0x9f, 0xfc, 0x2e, 0x8, 0x33, 0x69, 0x65, 0xbd, 0xfe, 0x42, 0x2, 0xec, 0x9e, 0x97, 0xc9, 0x3, 0x2c, 0x84, 0x78, 0xc, 0x88, 0xb2, 0x88, 0xfb, 0xf2, 0x74}}
	return a, nil
}

var _review_requestHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x54\x5d\x6b\xdb\x4a\x10\x7d\xbe\xfa\x15\x07\x73\xe1\xde\x82\x3f\x68\x1f\x83\x2a\x70\x13\x68\x4d\x53\x08\xb6\x9b\x92\xc7\xb1\x76\x2c\x6d\xb3\xda\x51\x67\xd7\x36\x26\xe4\xbf\x97\x55\xe4\x28\x69\x92\xa7\xbe\xce\x9e\x99\x39\x73\xe6\xcc\xe6\x6d\xf1\x85\x9d\x13\xac\x97\x8b\xd5\x1c\x73\xd3\x58\x1f\xc6\xf9\xac\x2d\xb2\x2c\x6f\x8b\x1f\x8c\x9a\xf6\x0c\xe5\x92\xed\x9e\x0d\x08\x9e\x0f\x3d\xf8\xb3\x93\x0d\x39\x5c\x58\xe5\x32\x8a\x1e\xa1\x5c\xd9\x10\x95\xa2\x15\x0f\xe5\x5f\x3b\x0e\x11\x5b\x95\x06\x84\xeb\xf9\xea\x0a\xb1\xa6\x98\x79\x66\x13\x10\x05\x9b\x54\x77\x6f\xf9\xc0\x66\x8a\x75\xcd\xa7\x14\x51\xd4\x14\xb0\x67\xb5\x5b\xcb\x06\xb1\x66\xab\xe0\x86\xac\x03\x19\xa3\x1c\xc2\xff\x1c\xde\x81\xbc\x79\x42\x2c\xbb\xfa\x7a\xbe\x7a\xff\x01\x2d\x85\x70\x10\x35\xa9\x83\xe1\x52\x8f\x6d\x04\xa1\x64\x8d\x76\x6b\x4b\x8a\x8c\x43\x6d\xcb\x1a\x07\xeb\x5c\xa2\x50\xb1\x67\xa5\xc8\x06\xe2\xdd\x11\x76\x8b\xa3\xec\x40\x6d\xab\xb2\x67\xc4\xda\x86\xac\xe7\x35\xc5\x8d\xec\x50\x92\x47\x22\x0d\xc3\x91\xac\x0b\xa0\x8d\xec\x62\x22\xf9\x7c\x7e\xf1\x5d\xac\x53\x14\xc1\x46\x3e\xcb\x7b\x51\x73\x42\xad\xbc\xfd\x38\xba\xbb\xc3\xb4\x7b\x5f\x76\x32\x7c\x5f\x5e\xe2\xfe\x7e\x54\xbc\x1a\xce\x67\x54\x3c\xae\xe5\x89\x58\xb8\x99\x7f\xbb\x84\x0d\x67\xa7\x47\xe5\xae\xc0\xb2\x7f\x4e\x99\x29\xf6\x90\x27\x0f\xaa\x1e\x21\x0a\xe5\x9f\x5c\xbe\x42\xbc\x2f\x3c\xc6\x2e\x70\xf7\xba\x15\xe7\xe4\x60\x7d\x85\x86\x23\x19\x8a\x74\x6a\xb6\x73\x45\xf6\x4f\xee\x6c\xb1\xb8\x38\x43\x1e\xa2\x8a\xaf\xba\xee\xd7\x8b\x8b\xae\x73\x1f\xca\x67\xce\xf6\xc8\xb5\xdc\xb2\x7f\x0e\xee\x42\x2f\xe1\x40\xc2\x2f\x3b\x6a\xac\x6c\x06\xa3\x3d\x4f\x1f\x10\x83\x13\x5f\x14\xcb\x67\x89\x6b\xaf\x01\x95\x25\xb7\x2f\x27\x1f\x26\x2e\xa5\x69\x92\xbb\x9c\xf5\x8c\x28\xe2\x40\xa1\x97\x61\x10\xba\xc8\x4b\x31\x5c\xfc\x8b\xca\x84\xde\xc8\x98\x58\x0c\xe3\x63\x12\x31\x7a\x3a\xe0\x08\x13\xca\x67\x5d\xd6\xb0\x4a\x79\x6b\x11\x63\x24\x0b\x5a\xc3\xa0\x1e\x92\xb6\xd3\x70\x08\x54\x71\x22\x94\x12\x1c\x85\x08\xd2\x6a\xd7\xb0\x8f\x7f\x47\x6d\x89\x49\x83\xd1\xe8\x4f\x82\x57\x8e\x29\x30\xbc\x44\xee\xee\x37\x9d\xc7\x7f\xce\x21\x9d\x71\xba\xb1\xd0\x72\x99\x1c\xf5\xa0\x9b\xa6\x0d\xc0\x3c\x2e\x82\xbd\x69\xc5\xfa\x88\xad\x68\x07\x39\xbf\x5c\x64\x27\x79\xa3\xe0\x20\x7a\x3b\x7d\xec\xf5\x29\x39\x76\xc9\x15\xa9\x09\xe3\x7c\xa3\x98\x15\xd9\x1b\x7f\xcd\x8a\x75\x6f\x4b\xc6\x9a\xa9\x49\xf9\xbf\x03\x00\x00\xff\xff\x0a\xef\x65\x72\xc9\x04\x00\x00")

func review_requestHtmlBytes() ([]byte, error) {
	return bindataRead(
		_review_requestHtml,
		"review_request.html",
	)
}

func review_requestHtml() (*asset, error) {
	bytes, err := review_requestHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "review_request.html", size: 1225, mode: os.FileMode(0644), modTime: time.Unix(1634832170, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x3a, 0xaf, 0xd0, 0xf, 0x60, 0x6c, 0x66, 0x89, 0xd8, 0x30, 0x12, 0x82, 0x14, 0xb6, 0xa2, 0x82, 0xa6, 0x1a, 0x2a, 0x94, 0xca, 0xa2, 0x2b, 0x7a, 0x3e, 0xe8, 0xa6, 0xa1, 0x7a, 0x1a, 0x6a, 0x81}}
	return a, nil
}

var _review_requestTxt = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x93\x41\x6b\x1b\x49\x10\x85\xef\xfd\x2b\x1e\x62\x61\x77\x41\x12\xec\x1e\x7d\x53\x6c\x48\x44\x1c\x30\x23\xc5\xc1\xc7\xf2\xf4\xd3\x4c\xc7\x3d\x5d\x93\xee\x96\x06\x61\xfc\xdf\x43\x8f\xc6\x91\x42\xe2\x43\x6e\xa2\xe6\x75\xd5\x57\xaf\x9e\x3e\xd0\x7b\xc5\xb6\x5a\x6f\x56\x58\xd9\xce\x85\x34\x37\xe6\x0b\xd1\xca\x81\x88\xac\xe9\x0e\xb4\x10\x04\x0e\x93\xea\xbd\xd7\x47\xf1\xb8\x71\x91\x75\xd6\x78\x44\x64\xe3\x52\x8e\x92\x9d\x06\x44\x7e\xdb\x33\x65\xec\xa2\x76\x10\xdc\xaf\x36\x77\xc8\xad\x64\x13\x48\x9b\x90\x15\x8f\xa5\xef\xc1\x71\xa0\x5d\x62\xdb\xf2\xf5\x89\x46\xb4\x92\x70\x60\x74\x3b\x47\x8b\xdc\xd2\x45\xb0\x13\xe7\x21\xd6\x46\xa6\xf4\x0f\xd3\xbf\x90\x60\x2f\xc0\xcc\xdd\xc7\xeb\xcd\x7f\xff\xa3\x97\x94\x06\x8d\xb6\x4c\xb0\xac\xe3\xb1\xcf\x10\xd4\x8c\xd9\xed\x5c\x2d\x99\x18\x5a\x57\xb7\x18\x9c\xf7\x05\xa1\x61\x60\x94\x4c\x0b\x0d\xfe\x08\xb7\xc3\x51\xf7\x90\xbe\x8f\x7a\x20\x72\xeb\x92\x99\xb8\x96\x78\xd0\x3d\x6a\x09\x28\xd0\xb0\xcc\xe2\x7c\x82\x3c\xea\x3e\x17\xc8\x9f\xf7\xd7\x30\xd6\x46\x2b\x91\x5c\xe6\x95\x31\xcf\xcf\x58\x8e\x85\x6a\xdc\xfb\x73\x75\x8b\x97\x17\x63\x2e\x76\xc7\xc3\xea\xd3\x2d\x5c\x9a\xc4\xd5\x54\x1d\x55\x7a\xb2\xe4\x08\x8d\x88\xfc\xca\xfa\x37\x53\xa7\x36\x73\xec\x13\xc7\xaf\x3b\xf5\x5e\x07\x17\x1a\x74\xcc\x62\x25\xcb\x95\x31\xeb\x9b\x2b\x94\xf6\xf7\xeb\x9b\xd2\x7a\xab\x4f\x0c\xa7\xca\xf8\xb3\xd4\xaa\xb1\x2b\x23\xed\xf9\xc0\x27\xc9\xf9\xcb\xf9\xf2\x13\x9f\xd4\x35\xfb\x5f\xa9\xce\x34\xb5\x76\x5d\x39\x9b\x77\x81\xc8\xaa\x1e\x92\x26\xc4\xb2\xf2\x5f\x68\x6c\x9a\x42\x81\x85\xbb\x60\xc4\x22\x63\x76\x09\x38\xc3\x42\xc6\x99\x6f\x38\x31\x47\x39\xa0\xb3\x84\x4c\x92\x62\x4f\xc7\x94\xa4\x61\x99\x5a\x1e\x78\x49\x19\x12\x9b\x7d\xc7\x90\xff\x78\x7e\x85\x45\x87\xd9\xcc\x98\x3b\x4f\x49\x44\xd0\xcc\x31\xe2\x25\x41\x7f\x7b\x8f\x92\xf4\x12\xc3\xd4\xb3\x2e\x77\x3b\x39\x10\x8b\x69\xb0\x3f\xbc\x63\xb0\xbd\xba\x90\xb1\xd3\x38\x4a\xae\x6f\xd7\xe6\xd5\xa8\xac\x18\x34\x3e\x2d\x8d\x79\x57\x72\x50\xb1\x91\x68\xd3\xdc\xbc\xf1\x0f\xdc\x30\x1e\x5c\x4d\x6c\x29\x9d\xf9\x1e\x00\x00\xff\xff\x86\xa5\x27\x98\xd2\x03\x00\x00")

func review_requestTxtBytes() ([]byte, error) {
	return bindataRead(
		_review_requestTxt,
		"review_request.txt",
	)
}

func review_requestTxt() (*asset, error) {
	bytes, err := review_requestTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "review_request.txt", size: 978, mode: os.FileMode(0644), modTime: time.Unix(1634832170, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x8a, 0x58, 0xb6, 0xd4, 0x10, 0xd8, 0x3d, 0xe1, 0xce, 0x64, 0x58, 0x97, 0x33, 0xa8, 0xef, 0x45, 0xf5, 0xa7, 0xc7, 0xd5, 0x9f, 0xe4, 0x56, 0xdb, 0xc2, 0xa5, 0xc5, 0xec, 0x21, 0xc1, 0xa3, 0xe5}}
	return a, nil
}

var _verify_contactHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x92\xcd\x8e\x9b\x40\x0c\x80\xef\x3c\x85\x95\x4b\x2f\x29\xdc\x23\x8a\xfa\x27\xb5\x2b\x55\x3d\x6c\xd2\xde\xcd\xe0\x84\x11\xc3\x98\xda\x26\x11\x5a\xe5\xdd\xab\x99\x91\xd2\xed\x61\x6f\xe0\xdf\xcf\x1f\xb4\x4b\xf7\x9d\x42\x60\x78\x79\x81\xfa\x27\xce\x04\xf7\xfb\xbe\x6d\x96\xae\xaa\xda\xa5\x3b\x8d\x18\x27\xd8\x78\x85\x33\x0b\xe8\xda\xcf\xde\xcc\xc7\x0b\x20\x9c\x9e\x9f\x8e\x9f\x20\x92\xdd\x58\x26\x10\xba\x78\x35\x41\xf3\x1c\x41\xe8\xcf\x4a\x6a\x35\x9c\x18\x7a\xba\xf8\x08\x36\x12\x08\x5d\x3d\xdd\xaa\x45\xd8\x91\xea\x1e\x96\x40\xa8\x04\x2d\xc2\x28\x74\xfe\xb0\x4b\x00\xbf\x49\xfc\x79\xfb\xc2\xd1\xd0\xd9\xaf\xe7\x1f\x70\xbf\xef\xba\x6b\x0e\x26\x0a\x01\x9a\xd1\x07\xe8\x37\x70\xc1\xbb\x29\x91\xd8\xe8\xb5\x0a\x3e\x4e\x75\xdb\x60\x57\xc8\x9f\xce\x99\xd9\x61\x8c\x6c\xa5\x14\xb8\x50\xa4\xca\xc7\x6e\xc7\xcb\x06\x18\x07\x58\x50\x8d\x1e\x79\xe8\x29\xf0\x0d\x7c\x34\xce\x5b\xab\x5e\xf8\xa6\x24\x80\xc3\x20\xa4\x0a\x3d\xca\xe1\xe1\xa8\x75\x3c\x50\xf7\x06\x7d\xdb\xe4\xec\x6b\xac\x77\x42\x30\xe2\x35\xb3\x0b\xaf\x7d\x20\x28\x17\xa6\x48\x3e\xd2\x95\x09\xe5\xd8\x57\xb0\x25\x9a\x28\x8b\x7c\x1c\x66\x1f\xb5\x42\xfb\x27\x31\x75\x18\x1f\x72\xe6\xa3\x89\x57\xac\x3d\xef\xba\xff\xdf\x93\x28\x60\x49\x4a\xca\xba\x79\x5e\xa3\xb7\xad\x7a\x8c\x19\xcd\x16\x3d\x34\x4d\xee\x78\x9f\xbe\xb0\x2e\xe8\xa8\xd6\x80\x6e\xaa\x1d\xcf\xbb\xee\x98\x1e\xc1\x8d\x18\x23\x85\x34\xb1\xce\x36\x3e\x93\x5a\xfa\x1b\x50\x06\xdd\xb7\xbd\x40\xd3\x55\x85\xf6\x5b\xe0\x1e\x03\x7c\xf5\x42\xce\x58\x36\x38\x92\x5c\xbd\x23\x38\x11\xce\x59\xd0\xdf\x00\x00\x00\xff\xff\x76\x5a\xcf\xdb\x8c\x02\x00\x00")

func verify_contactHtmlBytes() ([]byte, error) {
	return bindataRead(
		_verify_contactHtml,
		"verify_contact.html",
	)
}

func verify_contactHtml() (*asset, error) {
	bytes, err := verify_contactHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "verify_contact.html", size: 652, mode: os.FileMode(0644), modTime: time.Unix(1622582599, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xb2, 0xbc, 0x1a, 0xa9, 0x34, 0x54, 0x24, 0x29, 0xe9, 0x94, 0x3c, 0xc1, 0xaf, 0xcd, 0xbe, 0x61, 0xa2, 0xcc, 0x6a, 0xbb, 0xa9, 0xe4, 0xda, 0x44, 0x6e, 0x20, 0x3f, 0xe1, 0x28, 0xe6, 0xba, 0xed}}
	return a, nil
}

var _verify_contactTxt = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x91\xbd\x8e\xdc\x30\x0c\x84\x7b\x3d\xc5\x74\x69\x36\x7e\x80\x54\xf9\x03\x92\x03\x82\x14\xb7\x9b\xf4\xb4\xcc\x5d\x13\x96\x45\x87\xa2\x77\x61\x1c\xfc\xee\x81\x2c\xe4\x90\xe6\x3a\x63\x38\xe6\x37\x23\x7e\xe7\x94\x14\x2f\x2f\xe8\x7e\xd2\xcc\xd8\xf7\x53\x08\x97\x91\xf2\x84\x4d\x57\x5c\xd5\x50\xd6\x7e\x16\x77\xc9\x37\x10\x2e\xcf\x4f\xe7\x4f\xc8\xec\x0f\xb5\x09\xc6\x37\x29\x6e\xe4\xa2\x19\xc6\x7f\x56\x2e\xde\xe1\xa2\xe8\xf9\x26\x19\x3e\x32\x8c\xef\xc2\x0f\x2c\xa6\x91\x4b\x39\x61\x49\x4c\x85\x0f\xe0\x6f\x36\xb9\x6e\x5f\x34\x3b\x45\xff\xf5\xfc\x03\xfb\x7e\x3f\xa4\x8a\x36\xf0\x4c\x92\xd0\x6f\x88\x49\xe2\x54\xf1\x3e\x4a\x41\x92\x3c\x75\x21\x3c\x5d\x8f\x80\x91\x72\x56\x6f\x16\x68\x43\x56\xc7\x2b\x28\xea\xb2\x81\xf2\x80\x85\x8a\xf3\xeb\x1c\x3d\x27\x7d\x40\xb2\x6b\xa3\xf5\xa6\x8f\xc2\x06\x1a\x06\xe3\x52\xd0\x93\x7d\x08\xe1\x8d\x98\xff\xf0\xef\x8c\x31\xd2\xfd\xc8\x66\xba\xf6\x89\xd1\x1a\x54\xe5\x58\x1b\xdb\x6f\xad\xcc\x7f\xa1\x9a\x5a\xd3\xb4\x17\xa5\x61\x96\x5c\x40\xde\xbe\x3e\xba\x49\xa1\x4e\x14\x6a\xb5\x56\x5b\x35\xcf\x6b\x16\xdf\x70\x4e\x14\x27\xc4\x91\x72\xe6\x84\xc3\xfa\xbe\xde\xa3\x2c\x14\xb9\x2b\x75\xda\x45\x9d\xbb\x10\x3e\x73\xf1\x7a\x25\xb2\xa1\x9c\x42\x43\x7d\x4b\xda\x53\xc2\x57\x31\x8e\xae\xb6\xe1\xcc\x76\x97\xc8\xb8\x30\xcd\xe1\x6f\x00\x00\x00\xff\xff\xb4\xa7\x17\x83\x10\x02\x00\x00")

func verify_contactTxtBytes() ([]byte, error) {
	return bindataRead(
		_verify_contactTxt,
		"verify_contact.txt",
	)
}

func verify_contactTxt() (*asset, error) {
	bytes, err := verify_contactTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "verify_contact.txt", size: 528, mode: os.FileMode(0644), modTime: time.Unix(1622582599, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x1a, 0xe1, 0x94, 0x81, 0xf5, 0x2a, 0x4c, 0x11, 0x5d, 0xdb, 0xfb, 0xfa, 0x72, 0xc, 0x8b, 0x3f, 0xf, 0x52, 0x18, 0x43, 0x74, 0x20, 0x38, 0x6d, 0xc8, 0x3e, 0x62, 0xf, 0x6f, 0x7f, 0xd8, 0xc}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"deliver_certs.html":       deliver_certsHtml,
	"deliver_certs.txt":        deliver_certsTxt,
	"reject_registration.html": reject_registrationHtml,
	"reject_registration.txt":  reject_registrationTxt,
	"review_request.html":      review_requestHtml,
	"review_request.txt":       review_requestTxt,
	"verify_contact.html":      verify_contactHtml,
	"verify_contact.txt":       verify_contactTxt,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"deliver_certs.html": {deliver_certsHtml, map[string]*bintree{}},
	"deliver_certs.txt": {deliver_certsTxt, map[string]*bintree{}},
	"reject_registration.html": {reject_registrationHtml, map[string]*bintree{}},
	"reject_registration.txt": {reject_registrationTxt, map[string]*bintree{}},
	"review_request.html": {review_requestHtml, map[string]*bintree{}},
	"review_request.txt": {review_requestTxt, map[string]*bintree{}},
	"verify_contact.html": {verify_contactHtml, map[string]*bintree{}},
	"verify_contact.txt": {verify_contactTxt, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
