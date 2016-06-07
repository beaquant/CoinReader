package BterReader

import(
    "fmt"
    "time"
    "strconv"
    "strings"
)


// ReadHistory readout histroy datas from btc38.com, datas saved in History datas
func (ths *BterReader) ReadHistory() bool {
    return ths.ReadOrderbook()
}

// PrintHistory print histroy datas to string
func (ths *BterReader) PrintHistory(length int) string {
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

func (ths *BterReader) decodeHistory(ret map[string]interface{}, ch chan bool) {
    if ret == nil {
        ch <- false
        return
    }
    
	trade,ok := ret["history"]
    
    if ok && trade != nil {
        historymap := trade.(map[string]interface{})
        rslt,ok := historymap["result"]
        if !ok || (rslt.(bool) == false) {
            ch <- false
            return
        }
        histy,ok := historymap["history"]
        if ok && histy != nil {
            tradeMap := histy.([]interface{})
            relLen := len(tradeMap)
            ths.Historys = make([]*BterHistory, relLen)

            for index := 0 ; index < relLen ; index++ {
                tmpMap := (tradeMap[index]).([]interface{})
                ph := &BterHistory{}
                ph.Price, _ = strconv.ParseFloat(tmpMap[2].(string), 64)
                ph.Amount, _ = strconv.ParseFloat(tmpMap[3].(string), 64)
                if strings.HasPrefix(tmpMap[6].(string),"buy") {
                    ph.Type = "Buy"
                } else {
                    ph.Type = "Sell"
                }
                ph.DateTime,_ = time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%d-%s",time.Now().Year(),tmpMap[0].(string)))
                
                ph.Total = ph.Price * ph.Amount
                ths.Historys[index] = ph
            }
            ch <- true
        }
    }
    
    ch <- false
}