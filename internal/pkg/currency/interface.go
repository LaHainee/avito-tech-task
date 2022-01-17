package currency

//go:generate moq -out ./mock/currency_mock.go -pkg mock . ConverterIface:MockConverterIface
type ConverterIface interface {
	Update()
	Get(string) (float64, error)
}
