package main

import (
	"fmt"
	"github.com/Max121279/routing_ripe/src/lib"
	"reflect"
	"testing"
)

// Тест для ipRangeToCIDR
func TestIpRangeToCIDR(t *testing.T) {
	tests := []struct {
		start    string
		end      string
		expected []string
	}{
		{"192.168.0.0", "192.168.0.255", []string{"192.168.0.0/24"}},
		{"192.168.0.0", "192.168.1.255", []string{"192.168.0.0/23"}},
		{"10.0.0.0", "10.0.0.15", []string{"10.0.0.0/28"}},
		{"10.0.0.0", "10.0.0.0", []string{"10.0.0.0/32"}},
	}

	for _, test := range tests {
		result, err := ipRangeToCIDR(test.start, test.end)
		if err != nil {
			t.Errorf("Ошибка ipRangeToCIDR(%s, %s): %v", test.start, test.end, err)
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ipRangeToCIDR(%s, %s) = %v; ожидается %v", test.start, test.end, result, test.expected)
		}
	}
}

// Тест для summarizeSubnets
func TestSummarizeSubnets(t *testing.T) {
	tests := []struct {
		subnets  []string
		expected []string
	}{
		{
			subnets:  []string{"192.168.0.0/24", "192.168.1.0/24"},
			expected: []string{"192.168.0.0/23"},
		},
		{
			subnets:  []string{"10.0.0.0/28", "10.0.0.16/28"},
			expected: []string{"10.0.0.0/27"}, // Ожидается объединение
		},
		{
			subnets: []string{"10.0.0.0/32", "10.0.0.1/32", "10.0.0.2/32", "10.0.0.3/32", "10.0.0.4/32"},
			//subnets:  []string{"10.0.0.1/32", "10.0.0.2/32", "10.0.0.3/32", "10.0.0.4/32"},
			expected: []string{"10.0.0.0/30", "10.0.0.4/32"},
		},
		{
			subnets:  []string{"192.168.0.0/24", "192.168.5.0/24"},
			expected: []string{"192.168.0.0/24", "192.168.5.0/24"},
		},
	}

	for _, test := range tests {
		result := summarizeSubnets(test.subnets)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("summarizeSubnets(%v) = %v; ожидается %v", test.subnets, result, test.expected)
		}
	}
}

// Тест для filterSubnets
func TestFilterSubnets(t *testing.T) {
	tests := []struct {
		cidr     string
		ignored  map[string]bool
		expected []string
	}{
		{
			cidr: "192.168.0.0/30",
			ignored: map[string]bool{
				"192.168.0.1": true,
				"192.168.0.3": true,
			},
			expected: []string{"192.168.0.0/32", "192.168.0.2/32"},
		},
		{
			cidr: "10.0.0.0/29",
			ignored: map[string]bool{
				"10.0.0.1": true,
				"10.0.0.5": true,
			},
			expected: []string{"10.0.0.0/32", "10.0.0.2/32", "10.0.0.3/32", "10.0.0.4/32", "10.0.0.6/32", "10.0.0.7/32"},
		},
		{
			cidr: "192.168.1.0/29",
			ignored: map[string]bool{
				"192.168.1.2": true,
				"192.168.1.6": true,
			},
			expected: []string{"192.168.1.0/32", "192.168.1.1/32", "192.168.1.3/32", "192.168.1.4/32", "192.168.1.5/32", "192.168.1.7/32"},
		},
	}

	for _, test := range tests {
		result, err := filterSubnets(test.cidr, test.ignored)
		if err != nil {
			t.Errorf("Ошибка filterSubnets(%s, %v): %v", test.cidr, test.ignored, err)
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("filterSubnets(%s, %v) = %v; ожидается %v", test.cidr, test.ignored, result, test.expected)
		}
	}
}

