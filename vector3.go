package main

import "math"

type Vector3 [3]float64

func (v Vector3) add(o Vector3) Vector3 {
	return Vector3{v[0] + o[0], v[1] + o[1], v[2] + o[2]}
}

func (v Vector3) sub(o Vector3) Vector3 {
	return Vector3{v[0] - o[0], v[1] - o[1], v[2] - o[2]}
}

func (v Vector3) dot(o Vector3) float64 {
	return v[0]*o[0] + v[1]*o[1] + v[2]*o[2]
}

func (v Vector3) mult(o Vector3) Vector3 {
	return Vector3{v[0] * o[0], v[1] * o[1], v[2] * o[2]}
}

func (v Vector3) multScalar(x float64) Vector3 {
	return Vector3{v[0] * x, v[1] * x, v[2] * x}
}

func (v Vector3) len() float64 {
	return math.Sqrt(v.dot(v))
}

func (v Vector3) norm() Vector3 {
	return v.multScalar(1.0 / v.len())
}

func (v Vector3) reflect(normal Vector3) Vector3 {
	return v.sub(normal.multScalar(2.0 * v.dot(normal)))
}
