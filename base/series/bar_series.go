package series

import (
	"errors"
	"fmt"
	"quant/base/bar"
	"quant/base/xbase"
	_ "reflect"
	"time"
)

var (
	ErrInvalidUseOfBarSeries = errors.New("invalid use of BarSeries")
)

type BarSeries struct {
	Symbol         string
	StartTime      time.Time
	EndTime        time.Time
	DateTime       []time.Time
	bars           []bar.Bar
	barField       bar.BarField
	mapDatetimeBar map[int]bar.Bar

	InnerChilds []xbase.ISeries `json:"-" `
	InnerParent xbase.ISeries   `json:"-" ` // always nil

	// use between OnMeasure() && OnDraw()
	drawStartTime time.Time
	drawEndTime   time.Time
	drawStartIdx  int
	drawEndIdx    int
	//
	IndicatorType int
}

func (this *BarSeries) Keys() []time.Time {
	all := []time.Time{}
	for _, item := range this.bars {
		all = append(all, item.DateTime)
	}
	return all
}

func (this *BarSeries) Values() []float64 {
	all := []float64{}
	for _, item := range this.bars {
		all = append(all, item.Get(this.barField))
	}
	return all

}

func (this *BarSeries) Count() int {
	return len(this.bars)
}

func (this *BarSeries) Index(datetime *time.Time) int {
	for idx, item := range this.bars {
		if item.DateTime.Equal(*datetime) {
			return idx
		}
	}
	return -1
}

func (this *BarSeries) Now() time.Time {
	return this.EndTime
}

func (this *BarSeries) ValueAtTime(datetime *time.Time) float64 {

	idx := this.Index(datetime)
	if idx >= 0 {
		return this.ValueAtIndex(idx)
	} else {
		fmt.Println("ValueAtTime invalid datetime: %v %v", datetime, idx)
		panic(datetime)
		return -1
	}
}

func (this *BarSeries) ValueAtIndex(index int) float64 {
	if index >= len(this.bars) || index < 0 {
		fmt.Println("OutOfArray: %v %v", len(this.bars), index)
		panic(index)
	}
	return this.bars[index].Get(this.barField)
}

func NewBarSeries() *BarSeries {
	s := &BarSeries{}
	s.Init(nil)
	return s
}

func (this *BarSeries) GetIndicatorType() int {
	return this.IndicatorType
}

func (this *BarSeries) Init(parent xbase.ISeries) {
	this.DateTime = []time.Time{}
	this.bars = []bar.Bar{}
	this.mapDatetimeBar = map[int]bar.Bar{}

	this.InnerParent = parent
	this.barField = bar.Close // default use close
	this.IndicatorType = xbase.IndicatorTypeDayBar
}

func (this *BarSeries) Match(symbol string) bool {
	if this.Symbol == symbol {
		return true
	} else {
		return false
	}
}

func (this *BarSeries) AddChild(child xbase.ISeries) {
	this.InnerChilds = append(this.InnerChilds, child)
}

func (this *BarSeries) Append(datetime *time.Time, value float64) {
	if debug {
		fmt.Println("BarSeries.Append dummy function")
	}
	panic(ErrInvalidUseOfBarSeries)
}

func (this *BarSeries) AppendBar(bar_ bar.Bar) {
	if debug {
		fmt.Println("BarSeries.Append:", bar_.Get(this.barField))
	}

	datetime := bar_.DateTime

	if len(this.bars) == 0 {
		this.StartTime = datetime
		this.EndTime = datetime
	} else {
		this.EndTime = datetime
	}

	sec := int(datetime.Unix())
	oldBar, ok := this.mapDatetimeBar[sec]
	if ok {
		// can not append duplicate record
		fmt.Println("can not append duplicate record: %v %v", oldBar, ok)
		panic(datetime)
	}

	this.mapDatetimeBar[sec] = bar_
	this.DateTime = append(this.DateTime, datetime)
	this.bars = append(this.bars, bar_)

	for _, child := range this.InnerChilds {
		// fmt.Println("child iseries append", this.Symbol, bar_.Get(this.barField), reflect.TypeOf(child), child.Count())
		child.Append(&datetime, bar_.Get(this.barField))
	}
}

// indicator.IIndicator.OnMeasure

// num of table, X cordirate [datetime count], Y cordirate [min, max, num, percent]
func (this *BarSeries) OnMeasure(start, end time.Time) (int, float64, float64, int, float64) {
	this.drawStartTime = start
	this.drawEndTime = end
	this.drawStartIdx = 0
	this.drawEndIdx = 0

	// put it here now.
	var min, max float64
	var num int

	min = 10000
	max = -10000
	for idx, bar_ := range this.bars {
		if bar_.DateTime.Before(this.drawStartTime) {
			continue
		} else if bar_.DateTime.After(this.drawEndTime) {
			continue
		}
		if num == 0 {
			this.drawStartIdx = idx
		}
		this.drawEndIdx = idx + 1
		num++

		value := bar_.Get(this.barField)
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	if debug {
		fmt.Println("BarSeries OnMeasure", num, min, max, 100, 100)
	}
	return num, min, max, 100, 100
}

func (this *BarSeries) OnLayout() []time.Time {

	/*
	 * Note:
	 *  For a = [0,1,2,3]:
	 *    a[0:1] means [0]
	 *    a[0:4] means [0,1,2,3]
	 */
	return this.DateTime[this.drawStartIdx:this.drawEndIdx]
}

// indicator.IIndicator.OnDraw(ICanvas)
func (this *BarSeries) OnDraw(canvas xbase.ICanvas) {
	if debug {
		fmt.Println("symbol", this.Symbol, "Bars onDraw")
	}

	canvas.DrawBar(this.bars[this.drawStartIdx:this.drawEndIdx], 1)

	// canvas.DrawBuy(table,  []time.Time, []float64,color)
	// canvas.DrawSell(table, []time.Time, []float64,color)
	// canvas.DrawSpark(table,[]time.Time, []float64,color)
	// canvas.DrawShit(table, []time.Time, []float64,color)
}
