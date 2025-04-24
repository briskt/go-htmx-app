package saml

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	gosaml2 "github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	goxmldsig "github.com/russellhaering/goxmldsig"
)

type Config struct {
	AssertionConsumerServiceURL string `json:"AssertionConsumerServiceURL"`
	AudienceURI                 string `json:"AudienceURI"`
	IDPMetadataURL              string `json:"IDPMetadataURL"`
	SPEntityID                  string `json:"SPEntityID"`
	SPPublicCert                string `json:"SPPublicCert"`
	SPPrivateKey                string `json:"SPPrivateKey"`
}

type Provider struct {
	gosaml2.SAMLServiceProvider
}

// GetKeyPair implements dsig.X509KeyStore interface
func (c *Config) GetKeyPair() (privateKey *rsa.PrivateKey, cert []byte, err error) {
	rsaKey, err := getRsaPrivateKey(c.SPPrivateKey, c.SPPublicCert)
	if err != nil {
		return &rsa.PrivateKey{}, []byte{}, err
	}

	return rsaKey, []byte(c.SPPublicCert), nil
}

func New(config Config) (*Provider, error) {
	res, err := http.Get(config.IDPMetadataURL)
	if err != nil {
		return nil, fmt.Errorf("error calling IdP metadata URL: %w", err)
	}

	rawMetadata, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading IdP metadata response: %w", err)
	}

	metadata := &types.EntityDescriptor{}
	err = xml.Unmarshal(rawMetadata, metadata)
	if err != nil {
		return nil, fmt.Errorf("error parsing IdP Metadata response: %w", err)
	}

	certStore := goxmldsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for i, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				return nil, fmt.Errorf("IdP metadata certificate #%d is empty", i)
			}
			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				return nil, fmt.Errorf("error decoding IdP cert data: %w", err)
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {
				return nil, fmt.Errorf("error parsing IdP cert: %w", err)
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}

	p := &Provider{
		gosaml2.SAMLServiceProvider{
			IdentityProviderSSOURL:      metadata.IDPSSODescriptor.SingleSignOnServices[0].Location,
			IdentityProviderSLOURL:      metadata.IDPSSODescriptor.SingleLogoutServices[0].Location,
			IdentityProviderIssuer:      metadata.EntityID,
			AssertionConsumerServiceURL: config.AssertionConsumerServiceURL,
			ServiceProviderIssuer:       config.SPEntityID,
			SignAuthnRequests:           true,
			AudienceURI:                 config.AudienceURI,
			IDPCertificateStore:         &certStore,

			// since Config implements goxmldsig.X509KeyStore interface, just pass in the pointer to our config
			SPKeyStore:        &config,
			SPSigningKeyStore: &config,
		},
	}
	return p, nil
}

func (p *Provider) GetUser(c echo.Context) (string, error) {
	samlResp := c.FormValue("SAMLResponse")
	if samlResp == "" {
		return "", fmt.Errorf("no SAML response provided in query")
	}

	info, err := p.RetrieveAssertionInfo(samlResp)
	if err != nil {
		return "", fmt.Errorf("invalid SAML assertion: %s", err)
	}

	if info.WarningInfo.InvalidTime {
		return "", fmt.Errorf("invalid SAML assertion time")
	}

	if info.WarningInfo.NotInAudience {
		return "", fmt.Errorf("invalid SAML assertion, not in audience")
	}
	attributes := info.Assertions[0].AttributeStatement.Attributes
	return getFirstValue("employeeNumber", attributes), nil
}

func getFirstValue(attrName string, attributes []types.Attribute) string {
	for _, attr := range attributes {
		if attr.Name != attrName {
			continue
		}

		if len(attr.Values) > 0 {
			return attr.Values[0].Value
		}
		return ""
	}
	return ""
}

func getRsaPrivateKey(privateKey, publicCert string) (*rsa.PrivateKey, error) {
	var rsaKey *rsa.PrivateKey

	if privateKey == "" {
		return rsaKey, errors.New("A valid PEM or base64 encoded privateKey is required")
	}

	if publicCert == "" {
		return rsaKey, errors.New("A valid PEM or base64 encoded publicCert is required")
	}

	privateKeyBytes, err := decodeKey(privateKey, "PRIVATE KEY")
	if err != nil {
		return nil, fmt.Errorf("problem with RSA private key: %w", err)
	}

	var parsedKey any
	if parsedKey, err = x509.ParsePKCS8PrivateKey(privateKeyBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS1PrivateKey(privateKeyBytes); err != nil {
			return rsaKey, fmt.Errorf("unable to parse RSA private key: %s", err)
		}
	}

	var ok bool
	rsaKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		return rsaKey, errors.New("unable to assert parsed key type")
	}

	publicCertBytes, err := decodeKey(publicCert, "CERTIFICATE")
	if err != nil {
		return nil, fmt.Errorf("problem with RSA public cert: %w", err)
	}

	cert, err := x509.ParseCertificate(publicCertBytes)
	if err != nil {
		return rsaKey, fmt.Errorf("unable to parse RSA public cert: %s", err)
	}

	var pubKey *rsa.PublicKey
	if pubKey, ok = cert.PublicKey.(*rsa.PublicKey); !ok {
		return rsaKey, errors.New("unable to assert RSA public cert type")
	}

	rsaKey.PublicKey = *pubKey

	return rsaKey, nil
}

// decodeKey decodes a key from either a PEM-encoded string or a base64 string
func decodeKey(key, expectedType string) ([]byte, error) {
	block, _ := pem.Decode([]byte(key))
	if block != nil {
		if block.Type != expectedType {
			return nil, fmt.Errorf("key is of the wrong type, expected %s but found %s", expectedType, block.Type)
		}
		return block.Bytes, nil
	}

	bytes := make([]byte, base64.StdEncoding.DecodedLen(len(key)))
	n, err := base64.StdEncoding.Decode(bytes, []byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return bytes[:n], nil
}
