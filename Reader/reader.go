package Reader

import (
    "fmt"
    "sync"
)

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
    
    Orders      map[string][]*OrderBook
    OrderLocker *sync.Mutex
    
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
    ths.OrderLocker = new(sync.Mutex)
    
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


// PrintOrderBook print order datas to string
func (ths *ReaderDef) PrintOrderBook(length int) string {
    if ths.Orders == nil {
        return "> No datas!!\r\n"
    }
    
    buyList,_ := ths.Orders[OrderBuyStringKey]
    sellList,_ := ths.Orders[OrderSellStringKey]
    
    relLenBuy := len(buyList)
    relLenSell := len(sellList)
    
    if length != -1{
        if length < relLenBuy {
            relLenBuy = length
        }
        if length < relLenSell {
            relLenSell = length
        }
    }
    
    ret := fmt.Sprintf("\r\n>  %s / %s Open orders (Records length = %d)\r\n",
        ths.MonetaryName, ths.CoinName, length)
    //> Price          Amount       Total             Price          Amount       Total
    //      > 0.00001071    868.80058877    0.00930485              0.00001074      15933.88623733  0.17112994
    ret += ">      ************ Buy ************                         ************ Sell ************ \r\n"
    ret += "> Price         Amount          Total                   Price           Amount          Total\r\n"
    indexBuy := 0
    indexSell := 0
    for ; indexBuy < relLenBuy || indexSell < relLenSell; {
        if (indexBuy < relLenBuy) && (indexSell < relLenSell) {
            bItm := buyList[indexBuy]
            sItm := sellList[indexSell]
            ret += fmt.Sprintf("> %.8f\t%.8f\t%.8f\t\t%.8f\t%.8f\t%.8f\r\n", 
                bItm.Price, bItm.Amount, bItm.Total,
                sItm.Price, sItm.Amount, sItm.Total)
        } else if (indexBuy >= relLenBuy) && (indexSell < relLenSell) {
            sItm := sellList[indexSell]
            //                  > 0.00001071    8.80058877    0.00930485              0.00001074      15933.88623733  0.17112994
            ret += fmt.Sprintf("> -         \t-         \t-         \t\t%.8f\t%.8f\t%.8f\r\n", 
                sItm.Price, sItm.Amount, sItm.Total)
        } else if (indexBuy < relLenBuy) && (indexSell >= relLenSell) {
            bItm := buyList[indexBuy]
            //                  > 0.00001071    8.80058877    0.00930485              0.00001074      15933.88623733  0.17112994
            ret += fmt.Sprintf("> %.8f\t%.8f\t%.8f\t\t-         \t-         \t-\r\n", 
                bItm.Price, bItm.Amount, bItm.Total)
        } else {
            break
        }
        indexBuy++
        indexSell++
    }
    return ret
}