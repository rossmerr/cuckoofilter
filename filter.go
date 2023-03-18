package cuckoofilter

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
)

type HashFunction func([]byte) uint

type Filter[T Hash] struct {
	buckets           [][]byte
	fingerprintLength uint
	length            uint
	hash              HashFunction
	bytesPerBucket    uint
	maxNumKicks       uint
	entriesPerBucket  uint
}

type FilterOption[T Hash] func(*Filter[T])

func WithHash[T Hash](hash HashFunction) FilterOption[T] {
	return func(s *Filter[T]) {
		s.hash = hash
	}

}

// New Cuckoo filter
func NewFilter[T Hash](length, entriesPerBucket, fingerprintLength, maxNumKicks uint, opts ...FilterOption[T]) *Filter[T] {

	buckets := make([][]byte, length)
	for i := uint(0); i < length; i++ {
		buckets[i] = make([]byte, entriesPerBucket)
	}

	bytesPerBucket := entriesPerBucket * fingerprintLength

	var hasher = sha1.New()
	hash := func(data []byte) uint {
		hasher.Write([]byte(data))
		hash := hasher.Sum(nil)
		hasher.Reset()
		return uint(binary.LittleEndian.Uint64(hash))
	}

	filter := &Filter[T]{
		buckets:           buckets,
		fingerprintLength: fingerprintLength,
		length:            length,
		hash:              hash,
		maxNumKicks:       maxNumKicks,
		bytesPerBucket:    bytesPerBucket,
		entriesPerBucket:  entriesPerBucket,
	}

	for _, opt := range opts {
		opt(filter)
	}

	return filter
}

// n The capacity
// fp The false positive rate
func NewFilterFalsePositiveRate[T Hash](n uint, fp float64, opts ...FilterOption[T]) *Filter[T] {
	// For two Buckets and each bucket has up to four fingerprints
	entriesPerBucket := uint(4)
	fingerprintLength := fingerprintLength(entriesPerBucket, fp)
	length := upperPower2(n / fingerprintLength * 8)
	if float64(n/length/entriesPerBucket) > 0.96 {
		length <<= 1
	}

	buckets := make([][]byte, int(length))
	for i := uint(0); i < length; i++ {
		buckets[i] = make([]byte, entriesPerBucket)
	}

	return NewFilter(length, entriesPerBucket, fingerprintLength, length, opts...)
}

func (s *Filter[T]) Add(item T) error {
	key := item.Sum()
	f := s.fingerprint(key)
	i1 := s.index(f, key)
	i2 := s.swapIndex(f, i1)

	b1 := s.buckets[i1%s.length]
	for i, v := range b1 {
		if v == 0 {
			copy(b1[i:], f)
			return nil
		}
	}
	b2 := s.buckets[i2%s.length]
	for i, v := range b2 {
		if v == 0 {
			copy(b2[i:], f)
			return nil
		}
	}

	//randomly pick i1 or i2
	i := uint(rand.Intn(2))
	if i == 0 {
		i = i1
	} else {
		i = i2
	}

	for n := uint(0); n < s.maxNumKicks; n++ {
		// randomly select an entry e from bucket;
		e := uint(rand.Intn(int(s.entriesPerBucket)))
		f = s.buckets[i%s.length][e : e+s.fingerprintLength]
		swap := s.swapIndex(f, i)
		for r, v := range s.buckets[swap%s.length] {
			// swap entry to alternate position
			if v == 0 {
				copy(s.buckets[swap%s.length][r:], f)
				copy(s.buckets[i%s.length][e:], make([]byte, s.fingerprintLength))
				return nil
			}
		}
		i = swap
	}

	return fmt.Errorf("cuckoo filter full")
}

func (s *Filter[T]) Contains(item T) bool {
	key := item.Sum()
	f := s.fingerprint(key)
	i1 := s.index(f, key)
	i2 := s.swapIndex(f, i1)

	b1 := s.buckets[i1%s.length]
	for i := uint(0); i < uint(len(b1)); i += uint(s.fingerprintLength) {
		if bytes.Equal(b1[i:i+s.fingerprintLength], f) {
			return true
		}
	}

	b2 := s.buckets[i2%s.length]
	for i := uint(0); i < uint(len(b2)); i += uint(s.fingerprintLength) {
		if bytes.Equal(b2[i:i+s.fingerprintLength], f) {
			return true
		}
	}

	return false
}

func (s *Filter[T]) Remove(item T) {
	key := item.Sum()
	f := s.fingerprint(key)
	i1 := s.index(f, key)
	i2 := s.swapIndex(f, i1)

	b1 := s.buckets[i1%s.length]
	for i := uint(0); i < uint(len(b1)); i += uint(s.fingerprintLength) {
		if bytes.Equal(b1[i:i+s.fingerprintLength], f) {
			copy(b1[i:i+s.fingerprintLength], make([]byte, s.fingerprintLength))
		}
	}

	b2 := s.buckets[i2%s.length]
	for i := uint(0); i < uint(len(b2)); i += uint(s.fingerprintLength) {
		if bytes.Equal(b2[i:i+s.fingerprintLength], f) {
			copy(b2[i:i+s.fingerprintLength], make([]byte, s.fingerprintLength))
		}
	}
}

func (s *Filter[T]) Length() int {
	return int(s.length)
}

func (s *Filter[T]) swapIndex(f []byte, index uint) uint {
	return index ^ (s.hash(f) & (s.length - 1))
}

func (s *Filter[T]) index(f []byte, key uint) uint {
	return key ^ (s.hash(f) & (s.length - 1))
}

func upperPower2(x uint) uint {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32
	x++
	return x
}

func (s *Filter[T]) fingerprint(value uint) []byte {
	arr := intToBytes(value)
	return arr[0:s.fingerprintLength]
}

func fingerprintLength(b uint, fp float64) uint {
	f := uint(math.Ceil((math.Log(2 * (float64(b) / fp)))))
	f /= 8
	if f < 1 {
		return 1
	}
	return f
}

func intToBytes(num uint) []byte {
	buff := new(bytes.Buffer)
	bigOrLittleEndian := binary.LittleEndian
	err := binary.Write(buff, bigOrLittleEndian, uint64(num))
	if err != nil {
		panic(err)
	}

	return buff.Bytes()
}
