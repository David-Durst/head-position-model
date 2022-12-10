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

func deg2Rad(angleInDegrees float64) float64 {
	return (angleInDegrees) * math.Pi / 180.0
}
func rad2Degree(angleInRadians float64) float64 {
	return (angleInRadians) * 180.0 / math.Pi
}

func angleVectors(angles r2.Point) r3.Vector {
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
	viewVec := angleVectors(viewAngles)
	viewVec.Z = 0.
	unitViewVec := viewVec.Normalize()
	neckDownAmount := duckAmount*DuckingNeckDownAmount + (1-duckAmount)*StandingNeckDownAmount
	headForwardAmount := duckAmount*DuckingHeadForwardAmount + (1-duckAmount)*StandingHeadForwardAmount
	return r3.Vector{
		X: eyePosition.X + math.Cos(deg2Rad(adjustedPitch))*unitViewVec.X*headForwardAmount,
		Y: eyePosition.Y + math.Cos(deg2Rad(adjustedPitch))*unitViewVec.Y*headForwardAmount,
		Z: eyePosition.Z - neckDownAmount + math.Sin(deg2Rad(adjustedPitch))*headForwardAmount,
	}
}

type AABB struct {
	minCorner r3.Vector
	maxCorner r3.Vector
}

const PlayerStandingHeight = 72
const PlayerCrouchedHeight = 54
const PlayerWidth = 32

func getAABBForPlayer(footPosition r3.Vector, duckAmount float64) AABB {
	//https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Mapper%27s_Reference
	// looks like coordinates are center at feet - tested using getpos_exact and box commands from
	//https://old.reddit.com/r/csmapmakers/comments/58ch3f/useful_console_commands_for_map_making_csgo/
	//making box with these coordinates wraps player perfectly
	// https://developer.valvesoftware.com/wiki/Dimensions#Eyelevel
	// eye level is 64 units when standing, 46 when crouching
	// these numbers are weird, use mappers reference not this
	// getpos is eye level, getpos_exact is foot, both are center of model
	height := duckAmount*PlayerCrouchedHeight + (1.-duckAmount)*PlayerStandingHeight
	return AABB{
		minCorner: r3.Vector{
			footPosition.X - PlayerWidth/2, footPosition.Y - PlayerStandingHeight/2, footPosition.Z,
		},
		maxCorner: r3.Vector{
			footPosition.X - PlayerWidth/2, footPosition.Y - PlayerStandingHeight/2, footPosition.Z + height,
		},
	}
}
