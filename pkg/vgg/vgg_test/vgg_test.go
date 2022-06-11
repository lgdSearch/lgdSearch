package vgg_test

import (
	"testing"
	"lgdSearch/pkg/vgg"
	"os"
	"io/ioutil"
)

func TestSearch(t *testing.T) {
	fp, err := os.Open("char.jpg")
	if err != nil {
		t.Error(err.Error())
	}
	defer fp.Close()
	bytes, err := ioutil.ReadAll(fp)
	if err != nil {
		t.Error(err.Error())
	}
	result, err := vgg.Search(bytes)
	if err != nil {
		t.Error(err.Error())
	}
	fd, err := os.Create("./test.jpg")
	if err != nil {
		t.Error(err.Error())
	}
	defer fd.Close()
	_, err = fd.Write(result[0])
	if err != nil {
		t.Error(err.Error())
	}
}