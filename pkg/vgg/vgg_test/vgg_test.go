package vgg_test

import (
	"testing"
	"lgdSearch/pkg/vgg"
	"os"
	"io/ioutil"
)

func TestGetFeature(t *testing.T) {
	fp, err := os.Open("char.jpg")
	if err != nil {
		t.Error(err.Error())
	}
	defer fp.Close()
	bytes, err := ioutil.ReadAll(fp)
	if err != nil {
		t.Error(err.Error())
	}
	result, err := vgg.GetFeature(bytes)
	if err != nil {
		t.Error(err.Error())
	}
	println(result)
}