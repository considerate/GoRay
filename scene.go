package main

import "github.com/unit3/kdtree"

type Scene struct {
	objects   []Object
	lights    []Light
	photonMap *kdtree.Tree
}

func setupScene() Scene {
	materials := []Material{
		Material{
			Vector3{1, 0, 0}, //color
			0.6,              //diffuse
			0.4,              //specular
			0.0,              //absorb
			1.0,              //reflection
			0.3,              //refraction
			1.3,              //refraction index
			true,
		},
		Material{
			Vector3{0, 1, 0}, //color
			0.8,              //diffuse
			0.2,              //specular
			0.0,              //absorb
			2.0,              //reflection
			0.0,              //refraction
			1.0,              //refraction index
			false,
		},
		Material{
			Vector3{0, 0, 1}, //color
			0.8,              //diffuse
			0.0,              //specular
			0.2,              //absorb
			0.0,              //reflection
			0.0,              //refraction
			1.0,              //refraction index
			false,
		},
		Material{
			Vector3{0, 1, 0}, //color
			0.8,              //diffuse
			0.0,              //specular
			0.2,              //absorb
			0.0,              //reflection
			0.0,              //refraction
			1.0,              //refraction index
			false,
		},
		Material{
			Vector3{1, 0, 0}, //color
			0.8,              //diffuse
			0.0,              //specular
			0.2,              //absorb
			0.0,              //reflection
			0.0,              //refraction
			1.0,              //refraction index
			false,
		},
		Material{
			Vector3{1, 1, 1}, //color
			0.8,              //diffuse
			0.0,              //specular
			0.2,              //absorb
			0.0,              //reflection
			0.0,              //refraction
			1.0,              //refraction index
			false,
		},
	}
	spheres := []Sphere{
		Sphere{
			Vector3{0, 1, 10},
			2,
			&materials[0],
		},
		Sphere{
			Vector3{-8, 1, 12},
			2,
			&materials[1],
		},
	}
	lights := []Light{
		Light{
			Vector3{0, 5, -5},
		},
		Light{
			Vector3{0, 5, 12},
		},
	}
	planes := []Plane{
		Plane{
			Vector3{0, -2, 0},
			Vector3{0, 1, 0},
			&materials[2],
		},
		Plane{
			Vector3{20, 0, 0},
			Vector3{-1, 0, 0},
			&materials[3],
		},
		Plane{
			Vector3{-20, 0, 0},
			Vector3{1, 0, 0},
			&materials[4],
		},
		Plane{
			Vector3{0, 0, 30},
			Vector3{0, 0, -1},
			&materials[5],
		},
		Plane{
			Vector3{0, 0, -10},
			Vector3{0, 0, 1},
			&materials[5],
		},
	}

	objects := make([]Object, len(spheres)+len(planes))
	count := 0
	for _, sphere := range spheres {
		objects[count] = sphere
		count++
	}
	for _, plane := range planes {
		objects[count] = plane
		count++
	}
	return Scene{
		objects,
		lights,
		nil,
	}
}
