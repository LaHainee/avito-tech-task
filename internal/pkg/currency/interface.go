package currency

type Service interface {
	Update()
	Get(string) (float64, error)
}
