package repository

import (
	"errors"
	"product-api/internal/model"
	"sync"
	"time"
)

type state int

const (
	StateClosed state = iota
	StateOpen
	StateHalfOpen
)

var ErrCircuitOpen = errors.New("circuit breaker terbuka: database sedang bermasalah")

type CircuitBreakerRepository struct {
	next ProductRepository

	mu               sync.Mutex
	state            state
	failureCount     int
	failureThreshold int
	openUntil        time.Time
	openDuration     time.Duration
}

func NewCircuitBreakerRepository(next ProductRepository) *CircuitBreakerRepository {
	return &CircuitBreakerRepository{
		next:             next,
		state:            StateClosed,
		failureThreshold: 5,
		openDuration:     10 * time.Second,
	}
}

func (c *CircuitBreakerRepository) beforeRequest() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch c.state {
	case StateOpen:
		if time.Now().After(c.openUntil) {
			c.state = StateHalfOpen
			return nil 
		}
		return ErrCircuitOpen

	case StateHalfOpen:
		return ErrCircuitOpen

	default: 
		return nil
	}
}

func (c *CircuitBreakerRepository) afterReq(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err != nil && !errors.Is(err, ErrProductNotFound) {
		c.failureCount++
		if c.state == StateHalfOpen || c.failureCount >= c.failureThreshold {
			c.state = StateOpen
			c.openUntil = time.Now().Add(c.openDuration)
		}
		return
	}

	c.state = StateClosed
	c.failureCount = 0
}

func (c *CircuitBreakerRepository) GetAll()([]model.Product, error){
	if err := c.beforeRequest(); err != nil{
		return nil, err
	}

	products, err := c.next.GetAll()
	c.afterReq(err)
	return products, err
}

func (c *CircuitBreakerRepository) GetByID(id int) (model.Product, error) {
	if err := c.beforeRequest(); err != nil {
		return model.Product{}, err
	}

	product, err := c.next.GetByID(id)
	c.afterReq(err)
	return product, err
}

func (c *CircuitBreakerRepository) Create(p model.Product)(model.Product, error){
	if err := c.beforeRequest(); err != nil {
		return model.Product{}, err
	}

	created, err := c.next.Create(p)
	c.afterReq(err)
	return created, err
}

func (c *CircuitBreakerRepository) Update(id int, updated model.Product)(model.Product, error){
	if err := c.beforeRequest(); err != nil {
		return model.Product{}, err
	}

	updated, err := c.next.Update(id, updated)
	c.afterReq(err)
	return updated, err
}

func (c *CircuitBreakerRepository) Delete(id int) error{
	if err := c.beforeRequest(); err != nil{
		return err
	}

	err := c.next.Delete(id)
	c.afterReq(err)
	return err
}