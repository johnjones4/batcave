package iface

// type GPIOController struct {
// 	controller
// 	modePins  []rpio.Pin
// 	togglePin rpio.Pin
// 	debounce  time.Duration
// }

// func NewGPIOController(log logrus.FieldLogger, modePins []int, togglePin int, debounce time.Duration) (*GPIOController, error) {
// 	err := rpio.Open()
// 	if err != nil {
// 		return nil, err
// 	}

// 	mpins := make([]rpio.Pin, len(modePins))
// 	for i, p := range modePins {
// 		mpins[i] = rpio.Pin(p)
// 		mpins[i].PullUp()
// 	}

// 	return &GPIOController{
// 		controller: controller{
// 			log:           log,
// 			signalChannel: make(chan core.ControllerSignal),
// 		},
// 		modePins:  mpins,
// 		debounce:  debounce,
// 		togglePin: rpio.Pin(togglePin),
// 	}, nil
// }

// func (c *GPIOController) Close() error {
// 	return rpio.Close()
// }

// func (c *GPIOController) Start(ctx context.Context) {
// 	var lastPress time.Time
// 	var isTogglingDown bool
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			if isTogglingDown && c.togglePin.Read() == rpio.High {
// 				isTogglingDown = false
// 				c.signalChannel <- core.ControllerSignal{
// 					Type: core.SignalTypeToggle,
// 				}
// 				continue
// 			}

// 			if !isTogglingDown && c.togglePin.Read() == rpio.Low {
// 				isTogglingDown = true
// 				c.signalChannel <- core.ControllerSignal{
// 					Type: core.SignalTypeToggle,
// 				}
// 				continue
// 			}

// 		modePinLoop:
// 			for i, pin := range c.modePins {
// 				if pin.Read() == rpio.Low && time.Since(lastPress) > c.debounce {
// 					lastPress = time.Now()
// 					c.signalChannel <- core.ControllerSignal{
// 						Type: core.SignalTypeMode,
// 						Mode: i,
// 					}
// 					break modePinLoop
// 				}
// 			}
// 		}
// 	}
// }
