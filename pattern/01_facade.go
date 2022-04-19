package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

type MagicDust struct{}

func (obj *MagicDust) MagicDustUse() {
	fmt.Println("*MagicDust Use*")
}

type Spell struct{}

func (obj *Spell) SpellUse() {
	fmt.Println("*Spell Use*")
}

type PoisonExtraction struct{}

func (obj *PoisonExtraction) PoisonExtractionUse() {
	fmt.Println("*PoisonExtraction Use*")
}

type MushroomCollection struct{}

func (obj *MushroomCollection) MushroomCollectionUse() {
	fmt.Println("*MushroomCollection Use*")
}

type GwitichDrink struct{}

func (obj *GwitichDrink) GwitichDrinkUse() {
	fmt.Println("*GwitichDrink Use*")
}

type MagicWorldFacade struct {
	magicDust          *MagicDust
	spell              *Spell
	poisonExtraction   *PoisonExtraction
	mushroomCollection *MushroomCollection
	gwitichDrink       *GwitichDrink
}

func newMagicWorldFacade() *MagicWorldFacade {
	return &MagicWorldFacade{
		magicDust:          &MagicDust{},
		spell:              &Spell{},
		poisonExtraction:   &PoisonExtraction{},
		mushroomCollection: &MushroomCollection{},
		gwitichDrink:       &GwitichDrink{},
	}
}

func (obj *MagicWorldFacade) MagicDustUse() {
	obj.magicDust.MagicDustUse()
}

func (obj *MagicWorldFacade) SpellUse() {
	obj.spell.SpellUse()
}

func (obj *MagicWorldFacade) PoisonExtractionUse() {
	obj.poisonExtraction.PoisonExtractionUse()
}

func (obj *MagicWorldFacade) MushroomCollectionUse() {
	obj.mushroomCollection.MushroomCollectionUse()
}

func (obj *MagicWorldFacade) GwitichDrinkUse() {
	obj.gwitichDrink.GwitichDrinkUse()
}

func (obj *MagicWorldFacade) PrepareForWar() {
	fmt.Println("PREP FOR WAR")
	obj.MagicDustUse()
	obj.SpellUse()
	obj.PoisonExtractionUse()
	obj.MushroomCollectionUse()
	obj.GwitichDrinkUse()

}
