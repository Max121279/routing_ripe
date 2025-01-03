package main

// Тест для ipRangeToCIDR
//func TestIpRangeToCIDR(t *testing.T) {
//	tests := []struct {
//		start    string
//		end      string
//		expected []string
//	}{
//		{"192.168.0.0", "192.168.0.255", []string{"192.168.0.0/24"}},
//		{"192.168.0.0", "192.168.1.255", []string{"192.168.0.0/23"}},
//		{"10.0.0.0", "10.0.0.15", []string{"10.0.0.0/28"}},
//		{"10.0.0.0", "10.0.0.0", []string{"10.0.0.0/32"}},
//	}
//
//	for _, test := range tests {
//		result, err := ipRangeToCIDR(test.start, test.end)
//		if err != nil {
//			t.Errorf("Ошибка ipRangeToCIDR(%s, %s): %v", test.start, test.end, err)
//		}
//		if !reflect.DeepEqual(result, test.expected) {
//			t.Errorf("ipRangeToCIDR(%s, %s) = %v; ожидается %v", test.start, test.end, result, test.expected)
//		}
//	}
//}
//
//// Тест для summarizeSubnets
//func TestSummarizeSubnets(t *testing.T) {
//	tests := []struct {
//		subnets  []string
//		expected []string
//	}{
//		{
//			subnets:  []string{"192.168.0.0/24", "192.168.1.0/24"},
//			expected: []string{"192.168.0.0/23"},
//		},
//		{
//			subnets:  []string{"10.0.0.0/28", "10.0.0.16/28"},
//			expected: []string{"10.0.0.0/27"}, // Ожидается объединение
//		},
//		{
//			subnets:  []string{"10.0.0.0/32", "10.0.0.1/32"},
//			expected: []string{"10.0.0.0/31"},
//		},
//		{
//			subnets:  []string{"192.168.0.0/24"},
//			expected: []string{"192.168.0.0/24"},
//		},
//	}
//
//	for _, test := range tests {
//		result := summarizeSubnets(test.subnets)
//		if !reflect.DeepEqual(result, test.expected) {
//			t.Errorf("summarizeSubnets(%v) = %v; ожидается %v", test.subnets, result, test.expected)
//		}
//	}
//}
//
//// Тест для filterSubnets
//func TestFilterSubnets(t *testing.T) {
//	tests := []struct {
//		cidr     string
//		ignored  map[string]bool
//		expected []string
//	}{
//		{
//			cidr: "192.168.0.0/30",
//			ignored: map[string]bool{
//				"192.168.0.1": true,
//				"192.168.0.3": true,
//			},
//			expected: []string{"192.168.0.0/32", "192.168.0.2/32"},
//		},
//		{
//			cidr: "10.0.0.0/29",
//			ignored: map[string]bool{
//				"10.0.0.1": true,
//				"10.0.0.5": true,
//			},
//			expected: []string{"10.0.0.0/32", "10.0.0.2/32", "10.0.0.3/32", "10.0.0.4/32", "10.0.0.6/32", "10.0.0.7/32"},
//		},
//		{
//			cidr: "192.168.1.0/29",
//			ignored: map[string]bool{
//				"192.168.1.2": true,
//				"192.168.1.6": true,
//			},
//			expected: []string{"192.168.1.0/32", "192.168.1.1/32", "192.168.1.3/32", "192.168.1.4/32", "192.168.1.5/32", "192.168.1.7/32"},
//		},
//		{
//			cidr: "10.0.0.0/8",
//			ignored: map[string]bool{
//				"10.134.1.1": true,
//			},
//			expected: []string{
//				"10.0.0.0/32", "10.0.0.1/32", "...", // Все адреса до 10.134.1.0
//				"10.134.1.0/32", "10.134.1.2/32", "...", // Все после 10.134.1.1
//			},
//		},
//	}
//
//	for _, test := range tests {
//		result, err := filterSubnets(test.cidr, test.ignored)
//		if err != nil {
//			t.Errorf("Ошибка filterSubnets(%s, %v): %v", test.cidr, test.ignored, err)
//		}
//		if !reflect.DeepEqual(result, test.expected) {
//			t.Errorf("filterSubnets(%s, %v) = %v; ожидается %v", test.cidr, test.ignored, result, test.expected)
//		}
//	}
//}
//
//func TestSummarizeSubnetsWithExclusions(t *testing.T) {
//	tests := []struct {
//		subnets  []string
//		excluded []string
//		expected []string
//	}{
//		{
//			subnets:  []string{"10.0.0.0/8"},
//			excluded: []string{"10.134.1.24"},
//			expected: []string{
//				"10.0.0.0/9", "10.128.0.0/11", "10.132.0.0/15", "10.134.0.0/24",
//				"10.134.1.0/32", "10.134.1.2/31", "10.134.1.4/22", "10.134.2.0/15",
//				"10.136.0.0/13", "10.144.0.0/12", "10.160.0.0/11", "10.192.0.0/10",
//			},
//		},
//	}
//
//	for _, test := range tests {
//		result, _ := summarizeSubnetsWithExclusions(test.subnets, test.excluded)
//		if !reflect.DeepEqual(result, test.expected) {
//			t.Errorf("summarizeSubnetsWithExclusions(%v, %v) = %v; ожидается %v",
//				test.subnets, test.excluded, result, test.expected)
//		}
//	}
//}
