package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"os"
)

type Results struct {
	MD5    string
	SHA1   string
	SHA256 string
	SHA512 string
}

func Compute(path string) (Results, error) {
	f, err := os.Open(path)
	if err != nil {
		return Results{}, err
	}
	defer f.Close()

	hMD5 := md5.New()
	hSHA1 := sha1.New()
	hSHA256 := sha256.New()
	hSHA512 := sha512.New()

	w := io.MultiWriter(hMD5, hSHA1, hSHA256, hSHA512)

	if _, err := io.Copy(w, f); err != nil {
		return Results{}, err
	}

	return Results{
		MD5:    fmt.Sprintf("%x", hMD5.Sum(nil)),
		SHA1:   fmt.Sprintf("%x", hSHA1.Sum(nil)),
		SHA256: fmt.Sprintf("%x", hSHA256.Sum(nil)),
		SHA512: fmt.Sprintf("%x", hSHA512.Sum(nil)),
	}, nil
}
