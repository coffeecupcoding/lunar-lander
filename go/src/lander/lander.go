package lander

import (
	"math"
)

// Specific Impulse of fuel, in miles/second
//   These units mean fuel is measured by mass, not weight
const fuelIsp = 1.8

// Lunar gravity in miles/(second^2)
const Gravity = 0.001

type Kinematics struct {
	Velocity    float64
	Altitude    float64
	ElapsedTime float64
}

func velocityMassFactor(massChange float64) float64 {
	return ((-1.0 * massChange) +
		(-1.0 * (math.Pow(massChange, 2.0) / 2.0)) +
		(-1.0 * (math.Pow(massChange, 3.0) / 3.0)) +
		(-1.0 * (math.Pow(massChange, 4.0) / 4.0)) +
		(-1.0 * (math.Pow(massChange, 5.0) / 5.0)))
}

func altitudeMassFactor(massChange float64) float64 {
	return ((massChange / 2.0) +
		(math.Pow(massChange, 2.0) / 6.0) +
		(math.Pow(massChange, 3.0) / 12.0) +
		(math.Pow(massChange, 4.0) / 20.0) +
		(math.Pow(massChange, 5.0) / 30.0))
}

type Lander struct {
	CapsuleMass float64
	Fuel        float64
	TotalMass   float64
	Kinematics
}

func (l *Lander) OutOfFuel() bool {
	return l.Fuel < 0.001
}

func (l *Lander) calcMassChange(burnRate float64, burnTime float64) float64 {
	return (burnRate * burnTime) / l.TotalMass
}

func (l *Lander) calcVelocity(burnTime float64, massChange float64) float64 {
	return l.Velocity + (Gravity * burnTime) +
		(fuelIsp * velocityMassFactor(massChange))
}

func (l *Lander) calcAltitude(burnTime float64, massChange float64) float64 {
	return l.Altitude +
		(-1.0 * (Gravity * (math.Pow(burnTime, 2.0) / 2.0))) +
		(-1.0 * l.Velocity * burnTime) +
		(fuelIsp * burnTime * altitudeMassFactor(massChange))
}

func (l *Lander) CalcDynamics(burnRate float64, burnTime float64) Kinematics {
	var newK Kinematics
	massChange := l.calcMassChange(burnRate, burnTime)
	newK.Velocity = l.calcVelocity(burnTime, massChange)
	newK.Altitude = l.calcAltitude(burnTime, massChange)
	newK.ElapsedTime = l.ElapsedTime
	return newK
}

func (l *Lander) ActualBurnTime(burnRate float64, burnTime float64) float64 {
	if l.Fuel < (burnRate * burnTime) {
		return l.Fuel / burnRate
	} else {
		return burnTime
	}
}

func (l *Lander) UpwardBurnTime(burnRate float64) float64 {
	factor := (1.0 - ((l.TotalMass * Gravity) /
		(fuelIsp * burnRate))) / 2.0
	return (((l.TotalMass * l.Velocity) / (fuelIsp * burnRate * (factor + math.Sqrt(
		(factor*factor)+(l.Velocity/fuelIsp))))) + 0.05)
}

func (l *Lander) UpdateLander(burnRate float64, burnTime float64, newPhys Kinematics) {
	l.Velocity = newPhys.Velocity
	l.Altitude = newPhys.Altitude
	l.ElapsedTime = l.ElapsedTime + burnTime
	l.Fuel -= (burnRate * burnTime)
	l.TotalMass = l.CapsuleMass + l.Fuel
}

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
