package main

type Scene struct {
	objects []Object
	lights  []Light
}

func setupScene() Scene {
	spheres := make([]Sphere, 2)
	material := &Material{
		Vector3{1, 0, 0},
		0.8,
		0.3,
		0.5,
	}
	material2 := &Material{
		Vector3{0, 1, 0},
		0.8,
		0.3,
		0.9,
	}

	spheres[0] = Sphere{
		Vector3{0, 1, 10},
		2,
		material,
	}
	spheres[1] = Sphere{
		Vector3{-8, 1, 12},
		2,
		material2,
	}
	lights := make([]Light, 2)
	lights[0] = Light{
		Vector3{0, 0, 5},
	}
	lights[1] = Light{
		Vector3{0, 10, 5},
	}

	planes := make([]Plane, 3)

	material3 := &Material{
		Vector3{0, 0, 1},
		0.8,
		0.3,
		0,
	}

	material4 := &Material{
		Vector3{0, 1, 0},
		0.8,
		0.3,
		0,
	}

	material5 := &Material{
		Vector3{1, 0, 0},
		0.8,
		0.3,
		0,
	}
	planes[0] = Plane{
		Vector3{0, -2, 0},
		Vector3{0, 1, 0},
		material3,
	}

	planes[1] = Plane{
		Vector3{20, 0, 0},
		Vector3{-1, 0, 0},
		material4,
	}

	planes[2] = Plane{
		Vector3{-20, 0, 0},
		Vector3{1, 0, 0},
		material5,
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
	}
}
