package types

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/cloudfoundry-community/bosh-vault/secret"
	"golang.org/x/crypto/ssh"
)

const SshKeypairType = "ssh"

type SshKeypairRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (record SshKeypairRecord) Store(secretStore secret.Store, name string) (CredentialResponse, error) {
	var resp SshKeypairResponse
	secretId, err := secretStore.Set(name, map[string]interface{}{
		"public_key":             record.PublicKey,
		"private_key":            record.PrivateKey,
		"public_key_fingerprint": record.PublicKeyFingerprint,
	})
	if err != nil {
		return resp, err
	}

	resp = SshKeypairResponse{
		Name:  name,
		Id:    secretId,
		Value: record,
	}

	return resp, nil
}

func (r *SshKeypairRequest) Validate() bool {
	return r.Type == SshKeypairType
}

func (r *SshKeypairRequest) Generate(secretStore secret.Store) (CredentialRecordInterface, error) {

	privKey, err := rsa.GenerateKey(rand.Reader, RsaKeySizeBits)
	if err != nil {
		return nil, err
	}

	pemPriv := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})

	pubKey, err := ssh.NewPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(pubKey)

	sshKeyPair := SshKeypairRecord{
		PrivateKey:           string(pemPriv),
		PublicKey:            string(pubKeyBytes),
		PublicKeyFingerprint: ssh.FingerprintLegacyMD5(pubKey),
	}

	return sshKeyPair, nil
}

func (r *SshKeypairRequest) CredentialType() string {
	return r.Type
}

func (r *SshKeypairRequest) CredentialName() string {
	return r.Name
}

type SshKeypairRecord struct {
	PublicKey            string `json:"public_key"`
	PrivateKey           string `json:"private_key"`
	PublicKeyFingerprint string `json:"public_key_fingerprint"`
}

type SshKeypairResponse struct {
	Name  string           `json:"name"`
	Id    string           `json:"id"`
	Value SshKeypairRecord `json:"value"`
}
