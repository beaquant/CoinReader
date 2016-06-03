package PoloniexReader

import(
    "fmt"
    "github.com/jojopoper/CoinReader/rhttp"
    "time"
    "strconv"
)


// ReadHistory readout histroy datas from poloniex.com, datas saved in History datas
func (ths *PReader) ReadHistory() bool {
    ret, err := ths.readHistory(ths.MonetaryName,ths.CoinName)
    if err == nil {
        return ths.decodeHistory(ret)
    }
    return false
}

// PrintHistory print histroy datas to string
func (ths *PReader) PrintHistory(length int) string {
    if ths.Historys == nil {
        return "> No datas!!\r\n"
    }
    
    relLen := len(ths.Historys)
    if length != -1 && length < relLen {
        relLen = length
    }
    
    ret := fmt.Sprintf("\r\n>  %s / %s Trade history datas (Records length = %d)\r\n",
        ths.MonetaryName, ths.CoinName, relLen)
    //> 2016-06-02 09:58:21   buy     0.00001069      187.09073900    0.00199999
    ret += "> DateTime              Type    Price           Amount          Total\r\n"
    
    for index := 0; index < relLen; index++ {
        his := ths.Historys[index]
        ret += fmt.Sprintf("> %s\t%s\t%.8f\t%.8f\t%.8f\r\n", 
            his.DateTime.Format("2006-01-02 15:04:05"),
            his.Type, his.Price, his.Amount, his.Total)
    }
    return ret
}

func (ths *PReader) readHistory(mon, coin string) ([]interface{},error) {
    address := ths.BaseAddress + fmt.Sprintf("public?command=returnTradeHistory&currencyPair=%s_%s",mon,coin)
    ret,err := rhttp.HttpGet(address,rhttp.HTTP_RETURN_TYPE_SLICE)
    if err != nil {
        return nil,err
    }
    return ret.([]interface{}),nil
}

func (ths *PReader) decodeHistory(ret []interface{}) bool {
    if ret == nil {
        return false
    }
    retLen := len(ret)
    ths.Historys = make([]*PHistory, retLen)
    
    for index := 0 ; index < retLen ; index++ {
        retMap := ret[index].(map[string]interface{})
        ph := &PHistory{}
        ph.GlobalTradeID,_ = strconv.ParseUint(fmt.Sprintf("%.0f",retMap["globalTradeID"].(float64)),10,64)
        ph.TradeID,_ = strconv.ParseUint(fmt.Sprintf("%.0f",retMap["tradeID"].(float64)),10,64)
        ph.DateTime,_ = time.Parse("2006-01-02 15:04:05",retMap["date"].(string))
        ph.Type = retMap["type"].(string)
        ph.Price,_ = strconv.ParseFloat(retMap["rate"].(string),64)
        ph.Amount,_ = strconv.ParseFloat(retMap["amount"].(string),64)
        ph.Total,_ = strconv.ParseFloat(retMap["total"].(string),64)
        ths.Historys[index] = ph
    }
    
    return true
}