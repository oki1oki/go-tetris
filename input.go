package main

const (
	repeatDelay    = 12 // тиков до начала повтора
	repeatInterval = 3  // тиков между повторами
)

type keyRepeat struct {
	held int // сколько тиков подряд клавиша зажата
}

func (k *keyRepeat) fire(pressed bool) bool {
	if !pressed {
		k.held = 0
		return false
	}

	k.held++
	if k.held == 1 {
		return true
	}

	return k.held > repeatDelay && (k.held-repeatDelay)%repeatInterval == 0
}
