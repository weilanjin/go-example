package strategr

import "testing"

func TestPayByCash(t *testing.T) {
	payment := NewPayment("Tom", "", 123, &Cash{})
	payment.Pay()
}

func TestPayByBank(t *testing.T) {
	payment := NewPayment("Tom", "", 123, &Bank{})
	payment.Pay()
}
