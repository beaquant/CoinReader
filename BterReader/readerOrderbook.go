package BterReader

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
func (ths *BterReader) ReadOrderbook() bool {
    ret, err := ths.readOrderBook(ths.MonetaryName,ths.CoinName)
    if err == nil {
        orderCh := make(chan bool)
        historyCh := make(chan bool)
        go ths.decodeOrderBook(ret, orderCh)
        go ths.decodeHistory(ret, historyCh)
        bOrder := <- orderCh
        bHistory := <- historyCh
        if bOrder && bHistory {
            return true
        }
    }
    return false
}

func (ths *BterReader) readOrderBook(mon, coin string) (map[string]interface{},error) {
    address := ths.BaseAddress + "json_svr/query/?u=12" + ths.pageRand()
    data := fmt.Sprintf("type=ask_bid_list&symbol=%s_%s&tid=%s", coin, mon, ths.lastTradeID)
    var ret interface{}
    var err error
    if ths.UseProxy() {
        client := rhttp.GetProxyClient(ths.ProxyAddress, ths.ProxyPort)
        ret,err = rhttp.HttpClientPostForm(client, address, rhttp.HTTP_RETURN_TYPE_MAP, data)
        if err != nil {
            return nil,err
        }
    } else {
        ret,err = rhttp.HttpPostForm(address, rhttp.HTTP_RETURN_TYPE_MAP, data)
        if err != nil {
            return nil,err
        }
    }
    return ret.(map[string]interface{}),nil
}

func (ths *BterReader) decodeOrderBook(ret map[string]interface{}, ch chan bool) {
    if ret == nil {
        ch <- false
        return
    }
    rslt,ok := ret["result"]
    if !ok || (rslt.(bool) == false) {
        ch <- false
        return
    }
    
    orders,ok := ret["orders"]
    if !ok || orders == nil {
        ch <- false
        return
    }
    ordersmap := orders.(map[string]interface{})

	buyOrder, bok := ordersmap["bids"]
	sellOrder, sok := ordersmap["asks"]
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

func (ths *BterReader) decodeSellOB(flag *sync.WaitGroup, list []interface{}) {
    defer flag.Done()
    
    relLen := len(list)
    sellList := make([]*Reader.OrderBook, relLen)
    ths.OrderLocker.Lock()
    ths.Orders[Reader.OrderSellStringKey] = sellList
    ths.OrderLocker.Unlock()

    lastPriceIndex := -1
    
    for index := 0; index < relLen; index++ {
        tmpMap := (list[index]).([]interface{})
        p, _ := strconv.ParseFloat(tmpMap[5].(string), 64)
        c, _ := strconv.ParseFloat(tmpMap[6].(string), 64)
        if lastPriceIndex != -1 {
            if sellList[lastPriceIndex].Price == p {
                sellList[lastPriceIndex].Amount += c
                sellList[lastPriceIndex].Calc()
                continue
            }
        }
        ob := &Reader.OrderBook{}
        ob.Price = p
        ob.Amount = c
        ob.Calc()
        lastPriceIndex++
        sellList[lastPriceIndex] = ob
    }
}

func (ths *BterReader) decodeBuyOB(flag *sync.WaitGroup, list []interface{}) {
    defer flag.Done()
    
    relLen := len(list)
    buyList := make([]*Reader.OrderBook, relLen)
    ths.OrderLocker.Lock()
    ths.Orders[Reader.OrderBuyStringKey] = buyList
    ths.OrderLocker.Unlock()
    
    lastPriceIndex := -1
    
    for index := 0; index < relLen; index++ {
        tmpMap := (list[index]).([]interface{})
        p, _ := strconv.ParseFloat(tmpMap[5].(string), 64)
        c, _ := strconv.ParseFloat(tmpMap[6].(string), 64)
        if lastPriceIndex != -1 {
            if buyList[lastPriceIndex].Price == p {
                buyList[lastPriceIndex].Amount += c
                buyList[lastPriceIndex].Calc()
                continue
            }
        }
        ob := &Reader.OrderBook{}
        ob.Price = p
        ob.Amount = c
        ob.Calc()
        lastPriceIndex++
        buyList[lastPriceIndex] = ob
    }
}


// pageRand
func (ths *BterReader) pageRand() string {
    curRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("&c=%d", curRand.Intn(299999))
}