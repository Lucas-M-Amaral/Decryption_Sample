package flow

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
)

func Process(data []byte) (encData []byte) {

	word := "mamaco"
	pubKey, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		fmt.Println("ERROR(L17): ", err)
		return nil
	}
	fmt.Println("\nPublic Key: ", pubKey)

	pubKeyRsa := pubKey.(*rsa.PublicKey) //casting from any key to RSA

	encData, err = rsa.EncryptPKCS1v15(rand.Reader, pubKeyRsa, []byte(word)) //It could be this one as well: rsa.EncryptOAEP(hash, random, pub, msg, label)
	if err != nil {
		fmt.Println("ENCRYPTION FAILED!")
		return nil
	}
	fmt.Printf("\nEncrypted Data: %x\n", encData)

	return encData
}
