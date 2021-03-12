package timedata

import (
	"encoding/binary"
	"errors"
	"time"
)

const (
	DAY_HOURS = 24
	DAY_HOUR_MINUTES = 60

	// less than DAY_HOUR_MINUTES
	SEGMENTS_PER_HOUR = 6
	SEGMENTS_PREPARE = 1000
)

var (

	VALUE_EXPIRATION = 30 * time.Minute

	// errors
	ErrTooLongDuration = errors.New("Expiration should be in a day.")
	ErrValueExpired = errors.New("Value had expired.")
	ErrOutOfIndex = errors.New("Out of index of segement.")
)

type (
	TimeManager struct {
		perHour int
		prepare int
		expire time.Duration
		segmets [][]*Segment
	}
	Segment struct {
		length int
		expiredIndex int
		nodes []*valuenode
	}
	valuenode struct {
		Value interface{}
		expiration int64
	}
)

func NewTimeManager(d time.Duration) *TimeManager {
	return NewTimeManagerWithOption(SEGMENTS_PER_HOUR, SEGMENTS_PREPARE, d)
}

func NewTimeManagerWithOption(segNum, segPrepare int, d time.Duration) *TimeManager{
	if segNum > DAY_HOUR_MINUTES {
		segNum = DAY_HOUR_MINUTES
	}
	if segNum < SEGMENTS_PER_HOUR {
		segNum = SEGMENTS_PER_HOUR
	}
	if segPrepare < SEGMENTS_PREPARE {
		segPrepare = SEGMENTS_PREPARE
	}
	return &TimeManager{
		perHour: segNum,
		prepare: segPrepare,
		expire: d,
		segmets: make([][]*Segment, DAY_HOURS),
	}
}

func (t *TimeManager) Add(value interface{}) []byte {
	expiring := time.Now().Add(t.expire)
	h, m, _ := expiring.Clock()
	index := m * t.perHour / DAY_HOUR_MINUTES
	// TODO: optimize space alloc
	if t.segmets[h] == nil {
		t.segmets[h] = make([]*Segment, t.perHour)
		for i:=0;i<t.perHour;i++ {
			t.segmets[h][i] = NewSegmentWithSize(t.prepare)
		}
	}
	segIndex := t.segmets[h][index].Add(value, expiring.Unix())
	if segIndex == -1 {
		return nil
	}
	return indexToByte(h, index, segIndex)
}

func (t *TimeManager) GetValue(data []byte) (interface{}, error) {
	h, seg, index := byteToIndex(data)
	checkIndex(t, h, seg, index)
	return t.segmets[h][seg].Get(index)
}

func checkIndex(t *TimeManager, h, seg, index int) bool {
	if h < 0 || h >= 24 {
		return false
	}
	if seg < 0 || seg >= t.perHour {
		return false
	}
	if index<0 {
		return false
	}
	return true
}

func indexToByte(h, seg, index int) []byte {
	res := make([]byte, 6)
	res[0] = byte(h)
	res[1] = byte(seg)
	binary.BigEndian.PutUint32(res[2:], uint32(index))
	return res
}

func byteToIndex(data []byte) (h, seg, index int) {
	if len(data) != 6 {
		return -1, -1, -1
	}
	h = int(data[0])
	seg = int(data[1])
	index = int(binary.BigEndian.Uint32(data[2:]))
	return
}

// TODO: timer to purge segments

func NewSegment() *Segment {
	return NewSegmentWithSize(SEGMENTS_PREPARE)
}

func NewSegmentWithSize(size int) *Segment {
	return &Segment{
		length: 0,
		expiredIndex: -1,
		nodes: make([]*valuenode, 0, size),
	}
}

func (s *Segment) Add(value interface{}, expiration int64) int {
	if expiration < time.Now().Unix() {
		return -1
	}
	s.nodes = append(s.nodes, &valuenode{value, expiration})
	index := s.length
	s.length++
	return index
}

func (s *Segment) Get(index int) (interface{}, error) {
	if index < s.expiredIndex {
		return nil, ErrValueExpired
	}
	if index >= s.length {
		return nil, ErrOutOfIndex
	}
	if s.nodes[index].expiration < time.Now().Unix() {
		s.expiredIndex = index
		return nil, ErrValueExpired
	}
	return s.nodes[index].Value, nil
}

