package main

type Plane struct {
	position  Vector3
	n         Vector3
	_material *Material
}

func (plane Plane) normal(Vector3) Vector3 {
	return plane.n
}

func (plane Plane) intersect(ray *Ray) (t1, t2 float64) {
	dir := ray.dir
	o := ray.origin
	n := plane.n
	p := plane.position
	t := (p.sub(o)).dot(n) / (dir.dot(n))
	return t, t
}

func (plane Plane) material() *Material {
	return plane._material
}
