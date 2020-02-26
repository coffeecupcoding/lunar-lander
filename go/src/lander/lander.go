// Package lander models the lunar lander's kinematics and dynamics
package lander

import (
	"math"
)

// Specific Impulse of fuel, in miles/second
// These units mean fuel is measured by mass, not weight
const fuelIsp = 1.8

// Gravity is the lunar gravity in miles/(second^2)
const Gravity = 0.001

// Kinematics models the motion of the lunar lander
// It also tracks elapsed time since that is directly tied to position and
// velocity changes
type Kinematics struct {
	Velocity    float64
	Altitude    float64
	ElapsedTime float64
}

// velocityMassFactor implements part of a series solution for lander
// velocity change based on thrust and mass change due to fuel usage
func velocityMassFactor(massChange float64) float64 {
	return ((-1.0 * massChange) +
		(-1.0 * (math.Pow(massChange, 2.0) / 2.0)) +
		(-1.0 * (math.Pow(massChange, 3.0) / 3.0)) +
		(-1.0 * (math.Pow(massChange, 4.0) / 4.0)) +
		(-1.0 * (math.Pow(massChange, 5.0) / 5.0)))
}

// altitudeMassFactor implements part of a series solution for lander
// altitude change based on thrust and mass change due to fuel usage
func altitudeMassFactor(massChange float64) float64 {
	return ((massChange / 2.0) +
		(math.Pow(massChange, 2.0) / 6.0) +
		(math.Pow(massChange, 3.0) / 12.0) +
		(math.Pow(massChange, 4.0) / 20.0) +
		(math.Pow(massChange, 5.0) / 30.0))
}

// Lander models the lunar lander via parameters and associated functions
type Lander struct {
	CapsuleMass float64
	Fuel        float64
	TotalMass   float64
	Kinematics
}

// OutOfFuel tests whether the lander fuel is 'close' to zero
func (l *Lander) OutOfFuel() bool {
	return l.Fuel < 0.001
}

// calcMassChange is an internal function that returns the percentage
// mass change due to a period of fuel usage
func (l *Lander) calcMassChange(burnRate float64, burnTime float64) float64 {
	return (burnRate * burnTime) / l.TotalMass
}

// calcVelocity is an internal function that returns a new velocity due to
// a period of fuel usage along with the effect of gravity
func (l *Lander) calcVelocity(burnTime float64, massChange float64) float64 {
	return l.Velocity + (Gravity * burnTime) +
		(fuelIsp * velocityMassFactor(massChange))
}

// calcAltitude is an internal function that returns a new altitude due to
// a period of fuel usage along with the effect of gravity
func (l *Lander) calcAltitude(burnTime float64, massChange float64) float64 {
	return l.Altitude +
		(-1.0 * (Gravity * (math.Pow(burnTime, 2.0) / 2.0))) +
		(-1.0 * l.Velocity * burnTime) +
		(fuelIsp * burnTime * altitudeMassFactor(massChange))
}

// CalcDynamics returns a new Kinematics based on // the current lander
// Kinematics with altitude and velocity adjusted for the current burn
func (l *Lander) CalcDynamics(burnRate float64, burnTime float64) Kinematics {
	var newK Kinematics
	massChange := l.calcMassChange(burnRate, burnTime)
	newK.Velocity = l.calcVelocity(burnTime, massChange)
	newK.Altitude = l.calcAltitude(burnTime, massChange)
	newK.ElapsedTime = l.ElapsedTime
	return newK
}

// ActualBurnTime returns a value equal to or less than the requested burn
// time based on the burn rate and the available fuel
func (l *Lander) ActualBurnTime(burnRate float64, burnTime float64) float64 {
	if l.Fuel < (burnRate * burnTime) {
		return l.Fuel / burnRate
	} else {
		return burnTime
	}
}

// UpwardBurnTime handles the dynamics case where the burn has caused the
// velocity of the lander to become negative for some part of the burn
// This is probably the most obscure part of the program.
func (l *Lander) UpwardBurnTime(burnRate float64) float64 {
	factor := (1.0 - ((l.TotalMass * Gravity) /
		(fuelIsp * burnRate))) / 2.0
	return (((l.TotalMass * l.Velocity) / (fuelIsp * burnRate * (factor + math.Sqrt(
		(factor*factor)+(l.Velocity/fuelIsp))))) + 0.05)
}

// UpdateLander updates the Lander structure based on a burn and the
// (externally calculated) changes to altitude and velocity
func (l *Lander) UpdateLander(burnRate float64, burnTime float64, newPhys Kinematics) {
	l.Velocity = newPhys.Velocity
	l.Altitude = newPhys.Altitude
	l.ElapsedTime = l.ElapsedTime + burnTime
	l.Fuel -= (burnRate * burnTime)
	l.TotalMass = l.CapsuleMass + l.Fuel
}

// CalcImpact iteratively determines the moment of impact along with the
// velocity at that moment.  It may fail (loop endlessly) if called for a burn
// that does not end on the surface
func (l *Lander) CalcImpact(burnRate float64, timeToImpact float64) {
	var calcVelocity float64
	var nextPhys Kinematics
	for timeToImpact >= 0.005 {
		calcVelocity = l.Velocity + math.Sqrt((math.Pow(l.Velocity, 2))+
			(2.0*l.Altitude*(Gravity-
				(fuelIsp*(burnRate/l.TotalMass)))))
		timeToImpact = 2.0 * (l.Altitude / calcVelocity)
		nextPhys = l.CalcDynamics(burnRate, timeToImpact)
		l.UpdateLander(burnRate, timeToImpact, nextPhys)
	}
}
