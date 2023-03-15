package cuckoofilter

type Hash interface {
	Sum() int
}
