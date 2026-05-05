package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"autoscaler"
)

func main() {
	if len(os.Args) == 1 {
		runAutoMode()
		return
	}

	if len(os.Args) != 3 {
		fmt.Println("Использование: autoscaler <cpu_load> <traffic_rate>")
		fmt.Println("Или запустите без аргументов для работы с локальными метриками ПК")
		fmt.Println("Пример: autoscaler 65.5 420")
		return
	}

	cpu, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Неверное значение CPU: %v\n", err)
		return
	}

	traffic, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Printf("Неверное значение трафика: %v\n", err)
		return
	}

	controller := autoscaler.NewScalingController()
	action := controller.ComputeScalingAction(cpu, traffic)
	fmt.Println(action)
}

func runAutoMode() {
	fmt.Println("Запущен режим постоянного мониторинга локального ПК.")
	fmt.Println("Для остановки нажмите Ctrl+C.")

	controller := autoscaler.NewScalingController()
	for {
		cpu, traffic, err := autoscaler.GetLocalMetrics()
		if err != nil {
			fmt.Printf("Ошибка мониторинга локальной системы: %v\n", err)
			return
		}

		action := controller.ComputeScalingAction(cpu, traffic)

		fmt.Printf("\n[Время: %s]\n", time.Now().Format("15:04:05"))
		fmt.Printf("Загрузка CPU: %.1f%%\n", cpu)
		fmt.Printf("Показатель трафика: %.0f\n", traffic)
		fmt.Printf("Решение по масштабированию: %d\n", action)

		if action > 0 {
			fmt.Println("Рекомендуется подключить дополнительный ресурс.")
		} else if action < 0 {
			fmt.Println("Рекомендуется уменьшить количество ресурсов.")
		} else {
			fmt.Println("Дополнительные ресурсы не требуются.")
		}

		time.Sleep(5 * time.Second)
	}
}
