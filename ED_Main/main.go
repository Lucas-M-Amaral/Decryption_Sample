package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/go-tpm/legacy/tpm2"
)

func main() {
	tpm, err := tpm2.OpenTPM()
	if err != nil {
		fmt.Println(err)
		fmt.Println((log.Lshortfile))
	}

	keyTemplate := tpm2.Public{ //we specify a template cause we were having policy issues with the default template
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagFixedParent | tpm2.FlagFixedTPM | tpm2.FlagSensitiveDataOrigin | tpm2.FlagUserWithAuth | tpm2.FlagDecrypt,
		AuthPolicy: nil,
		RSAParameters: &tpm2.RSAParams{
			KeyBits:    2048,
			ModulusRaw: make([]byte, 256),
		},
	}

	keyHandle, outPublic, err := tpm2.CreatePrimary(tpm, tpm2.HandleOwner, tpm2.PCRSelection{}, "", "", keyTemplate)
	if err != nil {
		fmt.Println(err)
		fmt.Println((log.Lshortfile))
	}

	fmt.Println(tpm2.ReadPublic(tpm, keyHandle))
	fmt.Println("\nPublic part: \n", outPublic)

	pk, err := x509.MarshalPKIXPublicKey(outPublic) //converts outPublic type to bytes
	if err != nil {
		fmt.Println("ERROR(L42): ", err)
	}
	payload, err := json.Marshal(pk)
	if err != nil {
		fmt.Println("ERROR(L48): ", err)
	}

	fmt.Println("\nPayload: \n", payload)

	//From here beyond, we are treating the response

	response, err := http.Post("http://localhost:8080/keyData", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("ERROR(L47): ", err)
	}

	defer response.Body.Close() //sets the json body to read only (closes it)

	dataReceived, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ERROR(L53): ", err)
	}

	var tResponse []byte                           //
	err = json.Unmarshal(dataReceived, &tResponse) // conversion of the received json to an understandable structure (treatedResponse)
	if err != nil {                                //
		fmt.Println("ERROR(L59): ", err)
	}

	fmt.Printf("\nTreated Response: %x\n", tResponse)

	decData, err := tpm2.RSADecrypt(tpm, keyHandle, "", tResponse, nil, "")
	if err != nil {
		fmt.Println("ERROR(L66): ", err)
	}

	fmt.Printf("\nDecrypted data: %x\n", decData[len(decData)-6:]) //starts printing after the padding (the size of the data minus the size of the word)
	fmt.Printf("\nDecrypted data (string): %s\n", string(decData))

	defer tpm.Close()
}
