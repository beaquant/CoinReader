package BterReader

import(
    "github.com/jojopoper/CoinReader/Reader"
)

// BterHistory Trade histroy data struct
type BterHistory struct {
    Reader.History
}

// BterReader Bter reader struct
type BterReader struct {
    Reader.ReaderDef
    
    Historys    []*BterHistory
    lastTradeID string
}

// Init init parameters
func (ths *BterReader) Init(m, c string, v ... interface{}) {
    ths.ReaderDef.Init(m,c,v...)
    ths.BaseAddress = "https://bter.com/"
    ths.Historys = nil
    ths.lastTradeID = "0"
}

// ReadAll Read all of history and order datas
// Bter.com, History&Order all data can be read in one request.
func (ths *BterReader) ReadAll() bool {
    return ths.ReadOrderbook()
}