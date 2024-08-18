package main

import (
	"fmt"
	"time"
	"xpc"

	"go.einride.tech/pid"
)

var (
	AltitudeHold = true
	AltitudeFeet = 2000.0
)

func main() {
	conn, err := xpc.Dial(xpc.Host{
		XPHost:  "localhost",
		XPPort:  49009,
		Timeout: time.Millisecond * 50,
	})

	if err != nil {
		panic(err)
	}

	samplingInterval := 100 * time.Millisecond

	ticker := time.NewTicker(samplingInterval)

	P := 0.05
	I := P / 5
	D := P / 5

	pitchSetpoint := 1.0
	pitch := pid.Controller{
		Config: pid.ControllerConfig{
			ProportionalGain: P,
			IntegralGain:     I,
			DerivativeGain:   D,
		},
	}

	rollSetpoint := 0.0
	roll := pid.Controller{
		Config: pid.ControllerConfig{
			ProportionalGain: P,
			IntegralGain:     I,
			DerivativeGain:   D,
		},
	}

	altitude := pid.Controller{
		Config: pid.ControllerConfig{
			ProportionalGain: P,
			IntegralGain:     I,
			DerivativeGain:   D,
		},
	}

	for {
		select {
		case <-ticker.C:
			posi, err := xpc.GetPOSI(conn, 0)
			if err != nil {
				panic(err)
			}

			if AltitudeHold {
				altitude.Update(pid.ControllerInput{
					ReferenceSignal:  AltitudeFeet,
					ActualSignal:     posi.Altitude * 3.281,
					SamplingInterval: samplingInterval,
				})

				pitchSetpoint = xpc.Clamp(altitude.State.ControlSignal, -15, 10)
			}

			pitch.Update(pid.ControllerInput{
				ReferenceSignal:  pitchSetpoint,
				ActualSignal:     float64(posi.Pitch),
				SamplingInterval: samplingInterval,
			})

			roll.Update(pid.ControllerInput{
				ReferenceSignal:  rollSetpoint,
				ActualSignal:     float64(posi.Roll),
				SamplingInterval: samplingInterval,
			})

			if err := xpc.SendCTRL(conn, &xpc.CTRL{
				Elevator:   float32(pitch.State.ControlSignal),
				Aileron:    float32(roll.State.ControlSignal),
				Rudder:     -998,
				Throttle:   -998,
				Flaps:      -998,
				Speedbrake: -998,
			}); err != nil {
				panic(err)
			}

			fmt.Printf("PITCH: %.2f -> %.2f	ROLL: %.2f -> %.2f	ALT: %.2f -> %.2f\n",
				posi.Pitch, pitch.State.ControlSignal, posi.Roll, roll.State.ControlSignal, posi.Altitude*3.281, altitude.State.ControlSignal)
		}
	}
}
