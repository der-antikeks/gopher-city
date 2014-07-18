package main

import (
	"fmt"
	"log"
	"math/rand"
)

const PrintDebug bool = false

func Debug(v ...interface{}) {
	if PrintDebug {
		log.Println(v...)
	}
}

/*
	economic constants (duh!)

	1 Gopher > -0,4 Goods
	1 Good > -1,5 Gopher, -0,5 Products
	1 Product > -2 Gopher

	5 Low Residential > 20 Gopher, -8 Goods
	3 Low Commercial > -12 Gopher, +8 Goods, -4 Products
	2 Low Industrial > -8  Gopher,           +4 Products
*/
const (
	gopherNeedsGoods       = 0.4
	workerProducesGoods    = 1.0 / 1.5
	goodNeedsProducts      = 0.5
	workerProducesProducts = 0.5
)

func main() {
	rand.Seed(42)
	log.SetFlags(0)

	// A group of wild gophers appears!
	gophers := Gophers{
		NewGopher("Klas"), NewGopher("Sture"), NewGopher("Verner"), NewGopher("Asbjörn"),
		NewGopher("Loke"), NewGopher("Vidar"), NewGopher("Markus"), NewGopher("Staffan"),
		NewGopher("Knut"), NewGopher("Stian"), NewGopher("Magnus"), NewGopher("Theodor"),
		NewGopher("Acke"), NewGopher("Stian"), NewGopher("Gunnar"), NewGopher("Halsten"),
		NewGopher("Noak"), NewGopher("Alvar"), NewGopher("Viktor"), NewGopher("Sigvard"),
	}

	// Setup environment
	s := SpatialSystem()
	s.AddResidentials(
		NewResidential(4, gophers[0:4]),
		NewResidential(4, gophers[4:8]),
		NewResidential(4, gophers[8:12]),
		NewResidential(4, gophers[12:16]),
		NewResidential(4, gophers[16:20]),
	)
	s.AddCommercials(
		NewCommercial(4, nil),
		NewCommercial(4, nil),
		NewCommercial(4, nil),
	)
	s.AddIndustrials(
		NewIndustrial(4, nil),
		NewIndustrial(4, nil),
	)

	// Simulate
	for step := 0; step < 10; step++ {
		gophers.Shuffle()

		gophers.Shop()
		gophers.Work()
		gophers.Sleep()
	}

	// Results
	fmt.Println(gophers.String())
	fmt.Println(SpatialSystem().String())
}

type spatialSystem struct {
	residentials []*Residential
	commercials  []*Commercial
	industrials  []*Industrial
}

var spatialSystemSingleton *spatialSystem

func SpatialSystem() *spatialSystem {
	if spatialSystemSingleton == nil {
		spatialSystemSingleton = &spatialSystem{
			residentials: []*Residential{},
			commercials:  []*Commercial{},
			industrials:  []*Industrial{},
		}
	}
	return spatialSystemSingleton
}

func (s *spatialSystem) String() string {
	r := "Residentials {\n"
	for _, b := range s.residentials {
		r += "\t" + b.String() + "\n"
	}

	r += "}\nCommercials {\n"
	for _, b := range s.commercials {
		r += "\t" + b.String() + "\n"
	}

	r += "}\nIndustrials {\n"
	for _, b := range s.industrials {
		r += "\t" + b.String() + "\n"
	}
	r += "}"

	return r
}

func (s *spatialSystem) AddResidentials(rs ...*Residential) {
	s.residentials = append(s.residentials, rs...)
}
func (s *spatialSystem) AddCommercials(cs ...*Commercial) {
	s.commercials = append(s.commercials, cs...)
}
func (s *spatialSystem) AddIndustrials(is ...*Industrial) {
	s.industrials = append(s.industrials, is...)
}

func (s *spatialSystem) Residentials() []*Residential {
	return s.residentials
}
func (s *spatialSystem) Commercials() []*Commercial {
	return s.commercials
}
func (s *spatialSystem) Industrials() []*Industrial {
	return s.industrials
}

type Gopher struct {
	name      string
	worked    bool
	shopped   bool
	happiness float64 // happiness is a float64!

	job  Building
	home *Residential
}

func NewGopher(name string) *Gopher {
	return &Gopher{
		name:      name,
		happiness: 0.5,
	}
}

func (g *Gopher) HasWorked() bool {
	return g.worked
}

func (g *Gopher) WorkDone() {
	g.worked = true
}

func (g *Gopher) HasShopped() bool {
	return g.shopped
}

func (g *Gopher) ShopDone() {
	g.shopped = true
}

