package Btc38Reader

import(
    "fmt"
    "github.com/jojopoper/CoinReader/Reader"
    "github.com/jojopoper/CoinReader/rhttp"
    "math/rand"
    "strconv"
	"sync"
    "time"
)

// ReadOrderbook readout open orders from poloniex.com
func (ths *B38Reader) ReadOrderbook() bool {
    ret, err := ths.readOrderBook(ths.MonetaryName,ths.CoinName)
    if err == nil {
        orderCh := make(chan bool)
        historyCh := make(chan bool)
        go ths.decodeOrderBook(ret,orderCh)
        go ths.decodeHistory(ret,historyCh)
        bOrder := <- orderCh
        bHistory := <- historyCh
        if bOrder && bHistory {
            return true
        }
    }
    return false
}

func (ths *B38Reader) readOrderBook(mon, coin string) (map[string]interface{},error) {
    curRand := rand.New(rand.NewSource(time.Now().UnixNano()))
    address := ths.BaseAddress + fmt.Sprintf("trade/getTradeList30.php?mk_type=%s&coinname=%s&n=0.00%d1",
        mon, coin, curRand.Int31())
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

func (ths *B38Reader) decodeOrderBook(ret map[string]interface{}, ch chan bool) {
    if ret == nil {
        ch <- false
        return
    }
    buyOrder, bok := ret["buyOrder"]
    sellOrder, sok := ret["sellOrder"]
    f := new(sync.WaitGroup)
    if bok && buyOrder != nil {
        f.Add(1)
        go ths.decodeBuyOB(f, buyOrder.([]interface{}))
    } else {
        ch <- false
        return 
    }
    if sok && sellOrder != nil {
        f.Add(1)
        go ths.decodeSellOB(f, sellOrder.([]interface{}))
    } else {
        f.Wait()
        ch <- false
        return 
    }
    f.Wait()
    ch <- true
    return
}

func (ths *B38Reader) decodeSellOB(flag *sync.WaitGroup, list []interface{}) {
    defer flag.Done()
    relLen := len(list)
    sellList := make([]*Reader.OrderBook, relLen)
    ths.OrderLocker.Lock()
    ths.Orders[Reader.OrderSellStringKey] = sellList
    ths.OrderLocker.Unlock()

    for index := 0; index < relLen; index++ {
        tmpMap := (list[index]).(map[string]interface{})
        ob := &Reader.OrderBook{}
        for k, v := range tmpMap {
            switch k {
            case "price":
                ob.Price, _ = strconv.ParseFloat(v.(string), 64)
            case "amount":
                ob.Amount, _ = strconv.ParseFloat(v.(string), 64)
            }
        }
        ob.Calc()
        sellList[index] = ob
    }
}

func (ths *B38Reader) decodeBuyOB(flag *sync.WaitGroup, list []interface{}) {
    defer flag.Done()
    relLen := len(list)
    buyList := make([]*Reader.OrderBook, relLen)
    ths.OrderLocker.Lock()
    ths.Orders[Reader.OrderBuyStringKey] = buyList
    ths.OrderLocker.Unlock()

    for index := 0; index < relLen; index++ {
        tmpMap := (list[index]).(map[string]interface{})
        ob := &Reader.OrderBook{}
        for k, v := range tmpMap {
            switch k {
            case "price":
                ob.Price, _ = strconv.ParseFloat(v.(string), 64)
            case "amount":
                ob.Amount, _ = strconv.ParseFloat(v.(string), 64)
            }
        }
        ob.Calc()
        buyList[index] = ob
    }
}