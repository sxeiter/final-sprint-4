package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65 // базовая длина шага, м
	mInKm                      = 1000 // метров в километре
	minInH                     = 60   // минут в часе
	stepLengthCoefficient      = 0.45 // коэффициент для расчёта длины шага на основе роста
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчёта калорий при ходьбе
)

// parseTraining парсит строку с данными о тренировке.
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных: ожидается три элемента через запятую")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, errors.New("количество шагов должно быть целым числом")
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}

	activityType := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, errors.New("некорректный формат продолжительности")
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше нуля")
	}

	return steps, activityType, duration, nil
}

// distance вычисляет дистанцию в километрах.
// Использует:
// - lenStep как базовую длину шага, ЕСЛИ рост не передан (height <= 0)
// - height * stepLengthCoefficient, ЕСЛИ рост передан (height > 0)
func distance(steps int, height float64) float64 {
	var stepLen float64
	if height > 0 {
		stepLen = height * stepLengthCoefficient // учёт роста пользователя
	} else {
		stepLen = lenStep // базовый шаг 0.65 м
	}

	totalDistanceMeters := float64(steps) * stepLen
	return totalDistanceMeters / mInKm // перевод в километры
}

// meanSpeed вычисляет среднюю скорость в км/ч.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distKm := distance(steps, height)
	durationHours := duration.Hours()

	return distKm / durationHours
}

// RunningSpentCalories рассчитывает калории для бега
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("длительность должна быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH
	return calories, nil
}

// WalkingSpentCalories рассчитывает калории для ходьбы
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("длительность должна быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)
	if speed == 0 {
		return 0, fmt.Errorf("невозможно рассчитать скорость")
	}

	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient // корректировка для ходьбы
	return calories, nil
}

// TrainingInfo формирует строку с подробной информацией о тренировке.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distKm := distance(steps, height)
	speedKmH := meanSpeed(steps, height, duration)

	var calories float64
	switch activityType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки: " + activityType)
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	durationHours := duration.Hours()

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType, durationHours, distKm, speedKmH, calories,
	), nil
}
