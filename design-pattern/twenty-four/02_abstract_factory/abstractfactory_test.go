package abstractfactory

import "testing"

func getMainAndDetail(factory DAOFactory) {
	factory.CreateOrderMainDAO().SaveOrderMain()
	factory.CreateOrderDetailDAO().SaveOrderDetail()
}

func TestRdbFactory(t *testing.T) {
	var factor DAOFactory = &RDBDAOFactory{}
	getMainAndDetail(factor)
}

func TestXMLFactory(t *testing.T) {
	var factor DAOFactory = &XMLDAOFactory{}
	getMainAndDetail(factor)
}
