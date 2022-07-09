package supertrend

import (
	"github.com/c9s/bbgo/pkg/bbgo"
	"github.com/c9s/bbgo/pkg/indicator"
	"github.com/c9s/bbgo/pkg/types"
)

type DoubleDema struct {
	DemaInterval types.Interval `json:"demaInterval"`

	// FastDEMAWindow DEMA window for checking breakout
	FastDEMAWindow int `json:"fastDEMAWindow"`
	// SlowDEMAWindow DEMA window for checking breakout
	SlowDEMAWindow int `json:"slowDEMAWindow"`
	fastDEMA       *indicator.DEMA
	slowDEMA       *indicator.DEMA
}

// getDemaSignal get current DEMA signal
func (dd *DoubleDema) getDemaSignal(openPrice float64, closePrice float64) types.Direction {
	var demaSignal types.Direction = types.DirectionNone

	if closePrice > dd.fastDEMA.Last() && closePrice > dd.slowDEMA.Last() && !(openPrice > dd.fastDEMA.Last() && openPrice > dd.slowDEMA.Last()) {
		demaSignal = types.DirectionUp
	} else if closePrice < dd.fastDEMA.Last() && closePrice < dd.slowDEMA.Last() && !(openPrice < dd.fastDEMA.Last() && openPrice < dd.slowDEMA.Last()) {
		demaSignal = types.DirectionDown
	}

	return demaSignal
}

// preloadDema preloads DEMA indicators
func (dd *DoubleDema) preloadDema(kLineStore *bbgo.MarketDataStore) {
	if klines, ok := kLineStore.KLinesOfInterval(dd.fastDEMA.Interval); ok {
		for i := 0; i < len(*klines); i++ {
			dd.fastDEMA.Update((*klines)[i].GetClose().Float64())
		}
	}
	if klines, ok := kLineStore.KLinesOfInterval(dd.slowDEMA.Interval); ok {
		for i := 0; i < len(*klines); i++ {
			dd.slowDEMA.Update((*klines)[i].GetClose().Float64())
		}
	}
}

// setupDoubleDema initializes double DEMA indicators
func (dd *DoubleDema) setupDoubleDema(kLineStore *bbgo.MarketDataStore, interval types.Interval) {
	dd.DemaInterval = interval

	// DEMA
	if dd.FastDEMAWindow == 0 {
		dd.FastDEMAWindow = 144
	}
	dd.fastDEMA = &indicator.DEMA{IntervalWindow: types.IntervalWindow{Interval: dd.DemaInterval, Window: dd.FastDEMAWindow}}
	dd.fastDEMA.Bind(kLineStore)

	if dd.SlowDEMAWindow == 0 {
		dd.SlowDEMAWindow = 169
	}
	dd.slowDEMA = &indicator.DEMA{IntervalWindow: types.IntervalWindow{Interval: dd.DemaInterval, Window: dd.SlowDEMAWindow}}
	dd.slowDEMA.Bind(kLineStore)

	dd.preloadDema(kLineStore)
}