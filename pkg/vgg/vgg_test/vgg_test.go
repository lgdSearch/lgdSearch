package vgg_test

import (
	"testing"
)

func TestMain(m *testing.M) {
	vgg.LoadModel()
	m.Run()
}

func TestGetFeature(t *testing.T) {
	vgg.GetFeature()
}