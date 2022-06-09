package vgg_test

import (
	"testing"
	"lgdSearch/pkg/vgg"
)

func TestGetFeature(t *testing.T) {
	result, err := vgg.GetFeature("./vgg_test/char.jpg")
	if err != nil {
		t.Error(err.Error())
	}
	println(result)
}