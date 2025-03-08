package cache

type Record struct {
	sum   float64
	count uint64
}

func NewRecord(sum float64, count uint64) Record {
	return Record{sum, count}
}

func (record *Record) IncrSum(value float64) float64 {
	record.sum += value
	return record.sum
}

func (record *Record) IncrCount(value uint64) uint64 {
	record.count += value
	return record.count
}
