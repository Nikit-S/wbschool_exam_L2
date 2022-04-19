package pattern

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/
type house struct {
	floors         int
	windowMaterial string
	doorMaterial   string
}

type iBuilder interface {
	setDoors()
	setFloor()
	setWindows()
	getHouse() house
}

type defaultBuilder struct {
	doorMaterial   string
	floors         int
	windowMaterial string
}

func newDefaultBuilder() *defaultBuilder {
	return &defaultBuilder{}
}

func (b *defaultBuilder) setDoors() {
	b.doorMaterial = "wood"
}

func (b *defaultBuilder) setFloor() {
	b.floors = 2
}

func (b *defaultBuilder) setWindows() {
	b.windowMaterial = "Steklopaket"
}

type izbaBuilder struct {
	doorMaterial   string
	floors         int
	windowMaterial string
}

func (b *defaultBuilder) getHouse() house {
	return house{
		doorMaterial:   b.doorMaterial,
		floors:         b.floors,
		windowMaterial: b.windowMaterial,
	}
}

func newIzbaBuilder() *izbaBuilder {
	return &izbaBuilder{}
}

func (b *izbaBuilder) setDoors() {
	b.doorMaterial = "wood"
}

func (b *izbaBuilder) setFloor() {
	b.floors = 1
}

func (b *izbaBuilder) setWindows() {
	b.windowMaterial = "Sluda"
}

func (b *izbaBuilder) getHouse() house {
	return house{
		doorMaterial:   b.doorMaterial,
		floors:         b.floors,
		windowMaterial: b.windowMaterial,
	}
}

func getBuilder(btype string) iBuilder {
	switch btype {
	default:
		return &defaultBuilder{}
	case "izba":
		return &izbaBuilder{}
	}
}

type prarab struct {
	builder iBuilder
}

func newPrarab(b iBuilder) *prarab {
	return &prarab{
		builder: b,
	}
}

func (p *prarab) setBuilder(b iBuilder) {
	p.builder = b
}

func (p *prarab) buildHouse() house {
	p.builder.setDoors()
	p.builder.setWindows()
	p.builder.setFloor()
	return p.builder.getHouse()
}
