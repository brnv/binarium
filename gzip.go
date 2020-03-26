package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func decompressIfGzipped(input []byte) []byte {
	var output []byte

	reader, err := gzip.NewReader(bytes.NewReader(input))
	if err != nil {
		logger.Infof("can't create reader: '%s', not gzipped", err.Error())
		output = input
	} else {
		logger.Info("reader ready, gzipped")
		output, _ = ioutil.ReadAll(reader)
	}

	return output
}
