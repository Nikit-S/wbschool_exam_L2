package pattern

import "fmt"

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

type iBaby interface {
	setName(name string)
	setEyeColor(red, green, blue int)
	getName() string
	getEyeColor() (int, int, int)
	String() string
}

type baby struct {
	name     string
	eyecolor struct {
		red   int
		green int
		blue  int
	}
}

func (b *baby) String() string {
	return fmt.Sprintf("Name: %s, Eye Color: (r: %d g: %d b: %d)", b.name, b.eyecolor.red, b.eyecolor.green, b.eyecolor.blue)
}

func (b *baby) setName(name string) {
	b.name = name
}

func (b *baby) setEyeColor(red, green, blue int) {
	b.eyecolor.red = red
	b.eyecolor.green = green
	b.eyecolor.blue = blue

}

func (b *baby) getName() string {
	return b.name
}

func (b *baby) getEyeColor() (int, int, int) {
	return b.eyecolor.red, b.eyecolor.green, b.eyecolor.blue
}

type sanya struct {
	baby
}

func newSanya() iBaby {
	return &sanya{baby: baby{
		name: "Sanya",
		eyecolor: struct {
			red   int
			green int
			blue  int
		}{red: 100,
			green: 10,
			blue:  10,
		},
	}}
}

type tamara struct {
	baby
}

func newTamara() iBaby {
	return &sanya{baby: baby{
		name: "Tamara",
		eyecolor: struct {
			red   int
			green int
			blue  int
		}{red: 10,
			green: 200,
			blue:  10,
		},
	}}
}

func getBaby(babyType string) (iBaby, error) {
	if babyType == "Sanya" {
		return newSanya(), nil
	}
	if babyType == "Tamara" {
		return newTamara(), nil
	}
	return nil, fmt.Errorf("No Such Baby Type")
}
