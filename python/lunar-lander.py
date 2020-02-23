#!/usr/bin/env python3

"""
A version of the old lunar lander game 
Rewritten from 101 BASIC Computer Games by David Ahl
"""

import math
import sys

# Globals
gravity = 0.001


class GameOver(Exception): pass


class Kinematics:
    """
    Kinematics-related data
    """
    def __init__(self, altitude, velocity, time=0.0):
        self.altitude = altitude
        self.velocity = velocity
        self.elapsed_time = time

    def altitude_mass_factor(q):
        """
        Effect on altitude from integrated change of mass
        q is the proportional change in mass due to fuel usage
        """
        factor = ((q / 2.0) +
                 ((q ** 2) / 6.0) +
                 ((q ** 3) / 12.0) +
                 ((q ** 4) / 20.0) +
                 ((q ** 5) / 30.0))
        return factor

    def velocity_mass_factor(q):
        """
        Effect on velocity from integrated change of mass
        q is the proportional change in mass due to fuel usage
        """
        factor = ((-q) +
                 (-1 * ((q ** 2) / 2.0)) +
                 (-1 * ((q ** 3) / 3.0)) +
                 (-1 * ((q ** 4) / 4.0)) +
                 (-1 * ((q ** 5) / 5.0)))
        return factor


class Lander:
    """
    Lander data and status
    """
    def __init__(self, initial_altitude=120.0, initial_velocity=1.0,
            capsule_mass=16500.0, initial_fuel=16500.0,
            fuel_specific_impulse = 1.8):
         self.capsule_mass = capsule_mass
         self.fuel = initial_fuel  # in pounds, um, mass
         self.total_mass = self.capsule_mass + self.fuel
         self.fuel_isp = fuel_specific_impulse
         self.phys = Kinematics(initial_altitude, initial_velocity)

    def out_of_fuel(self):
        if self.fuel < 0.001:
            return True
        return False

    def calc_mass_change(self, burn_rate, burn_time):
        # returns proportional change in mass due to fuel usage
        return (burn_rate * burn_time) / self.total_mass

    def calc_velocity(self, burn_rate, burn_time, q):
        """
        Caculate new velocity based on starting velocity and delta t
        """
        velocity = (self.phys.velocity +
                   (gravity * burn_time) +
                   (self.fuel_isp * Kinematics.velocity_mass_factor(q)))
        return velocity

    def calc_altitude(self, burn_rate, burn_time, q):
        """
        Caculate new altitude based on starting altitude and delta t
        """
        altitude = (self.phys.altitude +
                   ((-1) * (gravity * ((burn_time ** 2) / 2))) +
                   ((-1) * self.phys.velocity * burn_time) +
                   (self.fuel_isp * burn_time *
                       Kinematics.altitude_mass_factor(q)))
        return altitude

    def calc_dynamics(self, burn_rate, burn_time):
        """
        Return new altitude and velocity based on burn rate and time
        Returns the values as a Kinematics instance
        """
        q = self.calc_mass_change(burn_rate, burn_time)
        velocity = self.calc_velocity(burn_rate, burn_time, q)
        altitude = self.calc_altitude(burn_rate, burn_time, q)
        return Kinematics(altitude, velocity)

    def calc_burn_time(self, burn_rate, time_left):
        # Reduce burn time if there isn't enough fuel
        burn_time = time_left
        if self.fuel < burn_rate * burn_time:
            burn_time = self.fuel / burn_rate
        return burn_time

    def calc_upward_burn_time(self, burn_rate):
        """
        Handle the case where the burn takes the velocity negative
        This is the most obscure part of the original program
        """
        factor = (1 -
            (self.total_mass * gravity) / (self.fuel_isp * burn_rate)) / 2
        # the units work at least...
        # the '+ 0.05' at the end is for ?? good luck?
        new_burn_time = ((self.total_mass * self.phys.velocity) /
            (self.fuel_isp * burn_rate * (factor + math.sqrt(
            (factor * factor) + (self.phys.velocity / self.fuel_isp))))) + 0.05
        return new_burn_time

    def update_status(self, burn_rate, burn_time, new_phys, time_left):
        """
        Set lander data and elapsed time
        Returns the remaining time in this period
        """
        self.phys.elapsed_time = self.phys.elapsed_time + burn_time
        self.phys.altitude = new_phys.altitude
        self.phys.velocity = new_phys.velocity
        self.fuel = self.fuel - (burn_rate * burn_time)
        self.total_mass = self.capsule_mass + self.fuel
        return (time_left - burn_time)

    def calc_impact(self, burn_rate, iter_time):
        """
        Determine 'landing' parameters
        """
        while iter_time >= 0.005:
            iter_velocity = (self.phys.velocity +
                math.sqrt((self.phys.velocity ** 2) +
                (2 * self.phys.altitude * ((gravity) -
                (self.fuel_isp * (burn_rate / self.total_mass))))))
            iter_time = 2 * (self.phys.altitude / iter_velocity)
            new_phys = self.calc_dynamics(burn_rate, iter_time)
            # the return value is intentionally ignored here
            # because the end is near
            self.update_status(burn_rate, iter_time, new_phys, 0)



