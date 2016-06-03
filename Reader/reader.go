package Reader


type ReaderInterface interface {
    ReadAll() bool
    ReadHistory() bool
    PrintHistory(length int) string
    ReadOrderbook() bool
    PrintOrderBook(length int) string
}

// ReaderDef reader base struct
type ReaderDef struct {
    BaseAddress  string
    MonetaryName string
    CoinName     string
    
    Orders   map[string][]*OrderBook
}

// Init init parameters
// m is MonetaryName string
// c is CoinName string
func (ths *ReaderDef) Init(m, c string) {
    ths.MonetaryName = m
    ths.CoinName = c
    
    ths.Orders = make(map[string][]*OrderBook)
    ths.Orders[OrderBuyStringKey] = nil
    ths.Orders[OrderSellStringKey] = nil
}