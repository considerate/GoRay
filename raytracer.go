package main

import (
	"fmt"
	//"github.com/davecheney/profile"
	"github.com/skelterjohn/go.matrix"
	"github.com/unit3/kdtree"
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
const PHOTONS_PER_LIGHT = 1000
const PHOTONS_DEPTH = 3
const photonMap = true
const photonRadius = 5
const photonRadiusSq = photonRadius * photonRadius
const photonArea = math.Pi * photonRadiusSq
const vacuum = 1.0

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()
	scale := 1 / math.Tan(fov/2)
	focalLength := far - near
	camera := matrix.MakeDenseMatrix([]float64{
		scale * aspect, 0, 0, 0,
		0, -scale, 0, 0,
		0, 0, far / focalLength, 0,
		0, 0, 0, 1},
		4, 4)

	//Create scene (scene.go)
	fmt.Println("Creating scene")
	scene := setupScene()
	if photonMap {
		fmt.Println("Generating photons")
		scene.photonMap = shootPhotons(scene)
	}

	//Render image
	fmt.Println("Rendering image")
	render(camera, scene)
	fmt.Println("Done")
}

func shootPhotons(scene Scene) *kdtree.Tree {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	numLights := len(scene.lights)
	numPhotons := numLights * PHOTONS_PER_LIGHT * PHOTONS_DEPTH
	photons := make([]*kdtree.Node, numPhotons)
	count := 0
	photonPower := 600.0 / PHOTONS_PER_LIGHT
	var origin, dir Vector3
	var ray *Ray
	var material *Material
	var color, normal, probs Vector3
	var roulette, t float64
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
			color = Vector3{}
			for j = 0; j < PHOTONS_DEPTH; j++ {
				hit, t, object, _ = closestIntersection(ray, scene.objects)
				if !hit {
					break
				}
				position = origin.add(dir.multScalar(t))
				roulette = r.Float64()
				material = object.material()
				probs = Vector3{material.diffuse, material.specular, material.absorb}.norm()
				if roulette < probs[0] {
					//Diffuse reflection
					normal = object.normal(position)
					surfaceColor = shadePoint(position, normal, object, scene.lights)
					color = color.add(surfaceColor)
					photons[count] = kdtree.NewNode([]float64{position[0], position[1], position[2]})
					photon := &Photon{
						color.multScalar(photonPower),
						ray.dir.copy(),
					}
					photons[count].Data = new(interface{})
					*photons[count].Data = photon
					count++

					ray.dir = Vector3{
						1.0 - 2*r.Float64(),
						1.0 - 2*r.Float64(),
						1.0 - 2*r.Float64(),
					}.norm()
					if ray.dir.dot(normal) < 0 {
						//Invert rays reflected in wrong direction (opposite side of normal plane)
						ray.dir = ray.dir.multScalar(-1)
					}
					ray.origin = position.add(ray.dir.multScalar(rayEpsilon))
				} else if roulette >= probs[0] && roulette < probs[0]+probs[1] {
					//Specular reflection
					normal = object.normal(position)
					ray.dir = ray.dir.reflect(normal)
					ray.origin = position.add(ray.dir.multScalar(rayEpsilon))
				} else {
					//Absorb photon
					break
				}
			}
		}
	}
	treePhotons := photons[:count]
	tree := kdtree.BuildTree(treePhotons)
	return tree
}

func matrixToVector3(mat *matrix.DenseMatrix) Vector3 {
	arr := mat.Array()
	return Vector3{arr[0], arr[1], arr[2]}
}

func shadePoint(point, normal Vector3, object Object, lights []Light) Vector3 {
	material := object.material()

	resultingColor := Vector3{}
	if material.diffuse > 0 {
		for _, light := range lights {
			lightPos := light.position.copy()
			lightDir := lightPos.sub(point).norm()
			dot := normal.dot(lightDir)
			if dot > 0 {
				diffuse := material.diffuse * dot
				color := material.color.copy()
				resultingColor = resultingColor.add(color.multScalar(diffuse))
			}
		}
	}
	return resultingColor
}

func closestIntersection(ray *Ray, objects []Object) (bool, float64, Object, bool) {
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
	return hit, t, hitObject, inside
}

