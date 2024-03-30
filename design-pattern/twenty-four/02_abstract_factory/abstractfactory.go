package abstractfactory

// 用RDB和XML存储订单信息，抽象工厂分别能生成相关的主订单信息和订单详情信息。
import "log"

type OrderMainDAO interface {
	SaveOrderMain()
}

type OrderDetailDAO interface {
	SaveOrderDetail()
}

type DAOFactory interface {
	CreateOrderMainDAO() OrderMainDAO
	CreateOrderDetailDAO() OrderDetailDAO
}

// rdb 关系型数据库

type RDBMainDAO struct{}

func (*RDBMainDAO) SaveOrderMain() {
	log.Println("rdb main save")
}

type RDBDetailDAO struct{}

func (*RDBDetailDAO) SaveOrderDetail() {
	log.Println("rdb detail save")
}

// 获取实例

type RDBDAOFactory struct{}

func (*RDBDAOFactory) CreateOrderMainDAO() OrderMainDAO {
	return &RDBMainDAO{}
}

func (*RDBDAOFactory) CreateOrderDetailDAO() OrderDetailDAO {
	return &RDBDetailDAO{}
}

// xml 存储

type XMLMainDAO struct{}

func (*XMLMainDAO) SaveOrderMain() {
	log.Println("xml main save")
}

type XMLDetailDAO struct{}

func (*XMLDetailDAO) SaveOrderDetail() {
	log.Println("xml detail save")
}

// 获取实例

type XMLDAOFactory struct{}

func (*XMLDAOFactory) CreateOrderMainDAO() OrderMainDAO {
	return &XMLMainDAO{}
}

func (*XMLDAOFactory) CreateOrderDetailDAO() OrderDetailDAO {
	return &XMLDetailDAO{}
}
