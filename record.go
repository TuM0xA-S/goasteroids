package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func loadRecord() int {
	file, err := os.Open(RecordFileName)
	if err != nil {
		return 0
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return 0
	}

	Cipher.Decrypt(data, data)

	res, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}

	return res
}

func saveRecord(val int) {
	file, err := os.Create(RecordFileName)
	if err != nil {
		return
	}
	defer file.Close()

	data := []byte(fmt.Sprintf("%16d", val))

	Cipher.Encrypt(data, data)
	file.Write(data)
}
