package PoloniexReader

import(
    "github.com/jojopoper/CoinReader/Reader"
)

// PHistory Trade histroy data struct
type PHistory struct {
    Reader.History
    GlobalTradeID uint64
    TradeID uint64
}

// PReader Poloniex reader struct
type PReader struct {
    Reader.ReaderDef
    OrderDepth int
    
    Historys []*PHistory
}

// Init init parameters
func (ths *PReader) Init(m, c string) {
    ths.ReaderDef.Init(m,c)
    ths.BaseAddress = "https://poloniex.com/"
    ths.OrderDepth = 20
    ths.Historys = nil
}