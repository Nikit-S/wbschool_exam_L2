package pattern

import "fmt"

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

type state interface {
	addCharge(charge uint) error
	pressButton()
	unlock()
}

type phone struct {
	charging state
	unlocked state
	locked   state

	currentState state

	battery uint
}

func newPhone(b uint) *phone {
	ph := &phone{battery: b}
	ch := &charging{phone: ph}
	unl := &unlocked{phone: ph}
	l := &locked{phone: ph}

	ph.charging = ch
	ph.unlocked = unl
	ph.locked = l
	if ph.battery < 30 {
		ph.currentState = ch
	} else {
		ph.currentState = l
	}
	return ph
}

func (ph *phone) addCharge(charge uint) error {
	return ph.currentState.addCharge(charge)
}

func (ph *phone) pressButton() {
	ph.currentState.pressButton()
}
func (ph *phone) unlock() {
	ph.currentState.unlock()
}

type charging struct {
	*phone
}

func (st *charging) addCharge(charge uint) error {
	if st.battery+charge > 100 || charge == 0 {
		return fmt.Errorf("Wrong charge")
	}
	st.battery += charge
	fmt.Println("Charged up to: ", st.battery)
	if st.battery > 30 {
		st.currentState = st.locked
		fmt.Println("Turned On")
	}
	return nil
}

func (st *charging) pressButton() {
	fmt.Println("Showing charging screen")
	if st.battery < 10 {
		st.battery = 0
	} else {
		st.battery -= 10
	}
}

func (st *charging) unlock() {
	fmt.Println("Nothing happened")
}

type locked struct {
	*phone
}

func (st *locked) addCharge(charge uint) error {
	if st.battery+charge > 100 || charge == 0 {
		return fmt.Errorf("Wrong charge")
	}
	st.battery += charge
	fmt.Println("Charged up to: ", st.battery)
	return nil
}

func (st *locked) pressButton() {
	fmt.Println("Showing lock screen")
	if st.battery < 10 {
		st.battery = 0
		fmt.Println("Discharged")
		st.currentState = st.charging
	} else {
		st.battery -= 10
	}
}

func (st *locked) unlock() {
	fmt.Println("Unlocked")
	st.currentState = st.unlocked
}

type unlocked struct {
	*phone
}

func (st *unlocked) addCharge(charge uint) error {
	if st.battery+charge > 100 || charge == 0 {
		return fmt.Errorf("Wrong charge")
	}
	st.battery += charge
	fmt.Println("Charged up to: ", st.battery)
	return nil
}

func (st *unlocked) pressButton() {
	fmt.Println("Showing screen and apps ")
	if st.battery < 10 {
		st.battery = 0
		fmt.Println("Discharged")
		st.currentState = st.charging
	} else {
		st.battery -= 10
	}
}

func (st *unlocked) unlock() {
	fmt.Println("Nothing happened")
}

func (st *unlocked) lock() {
	st.currentState = st.locked
}
