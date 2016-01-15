package abm

import (
	"fmt"
	"log"
	"time"

	"github.com/benjamin-rood/abm-colour-polymorphism/calc"
)

/* Attack(Hunted) bool
Eat(Hunted) bool */

// RBB : Rule-Based-Behaviour for Visual Predator Agent
func (vp *VisualPredator) RBB(ctxt Context, cppPop []ColourPolymorphicPrey) (returning []VisualPredator) {
	jump := ""
	jump = vp.Age(ctxt)
	fmt.Println(jump)
	time.Sleep(time.Second)
	switch jump {
	case "PREY SEARCH":
		target, err := vp.PreySearch(cppPop, ctxt.VsrSearchChance)
		if err != nil {
			log.Println("vp.RBB:", err)
		}
		success := vp.Attack(target, ctxt.VpAttackChance)
		if success {
			goto Add
		}
		fallthrough
	case "PATROL":
		𝚯 := calc.RandFloatIn(-ctxt.VpTurn, ctxt.VpTurn)
		vp.Turn(𝚯)
		vp.Move()
	case "DEATH":
		goto End
	default:
		log.Println("vp.RBB Switch: FAIL: jump =", jump)
	}
Add:
	returning = append(returning, *vp)
End:
	return
}