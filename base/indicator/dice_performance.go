package indicator

import (
	"fmt"
	"quant/base/color"
	"quant/base/xbase"
	"time"
)

// Casino Dicing Game
type DicePerformance struct {
	Total float64
	BaseIndicator
}

func NewDicePerformance(amt float64) *DicePerformance {
	s := &DicePerformance{}
	s.Init(0)
	s.Total = amt

	// override
	s.IndicatorType = xbase.IndicatorTypePerformance
	return s
}

// pnl may be negative
func (this *DicePerformance) UpdateData(datetime *time.Time, pnl float64) {
	if debug {
		fmt.Println("DicePerformance.UpdateData:", pnl)
	}
	this.Total += pnl

	this.BaseIndicator.UpdateData(datetime, this.Total)
}

// indicator.IIndicator.OnDraw(ICanvas)
func (this *DicePerformance) OnDraw(canvas xbase.ICanvas) {
	if len(this.DateTime) == 0 {
		return
	}

	if debug {
		fmt.Println("DicePerformance", this.Name, "symbol", this.Symbol, "onDraw")
		fmt.Println("DicePerformance onDraw", " start:", this.drawStartIdx, " end:", this.drawEndIdx)
	}
	canvas.DrawLine(this.DateTime[this.drawStartIdx:this.drawEndIdx], this.Data[this.drawStartIdx:this.drawEndIdx], this.Color)

	if (len(this.DateTime)) >= 2 {
		startTime := this.DateTime[0]
		endTime := this.DateTime[len(this.DateTime)-1]
		dateTimeStartAndEnd := []time.Time{startTime, endTime}
		value := []float64{5000.0, 5000.0}
		canvas.DrawLine(dateTimeStartAndEnd, value, color.Gold)
	}

	// for i := 0; i < len(this.Data); i++ {
	// 	value := this.Data[i]
	// 	date_ := this.DateTime[i]
	// 	if value != 5000 {
	// 		fmt.Printf("index:%d, d:%s, v:%d\n", i, date_.Format("2006-01-02"), int(value))
	// 	}
	// }
}