func TestSummarizeSubnetsWithExclusions(t *testing.T) {
	tests := []struct {
		subnets  []string
		excluded map[string]bool
		expected []string
	}{
		{
			subnets:  []string{"10.0.0.0/8"},
			excluded: map[string]bool{"10.134.1.24": true, "10.134.1.25": true},
			expected: []string{
				"10.0.0.0/9", "10.128.0.0/14", "10.132.0.0/15", "10.134.0.0/24", "10.134.1.0/28",
				"10.134.1.16/29", "10.134.1.26/31", "10.134.1.28/30", "10.134.1.32/27", "10.134.1.64/26",
				"10.134.1.128/25", "10.134.2.0/23", "10.134.4.0/22", "10.134.8.0/21", "10.134.16.0/20",
				"10.134.32.0/19", "10.134.64.0/18", "10.134.128.0/17", "10.135.0.0/16", "10.136.0.0/13",
				"10.144.0.0/12", "10.160.0.0/11", "10.192.0.0/10",
			},
		},
	}

	for _, test := range tests {
		result, _ := summarizeSubnetsWithExclusions(test.subnets, test.excluded)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("summarizeSubnetsWithExclusions(%v, %v) = %v; ожидается %v",
				test.subnets, test.excluded, result, test.expected)
		}
	}
}

func TestSummarizeSubnetsWithExclusions2(t *testing.T) {
	tests := []struct {
		subnet   string
		excluded map[string]bool
		expected []string
	}{
		{
			subnet:   "10.0.0.0/8",
			excluded: map[string]bool{"10.134.1.24": true, "10.134.1.25": true},
			expected: []string{
				"10.0.0.0/9", "10.128.0.0/14", "10.132.0.0/15", "10.134.0.0/24", "10.134.1.0/28",
				"10.134.1.16/29", "10.134.1.26/31", "10.134.1.28/30", "10.134.1.32/27", "10.134.1.64/26",
				"10.134.1.128/25", "10.134.2.0/23", "10.134.4.0/22", "10.134.8.0/21", "10.134.16.0/20",
				"10.134.32.0/19", "10.134.64.0/18", "10.134.128.0/17", "10.135.0.0/16", "10.136.0.0/13",
				"10.144.0.0/12", "10.160.0.0/11", "10.192.0.0/10",
			},
		},
	}

	for _, test := range tests {
		result, _ := lib.SummarizeSubnetsWithExclusions(test.subnet, test.excluded)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("summarizeSubnetsWithExclusions(%v, %v) = %v; ожидается %v",
				test.subnet, test.excluded, result, test.expected)
		}
	}
}

func TestSummarizeSubnetsWithExclusions3(t *testing.T) {
	tests := []struct {
		subnet         string
		excluded       map[string]bool
		expected       []string
		excludedDomain []string
	}{
		{
			subnet:   "80.0.0.0/8",
			excluded: map[string]bool{"80.134.1.24": true, "80.134.1.25": true},
			expected: []string{
				"80.0.0.0/10", "80.64.0.0/13", "80.72.0.0/14", "80.76.0.0/16", "80.77.0.0/17", "80.77.128.0/19",
				"80.77.160.0/21", "80.77.168.0/27", "80.77.168.32/29", "80.77.168.40/30", "80.77.168.45/32",
				"80.77.168.46/31", "80.77.168.48/28", "80.77.168.64/26", "80.77.168.128/25", "80.77.169.0/24",
				"80.77.170.0/23", "80.77.172.0/22", "80.77.176.0/20", "80.77.192.0/18", "80.78.0.0/15", "80.80.0.0/12",
				"80.96.0.0/11", "80.128.0.0/14", "80.132.0.0/15", "80.134.0.0/24", "80.134.1.0/28", "80.134.1.16/29",
				"80.134.1.26/31", "80.134.1.28/30", "80.134.1.32/27", "80.134.1.64/26", "80.134.1.128/25", "80.134.2.0/23",
				"80.134.4.0/22", "80.134.8.0/21", "80.134.16.0/20", "80.134.32.0/19", "80.134.64.0/18", "80.134.128.0/17",
				"80.135.0.0/16", "80.136.0.0/13", "80.144.0.0/12", "80.160.0.0/11", "80.192.0.0/10",
			},
			excludedDomain: []string{"dlcache4.vibio.tv"},
		},
	}

	for _, test := range tests {
		for _, dom := range test.excludedDomain {
			ips, _ := lib.GetHostIPs(dom)
			for _, ip := range ips {
				test.excluded[ip] = true
			}
		}

		result, _ := lib.SummarizeSubnetsWithExclusions(test.subnet, test.excluded)

		fmt.Println("")
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("summarizeSubnetsWithExclusions(%v, %v) = %v; ожидается %v",
				test.subnet, test.excluded, result, test.expected)
		}
	}
}
