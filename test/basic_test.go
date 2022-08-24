package main

import (
	"github.com/David-Durst/head-position-model/pkg/headpositionmodel"
	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"testing"
)

func TestCallModel(t *testing.T) {
	headpositionmodel.ModelHeadPosition(r3.Vector{X: 50., Y: 50., Z: 50.}, r2.Point{Y: 20.}, 1.)
}
