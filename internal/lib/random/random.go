package random

import (
	"math/rand"
	"time"
)

// Глобальный генератор случайных чисел
var rnd *rand.Rand

// Инициализация генератора случайных чисел при старте программы
func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// TODO: Дополнительно проверять текущии alias на наличие в базе данных
func NewRandomString(size int) string {
	if size <= 0 {
		return ""
	}

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
