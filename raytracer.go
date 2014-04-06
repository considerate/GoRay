package main

import (
	"fmt"
	"github.com/skelterjohn/go.matrix"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

const WIDTH = 860
const HEIGHT = 640
const aspect float64 = float64(WIDTH) / float64(HEIGHT)
const fov float64 = math.Pi / 2.0
const near float64 = 1
const far float64 = 3
const MAX_TRACE_DEPTH uint8 = 5
const rayEpsilon = 0.001
const PHOTONS_PER_LIGHT = 100
const PHOTONS_DEPTH = 3
const photonMap = true
const photonRadius = 15.0

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	scale := 1 / math.Tan(fov/2)
	focalLength := far - near
	camera := matrix.MakeDenseMatrix([]float64{
		scale * aspect, 0, 0, 0,
		0, -scale, 0, 0,
		0, 0, far / focalLength, 0,
		0, 0, 0, 1},
		4, 4)
	fmt.Println("Generating photons")
	//Create scene (scene.go)
	scene := setupScene()
	if photonMap {
		scene.photons = shootPhotons(scene)
	}
	fmt.Println("Rendering image")
	//Render image
	render(camera, scene)
	fmt.Println("Done")
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

func raytrace(render bool, ray *Ray, scene Scene, depth ...uint8) (bool, Vector3, Vector3, Object) {
	tracedepth := depth[0]
	objects := scene.objects
	lights := scene.lights
	photons := scene.photons
	maxdepth := MAX_TRACE_DEPTH
	if len(depth) > 1 {
		maxdepth = depth[1]
	}
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
		var color Vector3
		if photonMap && render {
			color = Vector3{}

			photonArea := math.Pi * photonRadius * photonRadius
			photonSphere := Sphere{
				point,
				photonRadius,
				nil,
			}
			for _, photon := range photons {
				//fmt.Println(photon.position)
				if photonSphere.contains(photon.position) {

					weight := math.Max(0.0, -normal.dot(photon.incomingDirection))
					distV := photon.position.sub(point)
					distSq := distV.dot(distV)
					weight *= 1.0 - math.Sqrt(distSq)/photonRadius
					color = color.add(photon.color.multScalar(weight))
				}
			}
			color = color.multScalar(1.0 / photonArea * 3.0)
			//color = color.add(shadePoint(point, normal, object, lights).multScalar(0.3))
		} else {
			color = shadePoint(point, normal, object, lights)
		}

		material := object.material()
		if material.reflection > 0 && tracedepth < maxdepth-1 {
			r := ray.dir.reflect(normal)
			newRay := &Ray{point.add(r.multScalar(rayEpsilon)), r}
			nextHit, nextColor, _, _ := raytrace(true, newRay, scene, tracedepth+1, maxdepth)
			if nextHit {
				color = color.add(nextColor.multScalar(material.reflection))
			}
		}
		return true, color, point, object
	}
	return false, Vector3{}, Vector3{}, nil
}

func shootPhotons(scene Scene) []Photon {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	numLights := len(scene.lights)
	numPhotons := numLights * PHOTONS_PER_LIGHT * PHOTONS_DEPTH
	photons := make([]Photon, numPhotons)
	count := 0
	photonPower := 600.0 / PHOTONS_PER_LIGHT
	var origin, dir Vector3
	var ray *Ray
	var diffuse float64
	var color, normal Vector3
	var objectDiffuse float64
	var surfaceColor, position Vector3
	var hit bool
	var object Object
	var i, j int
	for _, light := range scene.lights {
		origin = light.position

		for i = 0; i < PHOTONS_PER_LIGHT; i++ {
			dir = Vector3{
				1.0 - 2*r.Float64(),
				1.0 - 2*r.Float64(),
				1.0 - 2*r.Float64(),
			}.norm()
			ray = &Ray{
				origin,
				dir,
			}
			diffuse = 1.0
			color = Vector3{}
			for j = 0; j < PHOTONS_DEPTH; j++ {
				hit, surfaceColor, position, object = raytrace(false, ray, scene, 0, 1)
				if !hit {
					break
				}
				objectDiffuse = object.material().diffuse
				diffuse *= objectDiffuse
				if objectDiffuse <= 0 {
					break
				}
				color = color.add(surfaceColor.multScalar(diffuse))
				photons[count] = Photon{
					position,
					color.multScalar(photonPower),
					ray.dir.copy(),
				}
				count++
				normal = object.normal(position)
				ray.dir = Vector3{
					1.0 - 2*r.Float64(),
					1.0 - 2*r.Float64(),
					1.0 - 2*r.Float64(),
				}.norm()
				if ray.dir.dot(normal) < 0 {
					ray.dir = ray.dir.multScalar(-1)
				}
				ray.origin = position.add(ray.dir.multScalar(rayEpsilon))
			}
		}
	}
	return photons
}

func renderpoint(camera *matrix.DenseMatrix, scene Scene, x, y float64) Vector3 {
	origin := Vector3{0, -1, 0}
	imgX := x / float64(WIDTH)
	imgY := y / float64(HEIGHT)
	dir := []float64{imgX, imgY, 1, 1}
	dirMatrix, _ := matrix.MakeDenseMatrix(dir, 1, 4).TimesDense(camera)
	rayDir := matrixToVector3(dirMatrix).norm()
	ray := &Ray{origin, rayDir}
	hit, surfaceColor, _, _ := raytrace(true, ray, scene, 0)
	if hit {
		return surfaceColor
	} else {
		return Vector3{}
	}
}

func gatherPhotons(point, normal Vector3, scene Scene) Vector3 {
	energy := Vector3{}
	return energy
}

func render(camera *matrix.DenseMatrix, scene Scene) {
	bounds := image.Rect(int(-WIDTH/2), int(-HEIGHT/2), int(WIDTH/2), int(HEIGHT/2))
	image := image.NewRGBA(bounds)

	wg := new(sync.WaitGroup)
	wg.Add(HEIGHT)
	for v := -HEIGHT / 2; v < HEIGHT/2; v++ {
		go func(v int) {
			for u := -WIDTH / 2; u < WIDTH/2; u++ {
				x := 2.0 * float64(u)
				y := 2.0 * float64(v)
				color1 := renderpoint(camera, scene, x, y)
				color2 := renderpoint(camera, scene, x+0.5, y)
				color3 := renderpoint(camera, scene, x, y+0.5)
				color4 := renderpoint(camera, scene, x+0.5, y+0.5)
				//Average 4 points for this pixel
				averageColor := color1.add(color2).add(color3).add(color4).multScalar(0.25)
				var red, green, blue, alpha uint8
				red = uint8(math.Min(averageColor[0], 1) * 255)
				green = uint8(math.Min(averageColor[1], 1) * 255)
				blue = uint8(math.Min(averageColor[2], 1) * 255)
				alpha = uint8(255)
				image.Set(u, v, color.RGBA{red, green, blue, alpha})
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	t := time.Now()
	filename := "output " + t.Format("2006-01-02 15:04:05") + ".png"
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	png.Encode(file, image)
}
