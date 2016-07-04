package PoloniexReader

import(
    "fmt"
    "github.com/jojopoper/CoinReader/Reader"
    "github.com/jojopoper/CoinReader/rhttp"
    "strconv"
	"sync"
)

// ReadOrderbook readout open orders from poloniex.com
func (ths *PReader) ReadOrderbook() bool {
    ret, err := ths.readOrderBook(ths.MonetaryName,ths.CoinName)
    if err == nil {
        return ths.decodeOrderBook(ret)
    }
    return false
}

func (ths *PReader) readOrderBook(mon, coin string) (map[string]interface{},error) {
    address := ths.BaseAddress + fmt.Sprintf("public?command=returnOrderBook&depth=%d&currencyPair=%s_%s",
        ths.OrderDepth, mon, coin)
    var ret interface{}
    var err error
    if ths.UseProxy() {
        ret,err = rhttp.HttpProxyGet(address,ths.ProxyAddress,ths.ProxyPort,rhttp.HTTP_RETURN_TYPE_MAP)
        if err != nil {
            return nil,err
        }
    } else {
        ret,err = rhttp.HttpGet(address,rhttp.HTTP_RETURN_TYPE_MAP)
        if err != nil {
            return nil,err
        }
    }
    return ret.(map[string]interface{}),nil
}

func (ths *PReader) decodeOrderBook(ret map[string]interface{}) bool {
    if ret == nil {
        return false
    }
    asks,_ := ret["asks"]
    bids,_ := ret["bids"]
    f := new(sync.WaitGroup)
    if asks != nil && bids != nil {
        f.Add(2)
        go ths.decodeAsksOB(f, asks.([]interface{}))
        go ths.decodeBidsOB(f, bids.([]interface{}))
        f.Wait()
        return true
    }
    return false
}

func (ths *PReader) decodeAsksOB(flag *sync.WaitGroup, list []interface{}) {
    defer flag.Done()
    relLen := len(list)
    sellList := make([]*Reader.OrderBook, relLen)
    ths.OrderLocker.Lock()
    ths.Orders[Reader.OrderSellStringKey] = sellList
    ths.OrderLocker.Unlock()

    for index := 0; index < relLen; index++ {
        itm := list[index].([]interface{})
        ob := &Reader.OrderBook{}
        ob.Price,_ = strconv.ParseFloat(itm[0].(string),64)
        ob.Amount = itm[1].(float64)
        ob.Calc()
        sellList[index] = ob
    }
}

func (ths *PReader) decodeBidsOB(flag *sync.WaitGroup, list []interface{}) {
    defer flag.Done()
    relLen := len(list)
    buyList := make([]*Reader.OrderBook, relLen)
    ths.OrderLocker.Lock()
    ths.Orders[Reader.OrderBuyStringKey] = buyList
    ths.OrderLocker.Unlock()
    
    for index := 0; index < relLen; index++ {
        itm := list[index].([]interface{})
        ob := &Reader.OrderBook{}
        ob.Price,_ = strconv.ParseFloat(itm[0].(string),64)
        ob.Amount = itm[1].(float64)
        ob.Calc()
        buyList[index] = ob
    }
}