package cuckoofilter

type Hash interface {
	Sum() uint
}