// adorable little gopher takes a nap :D
func (g *Gopher) Sleep() {
	if g.worked {
		g.happiness += 0.5
		g.worked = false
	} else {
		g.happiness -= 0.5
	}

	if g.shopped {
		g.happiness += 0.5
		g.shopped = false
	} else {
		g.happiness -= 0.5
	}

	// awwww! bonus
	g.happiness += 0.05

	// limit
	g.happiness = clamp(g.happiness, 0, 1)
}

func clamp(v, min, max float64) float64 {
	if v > max {
		return max
	} else if v < min {
		return min
	}
	return v
}

type Gophers []*Gopher

func (gs *Gophers) String() string {
	var r string
	for _, g := range *gs {
		r += fmt.Sprintf("{%s %0.2f", g.name, g.happiness)
		if g.job != nil {
			r += ", works in the "
			switch g.job.(type) {
			case *Commercial:
				r += "industry"
			case *Industrial:
				r += "commercial"
			}
		} else {
			r += ", is unemployed"
		}

		r += "}\n"
	}
	return r
}

func (gs Gophers) Shuffle() {
	for i := 0; i < len(gs); i++ {
		j := rand.Intn(i + 1)
		gs[i], gs[j] = gs[j], gs[i]
	}
}

func (gs Gophers) Sleep() {
	for _, g := range gs {
		g.Sleep()
	}
}

func (gs Gophers) Shop() {
	commercials := SpatialSystem().Commercials()
	for _, g := range gs {
		Debug("{G", g.name, "} goes shopping")
		for _, c := range commercials {
			if c.GetGoods(gopherNeedsGoods) {
				g.ShopDone()
				break
			}
		}
	}
}

func (gs Gophers) Work() {
	commercials := SpatialSystem().Commercials()
	industrials := SpatialSystem().Industrials()

	for _, g := range gs {
		if g.HasWorked() {
			continue
		}

		Debug("{G", g.name, "} goes working")
		for _, c := range commercials {
			if c.DoWork(g) {
				break
			}
		}
		if !g.HasWorked() {
			for _, i := range industrials {
				if i.DoWork(g) {
					break
				}
			}
		}
	}
}

type Residential struct {
	capacity  int
	residents []*Gopher
}

func NewResidential(size int, residents []*Gopher) *Residential {
	max := size
	if max > len(residents) {
		max = len(residents)
	}
	r := &Residential{
		capacity:  size,
		residents: residents[:max],
	}
	return r
}

func (r *Residential) String() string {
	return fmt.Sprintf("{R %d}", len(r.residents))
}

func (r *Residential) GetWorker() *Gopher {
	// first, look for gophers without job
	for _, g := range r.residents {
		if g.job == nil && !g.HasWorked() {
			Debug(r, "unemployed gopher found")
			return g
		}
	}

	// try employee poaching
	for _, g := range r.residents {
		if g.job != nil && !g.HasWorked() {
			Debug(r, "unoccupied gopher found")
			return g
		}
	}

	Debug(r, "no unemployed gophers")
	return nil
}

type Building interface {
	DoWork(*Gopher) bool
	RemoveWorker(*Gopher)
}

type Commercial struct {
	capacity int
	workers  []*Gopher

	products float64
	goods    float64
}

func NewCommercial(size int, workers []*Gopher) *Commercial {
	max := size
	if max > len(workers) {
		max = len(workers)
	}
	r := &Commercial{
		capacity: size,
		workers:  workers[:max],
	}
	return r
}

func (c *Commercial) String() string {
	return fmt.Sprintf("{C %d (%.4g/%.4g)}", len(c.workers), c.products, c.goods)
}

func (c *Commercial) RemoveWorker(worker *Gopher) {
	for i, g := range c.workers {
		if g == worker {
			copy(c.workers[i:], c.workers[i+1:])
			c.workers[len(c.workers)-1] = nil
			c.workers = c.workers[:len(c.workers)-1]

			return
		}
	}
}

func (c *Commercial) DoWork(worker *Gopher) bool {
	if len(c.workers) >= c.capacity {
		Debug(c, "no more capacity")
		return false
	}

	// hire worker
	c.workers = append(c.workers, worker)
	if worker.job != nil {
		worker.job.RemoveWorker(worker) // TODO: make less ugly!
	}
	worker.job = c

	// TODO: should be after production
	// but than the industry can poach the worker in progress...
	worker.WorkDone()

	// produce
	neededProducts := workerProducesGoods * goodNeedsProducts
	if c.products < neededProducts {
		Debug(c, "get products from industrials")
		industrials := SpatialSystem().Industrials()
		var gotProducts bool
		for _, i := range industrials {
			if gotProducts = i.GetProducts(neededProducts); gotProducts {
				c.products += neededProducts
				break
			}
		}
		if !gotProducts {
			Debug(c, "no products available")
			return false
		}
	} else {
		Debug(c, "products in stock")
	}

	Debug(c, "produce goods")
	c.products -= neededProducts
	c.goods += workerProducesGoods
	return true
}

