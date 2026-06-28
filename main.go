package main

import (
	"fmt"
	"sync"
	"time"
)

var keyStore map[string]string

type Transaction struct {
	ID        string
	Amount    float64
	Timestamp time.Time
}

type TransactionLookupCache struct {
	mu           sync.RWMutex
	transactions map[string]Transaction
}

func init() {
	println("Initializing the key store...")
	keyStore = make(map[string]string)
}

func NewTransactionLookupCache() *TransactionLookupCache {
	return &TransactionLookupCache{
		transactions: make(map[string]Transaction),
	}
}

func addTransaction(cache *TransactionLookupCache, transaction Transaction) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.transactions[transaction.ID] = transaction
	println("Added transaction:", transaction.ID)
}

func getTransaction(cache *TransactionLookupCache, id string) (Transaction, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	if transaction, exists := cache.transactions[id]; exists {
		println("Retrieved transaction:", id, "with amount:", transaction.Amount)
		return transaction, true
	}
	println("Transaction not found:", id)
	return Transaction{}, false
}

func createSampleTransactions(cache *TransactionLookupCache, workerID int) {

	transaction := Transaction{
		ID:        fmt.Sprintf("txn%d", workerID),
		Amount:    float64(workerID) * 100.0,
		Timestamp: time.Now(),
	}
	addTransaction(cache, transaction)

}

func getSampleTransaction(cache *TransactionLookupCache, workerID int) {

	time.Sleep(5 * time.Millisecond)
	id := fmt.Sprintf("txn%d", workerID)
	if transaction, found := getTransaction(cache, id); found {
		println("Sample transaction found:", transaction.ID, "with amount:", transaction.Amount)
	} else {
		println("Sample transaction not found:", id)
	}

}

func addKey(key, value string) {
	keyStore[key] = value
	println("Added key:", key)
}

func getKey(key string) string {
	value, exists := keyStore[key]
	if !exists {
		println("Key not found:", key)
		return ""
	}
	println("Retrieved key:", key, "with value:", value)
	return value
}

func time_to_live() int {
	return 60 // seconds
}

func remove_expired_transactions(cache *TransactionLookupCache) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	for id, transaction := range cache.transactions {
		if time.Since(transaction.Timestamp) > time.Duration(time_to_live())*time.Second {
			delete(cache.transactions, id)
			println("Removed expired transaction:", id)
		}
	}
}

func main() {
	println("Lookup cache test started...")

	var wg sync.WaitGroup
	cache := NewTransactionLookupCache()
	for i := range 2000 {
		wg.Go(func() {
			createSampleTransactions(cache, i)
		})
	}
	for i := range 200 {
		wg.Go(func() {
			getSampleTransaction(cache, i)
		})
	}
	wg.Go(func() {
		for {
			if len(cache.transactions) > 0 {
				println("Current number of transactions in cache:", len(cache.transactions))
				remove_expired_transactions(cache)
			} else {
				println("No transactions in cache.")
				break
			}
		}
	})

	wg.Wait()

	println("Final number of transactions in cache:", len(cache.transactions))
}
