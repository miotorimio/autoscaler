package autoscaler

import (
	"errors"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func GetLocalMetrics() (float64, float64, error) {
	cpu, err := GetLocalCPULoad()
	if err != nil {
		return 0, 0, err
	}

	traffic, err := GetLocalTrafficRate()
	if err != nil {
		return 0, 0, err
	}

	return cpu, traffic, nil
}

func GetLocalCPULoad() (float64, error) {
	if runtime.GOOS != "windows" {
		return 0, errors.New("local monitoring поддерживается только на Windows")
	}

	cpu, err := getCPUWithWmic()
	if err == nil {
		return clamp(cpu, 0, 100), nil
	}

	cpu, err = getCPUWithPowerShell()
	if err != nil {
		return 0, err
	}

	return clamp(cpu, 0, 100), nil
}

func getCPUWithWmic() (float64, error) {
	cmd := exec.Command("cmd", "/C", "wmic cpu get loadpercentage /value")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	text := strings.TrimSpace(string(output))
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "LoadPercentage=") {
			value := strings.TrimPrefix(line, "LoadPercentage=")
			return strconv.ParseFloat(strings.TrimSpace(value), 64)
		}
	}

	return 0, errors.New("не удалось прочитать загрузку CPU через wmic")
}

func getCPUWithPowerShell() (float64, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "(Get-CimInstance Win32_Processor | Measure-Object -Property LoadPercentage -Average).Average")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	text := strings.TrimSpace(strings.ReplaceAll(string(output), "\r", ""))
	if text == "" {
		return 0, errors.New("пустой результат PowerShell для CPU")
	}

	return strconv.ParseFloat(text, 64)
}

func GetLocalTrafficRate() (float64, error) {
	if runtime.GOOS != "windows" {
		return 0, errors.New("local monitoring поддерживается только на Windows")
	}

	traffic, err := getTrafficWithPowerShell()
	if err == nil {
		return clamp(traffic, 0, 1000), nil
	}

	traffic, err = getTrafficWithNetstat()
	if err != nil {
		return 0, err
	}

	return clamp(traffic, 0, 1000), nil
}

func getTrafficWithPowerShell() (float64, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "(Get-NetTCPConnection | Measure-Object).Count")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	text := strings.TrimSpace(strings.ReplaceAll(string(output), "\r", ""))
	if text == "" {
		return 0, errors.New("пустой результат PowerShell для трафика")
	}

	return strconv.ParseFloat(text, 64)
}

func getTrafficWithNetstat() (float64, error) {
	cmd := exec.Command("cmd", "/C", "netstat -n")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.ReplaceAll(string(output), "\r", ""), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, "TCP") {
			count++
		}
	}

	return float64(count), nil
}
