package generator

// mainly copied from https://golang.org/src/crypto/tls/generate_cert.go

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/trusch/pkid/entity"
	"github.com/trusch/pkid/types"
)

type Options struct {
	Name      string
	NotBefore time.Time
	ValidFor  time.Duration
	IsCA      bool
	RsaBits   int
	Curve     string
	Usage     x509.ExtKeyUsage
}

func (options *Options) fillDefaults() {
	if options.NotBefore.IsZero() {
		options.NotBefore = time.Now()
	}
	if options.ValidFor == 0 {
		options.ValidFor = 365 * 24 * time.Hour
	}
	if options.Curve == "" && options.RsaBits == 0 {
		options.Curve = "P521"
	}
	if options.Usage == 0 && !options.IsCA {
		options.Usage = x509.ExtKeyUsageAny
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			return nil
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func generateKey(rsaBits int, curve string) (interface{}, error) {
	var (
		priv interface{}
		err  error
	)
	switch curve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, rsaBits)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		err = errors.New("unknown elliptic curve (try P224 P256 P384 or P521)")
	}
	return priv, err
}

func getSerial(ca *types.CAEntity) (*big.Int, error) {
	var serialNumber *big.Int
	if ca == nil {
		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
		serial, err := rand.Int(rand.Reader, serialNumberLimit)
		if err != nil {
			return nil, fmt.Errorf("failed to generate serial number: %s", err)
		}
		serialNumber = serial
	} else {
		serialNumber = ca.Serial
	}
	return serialNumber, nil
}

func getSignerCertAndKey(template x509.Certificate, priv interface{}, caEntity *types.CAEntity) (*x509.Certificate, interface{}, error) {
	signerCert := &template
	signerKey := priv
	if caEntity != nil {
		ca, err := entity.NewEntityFromPEM([]byte(caEntity.Cert), []byte(caEntity.Key))
		if err != nil {
			return nil, nil, err
		}
		signerCert = ca.Cert
		signerKey = ca.Key
	}
	return signerCert, signerKey, nil
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func Generate(ca *types.CAEntity, options *Options) (*types.Entity, error) {
	options.fillDefaults()
	priv, err := generateKey(options.RsaBits, options.Curve)
	if err != nil {
		return nil, err
	}
	serial, err := getSerial(ca)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
			CommonName:   options.Name,
		},
		NotBefore:             options.NotBefore,
		NotAfter:              options.NotBefore.Add(options.ValidFor),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{options.Usage},
		BasicConstraintsValid: true,
	}
	if options.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}
	signerCert, signerKey, err := getSignerCertAndKey(template, priv, ca)
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, signerCert, publicKey(priv), signerKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to create certificate: %s", err)
	}

	certOut := &bytes.Buffer{}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut := &bytes.Buffer{}
	pem.Encode(keyOut, pemBlockForKey(priv))
	entity := &types.Entity{
		Name: options.Name,
		Cert: certOut.String(),
		Key:  keyOut.String(),
	}
	return entity, nil
}
