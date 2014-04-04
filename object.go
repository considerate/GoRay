package main

type Object interface {
	normal(Vector3) Vector3
	intersect(*Ray) (float64, float64)
	material() *Material
}