def landing(lem, burn_rate, burn_time):
    lem.calc_impact(burn_rate, burn_time)
    end_game(lem)

def out_of_fuel(lem):
    """
    Print message and set final velocity and time
    """
    print("\nFUEL OUT AT %d SECONDS" % lem.phys.elapsed_time)
    seconds_to_impact = (((-1.0) * lem.phys.velocity) +
        math.sqrt((lem.phys.velocity * lem.phys.velocity) +
            (2.0 * lem.phys.altitude * gravity))) / gravity
    lem.phys.velocity = lem.phys.velocity + (gravity * seconds_to_impact)
    lem.phys.elapsed_time = lem.phys.elapsed_time + seconds_to_impact
    end_game(lem)

def end_game(lem):
    velocity_mph = lem.phys.velocity * 3600.0
    print("\nON THE MOON AT %d SECONDS\nIMPACT VELOCITY %d MPH\n" % 
        (lem.phys.elapsed_time, velocity_mph))
    if velocity_mph < 1.2:
        print("PERFECT LANDING!!")
    elif velocity_mph <= 10.0:
        print("GOOD LANDING (COULD BE BETTER)")
    elif velocity_mph <= 60.0:
        print("CRAFT DAMAGE... YOU'RE STRANDED HERE")
        print("UNTIL A RESCUE PARTY ARRIVES.")
        print("I HOPE YOU HAVE ENOUGH OXYGEN!")
    else:
        print("THAT'S ONE SMALL IMPACT FOR THE MOON,")
        print("ONE GIANT BOOM FOR YOUR LANDER!")
        print("YOU BLASTED A NEW CRATER %d FEET DEEP!" %
            (velocity_mph * 0.227))
    raise GameOver

def intro():
    print("\n"
        "              LUNAR\n"
        "CREATIVE COMPUTING MORRISTOWN, NJ\n\n"
        "THIS IS A COMPUTER SIMULATION OF AN\n"
        "APOLLO LUNAR LANDING CAPSULE.\n\n"
        "THE ON-BOARD COMPUTER HAS FAILED (IT WAS\n"
        "MADE BY XEROX) SO YOU HAVE TO LAND THE\n"
        "CAPSULE MANUALLY.\n\n"
    )

def output_header():
    print(" SEC  MILES  FEET    MPH   FUEL  RATE")

def output_status(lem):
    miles = int(lem.phys.altitude)
    feet = 5280 * (lem.phys.altitude - miles)
    mph = 3600 * lem.phys.velocity
    print("%4d    %3d  %4d  %5d  %5d" %
        (lem.phys.elapsed_time, miles, feet, mph, lem.fuel), end='  ')

def get_burn_rate():
    user_input = ""
    while not user_input:
        user_input = input()
        if user_input == "":
            print("PLEASE ENTER A BURN RATE:", end=' ')
        else:
            try:
                burn_rate = float(user_input)
                if (burn_rate < 0.0) or (burn_rate > 200.0):
                    print("PLEASE ENTER A BURN RATE")
                    print("BETWEEN 0 AND 200 :", end=' ')
                    user_input = ""
            except ValueError:
                print("PLEASE ENTER A BURN RATE:", end=' ')
                user_input = ""
    return burn_rate

def run_game():
    print("SET THE BURN RATE OF THE RETRO ROCKETS\n"
        "TO ANY VALUE BETWEEN 0 (FREE FALL) AND\n"
        "200 (MAXIMUM BURN) IN POUNDS PER SECOND.\n"
        "SET A NEW BURN RATE EVERY 10 SECONDS.\n\n"
        "CAPSULE DRY WEIGHT IS 16,500 LBS;\n"
        "INITIAL FUEL IS 16,500 LBS.\n\n\n"
        "GOOD LUCK!\n\n"
    )
    
    lem = Lander()
    output_header()
    while True:
        burn_rate = 0.0
        output_status(lem)
        burn_rate = get_burn_rate()
        this_period = 10.0
        while True:
            if lem.out_of_fuel():
                out_of_fuel(lem)
            if this_period < 0.001:
                break
            burn_time = lem.calc_burn_time(burn_rate, this_period)
            lem.calc_dynamics(burn_rate, burn_time)
            new_phys = lem.calc_dynamics(burn_rate, burn_time)
            if new_phys.altitude <= 0:
                landing(lem, burn_rate, burn_time)
            if lem.phys.velocity > 0:
                if new_phys.velocity < 0:
                    burn_time = lem.calc_upward_burn_time(burn_rate)
                    new_phys = lem.calc_dynamics(burn_rate, burn_time)
                    if new_phys.altitude <= 0:
                        landing(lem, burn_rate, burn_time)
            this_period = lem.update_status(burn_rate, burn_time,
                              new_phys, this_period)

def run():
    intro()
    another_game = True
    while another_game:
        try:
            run_game()
        except GameOver:
            reply = input("\nTRY AGAIN?? ")
            if not (reply.startswith(('y', 'Y'))):
                another_game = False
        except KeyboardInterrupt:
            print("\nEXITING GAME")
            sys.exit(0)



if __name__ == "__main__":
    run()
