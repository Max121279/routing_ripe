package lib

import (
	"fmt"
	"os"
)

// Функция для обновления файла подсетей
func UpdateSubnetsFile(subnets []string, filePath string) error {

	// Записываем подсети в файл
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	for _, subnet := range subnets {
		_, err = file.WriteString(subnet + "\n")
		if err != nil {
			return fmt.Errorf("ошибка записи в файл: %v", err)
		}
	}

	fmt.Printf("Файл %s успешно обновлен\n", filePath)
	return nil
}
