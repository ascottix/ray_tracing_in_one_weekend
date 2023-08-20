package main

type HittableList struct {
	objects []Hittable
}

func NewHittableList() HittableList {
	return HittableList{}
}

func (hl *HittableList) Add(object Hittable) {
	hl.objects = append(hl.objects, object)
}

func (hl *HittableList) Clear() {
	hl.objects = nil
}

func (hl HittableList) Hit(ray Ray, rayTmin, rayTmax float64, rec *HitRecord) bool {
	tempRec := HitRecord{}
	hitAnything := false
	closestSoFar := rayTmax
	for _, object := range hl.objects {
		if object.Hit(ray, rayTmin, closestSoFar, &tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			*rec = tempRec
		}
	}

	return hitAnything
}
