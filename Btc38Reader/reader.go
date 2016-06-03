package Btc38Reader

import(
    "github.com/jojopoper/CoinReader/Reader"
)

// B38History Trade histroy data struct
type B38History struct {
    Reader.History
}

// B38Reader Btc38 reader struct
type B38Reader struct {
    Reader.ReaderDef
    
    Historys []*B38History
}

// Init init parameters
func (ths *B38Reader) Init(m, c string) {
    ths.ReaderDef.Init(m,c)
    ths.BaseAddress = "http://www.btc38.com/"
    ths.Historys = nil
}

// ReadAll Read all of history and order datas
// Btc38.com, History&Order all data can be read in one request.
func (ths *B38Reader) ReadAll() bool {
    return ths.ReadOrderbook()
}