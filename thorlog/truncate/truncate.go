package truncate

import (
	"bytes"
	"errors"
	"io"
	"sort"
	"strings"
)

const truncateSequence = `[...]`

type interval struct {
	from int
	to   int
}

// TruncateWithNewlines truncates the data smart by removing
// all lines that don't contain interesting strings from data.
// Interesting strings can be passed via mustContain argument.
// All lines containing mustContain arguments will be contained in data.
// The resulting string might still be larger than l.truncate! But this
// won't be a problem because the logger will truncate each value again.
func TruncateWithNewlines(data string, mustContain []Match, truncateLimit int, context int) string {
	if len(mustContain) == 0 {
		return strings.Replace(data, "\n", truncateSequence, -1)
	}
	var intervals []interval
	for _, str := range mustContain {
		i := int(str.Offset)
		previousNewlineIndex := strings.LastIndex(data[:i], "\n")
		nextNewlineIndex := strings.Index(data[i+len(str.Data):], "\n")
		if nextNewlineIndex < 0 {
			nextNewlineIndex = len(data)
		} else {
			nextNewlineIndex = nextNewlineIndex + i + len(str.Data)
		}
		intervals = append(intervals, interval{
			previousNewlineIndex + 1, nextNewlineIndex,
		})
	}
	intervals = reduceIntervals(intervals)
	var intervalStrings []string
	for _, interval := range intervals {
		if interval.to > interval.from {
			// Truncate the interval data to reduce it to the interesting part
			var offsetYaraStrings []Match
			for _, str := range mustContain {
				i := int(str.Offset)
				if i >= interval.from && i+len(str.Data) <= interval.to {
					// Match string occurs in this interval, include it in the truncation data
					offsetYaraStrings = append(offsetYaraStrings, Match{
						Offset: uint64(i - interval.from),
						Data:   str.Data,
					})
				}
			}
			truncatedIntervalData := SmartTruncate(data[interval.from:interval.to], offsetYaraStrings, truncateLimit, context)
			intervalStrings = append(intervalStrings, truncatedIntervalData)
		}
	}
	totalString := strings.Join(intervalStrings, truncateSequence)
	return strings.Replace(totalString, "\n", truncateSequence, -1)
}

// reduceIntervals minimizes the list of intervals by combining overlapping intervals.
// The resulting list is guaranteed to:
//   - be sorted
//   - have no overlapping intervals
//   - contain the same indices as the given intervals.
func reduceIntervals(intervals []interval) []interval {
	// sort intervals by their start position.
	// if the start position is the same, then sort by their end position.
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].from == intervals[j].from {
			return intervals[i].to < intervals[j].to
		}
		return intervals[i].from < intervals[j].from
	})
	var reduced []interval
	for _, nextInterval := range intervals {
		if nextInterval.to <= nextInterval.from { // Empty interval
			continue
		}
		if len(reduced) == 0 {
			reduced = append(reduced, nextInterval)
			continue
		}
		lastReduced := &reduced[len(reduced)-1]
		if nextInterval.from <= lastReduced.to { // Interval overlaps with existing interval
			if nextInterval.to > lastReduced.to { // Interval only overlaps partially
				lastReduced.to = nextInterval.to
			}
			continue
		}
		// Intervals does not overlap and can be added
		reduced = append(reduced, nextInterval)
	}
	return reduced
}

// SmartTruncate truncates the data smart by removing
// uninteresting strings from data. Interesting strings can be
// passed via mustContain argument. All mustContain arguments
// will be contained in data including 10 chars before and after each
// mustContain element.
// If the resulting string might still be larger than truncateLimit, some
// mustContain strings will be removed, starting with the largest ones.
func SmartTruncate(data string, mustContain []Match, truncateLimit int, context int) string {
	// Reduce intervals to an amount we can actually support by removing the largest ones
	sort.Slice(mustContain, func(i, j int) bool {
		return len(mustContain[i].Data) < len(mustContain[j].Data)
	})
	for {
		truncatedString := truncateContainingStrings(data, mustContain, truncateLimit, context)
		if truncateLimit <= 0 || len(truncatedString) <= truncateLimit || len(mustContain) == 0 {
			return truncatedString
		}

		// Reduce context, if possible, to get within range
		if context > 0 && len(truncatedString)-len(mustContain)*context*2 <= truncateLimit {
			contextReduction := (len(truncatedString) - truncateLimit) / (len(mustContain) * 2)
			if contextReduction <= 0 {
				contextReduction = 1
			}
			context -= contextReduction
			continue
		}
		mustContain = mustContain[:len(mustContain)-1]
	}
}