func (c *Commercial) GetGoods(amount float64) bool {
	// goods in stock
	if c.goods >= amount {
		Debug(c, "goods in stock")
		c.goods -= amount
		return true
	} else {
		Debug(c, "not enough goods in stock")
	}

	for c.goods < amount {
		// fetch a hired gopher
		Debug(c, "fetch a worker")
		var worker *Gopher
		for _, g := range c.workers {
			if !g.HasWorked() {
				worker = g
				break
			}
		}

		// all have worked, but we still have vacancies
		if worker == nil {
			Debug(c, "all workers are busy")
			if len(c.workers) < c.capacity {
				Debug(c, "hire new gopher from residentials")
				residentials := SpatialSystem().Residentials()
				for _, r := range residentials {
					worker = r.GetWorker()
					if worker != nil {
						break
					}
				}
				// no workers, no goods!
				if worker == nil {
					Debug(c, "no worker available")
					return false
				}
				// hire worker
				c.workers = append(c.workers, worker)
				if worker.job != nil {
					worker.job.RemoveWorker(worker) // TODO: make less ugly!
				}
				worker.job = c
			} else {
				Debug(c, "no more capacity")
				return false
			}
		}

		// TODO: should be after production
		worker.WorkDone()

		neededProducts := workerProducesGoods * goodNeedsProducts
		if c.products < neededProducts {
			Debug(c, "get products from industrials")
			industrials := SpatialSystem().Industrials()
			var gotProducts bool
			for _, i := range industrials {
				if gotProducts = i.GetProducts(neededProducts); gotProducts {
					c.products += neededProducts
					break
				}
			}
			if !gotProducts {
				Debug(c, "no products available")
				return false
			}
		} else {
			Debug(c, "products in stock")
		}

		Debug(c, "produce goods")
		c.products -= neededProducts
		c.goods += workerProducesGoods
	}

	Debug(c, "hand over", amount, "goods")
	c.goods -= amount
	return true
}

type Industrial struct {
	capacity int
	workers  []*Gopher

	products float64
}

func NewIndustrial(size int, workers []*Gopher) *Industrial {
	max := size
	if max > len(workers) {
		max = len(workers)
	}
	r := &Industrial{
		capacity: size,
		workers:  workers[:max],
	}
	return r
}

func (i *Industrial) String() string {
	return fmt.Sprintf("{I %d (%.4g)}", len(i.workers), i.products)
}

func (i *Industrial) RemoveWorker(worker *Gopher) {
	for j, g := range i.workers {
		if g == worker {
			copy(i.workers[j:], i.workers[j+1:])
			i.workers[len(i.workers)-1] = nil
			i.workers = i.workers[:len(i.workers)-1]

			return
		}
	}
}

func (i *Industrial) DoWork(worker *Gopher) bool {
	if len(i.workers) >= i.capacity {
		Debug(i, "no more capacity")
		return false
	}

	// hire worker
	i.workers = append(i.workers, worker)
	if worker.job != nil {
		worker.job.RemoveWorker(worker) // TODO: make less ugly!
	}
	worker.job = i

	// produce
	Debug(i, "produce products")
	i.products += workerProducesProducts
	worker.WorkDone()
	return true
}

func (i *Industrial) GetProducts(amount float64) bool {
	// products in stock
	if i.products >= amount {
		Debug(i, "products in stock")
		i.products -= amount
		return true
	} else {
		Debug(i, "not enough products in stock")
	}

	for i.products < amount {
		// fetch a hired gopher
		Debug(i, "fetch a worker")
		var worker *Gopher
		for _, g := range i.workers {
			if !g.HasWorked() {
				worker = g
				break
			}
		}

		// all have worked, but we still have vacancies
		if worker == nil {
			Debug(i, "all workers are busy")
			if len(i.workers) < i.capacity {
				Debug(i, "hire new gopher from residentials")
				residentials := SpatialSystem().Residentials()
				for _, r := range residentials {
					worker = r.GetWorker()
					if worker != nil {
						break
					}
				}
				// no workers, no products!
				if worker == nil {
					Debug(i, "no worker available")
					return false
				}
				// hire worker
				i.workers = append(i.workers, worker)
				if worker.job != nil {
					worker.job.RemoveWorker(worker) // TODO: make less ugly!
				}
				worker.job = i
			} else {
				Debug(i, "no more capacity")
				return false
			}
		}

		Debug(i, "produce products")
		i.products += workerProducesProducts
		worker.WorkDone()
	}

	Debug(i, "hand over", amount, "products")
	i.products -= amount
	return true
}
