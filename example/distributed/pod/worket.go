package pod

type Worker struct {
	Name     string
	TaskList []string
}

func NewWorker(name string, taskList []string) *Worker {
	return &Worker{
		Name:     name,
		TaskList: taskList,
	}
}

func (w *Worker) Key() string {
	return "worker:" + w.Name
}

func (w *Worker) Task() {

}
