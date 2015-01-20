package strategy

import (
	"fmt"
	"quant/base/bar"
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
	Name      string
	Symbol    string
	BarSeries *bar.BarSeries
}

func (this *Strategy) Init(symbol string, barSeries *bar.BarSeries) {
	if debug {
		fmt.Println("Strategy.Init()")
	}
	this.Name = "Strategy"
	this.Symbol = symbol
	this.BarSeries = barSeries
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

type IStrategy interface {
	Init(string, *bar.BarSeries)
	Match(string) bool

	OnStrategyStart()
	OnStrategyStop()

	OnBarOpen(bar.Bar)
	OnBar(bar.Bar)
	OnBarSlice(int)

	// public virtual void OnTrade(Trade trade)
	// public virtual void OnQuote(Quote quote)
}
