package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/Max121279/routing_ripe/src/lib"
)

const baseURL = "https://stat.ripe.net/data/country-resource-list/data.json?resource="

func fetchSubnets(config *lib.Config) ([]string, error) {
	url := baseURL + config.CountryCode
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки данных: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	var result struct {
		Data struct {
			Resources struct {
				IPv4 []string `json:"ipv4"`
			} `json:"resources"`
		} `json:"data"`
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка разбора JSON: %v", err)
	}

	ignoredIPs := make(map[string]bool)
	for _, ip := range config.IgnoredIPs {
		ignoredIPs[ip] = true
	}

	var subnets []string
	for _, resource := range result.Data.Resources.IPv4 {
		if strings.Contains(resource, "-") {
			ips := strings.Split(resource, "-")
			if len(ips) == 2 {
				cidrs, _ := ipRangeToCIDR(ips[0], ips[1])
				for _, cidr := range cidrs {
					//filtered, _ := filterSubnets(cidr, ignoredIPs)
					filtered, _ := lib.SummarizeSubnetsWithExclusions(cidr, ignoredIPs)
					subnets = append(subnets, filtered...)
				}
			}
		} else {
			//filtered, _ := filterSubnets(resource, ignoredIPs)
			filtered, _ := lib.SummarizeSubnetsWithExclusions(resource, ignoredIPs)
			subnets = append(subnets, filtered...)
		}
	}

	return summarizeSubnets(subnets), nil
}

func ipRangeToCIDR(start, end string) ([]string, error) {
	startIP := net.ParseIP(start).To4()
	endIP := net.ParseIP(end).To4()
	if startIP == nil || endIP == nil {
		return nil, fmt.Errorf("некорректный формат IP: %s - %s", start, end)
	}

	var result []string
	for bytesCompare(startIP, endIP) <= 0 {
		mask := 32
		for mask > 0 {
			mask--
			network := startIP.Mask(net.CIDRMask(mask, 32))
			if bytesCompare(network, startIP) != 0 || bytesCompare(lastIPInCIDR(network, mask), endIP) > 0 {
				mask++
				break
			}
		}

		cidr := fmt.Sprintf("%s/%d", startIP.String(), mask)
		result = append(result, cidr)
		startIP = lib.NextIPInRange(lastIPInCIDR(startIP, mask), net.CIDRMask(32, 32))
	}

	return result, nil
}

func summarizeSubnets(subnets []string) []string {
	ipNets := make([]net.IPNet, len(subnets))
	for i, subnet := range subnets {
		_, ipNet, _ := net.ParseCIDR(subnet)
		ipNets[i] = *ipNet
	}

	sort.Sort(lib.ByNumericalValue(ipNets))

	var summarized []net.IPNet
	for _, ipNet := range ipNets {
		if len(summarized) == 0 {
			summarized = append(summarized, ipNet)
			continue
		}

		last := summarized[len(summarized)-1]
		if canMerge(last, ipNet) {
			merged := mergeSubnets(last, ipNet)
			summarized[len(summarized)-1] = merged
		} else {
			summarized = append(summarized, ipNet)
		}
	}

	// Преобразуем в строковый формат
	result := make([]string, len(summarized))
	for i, net := range summarized {
		result[i] = net.String()
	}

	return result
}

func lastIPInCIDR(ip net.IP, prefixSize int) net.IP {
	mask := net.CIDRMask(prefixSize, 32)
	network := ip.Mask(mask)
	broadcast := make(net.IP, len(network))
	copy(broadcast, network)
	for i := 0; i < len(broadcast); i++ {
		broadcast[i] |= ^mask[i]
	}
	return broadcast
}

func filterSubnets(cidr string, ignoredIPs map[string]bool) ([]string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора подсети %s: %v", cidr, err)
	}

	var result []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); ip = lib.NextIPInRange(ip, net.CIDRMask(32, 32)) {
		if !ignoredIPs[ip.String()] {
			result = append(result, ip.String()+"/32")
		}
	}

	return result, nil
}

// summarizeSubnetsWithExclusions - основная функция для суммаризации подсетей с исключениями
func summarizeSubnetsWithExclusions(subnets []string, excludedIPs map[string]bool) ([]string, error) {
	var result []*net.IPNet

	// Перебираем все подсети
	for _, subnet := range subnets {
		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			return nil, fmt.Errorf("ошибка разбора подсети %s: %v", subnet, err)
		}

		// Разделяем подсеть с учётом исключаемых IP
		result = append(result, splitSubnetWithExclusions(ipNet, excludedIPs)...)
	}

	// Преобразуем результат в строковый формат и возвращаем
	var summarized []string
	for _, net := range result {
		summarized = append(summarized, net.String())
	}

	return summarized, nil
}

