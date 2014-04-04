package main

import (
	"github.com/skelterjohn/go.matrix"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sync"
)

const WIDTH = 800
const HEIGHT = 640
const aspect float64 = float64(WIDTH) / float64(HEIGHT)
const fov float64 = math.Pi / 2
const near float64 = 0
const far float64 = 1
const MAX_TRACE_DEPTH = 5
const rayEpsilon = 0.001

func main() {
	scale := 1 / math.Tan(fov/2)
	focalLength := far - near
	camera := matrix.MakeDenseMatrix([]float64{
		scale * aspect, 0, 0, 0,
		0, -scale, 0, 0,
		0, 0, far / focalLength, 1,
		0, 0, (far * near) / focalLength, 0},
		4, 4)
	//Create scene (scene.go)
	scene := setupScene()

	//Render image
	render(camera, scene)
}

func matrixToVector3(mat *matrix.DenseMatrix) Vector3 {
	arr := mat.Array()
	return Vector3{arr[0], arr[1], arr[2]}
}

func shadePoint(point, normal Vector3, object Object, lights []Light) Vector3 {
	material := object.material()

	resultingColor := Vector3{0, 0, 0}
	if material.diffuse > 0 {
		for _, light := range lights {
			lightDir := light.position.sub(point).norm()
			dot := normal.dot(lightDir)
			if dot > 0 {
				diffuse := material.diffuse * dot
				color := material.color
				resultingColor = resultingColor.add(color.multScalar(diffuse))
			}
		}
	}
	return resultingColor
}

func raytrace(ray *Ray, objects []Object, lights []Light, tracedepth uint8) (bool, Vector3) {
	var t float64
	t = math.MaxFloat64
	hit := false
	inside := false
	var hitObject Object
	for _, object := range objects {
		t1, t2 := object.intersect(ray)

		if t1 >= 0 && t1 <= t {
			t = t1
			hit = true
			hitObject = object
		} else if t2 >= 0 && t2 <= t {
			t = t2
			inside = true
			hit = true
			hitObject = object
		}
	}
	if hit {
		object := hitObject
		point := ray.origin.add(ray.dir.multScalar(t))
		normal := object.normal(point)
		if inside {
			normal = normal.multScalar(-1)
		}
		color := shadePoint(point, normal, object, lights)

		material := object.material()
		if material.reflection > 0 && tracedepth < MAX_TRACE_DEPTH {
			r := ray.dir.reflect(normal)
			newRay := &Ray{point.add(r.multScalar(rayEpsilon)), r}
			nextHit, nextColor := raytrace(newRay, objects, lights, tracedepth+1)
			if nextHit {
				color = color.add(nextColor.multScalar(material.reflection))
			}
		}
		return true, color
	}
	return false, Vector3{}
}

func render(camera *matrix.DenseMatrix, scene Scene) {
	bounds := image.Rect(int(-WIDTH/2), int(-HEIGHT/2), int(WIDTH/2), int(HEIGHT/2))
	image := image.NewRGBA(bounds)
	origin := Vector3{0, 0, 0}
	wg := new(sync.WaitGroup)
	wg.Add(HEIGHT)
	for v := -HEIGHT / 2; v < HEIGHT/2; v++ {
		go func(v int) {
			for u := -WIDTH / 2; u < WIDTH/2; u++ {
				imgX, imgY := 2.0*float64(u)/float64(WIDTH), 2*float64(v)/float64(HEIGHT)
				dir := []float64{imgX, imgY, 1, 1}
				dirMatrix, _ := matrix.MakeDenseMatrix(dir, 1, 4).TimesDense(camera)
				rayDir := matrixToVector3(dirMatrix).norm()
				ray := &Ray{origin, rayDir}
				image.Set(u, v, color.Black)
				hit, hitColor := raytrace(ray, scene.objects, scene.lights, 0)
				if hit {
					red := uint8(math.Min(hitColor[0], 1) * 255)
					green := uint8(math.Min(hitColor[1], 1) * 255)
					blue := uint8(math.Min(hitColor[2], 1) * 255)
					alpha := uint8(255)
					image.Set(u, v, color.RGBA{red, green, blue, alpha})
				}
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	filename := "output.png"
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	png.Encode(file, image)
}
