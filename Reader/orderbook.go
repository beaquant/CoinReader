package Reader

const (
    OrderBuyStringKey  = "Buy"
    OrderSellStringKey = "Sell"
)

//OrderBook buy&sell data struct
type OrderBook struct {
    Price  float64
    Amount float64
    Total  float64
}

// Calc calculate Total value
func (ths *OrderBook) Calc() {
    ths.Total = ths.Price * ths.Amount
}