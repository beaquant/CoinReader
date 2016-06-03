package Btc38Reader

import(
    "fmt"
    "time"
    "strconv"
)


// ReadHistory readout histroy datas from btc38.com, datas saved in History datas
func (ths *B38Reader) ReadHistory() bool {
    return ths.ReadOrderbook()
}

// PrintHistory print histroy datas to string
func (ths *B38Reader) PrintHistory(length int) string {
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

func (ths *B38Reader) decodeHistory(ret map[string]interface{}, ch chan bool) {
    if ret == nil {
        ch <- false
        return
    }
    trade, ok := ret["trade"]
    if ok && trade != nil {
		tradeMap := trade.([]interface{})
		relLen := len(tradeMap)
		ths.Historys = make([]*B38History, relLen)

        for index := 0 ; index < relLen ; index++ {
			tmpMap := (tradeMap[index]).(map[string]interface{})
            ph := &B38History{}
			for k, v := range tmpMap {
				switch k {
				case "price":
					ph.Price, _ = strconv.ParseFloat(v.(string), 64)
				case "volume":
					ph.Amount, _ = strconv.ParseFloat(v.(string), 64)
				case "type":
					itype, _ := strconv.Atoi(v.(string))
                    if itype == 1 {
                        ph.Type = "Buy"
                    } else {
                        ph.Type = "Sell"
                    }
				case "time":
					ph.DateTime,_ = time.Parse("2006-01-02 15:04:05",v.(string))
				}
			}
			ph.Total = ph.Price * ph.Amount
            ths.Historys[index] = ph
		}
        ch <- true
    }
    
    ch <- false
}