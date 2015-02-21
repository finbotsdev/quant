package strategy

import (
	"fmt"
	"quant/base"
	"quant/base/bar"
	"quant/base/order"
	"quant/base/series"
	"quant/base/xbase"
	"quant/canvas"
)

const (
	debug = true
)

func init() {
	if debug {
		fmt.Println("quant/base/stratety init")
	}
}

type Strategy struct {
	Name         string
	Symbol       string
	BarSeries    *series.BarSeries          `json:"-" `
	mapIndicator map[int][]xbase.IIndecator `json:"-" `
	drawCanvas   []*canvas.Canvas
	drawed       bool
}

func (this *Strategy) Init(symbol string, barSeries *series.BarSeries) {
	if debug {
		fmt.Println("Strategy.Init()")
	}
	this.Name = "Strategy"
	this.Symbol = symbol
	this.BarSeries = barSeries

	this.mapIndicator = map[int][]xbase.IIndecator{}
	this.drawed = false
}

func (this *Strategy) Key() string {
	return this.Name + this.Symbol
}

func (this *Strategy) Match(symbol string) bool {
	if this.Symbol == symbol {
		return true
	} else {
		return false
	}
}

func (this *Strategy) OnStrategyStart() {

}

func (this *Strategy) OnStrategyStop() {

}

func (this *Strategy) OnBarOpen(bar bar.Bar) {

}

func (this *Strategy) OnBar(bar bar.Bar) {

}

func (this *Strategy) OnBarSlice(size int) {

}

func (this *Strategy) Draw(table int, indicator_ xbase.IIndecator) {

	indicatorSlice, ok := this.mapIndicator[table]
	if ok {
		indicatorSlice = append(indicatorSlice, indicator_)
		this.mapIndicator[table] = indicatorSlice
	} else {
		this.mapIndicator[table] = []xbase.IIndecator{indicator_}
	}
}

func (this *Strategy) DoSvgDrawing() []string {

	tables := []string{}

	if !this.drawed {
		this.drawed = true

		for _, indicatorSlice := range this.mapIndicator {

			newCanvas := canvas.NewCanvas()
			this.drawCanvas = append(this.drawCanvas, newCanvas)

			// do drawing
			for _, oneIndicator := range indicatorSlice {
				newCanvas.Draw(oneIndicator)
			}
		}
	}

	for _, canvas := range this.drawCanvas {
		tables = append(tables, canvas.GetResult())
	}
	return tables
}

// for easy sending an Market Sell Order.
// kepp a refer to the returned order. check it's status(StatusFilled/StatusRejected)
// on the next bar
func (this *Strategy) Sell(qty int, remark string) *order.Order {
	o := order.NewMarketOrder(order.SideBuy, qty, remark)
	o.Symbol = this.Symbol

	base.Send(o)
	return o
}

func (this *Strategy) Buy(qty int, remark string) *order.Order {
	o := order.NewMarketOrder(order.SideBuy, qty, remark)
	o.Symbol = this.Symbol
	base.Send(o)
	return o
}

type IStrategy interface {
	Init(string, *series.BarSeries)
	Key() string
	Match(string) bool

	OnStrategyStart()
	OnStrategyStop()

	OnBarOpen(bar.Bar)
	OnBar(bar.Bar)
	OnBarSlice(int)
	// public virtual void OnTrade(Trade trade)
	// public virtual void OnQuote(Quote quote)

	DoSvgDrawing() []string
}
