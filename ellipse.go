package main

import "github.com/skelterjohn/go.matrix"

type Ellipse struct {
	position  Vector3
	radius    float64
	transform *matrix.DenseMatrix
}

func (ellipse Ellipse) intersect(ray *Ray) (t1, t2 float64) {
	return t1, t2
}
