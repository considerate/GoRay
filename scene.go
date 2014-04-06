package main

type Scene struct {
	objects []Object
	lights  []Light
	photons []Photon
}

func setupScene() Scene {
	materials := []Material{
		Material{
			Vector3{1, 0, 0},
			0.8,
			0.3,
			0.5,
		},
		Material{
			Vector3{0, 1, 0},
			0.8,
			0.3,
			0.9,
		},
		Material{
			Vector3{0, 0, 1},
			0.8,
			0.3,
			0,
		},
		Material{
			Vector3{0, 1, 0},
			0.8,
			0.3,
			0,
		},
		Material{
			Vector3{1, 0, 0},
			0.8,
			0.3,
			0,
		},
		Material{
			Vector3{1, 1, 1},
			0.8,
			0.3,
			0,
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
			Vector3{0, 20, 5},
		},
		Light{
			Vector3{0, 10, 12},
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