func raytrace(ray *Ray, scene Scene, ni float64, depth ...uint8) (bool, Vector3, Vector3, Object) {
	tracedepth := depth[0]
	maxdepth := MAX_TRACE_DEPTH
	objects := scene.objects
	lights := scene.lights
	origin := ray.origin.copy()
	dir := ray.dir.copy()
	if len(depth) > 1 {
		maxdepth = depth[1]
	}
	hit, t, hitObject, inside := closestIntersection(ray, objects)
	if hit {
		object := hitObject
		point := origin.add(dir.copy().multScalar(t))
		normal := object.normal(point).copy()
		material := object.material()
		nt := material.refrIdx
		if inside {
			normal = normal.multScalar(-1)
			ni = material.refrIdx
			nt = 1.0
			if nt > 0 {

			}
		}

		var color Vector3
		if photonMap {
			//Direct illumination
			indirectColor := Vector3{}

			//Indirect illumination
			gatheredPhotons, _ := scene.photonMap.FindRange(map[int]kdtree.Range{
				0: {point[0] - photonRadius, point[0] + photonRadius},
				1: {point[1] - photonRadius, point[1] + photonRadius},
				2: {point[2] - photonRadius, point[2] + photonRadius},
			})
			for _, photonNode := range gatheredPhotons {
				data := photonNode.Data
				photon := (*data).(*Photon)
				photonPos := Vector3{photonNode.Coordinates[0], photonNode.Coordinates[1], photonNode.Coordinates[2]}
				distv := photonPos.sub(point)
				distSq := distv.dot(distv)
				if distSq > photonRadiusSq {
					continue
				}
				weight := math.Max(0.0, -normal.dot(photon.incomingDirection))
				weight *= 1.0 - math.Sqrt(distSq)/photonRadius
				indirectColor = indirectColor.add(photon.color.copy().multScalar(weight))
			}
			indirectColor = indirectColor.multScalar(1.0 / photonArea * 3.0)
			color = shadePoint(point, normal, object, lights)
			color = color.multScalar(0.5).add(indirectColor)
		} else {
			color = shadePoint(point, normal, object, lights)
		}

		if material.specular > 0 && tracedepth < maxdepth-1 {

			reflection, transmission := 1.0, 1.0
			if material.fresnel {
				reflection, transmission = fresnelReflectionRefraction(ni, object, point, normal, ray)
				transmission = math.Min(transmission, 1.0)
				reflection = math.Min(reflection, 1.0)
			}
			if material.reflection > 0 && reflection > 0 {
				r := dir.reflect(normal)
				newRay := &Ray{point.add(r.multScalar(rayEpsilon)), r}
				nextHit, nextColor, _, _ := raytrace(newRay, scene, ni, tracedepth+1, maxdepth)
				if nextHit {
					color = color.add(nextColor.multScalar(reflection * material.reflection))
				}
			}
			if material.refraction > 0 && transmission > 0 {

				r := dir.refract(normal, ni, nt).norm()
				newRay := &Ray{point.add(r.multScalar(rayEpsilon)), r}
				nextHit, nextColor, _, _ := raytrace(newRay, scene, nt, tracedepth+1, maxdepth)
				if nextHit {
					color = color.add(nextColor.multScalar(transmission * material.refraction))
				}
			}
		}
		return true, color, point, object
	}
	return false, Vector3{}, Vector3{}, nil
}

func fresnelReflectionRefraction(currentIndex float64, object Object, point, normal Vector3, ray *Ray) (float64, float64) {
	cosθi := math.Abs(ray.dir.norm().dot(normal)) //Cosine of incoming angle (absolute)
	sinθi := math.Sqrt(1 - cosθi*cosθi)           //Trigonometric 1 (Sine of incoming angle)
	ni := currentIndex                            //Incoming refraction index
	nt := object.material().refrIdx               //Transmission refraction index
	sinθt := (ni / nt) * sinθi
	if cosθi*sinθt > 0.999 {
		//Grazing angle... result in pure reflection
		return 1.0, 0.0
	} else {
		cosθt := math.Sqrt(1.0 - sinθt*sinθt) //Trigonometric 1
		reflectionParallel := (nt*cosθi - ni*cosθt) / (ni*cosθt + nt*cosθi)
		reflectionPerpendicular := (ni*cosθi - nt*cosθt) / (ni*cosθi + nt*cosθt)

		//Average the two polarizations
		reflection := 0.5 * (reflectionParallel + reflectionPerpendicular)
		reflection *= reflection //Square

		transmissionParallel := (2 * ni * cosθi) / (ni*cosθt + nt*cosθi)
		transmissionPerpendicular := (2 * ni * cosθi) / (ni*cosθi + nt*cosθt)

		//Average the two polarizations
		transmission := 0.5 * (transmissionParallel + transmissionPerpendicular)
		transmission *= transmission
		transmission *= cosθt / cosθi

		if reflection < 0 {
			reflection = 0
		}
		if transmission < 0 {
			transmission = 0
		}
		return reflection, transmission
	}
}

func renderpoint(camera *matrix.DenseMatrix, scene Scene, x, y float64) Vector3 {
	origin := Vector3{0, -1, 0}
	imgX := x / float64(WIDTH)
	imgY := y / float64(HEIGHT)
	dir := []float64{imgX, imgY, 1, 1}
	dirMatrix, _ := matrix.MakeDenseMatrix(dir, 1, 4).TimesDense(camera)
	rayDir := matrixToVector3(dirMatrix).norm()
	ray := &Ray{origin, rayDir}
	hit, surfaceColor, _, _ := raytrace(ray, scene, vacuum, 0)
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
	filename := "output/output " + t.Format("2006-01-02 15:04:05") + ".png"
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	png.Encode(file, image)
}
