package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/
type burokrat interface {
	execute(*man)
}

type receptionist struct {
	next burokrat
}

func (r *receptionist) execute(m *man) {
	if m.wearingMask {
		fmt.Println("Thank you for wearing mask, you may proceed")
		r.next.execute(m)
		return
	}
	fmt.Println("YOU ARE NOT WEARING MASK GET OUT!")
	m.wearingMask = true
}

type smallBurokrat struct {
	next burokrat
}

func (b *smallBurokrat) execute(m *man) {
	if m.smallBurokratCheckUp {
		fmt.Println("Oh, you already have FZ-721, good boy, proceed")
		b.next.execute(m)
		return
	}
	fmt.Println("Here, sign FZ-721 and proceed")
	m.smallBurokratCheckUp = true
	b.next.execute(m)
}

type bigBurokrat struct {
	next burokrat
}

func (b *bigBurokrat) execute(m *man) {
	if m.bigBurokratCheckUp {
		fmt.Println("Oh, you already have FU-4101, good boy, proceed")
		b.next.execute(m)
		return
	}
	fmt.Println("Here, sign FU-4101 and now you can get out")
	m.bigBurokratCheckUp = true
}

type man struct {
	wearingMask          bool
	smallBurokratCheckUp bool
	bigBurokratCheckUp   bool
}

func (m *man) askForMask(r *receptionist) {
	fmt.Println("please, give me the mask")
	m.wearingMask = true
}
