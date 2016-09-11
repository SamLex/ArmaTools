package unitcapturereduce

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
)

// ReduceUnitCapture reduces BIS_fnc_UnitCapture in a lossy manner using an error threshold
func ReduceUnitCapture(rawCaptureData string, probabilityThreshold float64) (string, int, int, error) {
	captureData, err := parseCaptureData(rawCaptureData)
	if err != nil {
		return "", 0, 0, err
	}

	before := captureData.numberOfFrames()
	if before < 3 {
		return "", 0, 0, fmt.Errorf("There must be atleast 3 capture frames")
	}

	captureData.reduce(probabilityThreshold)

	after := captureData.numberOfFrames()

	return captureData.SQFString(), before, after, nil
}

func parseCaptureData(rawCaptureData string) (*captureKeyFrames, error) {
	var unsafeCaptureData [][]interface{}
	err := json.Unmarshal([]byte(rawCaptureData), &unsafeCaptureData)
	if err != nil {
		return nil, err
	}

	data := list.New()

	for i, unsafeKeyFrame := range unsafeCaptureData {
		if len(unsafeKeyFrame) != 5 {
			return nil, fmt.Errorf("Invalid UnitCapture Output")
		}

		time, timeOK := unsafeKeyFrame[0].(float64)

		unsafePosition, unsafePositionOK := unsafeKeyFrame[1].([]interface{})
		unsafeDirection, unsafeDirectionOK := unsafeKeyFrame[2].([]interface{})
		unsafeUp, unsafeUpOK := unsafeKeyFrame[3].([]interface{})
		unsafeVelocity, unsafeVelocityOK := unsafeKeyFrame[4].([]interface{})

		position := unsafeSliceToVec3(unsafePosition)
		direction := unsafeSliceToVec3(unsafeDirection)
		up := unsafeSliceToVec3(unsafeUp)
		velocity := unsafeSliceToVec3(unsafeVelocity)

		if !timeOK ||
			!unsafePositionOK || !unsafeDirectionOK || !unsafeUpOK || !unsafeVelocityOK ||
			position == nil || direction == nil || up == nil || velocity == nil {

			return nil, fmt.Errorf("Invalid UnitCapture Output")
		}

		keyframe := &captureKeyFrame{
			OriginalFrameNumber: i,
			Time:                time,
			Position:            *position,
			Direction:           *direction,
			Up:                  *up,
			Velocity:            *velocity,
		}

		data.PushBack(keyframe)
	}

	return (*captureKeyFrames)(data), nil
}

func unsafeSliceToVec3(slice []interface{}) *vec3 {
	if len(slice) != 3 {
		return nil
	}

	a, aOK := slice[0].(float64)
	b, bOK := slice[1].(float64)
	c, cOK := slice[2].(float64)

	if !aOK || !bOK || !cOK {
		return nil
	}

	return &vec3{
		A: a,
		B: b,
		C: c,
	}
}

type captureKeyFrames list.List

func (ckf *captureKeyFrames) numberOfFrames() int {
	return (*list.List)(ckf).Len()
}

func (ckf *captureKeyFrames) reduce(probabilityThreshold float64) {
	count := 0

	startElm := (*list.List)(ckf).Front()
	considerElm := startElm.Next()
	endElm := considerElm.Next()

	start := startElm.Value.(*captureKeyFrame)
	consider := considerElm.Value.(*captureKeyFrame)
	end := endElm.Value.(*captureKeyFrame)

	for {
		if start.OriginalFrameNumber != count {
			panic("A frame has been skipped")
		}

		newFrame := start.lerp(end, consider.Time)

		if vectorisedNormalDistributionPDF(newFrame.toSlice(), consider.toSlice(), 1) > probabilityThreshold {
			(*list.List)(ckf).Remove(considerElm)
		}

		if endElm.Next() == nil {
			break
		}

		startElm = startElm.Next()
		considerElm = startElm.Next()
		endElm = considerElm.Next()

		if endElm == nil {
			break
		}

		start = startElm.Value.(*captureKeyFrame)
		consider = considerElm.Value.(*captureKeyFrame)
		end = endElm.Value.(*captureKeyFrame)

		count++
	}
}

func (ckf *captureKeyFrames) SQFString() string {
	buf := &bytes.Buffer{}

	fmt.Fprint(buf, "[")
	for elm := (*list.List)(ckf).Front(); elm != nil; elm = elm.Next() {
		fmt.Fprint(buf, elm.Value.(*captureKeyFrame).SQFString())

		if elm.Next() != nil {
			fmt.Fprint(buf, ",")
		}
	}
	fmt.Fprint(buf, "]")

	return buf.String()
}

type captureKeyFrame struct {
	OriginalFrameNumber int
	Time                float64
	Position            vec3
	Direction           vec3
	Up                  vec3
	Velocity            vec3
}

func (ckf *captureKeyFrame) lerp(end *captureKeyFrame, time float64) *captureKeyFrame {
	t := (time - ckf.Time) / end.Time

	return &captureKeyFrame{
		OriginalFrameNumber: -1,
		Time:                time,
		Position:            ckf.Position.lerp(end.Position, t),
		Direction:           ckf.Direction.lerp(end.Direction, t),
		Up:                  ckf.Up.lerp(end.Up, t),
		Velocity:            ckf.Velocity.lerp(end.Velocity, t),
	}
}

func (ckf *captureKeyFrame) toSlice() []float64 {
	slice := make([]float64, 12)

	slice[0] = ckf.Position.A
	slice[1] = ckf.Position.B
	slice[2] = ckf.Position.C

	slice[3] = ckf.Direction.A
	slice[4] = ckf.Direction.B
	slice[5] = ckf.Direction.C

	slice[6] = ckf.Up.A
	slice[7] = ckf.Up.B
	slice[8] = ckf.Up.C

	slice[9] = ckf.Velocity.A
	slice[10] = ckf.Velocity.B
	slice[11] = ckf.Velocity.C

	return slice
}

func (ckf *captureKeyFrame) SQFString() string {
	return fmt.Sprintf("[%v,%s,%s,%s,%s]", ckf.Time, ckf.Position.SQFString(), ckf.Direction.SQFString(), ckf.Up.SQFString(), ckf.Velocity.SQFString())
}
