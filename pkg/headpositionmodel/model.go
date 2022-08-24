package headpositionmodel

import (
	"github.com/golang/geo/r2"
	"github.com/golang/geo/r3"
	"math"
)

const StandingHeadForwardAmount = 12.5
const DuckingHeadForwardAmount = 11.
const StandingNeckDownAmount = 8.5
const DuckingNeckDownAmount = 3.75
const StandingHeadAngleAdjustmentAdd = 2.5
const CrouchingHeadAngleAdjustmentAdd = 17.5
const StandingHeadAngleAdjustmentMul = 1.1
const CrouchingHeadAngleAdjustmentMul = 0.75

func Deg2Rad(angleInDegrees float64) float64 {
	return (angleInDegrees) * math.Pi / 180.0
}
func Rad2Degree(angleInRadians float64) float64 {
	return (angleInRadians) * 180.0 / math.Pi
}

func AngleVectors(angles r2.Point) r3.Vector {
	// https://github.com/ValveSoftware/source-sdk-2013/blob/master/sp/src/mathlib/mathlib_base.cpp#L901-L914
	// https://developer.valvesoftware.com/wiki/QAngle - QAngle is just a regular Euler angle
	var forward r3.Vector

	sy, cy := math.Sincos(angles.X)
	sp, cp := math.Sincos(angles.Y)

	forward.X = cp * cy
	forward.Y = cp * sy
	forward.Z = -sp
	return forward
}

func ModelHeadPosition(eyePosition r3.Vector, viewAngles r2.Point, duckAmount float64) r3.Vector {
	// no z factor if standing and pitch is 90 (looking down), all z factor and no x/y factor if pitch is -90 (looking up)
	// scale looking down is 0 and up is 90, perfect for sin/cos function where head makes quarter circle
	// also adjust by head angle amount since head is flat when looking down (0 deg after transformation) but back a little
	// when looking up (90 deg after transformation)
	headAngleAdjustmentAdd := duckAmount*CrouchingHeadAngleAdjustmentAdd +
		(1-duckAmount)*StandingHeadAngleAdjustmentAdd
	headAngleAdjustmentMul := duckAmount*CrouchingHeadAngleAdjustmentMul +
		(1-duckAmount)*StandingHeadAngleAdjustmentMul
	adjustedPitch := (viewAngles.Y*-1.+90.)/2.*headAngleAdjustmentMul + headAngleAdjustmentAdd
	// get unit vec of just x and y (z already handled)
	viewVec := AngleVectors(viewAngles)
	viewVec.Z = 0.
	unitViewVec := viewVec.Normalize()
	neckDownAmount := duckAmount*DuckingNeckDownAmount + (1-duckAmount)*StandingNeckDownAmount
	headForwardAmount := duckAmount*DuckingHeadForwardAmount + (1-duckAmount)*StandingHeadForwardAmount
	return r3.Vector{
		X: eyePosition.X + math.Cos(Deg2Rad(adjustedPitch))*unitViewVec.X*headForwardAmount,
		Y: eyePosition.Y + math.Cos(Deg2Rad(adjustedPitch))*unitViewVec.Y*headForwardAmount,
		Z: eyePosition.Z - neckDownAmount + math.Sin(Deg2Rad(adjustedPitch))*headForwardAmount,
	}
}