// splitSubnetWithExclusions делит подсеть на два сегмента и исключает IP
func splitSubnetWithExclusions(ipNet *net.IPNet, excludedIPs map[string]bool) []*net.IPNet {
	queue := []*net.IPNet{ipNet}
	var result []*net.IPNet

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

		// Разделяем подсеть на два сегмента
		nextMask := maskSize + 1
		left := &net.IPNet{
			IP:   current.IP.Mask(net.CIDRMask(nextMask, bits)),
			Mask: net.CIDRMask(nextMask, bits),
		}
		right := &net.IPNet{
			IP:   lib.NextIPInRange(left.IP, left.Mask),
			Mask: net.CIDRMask(nextMask, bits),
		}

		// Если исключаемый адрес не входит в оба сегмента, добавляем поддиапазоны в результат
		shouldAddLeft := true
		shouldAddRight := true
		for ip := range excludedIPs {
			excludedIP := net.ParseIP(ip)
			if left.Contains(excludedIP) {
				shouldAddLeft = false
			}
			if right.Contains(excludedIP) {
				shouldAddRight = false
			}
		}

		// Добавляем только те сегменты, которые не содержат исключаемые адреса
		if shouldAddLeft {
			result = append(result, left)
		}
		if shouldAddRight {
			result = append(result, right)
		}

		// Добавляем оба сегмента в очередь для дальнейшей обработки
		if shouldAddLeft {
			queue = append(queue, left)
		}
		if shouldAddRight {
			queue = append(queue, right)
		}
	}

	return result
}

func canMerge(a, b net.IPNet) bool {
	sizeA, bits := a.Mask.Size()
	sizeB, _ := b.Mask.Size()
	if sizeA != sizeB || bits != 32 {
		return false
	}

	// Проверяем, являются ли подсети последовательными
	nextIP := lib.NextIPInRange(a.IP, a.Mask)
	// проверяем начало сети совпадает с А
	newMask := net.CIDRMask(sizeA-1, bits)
	if bytesCompare(a.IP.Mask(newMask), a.IP) != 0 {
		return false
	}
	return bytesCompare(nextIP, b.IP) == 0
}

func mergeSubnets(a, b net.IPNet) net.IPNet {
	size, bits := a.Mask.Size()
	newMask := net.CIDRMask(size-1, bits)
	return net.IPNet{
		IP:   a.IP.Mask(newMask),
		Mask: newMask,
	}
}

func bytesCompare(a, b net.IP) int {
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

// Главная функция
func main() {
	var configPath string = ""
	// Определение флагов
	removeOnly := flag.Bool("d", false, "Только удаление маршрутов")
	addOnly := flag.Bool("s", false, "Только запрос и добавление маршрутов")
	displayOnly := flag.Bool("p", false, "Только запрос и отображение данных")
	flag.Parse()

	// Загрузка конфигурации
	value := os.Getenv("DEBUG")
	if value == "" {
		configPath = "/opt/routing/config.json"
	} else {
		configPath = "config.json"
	}
	config, err := lib.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		return
	}

	// Преобразование игнорируемых подсетей в map для быстрого доступа
	ignoredSubnets := make(map[string]bool)
	for _, subnet := range config.IgnoredSubnets {
		ignoredSubnets[subnet] = true
	}

	// Выполнение действий в зависимости от флагов
	switch {
	case *removeOnly:
		// Запускаем процесс обновления и применения маршрутов
		fmt.Println("Очистка старых маршрутов...")
		err = lib.RemoveRoutes(config.FilePath, config.Interface)
		if err != nil {
			fmt.Printf("Ошибка при удалении старых маршрутов: %v\n", err)
		}
	case *addOnly:
		fmt.Println("Запрос данных RIPE...")
		subnets, err := fetchSubnets(config)
		if err != nil {
			panic(err)
		}

		fmt.Println("Обновление файла подсетей...")
		err = lib.UpdateSubnetsFile(subnets, config.FilePath)
		if err != nil {
			fmt.Printf("Ошибка обновления файла: %v\n", err)
			os.Exit(0)
		}

		// Добавление новых маршрутов
		err = lib.AddRoutes(config.FilePath, config.Interface)
		if err != nil {
			fmt.Printf("Ошибка при добавлении новых маршрутов: %v\n", err)
		}
	case *displayOnly:
		fmt.Println("Запрос данных RIPE...")
		subnets, err := fetchSubnets(config)
		if err != nil {
			panic(err)
		}
		fmt.Println("Полученные подсети:")
		for _, subnet := range subnets {
			fmt.Println(subnet)
		}
	default:
		// Запускаем процесс обновления и применения маршрутов
		fmt.Println("Очистка старых маршрутов...")
		err = lib.RemoveRoutes(config.FilePath, config.Interface)
		if err != nil {
			fmt.Printf("Ошибка при удалении старых маршрутов: %v\n", err)
		}

		fmt.Println("Запрос данных RIPE...")
		subnets, err := fetchSubnets(config)
		if err != nil {
			panic(err)
		}

		fmt.Println("Обновление файла подсетей...")
		err = lib.UpdateSubnetsFile(subnets, config.FilePath)
		if err != nil {
			fmt.Printf("Ошибка обновления файла: %v\n", err)
			os.Exit(0)
		}

		// Добавление новых маршрутов
		err = lib.AddRoutes(config.FilePath, config.Interface)
		if err != nil {
			fmt.Printf("Ошибка при добавлении новых маршрутов: %v\n", err)
		}
		fmt.Println("Ожидание следующего обновления...")
	}
}
