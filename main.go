package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	binarycodec "tx_decoder/binary-codec"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("The script expects exactly one input")
		fmt.Println("Usage: decoder [blob] where blob is either hex(blob) encoded object, base64(blob) encoded object or base64(hex(blob)) encoded object")
		return
	}

	blob := os.Args[1]

	if isHex(blob) {
		err := decodeHexAndPrint(blob)
		if err != nil {
			log.Printf("Failed to decode blob: %v\n", err)
			return
		}

		return
	}

	log.Println("blob is not in hex format; trying base64...")
	blob, ok := decodeBase64(blob)
	if ok {
		// if the string is already hex encoded
		if isHex(blob) {
			err := decodeHexAndPrint(blob)
			if err != nil {
				log.Printf("Failed to decode blob: %v\n", err)
				return
			}

			return
		}

		// last try - encode in hex and try to decode
		blob = hex.EncodeToString([]byte(blob))
		err := decodeHexAndPrint(blob)
		if err != nil {
			log.Printf("Failed to decode blob: %v\n", err)
			return
		}

		return
	}

	log.Println("blob is not in hex or base64")

	os.Exit(1)
}

func decodeHexAndPrint(str string) error {
	if !isHex(str) {
		return errors.New("string is not in hex format")
	}

	m, err := binarycodec.Decode(str)
	if err != nil {
		return err
	}

	printMap(m)
	return nil

}

func isHex(str string) bool {
	_, err := hex.DecodeString(str)

	return err == nil
}

func decodeBase64(str string) (string, bool) {
	decoded, err := base64.StdEncoding.DecodeString(str)

	return string(decoded), err == nil
}

func printMap(m map[string]any) {
	str, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Println("failed to marshal json: ", err)
		os.Exit(1)
	}

	fmt.Println(string(str))
}
