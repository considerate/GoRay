package main

import "math"

type Vector3 [3]float64

func (v Vector3) add(o Vector3) Vector3 {
	v[0] += o[0]
	v[1] += o[1]
	v[2] += o[2]
	return v
}

func (v Vector3) sub(o Vector3) Vector3 {
	v[0] -= o[0]
	v[1] -= o[1]
	v[2] -= o[2]
	return v
}

func (v Vector3) dot(o Vector3) float64 {
	return v[0]*o[0] + v[1]*o[1] + v[2]*o[2]
}

func (v Vector3) mult(o Vector3) Vector3 {
	v[0] *= o[0]
	v[1] *= o[1]
	v[2] *= o[2]
	return v
}

func (v Vector3) multScalar(x float64) Vector3 {
	v[0] *= x
	v[1] *= x
	v[2] *= x
	return v
}

func (v Vector3) len() float64 {
	return math.Sqrt(v.dot(v))
}

func (v Vector3) norm() Vector3 {
	v = v.multScalar(1.0 / v.len())
	return v
}

func (v Vector3) reflect(normal Vector3) Vector3 {
	v = v.sub(normal.multScalar(2.0 * v.dot(normal)))
	return v
}

func (v Vector3) refract(normal Vector3, ni, nt float64) Vector3 {
	n := ni / nt
	cosθi := v.dot(normal)              //Cosine of incoming angle
	sinθi := math.Sqrt(1 - cosθi*cosθi) //Trigonometric 1 (Sine of incoming angle)
	sinθt := n * sinθi
	cosθt := math.Sqrt(1.0 - sinθt*sinθt)
	v = v.multScalar(n).sub(normal.multScalar(n + cosθt))
	return v
}

func (v Vector3) copy() Vector3 {
	return Vector3{v[0], v[1], v[2]}
}
