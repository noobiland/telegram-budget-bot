package expence

type transactionStatus struct {
	sum      int
	category string
	payment  string
}

func (ts *transactionStatus) setCategory(category string) {
	ts.category = category
}

func (ts *transactionStatus) setPayment(payment string) {
	ts.payment = payment
}

func (ts *transactionStatus) reset() {
	ts.sum = 0
	ts.category = ""
	ts.payment = ""
}