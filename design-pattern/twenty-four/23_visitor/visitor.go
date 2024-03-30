package visitor

import "log"

type Customer interface {
	Accept(Visitor)
}

type Visitor interface {
	Visit(Customer)
}

type EnterpriseCustomer struct {
	name string
}

func NewEnterpriseCustomer(name string) *EnterpriseCustomer {
	return &EnterpriseCustomer{
		name: name,
	}
}

func (e *EnterpriseCustomer) Accept(visitor Visitor) {
	visitor.Visit(e)
}

type CustomerCol struct {
	customers []Customer
}

func (c *CustomerCol) Accept(visitor Visitor) {
	for _, customer := range c.customers {
		customer.Accept(visitor)
	}
}

func (c *CustomerCol) Add(customer Customer) {
	c.customers = append(c.customers, customer)
}

type IndividualCustomer struct {
	name string
}

func NewIndividualCustomer(name string) *IndividualCustomer {
	return &IndividualCustomer{
		name: name,
	}
}

func (i *IndividualCustomer) Accept(visitor Visitor) {
	visitor.Visit(i)
}

type ServiceRequestVisitor struct{}

func (*ServiceRequestVisitor) Visit(customer Customer) {
	switch c := customer.(type) {
	case *EnterpriseCustomer:
		log.Printf("serving enterprise customer %s", c.name)
	case *IndividualCustomer:
		log.Printf("serving individual customer %s", c.name)
	}
}

type AnalysisVisitor struct{}

func (*AnalysisVisitor) Visit(customer Customer) {
	switch c := customer.(type) {
	case *EnterpriseCustomer:
		log.Printf("Analysis enterprise customer %s", c.name)
	}
}
