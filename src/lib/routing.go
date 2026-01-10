package lib

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Функция для удаления маршрутов
func RemoveRoutes(filePath, iface, gateway string) error {
	// Открываем предыдущий файл
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла для удаления маршрутов: %v", err)
	}
	defer file.Close()

	// Читаем файл построчно
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subnet := strings.TrimSpace(scanner.Text())
		if !strings.Contains(subnet, "/") {
			continue
		}

		if iface != "" {
			// Удаляем маршрут
			cmd := exec.Command("ip", "route", "del", subnet, "dev", iface)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Ошибка удаления маршрута %s: %s\n", subnet, string(output))
			} else {
				fmt.Printf("Маршрут для подсети %s удален\n", subnet)
			}
		} else if gateway != "" {
			// Удаляем маршрут
			cmd := exec.Command("ip", "route", "del", subnet, "via", gateway)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Ошибка удаления маршрута %s: %s\n", subnet, string(output))
			} else {
				fmt.Printf("Маршрут для подсети %s удален\n", subnet)
			}
		}

	}

	return nil
}

// Функция для добавления маршрута
func AddRoutes(filePath, iface, gateway string) error {
	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла для добавления маршрутов: %v", err)
	}
	defer file.Close()
	// Читаем файл построчно
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subnet := strings.TrimSpace(scanner.Text())
		if !strings.Contains(subnet, "/") {
			continue
		}
		if iface != "" {

			cmd := exec.Command("ip", "route", "add", subnet, "dev", iface)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Ошибка добавления маршрута %s: %s\n", subnet, string(output))
			} else {
				fmt.Printf("Маршрут для подсети %s добавлен\n", subnet)
			}

		} else if gateway != "" {

			cmd := exec.Command("ip", "route", "add", subnet, "via", gateway)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Ошибка добавления маршрута %s: %s\n", subnet, string(output))
			} else {
				fmt.Printf("Маршрут для подсети %s добавлен\n", subnet)
			}
		}
	}
	return nil
}
