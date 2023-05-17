package handler

type CalcService struct{}

func (s *CalcService) Add(request int, reply *int) error {
	*reply = request + 10
	return nil
}

func (s *CalcService) ServiceName() string {
	return "CalcService"
}
