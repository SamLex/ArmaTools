package unitcapturereduce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
)

// ReduceUnitCapture reduces BIS_fnc_UnitCapture in a lossy manner using an error threshold
func ReduceUnitCapture(rawCaptureData string, errorPercentThreshold float64) (string, int, int, error) {
	captureData, err := parseCaptureData(rawCaptureData)
	if err != nil {
		return "", 0, 0, err
	}

	before := captureData.numberOfFrames()

	if before < 3 {
		return "", 0, 0, fmt.Errorf("There must be atleast 3 capture frames")
	}

	captureData.reduce(errorPercentThreshold)

	after := captureData.numberOfFrames()

	return captureData.String(), before, after, nil
}

type captureFrame struct {
	Time float64

	Position  vec3
	Direction vec3
	Up        vec3
	Velocity  vec3

	Before *captureFrame
	Next   *captureFrame
}

func parseCaptureData(rawCaptureData string) (*captureFrame, error) {
	var unsafeCaptureData [][]interface{}
	err := json.Unmarshal([]byte(rawCaptureData), &unsafeCaptureData)
	if err != nil {
		return nil, err
	}

	var next *captureFrame
	first := (*captureFrame)(nil)
	before := (*captureFrame)(nil)

	for i, unsafeFrame := range unsafeCaptureData {
		if i == 0 {
			first = &captureFrame{}
			next = first
		}

		next.Time = unsafeFrame[0].(float64)
		next.Position = unsafeSliceToVec3(unsafeFrame[1].([]interface{}))
		next.Direction = unsafeSliceToVec3(unsafeFrame[2].([]interface{}))
		next.Up = unsafeSliceToVec3(unsafeFrame[3].([]interface{}))
		next.Velocity = unsafeSliceToVec3(unsafeFrame[4].([]interface{}))
		next.Before = before

		before = next
		next = &captureFrame{}

		if i != len(unsafeCaptureData)-1 {
			before.Next = next
		} else {
			before.Next = nil
		}
	}

	return first, nil
}

func (captureData *captureFrame) reduce(errorPercentThreshold float64) {
	start := captureData
	consider := start.Next
	end := consider.Next

	for {
		lerped := lerpCaptureFrame(start, end, consider.Time)

		if consider.maxPercentDifference(lerped) < errorPercentThreshold {
			start.Next = end
			end.Before = start
		}

		if end.Next == nil {
			break
		}

		start = start.Next
		consider = start.Next
		end = consider.Next

		if end == nil {
			break
		}
	}
}

func (captureData *captureFrame) String() string {
	buf := &bytes.Buffer{}
	current := captureData

	fmt.Fprint(buf, "[")
	for current != nil {
		fmt.Fprintf(buf, "[%v,%s,%s,%s,%s]", current.Time, current.Position.String(), current.Direction.String(), current.Up.String(), current.Velocity.String())

		if current.Next != nil {
			fmt.Fprint(buf, ",")
		}

		current = current.Next
	}
	fmt.Fprint(buf, "]")

	return buf.String()
}

func (captureData *captureFrame) numberOfFrames() int {
	count := 0
	current := captureData

	for current != nil {
		count++
		current = current.Next
	}

	return count
}

func lerpCaptureFrame(start *captureFrame, end *captureFrame, time float64) *captureFrame {
	t := (time - start.Time) / end.Time

	return &captureFrame{
		Time:      time,
		Position:  start.Position.lerp(end.Position, t),
		Direction: start.Direction.lerp(end.Direction, t),
		Up:        start.Up.lerp(end.Up, t),
		Velocity:  start.Velocity.lerp(end.Velocity, t),
	}
}

func (captureData *captureFrame) maxPercentDifference(other *captureFrame) float64 {
	positionDifference := captureData.Position.maxPercentDifference(other.Position)
	directionDifference := captureData.Direction.maxPercentDifference(other.Direction)
	upDifference := captureData.Up.maxPercentDifference(other.Up)
	velocityDifference := captureData.Velocity.maxPercentDifference(other.Velocity)

	return math.Max(positionDifference, math.Max(directionDifference, math.Max(upDifference, velocityDifference)))
}

func unsafeSliceToVec3(slice []interface{}) vec3 {
	return vec3{
		A: slice[0].(float64),
		B: slice[1].(float64),
		C: slice[2].(float64),
	}
}
