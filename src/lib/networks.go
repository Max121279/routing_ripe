package lib

import (
	"fmt"
	"net"
	"sort"
)

// sortByNumericalValue реализует сортировку подсетей по числовому представлению первого IP
type ByNumericalValue []net.IPNet

func (a ByNumericalValue) Len() int           { return len(a) }
func (a ByNumericalValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNumericalValue) Less(i, j int) bool { return ipToInt(a[i].IP) < ipToInt(a[j].IP) }

// intToIP преобразует целое число в IP-адрес
func intToIP(ipInt int) net.IP {
	return net.IPv4(byte(ipInt>>24), byte(ipInt>>16&0xFF), byte(ipInt>>8&0xFF), byte(ipInt&0xFF))
}

// ipToInt преобразует IP-адрес в целое число
func ipToInt(ip net.IP) int {
	ip = ip.To4()
	return int(ip[0])<<24 | int(ip[1])<<16 | int(ip[2])<<8 | int(ip[3])
}

// NextIPInRange находит следующий IP-адрес в диапазоне

// NextIPInRange находит следующий IP-адрес в диапазоне для IPv4
func NextIPInRange(ip net.IP, mask net.IPMask) net.IP {

	// Создаем копию IP-адреса, чтобы не изменять исходный
	retIP := make(net.IP, len(ip.To4()))
	copy(retIP, ip.To4())

	ones, _ := mask.Size()
	increment := 1 << uint(32-ones)

	// Находим следующий IP-адрес для IPv4
	for i := 3; i >= 0 && increment > 0; i-- {
		val := int(retIP[i]) + increment
		retIP[i] = byte(val % 256)
		increment = val / 256
	}

	return retIP
}

// getMidIP вычисляет средний IP-адрес для подсети
func getMidIP(startIP net.IP, mask net.IPMask) net.IP {
	// Преобразуем IP и маску в целые числа
	startIPInt := ipToInt(startIP)
	maskSize, _ := mask.Size()

	// Вычисляем количество адресов в подсети
	totalIPs := 1 << uint(32-maskSize)
	midIPInt := startIPInt + totalIPs/2

	// Преобразуем обратно в IP
	midIP := intToIP(midIPInt)
	return midIP
}

// containsExcludedInRange проверяет, содержится ли исключаемый IP в диапазоне подсети
func containsExcludedInRange(subnetRange net.IPNet, excludedIPs map[string]bool) bool {
	for ip := range excludedIPs {
		if subnetRange.Contains(net.ParseIP(ip)) {
			return true
		}
	}
	return false
}

// splitSubnetWithExclusions делит подсеть на два сегмента и исключает IP
func splitSubnetWithExclusions(ipNet net.IPNet, excludedIPs map[string]bool) []net.IPNet {
	var result []net.IPNet
	queue := []net.IPNet{ipNet}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Проверяем, содержится ли исключаемый IP в текущей подсети
		var containsExcluded bool
		for ip := range excludedIPs {
			if current.Contains(net.ParseIP(ip)) {
				containsExcluded = true
				break
			}
		}

		// Если в подсети нет исключаемых адресов, добавляем её в результат
		if !containsExcluded {
			result = append(result, current)
			continue
		}

		// Если подсеть минимальна (/32), добавляем её и продолжаем
		maskSize, bits := current.Mask.Size()
		if maskSize == bits {
			result = append(result, current)
			continue
		}

		// Разделяем подсеть на два сегмента, делим посередине
		nextMask := maskSize + 1
		midIP := getMidIP(current.IP, current.Mask)

		// Левый сегмент
		left := net.IPNet{
			IP:   current.IP.Mask(net.CIDRMask(nextMask, bits)),
			Mask: net.CIDRMask(nextMask, bits),
		}

		// Правый сегмент
		right := net.IPNet{
			IP:   midIP,
			Mask: net.CIDRMask(nextMask, bits),
		}

		// Добавляем поддиапазоны в очередь для дальнейшей обработки
		queue = append(queue, left, right)

	}

	return result
}

// SummarizeSubnetsWithExclusions - основная функция для суммаризации подсетей с исключениями
func SummarizeSubnetsWithExclusions(subnetCIDR string, excludedIPs map[string]bool) ([]string, error) {
	var result []net.IPNet

	_, ipNet, err := net.ParseCIDR(subnetCIDR)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга подсети %s: %v", subnetCIDR, err)
	}

	// Разделяем подсеть с учётом исключаемых IP
	result = append(result, splitSubnetWithExclusions(*ipNet, excludedIPs)...)

	// Сортируем результаты по строковому представлению
	sort.Sort(ByNumericalValue(result))

	// Преобразуем результат в строковый формат и возвращаем
	var summarized []string
	for _, subnet := range result {
		summarized = append(summarized, subnet.String())
	}

	return summarized, nil
}
