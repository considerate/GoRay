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

func (sphere Sphere) contains(point Vector3) bool {
	r := sphere.radius
	rSq := r * r
	c := sphere.center
	p := point
	d := c[0] - p[0]
	dSq := d * d
	if dSq > rSq {
		return false
	}
	d = c[1] - p[1]
	dSq += d * d
	if dSq > rSq {
		return false
	}
	d = c[2] - p[2]
	dSq += d * d
	if dSq > rSq {
		return false
	}
	return true
}

func (sphere Sphere) material() *Material {
	return sphere._material
}
