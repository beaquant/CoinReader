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
    
    proxyUse     bool
    ProxyAddress string // proxy server address. If not set, proxy is not be used.
    ProxyPort    string // proxy server port. If not set, 8181 is default.
}

// Init init parameters
// m is MonetaryName string
// c is CoinName string
// v is optional parameters, Set the parameter order described below:
//  v[0] -- proxyAddress string
//  v[1] -- proxyPort string
func (ths *ReaderDef) Init(m, c string, v ... interface{}) {
    ths.MonetaryName = m
    ths.CoinName = c
    
    ths.Orders = make(map[string][]*OrderBook)
    ths.Orders[OrderBuyStringKey] = nil
    ths.Orders[OrderSellStringKey] = nil
    
    vLen := len(v)
    if vLen >= 1 {
        switch v[0].(type) {
        case string:
            ths.ProxyAddress = v[0].(string)
        default:
            panic("The first parameter(ProxyAddress) must be of type string!")
        }
    }
    if vLen >= 2 {
        switch v[1].(type) {
        case string:
            ths.ProxyPort = v[1].(string)
        default:
            panic("The second parameter(ProxyPort) must be of type string!")
        }
    }
    if len(ths.ProxyAddress) > 0 {
        ths.proxyUse = true
        if len(ths.ProxyPort) == 0 {
            ths.ProxyPort = "8181"
        }
    } else {
        ths.proxyUse = false
    }
}

// UseProxy Return use or not use a proxy configuration.
func (ths *ReaderDef) UseProxy() bool {
    return ths.proxyUse
}