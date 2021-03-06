/*
Copyright 2016 Euan James Hunter

unitcapturereduce.go: Main file for armatools/unitcapturereduce

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package unitcapturereduce

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
)

/*
Reduction algorithm:
	Assume UnitCapture outputs keyframes that are linearly interpolated between by UnitPlay (an approximation of actual behaviour)
	Goal: Remove keyframes that do not greatly effect the followed 'path' ('path' here being the path through 12 demension space that the keyframes represent)

	1. Keep the first and last keyframe
	2. For each keyframe between the first and last consider removing it;
		2.1 Linearly interpolate between the one before and the one after using the time of the one under consideration
		2.2 Construct a normal distribution with the current keyframe as mean and a standard deviation of 1
		2.3 Calculate the probability of the linearly interpolated keyframe using this distribution
		2.4 If the probability of this keyframe is above the probability threshold, the current key frame can be removed
	3. Continue until all keyframes have been considered
*/

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

// Parse the SQF array capture data into a List of captureKeyFrame structs
func parseCaptureData(rawCaptureData string) (*captureKeyFrames, error) {
	// A SQF array is actually valid JSON, so parse as json
	// No typesafe way to represent inner slice though
	var unsafeCaptureData [][]interface{}
	err := json.Unmarshal([]byte(rawCaptureData), &unsafeCaptureData)
	if err != nil {
		return nil, err
	}

	data := list.New()

	// Unpack each parsed keyframe into a captureKeyFrame
	for i, unsafeKeyFrame := range unsafeCaptureData {
		if len(unsafeKeyFrame) != 5 {
			return nil, fmt.Errorf("Invalid UnitCapture Output")
		}

		// Typesafe assertions to the correct type

		time, timeOK := unsafeKeyFrame[0].(float64)

		unsafePosition, unsafePositionOK := unsafeKeyFrame[1].([]interface{})
		unsafeDirection, unsafeDirectionOK := unsafeKeyFrame[2].([]interface{})
		unsafeUp, unsafeUpOK := unsafeKeyFrame[3].([]interface{})
		unsafeVelocity, unsafeVelocityOK := unsafeKeyFrame[4].([]interface{})

		position := unsafeSliceToVec3(unsafePosition)
		direction := unsafeSliceToVec3(unsafeDirection)
		up := unsafeSliceToVec3(unsafeUp)
		velocity := unsafeSliceToVec3(unsafeVelocity)

		// Check everything asserted properly
		if !timeOK ||
			!unsafePositionOK || !unsafeDirectionOK || !unsafeUpOK || !unsafeVelocityOK ||
			position == nil || direction == nil || up == nil || velocity == nil {

			return nil, fmt.Errorf("Invalid UnitCapture Output")
		}

		// Pack into struct and add to list
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

// Unpack an unsafe slice of 3 elemenrs into a vec3
// Returns nil if the slice does not have 3 elements or are not all float64s
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

// Reduce the capture data (see algorithm above)
func (ckf *captureKeyFrames) reduce(probabilityThreshold float64) {
	startElm := (*list.List)(ckf).Front()
	considerElm := startElm.Next()
	endElm := considerElm.Next()

	start := startElm.Value.(*captureKeyFrame)
	consider := considerElm.Value.(*captureKeyFrame)
	end := endElm.Value.(*captureKeyFrame)

	for {
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
	}
}

// Convert the keyframes back to SQF array format
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

// Linearly interpolate between two keyframes
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

// Convert a keyframe into a slice with 12 elements
// Helpful for treating a keyframe as one big vector
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

// Convert a keyframe to SQF array format
func (ckf *captureKeyFrame) SQFString() string {
	return fmt.Sprintf("[%v,%s,%s,%s,%s]", ckf.Time, ckf.Position.SQFString(), ckf.Direction.SQFString(), ckf.Up.SQFString(), ckf.Velocity.SQFString())
}
