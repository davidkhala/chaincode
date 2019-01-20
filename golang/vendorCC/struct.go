package main

type Project struct {
	//for party A
	Title       string
	StackHolder []string
	Requirement []string
	Schedule    []Step
}
type Submit struct {
	DeliveryURL string
	ID          string
}
type Step struct {
	//for party A
	Installment int
	ID          string
	DeadLine    string
	Status      string
	lastSubmit  Submit
	lastAudit   Audit
	lastReview  Review
}
type Audit struct {
	ID      string
	Status  string
	Comment string
	Time    string
}
type Review struct {
	Status  string
	Comment string
	ID      string
	Time    string
}