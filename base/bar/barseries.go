package bar

import (
	"errors"
	"fmt"
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
	bars           []Bar
	barField       BarField
	mapDatetimeBar map[int]Bar

	InnerChilds []xbase.ISeries
	InnerParent xbase.ISeries // always nil
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
		all = append(all, item.get(this.barField))
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
	return this.bars[index].get(this.barField)
}

func NewBarSeries() *BarSeries {
	s := &BarSeries{}
	s.Init(nil)
	return s
}

func (this *BarSeries) Init(parent xbase.ISeries) {
	this.DateTime = []time.Time{}
	this.bars = []Bar{}
	this.mapDatetimeBar = map[int]Bar{}
	this.InnerParent = parent
	this.barField = Close // default use close
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

func (this *BarSeries) AppendBar(bar_ Bar) {
	if debug {
		fmt.Println("BarSeries.Append:", bar_.get(this.barField))
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
		// fmt.Println("child iseries append", this.Symbol, bar_.get(this.barField), reflect.TypeOf(child), child.Count())
		child.Append(&datetime, bar_.get(this.barField))
	}
}
