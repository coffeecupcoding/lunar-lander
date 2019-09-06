use std::io;
use std::io::Write;

const GRAVITY: f64 = 0.001;

struct Kinematics {
    velocity: f64,
    altitude: f64,
    elapsed_time: f64,
}

fn velocity_mass_factor(q: f64) -> f64 {
    ((-q) +
        (-1.0 * ((q.powi(2)) / 2.0)) +
        (-1.0 * ((q.powi(3)) / 3.0)) +
        (-1.0 * ((q.powi(4)) / 4.0)) +
        (-1.0 * ((q.powi(5)) / 5.0)))
}

fn altitude_mass_factor(q: f64) -> f64 {
    ((q / 2.0) +
        ((q.powi(2)) / 6.0) +
        ((q.powi(3)) / 12.0) +
        ((q.powi(4)) / 20.0) +
        ((q.powi(5)) / 30.0))
}

struct Lander {
    capsule_mass: f64,
    fuel: f64,
    total_mass: f64,
    fuel_isp: f64,
    phys: Kinematics,
}

impl Lander {
    fn out_of_fuel(&self) -> bool {
        self.fuel < 0.001
    }
    fn calc_mass_change(&self, burn_rate: f64, burn_time: f64) -> f64 {
        (burn_rate * burn_time ) / self.total_mass
    }
    fn calc_velocity(&self, burn_time: f64, q: f64) -> f64 {
        (self.phys.velocity + 
             (GRAVITY * burn_time) +
             (self.fuel_isp * velocity_mass_factor(q)))
    }
    fn calc_altitude(&self, burn_time: f64, q: f64) -> f64 {
        (self.phys.altitude +
            ((-1.0) * (GRAVITY * ((burn_time.powi(2)) / 2.0))) +
            ((-1.0) * self.phys.velocity * burn_time) +
            (self.fuel_isp * burn_time * altitude_mass_factor(q)))
    }
    fn calc_dynamics(&self, burn_rate: f64, burn_time: f64) -> Kinematics {
        let q = self.calc_mass_change(burn_rate, burn_time);
        let velocity = self.calc_velocity(burn_time, q);
        let altitude = self.calc_altitude(burn_time, q);
        Kinematics { velocity, altitude,
            elapsed_time: self.phys.elapsed_time }
    }
    fn calc_burn_time(&self, burn_rate: f64, burn_time: f64) -> f64 {
        if self.fuel < (burn_rate * burn_time) {
            self.fuel / burn_rate
        } else {
            burn_time
        }
    }
    fn calc_upward_burn_time(&self, burn_rate: f64) -> f64 {
        let factor: f64 = (1.0 - ((self.total_mass * GRAVITY) /
            (self.fuel_isp * burn_rate))) / 2.0;
        (((self.total_mass * self.phys.velocity) /
        (self.fuel_isp * burn_rate * (factor + f64::sqrt(
        (factor * factor) + (self.phys.velocity / self.fuel_isp))))) + 0.05)
    }
    fn update_status(&mut self, burn_rate: f64, burn_time: f64,
            new_phys: &Kinematics, time: f64) -> f64 {
        self.phys.velocity = new_phys.velocity;
        self.phys.altitude = new_phys.altitude;
        self.phys.elapsed_time = self.phys.elapsed_time + burn_time;
        self.fuel = self.fuel - (burn_rate * burn_time);
        self.total_mass = self.capsule_mass + self.fuel;
        (time - burn_time)
    }
    fn calc_impact(&mut self, burn_rate: f64, mut iter_time: f64) {
        let mut iter_velocity: f64;
        while iter_time >= 0.005 {
            iter_velocity = self.phys.velocity +
                f64::sqrt((self.phys.velocity.powi(2)) +
                (2.0 * self.phys.altitude * ((GRAVITY) -
                (self.fuel_isp * (burn_rate / self.total_mass)))));   
            iter_time = 2.0 * (self.phys.altitude / iter_velocity);
            let new_phys = self.calc_dynamics(burn_rate, iter_time);
            let _ = self.update_status(burn_rate, iter_time, &new_phys, 0.0);
        }
    }
}

fn landing(lem: &mut Lander, burn_rate: f64, burn_time: f64) {
    lem.calc_impact(burn_rate, burn_time);
    end_game(lem)
}

fn out_of_fuel(lem: &mut Lander) {
    println!("\nFUEL OUT AT {:.2} SECONDS", lem.phys.elapsed_time);
    let seconds_to_impact: f64 = (((-1.0) * lem.phys.velocity) +
        f64::sqrt((lem.phys.velocity.powi(2)) +
        (2.0 * lem.phys.altitude * GRAVITY))) / GRAVITY;
    lem.phys.velocity = lem.phys.velocity + (GRAVITY * seconds_to_impact);
    lem.phys.elapsed_time = lem.phys.elapsed_time + seconds_to_impact;
    end_game(lem)
}

