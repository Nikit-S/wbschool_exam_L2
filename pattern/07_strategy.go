package pattern

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

type transStrategy interface {
	getToAirport(*client, *route)
}

type client struct {
	money float64
}

type route struct {
	distance int
	time     float64
}

type GenericTransport[T rentCar | bicycle | taxi] struct {
	cost  float64
	speed int
}

type rentCar struct{}
type bicycle struct{}
type taxi struct{}

func (gt *GenericTransport[T]) getToAirport(c *client, r *route) {
	c.money -= gt.cost * ((float64)(r.distance) / (float64)(gt.speed))
	r.time = ((float64)(r.distance) / (float64)(gt.speed))
}
