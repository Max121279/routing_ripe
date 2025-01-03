package lib

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Функция для удаления маршрутов
func RemoveRoutes(filePath, iface string) error {
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

		// Удаляем маршрут
		cmd := exec.Command("ip", "route", "del", subnet, "dev", iface)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Ошибка удаления маршрута %s: %s\n", subnet, string(output))
		} else {
			fmt.Printf("Маршрут для подсети %s удален\n", subnet)
		}
	}

	return nil
}

// Функция для добавления маршрута
func AddRoutes(filePath, iface string) error {
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
		cmd := exec.Command("ip", "route", "add", subnet, "dev", iface)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Ошибка добавления маршрута %s: %s\n", subnet, string(output))
		} else {
			fmt.Printf("Маршрут для подсети %s добавлен\n", subnet)
		}
	}
	return nil
}
