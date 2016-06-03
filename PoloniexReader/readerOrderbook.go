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
    ret,err := rhttp.HttpGet(address,rhttp.HTTP_RETURN_TYPE_MAP)
    if err != nil {
        fmt.Println(err)
        return nil,err
    }
    return ret.(map[string]interface{}),nil
}

// PrintOrderBook print order datas to string
func (ths *PReader) PrintOrderBook(length int) string {
    if ths.Orders == nil {
        return "> No datas!!\r\n"
    }
    
    buyList,_ := ths.Orders[Reader.OrderBuyStringKey]
    sellList,_ := ths.Orders[Reader.OrderSellStringKey]
    
    relLen := len(buyList)
    relLenSell := len(sellList)
    if relLen > relLenSell {
        relLen = relLenSell
    }
    if length != -1 && length < relLen {
        relLen = length
    }
    
    ret := fmt.Sprintf("\r\n>  %s / %s Open orders (Records length = %d)\r\n",
        ths.MonetaryName, ths.CoinName, relLen)
    //> Price          Amount       Total             Price          Amount       Total
    //      > 0.00001071    868.80058877    0.00930485              0.00001074      15933.88623733  0.17112994
    ret += ">      ************ Buy ************                         ************ Sell ************ \r\n"
    ret += "> Price         Amount          Total                   Price           Amount          Total\r\n"
    
    for index := 0; index < relLen; index++ {
        bItm := buyList[index]
        sItm := sellList[index]
        ret += fmt.Sprintf("> %.8f\t%.8f\t%.8f\t\t%.8f\t%.8f\t%.8f\r\n", 
            bItm.Price, bItm.Amount, bItm.Total,
            sItm.Price, sItm.Amount, sItm.Total)
    }
    return ret
}

func (ths *PReader) decodeOrderBook(ret map[string]interface{}) bool {
    if ret == nil {
        return false
    }
    asks,_ := ret["asks"]
    bids,_ := ret["bids"]
    f := new(sync.WaitGroup)
    f.Add(2)
    go ths.decodeAsksOB(f, asks.([]interface{}))
    go ths.decodeBidsOB(f, bids.([]interface{}))
    f.Wait()
    return true
}

func (ths *PReader) decodeAsksOB(flag *sync.WaitGroup, list []interface{}) {
    defer flag.Done()
    relLen := len(list)
    sellList := make([]*Reader.OrderBook, relLen)
    ths.Orders[Reader.OrderSellStringKey] = sellList
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
    ths.Orders[Reader.OrderBuyStringKey] = buyList
    for index := 0; index < relLen; index++ {
        itm := list[index].([]interface{})
        ob := &Reader.OrderBook{}
        ob.Price,_ = strconv.ParseFloat(itm[0].(string),64)
        ob.Amount = itm[1].(float64)
        ob.Calc()
        buyList[index] = ob
    }
}