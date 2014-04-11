package main

type Material struct {
	color      Vector3
	diffuse    float64
	specular   float64
	absorb     float64
	reflection float64
	refraction float64
	refrIdx    float64
	fresnel    bool
}
