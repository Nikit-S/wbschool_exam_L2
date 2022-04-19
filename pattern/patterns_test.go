package pattern

import (
	"fmt"
	"testing"
)

func TestFacade(t *testing.T) {
	mwf := newMagicWorldFacade()
	mwf.MagicDustUse()
	mwf.SpellUse()
	mwf.PoisonExtractionUse()
	mwf.MushroomCollectionUse()
	mwf.GwitichDrinkUse()
	mwf.PrepareForWar()
}

func TestBuilder(t *testing.T) {
	defBuilder := getBuilder("")
	izbaBuilder := getBuilder("izba")
	prarab := newPrarab(defBuilder)
	defHouse := prarab.buildHouse()

	prarab.setBuilder(izbaBuilder)
	izbaHouse := prarab.buildHouse()
	fmt.Printf("Default House Door Type: %s\n", defHouse.doorMaterial)
	fmt.Printf("Default House Window Type: %s\n", defHouse.windowMaterial)
	fmt.Printf("Default House Num Floor: %d\n", defHouse.floors)

	fmt.Printf("Izba Door Type: %s\n", izbaHouse.doorMaterial)
	fmt.Printf("Izba Window Type: %s\n", izbaHouse.windowMaterial)
	fmt.Printf("Izba Num Floor: %d\n", izbaHouse.floors)
}

func TestVisitor(t *testing.T) {
	(&circle{}).accept(&areaSize{})
	(&square{}).accept(&areaSize{})
	(&rectangle{}).accept(&areaSize{})
}

func TestCommand(t *testing.T) {
	tv := &tv{}

	onCommand := &onCommand{
		device: tv,
	}

	offCommand := &offCommand{
		device: tv,
	}

	onButton := &button{
		command: onCommand,
	}
	onButton.press()

	offButton := &button{
		command: offCommand,
	}
	offButton.press()
}

func TestChain(t *testing.T) {
	m := &man{}
	bb := &bigBurokrat{}
	sb := &smallBurokrat{next: bb}
	r := &receptionist{next: sb}

	r.execute(m)
	m.askForMask(r)
	r.execute(m)
}

func TestFactory(t *testing.T) {
	s, _ := getBaby("Sanya")
	tam, _ := getBaby("Tamara")
	fmt.Println(s)
	fmt.Println(tam)
}

func TestStrategy(t *testing.T) {
	cli := &client{money: 10000.0}
	cli1 := &client{money: 5000.0}
	cli2 := &client{money: 10.0}
	rou := &route{distance: 50}

	rcar := GenericTransport[rentCar]{cost: 5000, speed: 100}
	tax := GenericTransport[taxi]{cost: 1000, speed: 100}
	bi := GenericTransport[bicycle]{cost: 0, speed: 20}

	fmt.Println(cli, rou)
	rcar.getToAirport(cli, rou)
	fmt.Println(cli, rou)

	fmt.Println(cli1, rou)
	tax.getToAirport(cli1, rou)
	fmt.Println(cli1, rou)

	fmt.Println(cli2, rou)
	bi.getToAirport(cli2, rou)
	fmt.Println(cli2, rou)
}

func TestState(t *testing.T) {
	ph := newPhone(10)

	ph.addCharge(10)
	ph.unlock()
	ph.pressButton()
	ph.pressButton()
	ph.pressButton()
	ph.addCharge(40)
	ph.pressButton()
	ph.unlock()
	ph.pressButton()
	ph.pressButton()
	ph.pressButton()
	ph.pressButton()
	ph.pressButton()

}
