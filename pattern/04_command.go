package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

type button struct {
	command
}

func (b *button) press() {
	b.execute()
}

type command interface {
	execute()
}

type onCommand struct {
	device
}

func (c *onCommand) execute() {
	c.on()
}

type offCommand struct {
	device
}

func (c *offCommand) execute() {
	c.off()
}

type device interface {
	on()
	off()
}

type tv struct {
	isRunning bool
}

func (t *tv) on() {
	t.isRunning = true
	fmt.Println("Turning tv on")
}

func (t *tv) off() {
	t.isRunning = false
	fmt.Println("Turning tv off")
}
