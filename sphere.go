package main

import "math"

type Sphere struct {
	center    Vector3
	radius    float64
	_material *Material
}

func (sphere Sphere) intersect(ray *Ray) (t1, t2 float64) {
	dir := ray.dir
	c := sphere.center
	o := ray.origin
	r := sphere.radius
	sphereDist := o.sub(c)
	a := dir.dot(sphereDist)
	b := sphereDist.dot(sphereDist)
	t1 = -a - math.Sqrt(a*a-b+r*r)
	t2 = -a + math.Sqrt(a*a-b+r*r)
	return t1, t2
}

func (sphere Sphere) normal(point Vector3) Vector3 {
	return point.sub(sphere.center).norm()
}

func (sphere Sphere) material() *Material {
	return sphere._material
}
