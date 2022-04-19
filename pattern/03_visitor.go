package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

type shape interface {
	getType() string
	accept(visitor)
	*square | *circle | *rectangle
}

type square struct {
	side int
}

type circle struct {
	radius int
}

type rectangle struct {
	a int
	b int
}

func (square) getType() string {
	return "square"
}

func (circle) getType() string {
	return "circle"
}

func (rectangle) getType() string {
	return "rectangle"
}

func (obj *square) accept(v visitor) {
	v.visitForSquare(obj)
}
func (obj *circle) accept(v visitor) {
	v.visitForCircle(obj)
}
func (obj *rectangle) accept(v visitor) {
	v.visitForRectangle(obj)
}

type visitor interface {
	visitForSquare(*square)
	visitForCircle(*circle)
	visitForRectangle(*rectangle)
}

type areaSize struct{}

func (a *areaSize) visitForSquare(s *square) {
	fmt.Println("area of the square is")
}
func (a *areaSize) visitForCircle(c *circle) {
	fmt.Println("area of the circle is")
}
func (a *areaSize) visitForRectangle(s *rectangle) {
	fmt.Println("area of the rectangle is")
}
