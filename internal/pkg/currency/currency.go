package currency

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"avito-tech-task/config"
	"avito-tech-task/internal/pkg/constants"
	createdErrors "avito-tech-task/internal/pkg/errors"
)

type Converter struct {
	Rates map[string]float64 `json:"rates,omitempty"`

	logger *logrus.Logger
	config *config.Config
	mutex  *sync.RWMutex
}

func NewConverter(config *config.Config, logger *logrus.Logger) *Converter {
	currency := &Converter{
		config: config,
		logger: logger,
	}
	currency.mutex = new(sync.RWMutex)

	currency.logger.Info("Initializing currency data")

	resp, err := http.Get(currency.config.CurrencyApiURL)
	if err != nil {
		currency.logger.Fatalf("Could not get actual currency data: %s", err)
		return nil
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			currency.logger.Errorf("Could not close response body: %s", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		currency.logger.Fatalf("Bad response status: %d", resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		currency.logger.Fatalf("Could not read response body: %s", err)
		return nil
	}

	if err = json.Unmarshal(body, &currency); err != nil {
		currency.logger.Fatalf("Could not umarshal response body into struct: %s", err)
		return nil
	}
	currency.Rates["RUB"] = 1

	return currency
}

func (c *Converter) Update() {
	c.logger.Info("Updating currency data")

	resp, err := http.Get(c.config.CurrencyApiURL)
	if err != nil {
		c.logger.Errorf("Could not get actual currency data: %s", err)
		return
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			c.logger.Errorf("Could not close response body: %s", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		c.logger.Errorf("Bad response status: %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Errorf("Could not read response body: %s", err)
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	if err = json.Unmarshal(body, &c); err != nil {
		c.logger.Errorf("Could not umarshal response body into struct: %s", err)
		return
	}
	c.Rates["RUB"] = 1
}

func (c *Converter) Get(currency string) (float64, error) {
	var (
		value float64
		ok    bool
	)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if value, ok = c.Rates[currency]; !ok {
		return 0, createdErrors.ErrNotSupportedCurrency
	}

	return value, nil
}

func UpdateCurrency(converter *Converter, cancel <-chan struct{}) {
	for {
		select {
		case <-cancel:
			return
		default:
			converter.Update()
		}

		time.Sleep(constants.CurrencyAPIUpdatePeriod)
	}
}
