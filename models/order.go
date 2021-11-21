package models

import "sync"

type Order struct {
	ID         int   `json:"id"`
	Items      []int `json:"items"`
	TableID    int   `json:"table_id"`
	WaiterID   int   `json:"waiter_id"`
	Priority   int   `json:"priority"`
	MaxWait    int   `json:"max_wait"`
	PickUpTime int   `json:"pick_up_time"`
}
type ReadyOrder struct {
	ID          int   `json:"id"`
	Items       []int `json:"items"`
	TableID     int   `json:"table_id"`
	WaiterID    int   `json:"waiter_id"`
	Priority    int   `json:"priority"`
	MaxWait     int   `json:"max_wait"`
	PickUpTime  int   `json:"pick_up_time"`
	CookingTime int   `json:"cooking_time"`
}

type OrdersList struct {
	sync.RWMutex
	AllOrders []Order
}

type ByPriorityAndMaxWait []Order

func (a ByPriorityAndMaxWait) Len() int {
	return len(a)
}
func (a ByPriorityAndMaxWait) Less(i, j int) bool {
	if a[i].Priority > a[j].Priority {
		return true
	}

	if a[i].Priority < a[j].Priority {
		return false
	}

	return a[i].MaxWait < a[j].MaxWait
}
func (a ByPriorityAndMaxWait) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