func truncateContainingStrings(data string, mustContain []Match, truncateLimit int, context int) string {
	var dataLen = len(data)
	if dataLen <= truncateLimit || truncateLimit <= 0 {
		return data
	}
	var res string
	var intervals []interval
	// get all intervals of the strings
	// that has to be in the resulting string
	for _, str := range mustContain {
		i := int(str.Offset)
		interval := interval{
			from: i - context,                 // the interval starts where the string was found
			to:   i + len(str.Data) + context, // the interval ends where the interval starts + len of string
		}
		if interval.from < 0 {
			interval.from = 0 // fix boundary
		}
		if interval.to > dataLen {
			interval.to = dataLen // fix boundary
		}
		intervals = append(intervals, interval)
	}
	intervals = reduceIntervals(intervals)

	// Add start and end intervals
	remainingLength := truncateLimit - truncatedLength(intervals, dataLen)
	if remainingLength < 0 {
		remainingLength = 0
	}
	intervals = append(intervals, interval{from: 0, to: remainingLength / 2}, interval{from: dataLen - remainingLength/2, to: dataLen})
	intervals = reduceIntervals(intervals)

	if len(intervals) == 0 {
		return truncateSequence
	}

	// maybe we truncated too much .. try to increase some intervals
	remainingLength = truncateLimit - truncatedLength(intervals, dataLen)
	for i := 1; i < len(intervals) && remainingLength > 0; i++ {
		if diff := intervals[i].from - intervals[i-1].to; diff > 0 {
			if remainingLength >= diff {
				intervals[i-1].to += diff
				remainingLength -= diff
			} else if remainingLength > 0 && remainingLength < diff {
				intervals[i-1].to += remainingLength
				remainingLength = 0
			}
		}
	}
	if remainingLength > 0 {
		diff := dataLen - intervals[len(intervals)-1].to
		if diff > 0 && diff <= remainingLength {
			intervals[len(intervals)-1].to += diff
		} else if diff > 0 && diff > remainingLength {
			intervals[len(intervals)-1].to += remainingLength
		}
	}
	intervals = reduceIntervals(intervals)
	if len(intervals) == 0 {
		return truncateSequence
	}

	// now iterate over the intervals and build the result
	if intervals[0].from > 0 {
		res += truncateSequence
	}
	res += data[intervals[0].from:intervals[0].to]
	for i := 1; i < len(intervals); i++ {
		if intervals[i].from > intervals[i-1].to {
			res += truncateSequence
		}
		res += data[intervals[i].from:intervals[i].to]
	}
	if intervals[len(intervals)-1].to < dataLen {
		res += truncateSequence
	}
	return res
}

func truncatedLength(intervals []interval, dataLen int) int {
	// now check where we have to start truncating our data.
	var length = 0
	for i := 0; i < len(intervals); i++ {
		length += intervals[i].to - intervals[i].from
		// if there is a gap between two intervals, we add a truncate sequence
		// so we have to reduce our start sequence by the length of the truncate sequence
		if i != 0 && intervals[i].from > intervals[i-1].to {
			length += len(truncateSequence)
		}
	}
	if len(intervals) > 0 {
		if intervals[0].from > 0 && intervals[0].from > length/2 {
			// if the first interval starts after our start sequence,
			// we add another truncate sequence
			length += len(truncateSequence)
		}
		if intervals[len(intervals)-1].to < dataLen && intervals[len(intervals)-1].to < dataLen-length/2 {
			// if the last interval ends before our end sequence,
			// we add another truncate sequence
			length += len(truncateSequence)
		}
	} else {
		// We need at least one truncate sequence between the initial and then end sequence
		length += len(truncateSequence)
	}
	return length
}

type Match struct {
	Offset uint64
	Data   []byte
}

func GetStringMatchContext(context uint64, rawData io.ReaderAt, match Match) (Match, error) {
	var startIndex int64
	var preDataContext uint64
	if context < match.Offset { // Offset is large enough to read whole context at the start
		startIndex = int64(match.Offset - context)
		preDataContext = context
	} else { // Reduce bytes to be read accordingly since the context at the start is truncated
		preDataContext = match.Offset
	}
	data := make([]byte, preDataContext+uint64(len(match.Data))+context)

	readBytes, err := rawData.ReadAt(data, startIndex)
	// Check: Read successful and match data appears as expected
	if err != nil && err != io.EOF {
		// No context available
		return Match{}, err
	}
	if readBytes < int(preDataContext)+len(match.Data) {
		return Match{}, errors.New("could not read sufficient context data")
	}
	if readBytes == len(match.Data) {
		return Match{}, errors.New("no context available")
	}
	if !bytes.Equal(data[preDataContext:preDataContext+uint64(len(match.Data))], match.Data) {
		return Match{}, errors.New("match data does not appear as expected")
	}
	return Match{
		Data:   data[:readBytes],
		Offset: uint64(startIndex),
	}, nil
}