fn end_game(lem: &Lander) {
    let velocity_mph = lem.phys.velocity * 3600.0;
    println!("\nON THE MOON AT {} SECONDS\nIMPACT VELOCITY {} MPH\n",
        lem.phys.elapsed_time as u64, velocity_mph as i64);
    if velocity_mph < 1.2 {
        println!("PERFECT LANDING!!");
    } else if velocity_mph <= 10.0 {
        println!("GOOD LANDING (COULD BE BETTER)");
    } else if velocity_mph <= 60.0 {
        println!("CRAFT DAMAGE... YOU'RE STRANDED HERE\n\
            UNTIL A RESCUE PARTY ARRIVES.\n\
            I HOPE YOU HAVE ENOUGH OXYGEN!");
    } else {
        println!("THAT'S ONE SMALL IMPACT FOR THE MOON,\n\
            ONE GIANT BOOM FOR YOUR LANDER!\n\
            YOU BLASTED A NEW CRATER {} FEET DEEP!",
            (velocity_mph * 0.227) as u64);
    }
}

fn intro() {
    println!("\n              LUNAR\n\
        CREATIVE COMPUTING MORRISTOWN, NJ\n\n\
        THIS IS A COMPUTER SIMULATION OF AN\n\
        APOLLO LUNAR LANDING CAPSULE.\n\n\
        THE ON-BOARD COMPUTER HAS FAILED (IT WAS\n\
        MADE BY XEROX) SO YOU HAVE TO LAND THE\n\
        CAPSULE MANUALLY.\n");
}

fn output_header() {
    println!(" SEC  MILES  FEET    MPH   FUEL  RATE");
}

fn output_status(lem: &Lander) {
    let elapsed_time = lem.phys.elapsed_time.trunc() as u64;
    let miles = lem.phys.altitude.trunc() as i64;
    let feet = (5280.0 * lem.phys.altitude.fract()) as i64;
    let mph = (3600.0 * lem.phys.velocity) as i64;
    let fuel = lem.fuel as i64;
    print!("{:4}    {:3}  {:4}  {:5}  {:5}  ", elapsed_time, miles, feet,
        mph, fuel);
    io::stdout().flush().unwrap();
}

fn prompt_for_input(prompt: &str) -> Result<String, io::Error> {
    let mut response = String::new();
    print!("{}", prompt);
    io::stdout().flush().unwrap();
    io::stdin().read_line(&mut response)?;
    Ok(response)
}

fn get_burn_rate() -> f64 {
    let mut burn_rate: f64;
    loop {
        let rate_input = match prompt_for_input(&"") {
            Ok(rate) => rate,
            Err(_) => {
                print!("PLEASE ENTER A BURN RATE: ");
                continue
            },
        };
        burn_rate = match rate_input.trim().parse() {
            Ok(rate) => rate,
            Err(_) => {
                print!("PLEASE ENTER A BURN RATE: ");
                continue
            },
        };
        if burn_rate < 0.0 || burn_rate > 200.0 {
            println!("PLEASE ENTER A BURN RATE");
            print!("BETWEEN 0 AND 200: ");
            continue
        }
        break
    }
    burn_rate
}

fn run_game() {
    println!("\
        SET THE BURN RATE OF THE RETRO ROCKETS\n\
        TO ANY VALUE BETWEEN 0 (FREE FALL) AND\n\
        200 (MAXIMUM BURN) IN POUNDS PER SECOND.\n\
        SET A NEW BURN RATE EVERY 10 SECONDS.\n\n\
        CAPSULE DRY WEIGHT IS 16,500 LBS;\n\
        INITIAL FUEL IS 16,500 LBS.\n\n\n\
        GOOD LUCK!\n");

    // Defaults
    let altitude: f64 = 120.0;
    let velocity: f64 = 1.0;
    let capsule_mass: f64 = 16500.0;
    let fuel: f64 = 16500.0;
    let total_mass: f64 = capsule_mass + fuel;
    let fuel_isp: f64 = 1.8;
    let phys = Kinematics { altitude, velocity, elapsed_time: 0.0 };

    let mut lem = Lander {
        capsule_mass,
        fuel,
        total_mass,
        fuel_isp,
        phys,
    };

    output_header();
    'game: loop {
        output_status(&lem);
        let burn_rate = get_burn_rate();
        let mut this_period: f64 = 10.0;
        loop {
            if lem.out_of_fuel() {
                out_of_fuel(&mut lem);
                break 'game;
            };
            if this_period < 0.001 {
                break
            };
            let mut burn_time = lem.calc_burn_time(burn_rate, this_period);
            lem.calc_dynamics(burn_rate, burn_time);
            let new_phys = lem.calc_dynamics(burn_rate, burn_time);
            if new_phys.altitude <= 0.0 {
                landing(&mut lem, burn_rate, burn_time);
                break 'game;
            }
            if lem.phys.velocity > 0.0 {
                if new_phys.velocity < 0.0 {
                    burn_time = lem.calc_upward_burn_time(burn_rate);
                    let new_phys = lem.calc_dynamics(burn_rate, burn_time);
                    if new_phys.altitude <= 0.0 {
                        landing(&mut lem, burn_rate, burn_time);
                        break 'game;
                    }
                }
            }
            this_period = lem.update_status(burn_rate, burn_time,
                              &new_phys, this_period);
        }
    }
}

fn main() {
    intro();
    loop {
        run_game();
        match prompt_for_input(&"\nTRY AGAIN?? ") {
            Ok(response) => {
                if !(response.starts_with('Y') | response.starts_with('y')) {
                    break
                }
            },
            Err(_) => {
                break
            },
        };
    }
}

