package account

import (
	"fmt"
	"github.com/milliyang/dice"
	"quant/base/bar"
	"quant/base/order"
	"quant/base/xbase"
)

func (this *Account) DiceHandleOrdersWithBar(orders []*order.Order, bar_ *bar.Bar) {
	pnl := 0
	for _, oneOrder := range orders {
		pnl += this.diceHandleOrderWithDiceRoll(oneOrder, bar_.Dice)
	}

	if len(orders) > 0 {
		fmt.Println(bar_.DateTime, "PNL:", pnl, bar_.Dice.Total)

		this.PnL += float64(pnl)
		this.PnlIndicator.UpdatePnl(&bar_.DateTime, float64(pnl))
	}
}

func (this *Account) diceHandleOrderWithDiceRoll(order_ *order.Order, diceRoll *dice.DiceRoll) int {

	pnl := 0

	switch order_.DiceBetType {
	case order.DiceBetTypeBig:
		if diceRoll.IsBig() {
			pnl += order_.DiceBetAmt
		} else {
			pnl -= order_.DiceBetAmt
		}
		break
	case order.DiceBetTypeSmall:
		if diceRoll.IsSmall() {
			pnl += order_.DiceBetAmt
		} else {
			pnl -= order_.DiceBetAmt
		}
		break
	case order.DiceBetTypeSingle:
		if diceRoll.IsSingle() {
			pnl += order_.DiceBetAmt
		} else {
			pnl -= order_.DiceBetAmt
		}
		break
	case order.DiceBetTypeEven:
		if diceRoll.IsEven() {
			pnl += order_.DiceBetAmt
		} else {
			pnl -= order_.DiceBetAmt
		}
		break
	case order.DiceBetTypeNumber:
		if !diceRoll.IsTriple() && diceRoll.Total == order_.DiceBetNumber {
			pnl += order_.DiceBetAmt * xbase.DiceNumberOddsMap[diceRoll.Total]
		} else {
			pnl -= order_.DiceBetAmt
		}
		break
	case order.DiceBetTypeTriple:
		if diceRoll.IsTriple() {
			pnl += order_.DiceBetAmt * xbase.DiceTrippleAnyOdds
		} else {
			pnl -= order_.DiceBetAmt
		}
		break
	case order.DiceBetTypeTripleNumber:
		if diceRoll.IsTriple() && diceRoll.Rolls[0] == order_.DiceBetNumber {
			pnl += order_.DiceBetAmt * xbase.DiceTrippleNumberOdds
		} else {
			pnl -= order_.DiceBetAmt
		}
		break
	case order.DiceBetTypeFaceNumber:
		match := diceRoll.ContainsFaceNumber(order_.DiceBetNumber)
		if match > 0 {
			pnl += order_.DiceBetAmt * match
		} else {
			pnl -= order_.DiceBetAmt
		}
		break
	default:
		panic(order_)
	}

	return pnl

}
