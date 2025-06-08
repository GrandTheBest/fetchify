package main

import (
	"fetchify/stdc"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

var ignoredMounts = []string{"/boot", "/efi", "/mnt", "/run", "/var", "/tmp", "/dev", "/sys", "/proc"}

const (
	Reset  = "\033[0m"
	Black  = "\033[30m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Purple = "\033[35m"
	Blue   = "\033[34m"
	White  = "\033[37m"
)

var koala = []string{
	"==== Fetchify ====",
	"   __         __  ",
	" /\"  \"\\     /\"  \"\\",
	"(  (\\  )___(  /)  )",
	" \\               / ",
	" /               \\ ",
	"/    () ___ ()    \\",
	"|      (   )      |",
	" \\      \\_/      / ",
	"   \\...__!__.../   ",
	"        \" mh       ",
	"   v1.0.4-generic   ",
}

func main() {
	helpFlag := flag.Bool("help", false, "Show help message")
	dontclearFlag := flag.Bool("dont-clear", false, "Don't clear screen before print fetch")

	flag.Parse()

	if *helpFlag {
		fmt.Println("Use:")
		fmt.Println("  --help   Show help message")
		fmt.Println("  --dont-clear  Don't clear screen before print fetch")

		return
	}

	if !*dontclearFlag {
		CMD := exec.Command("clear")
		CMD.Stdout = os.Stdout
		CMD.Run()
	}

	// Get info
	user, _ := user.Current()
	hostname, _ := os.Hostname()
	info, _ := host.Info()
	uptime := time.Duration(info.Uptime) * time.Second
	cpuInfo, _ := cpu.Info()

	_hostline := fmt.Sprintf("%s@%s", user.Username, hostname)
	var _separator string
	var separator string
	var hostline string
	var chars int
	os := info.Platform

	if info.Platform == "arch" {
		os = "i use arch, btw"
	}

	if user.Username == "root" {
		chars = stdc.CharsCount(_hostline) + 26
	} else {
		chars = stdc.CharsCount(_hostline)
	}

	for i := range chars {
		arr := make([]string, chars)

		(arr)[i] += "─"

		stdc.ArrayToString(&_separator, arr)
	}

	if user.Username == "root" {
		hostline = fmt.Sprintf("%s%s%s@%s%s%s - %sDon't login as root ^-^%s", Red, user.Username, Green, Cyan, hostname, Reset, Red, Reset)
		separator = fmt.Sprintf("%s%s%s", Red, _separator, Reset)
	} else {
		hostline = fmt.Sprintf("%s%s%s@%s%s%s", Red, user.Username, Green, Cyan, hostname, Reset)
		separator = fmt.Sprintf("%s%s%s", Green, _separator, Reset)
	}

	// Info for print
	infoLines := []string{
		hostline,
		separator,
		fmt.Sprintf("%sOS:%s %s %s", Yellow, Reset, os, info.PlatformVersion),
		fmt.Sprintf("%sKernel:%s %s", Yellow, Reset, info.KernelVersion),
		fmt.Sprintf("%sUptime:%s %s", Yellow, Reset, uptime),
		fmt.Sprintf("%sShell:%s %s", Yellow, Reset, getShell()),
		fmt.Sprintf("%sWM:%s %s", Yellow, Reset, getWM()),
		fmt.Sprintf("%sCPU:%s %s (%d cores)", Yellow, Reset, cpuInfo[0].ModelName, runtime.NumCPU()),
		fmt.Sprintf("%sGPU:%s %s", Yellow, Reset, getGPU()),
		fmt.Sprint(getRam()),
		fmt.Sprint(getSwap()),
		fmt.Sprint(getDisks()),
		fmt.Sprintf("%s███%s███%s███%s███%s███%s███%s███%s███", Black, Red, Green, Yellow, Blue, Purple, Cyan, White),
		fmt.Sprintf("%s███%s███%s███%s███%s███%s███%s███%s███\n", Black, Red, Green, Yellow, Blue, Purple, Cyan, White),
	}

	// ASCII art length
	maxArtWidth := 0
	for _, line := range koala {
		if len(line) > maxArtWidth {
			maxArtWidth = len(line)
		}
	}

	// Print ASCII art and info
	maxLines := len(koala)
	if len(infoLines) > maxLines {
		maxLines = len(infoLines)
	}

	for i := 0; i < maxLines; i++ {
		art := " "
		if i < len(koala) {
			art = koala[i]
		}

		text := ""
		if i < len(infoLines) {
			text = infoLines[i]
		}

		fmt.Printf("%-*s %s\n", maxArtWidth, art, text)
	}
}

// Get GPU info
func getGPU() string {
	cmd := exec.Command("sh", "-c", "lspci | grep VGA | cut -d ':' -f3")
	out, err := cmd.Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}
	return fmt.Sprintf("%sUnknown%s", Red, Reset)
}

// Get Shell info
func getShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return fmt.Sprintf("%sUnknown%s", Red, Reset)
	}
	parts := strings.Split(shell, "/")
	return parts[len(parts)-1]
}

// Get WM and DM
func getWM() string {
	var wm string

	cmd := exec.Command("sh", "-c", "xprop -root _NET_WM_NAME")
	out, err := cmd.Output()
	if err == nil && len(out) > 0 {
		parts := strings.Split(string(out), "=")
		if len(parts) > 1 {
			wm = strings.Trim(parts[1], " \"\n")
		}
	}

	if wm == "" {
		cmd = exec.Command("sh", "-c", "wmctrl -m | grep Name")
		out, err = cmd.Output()
		if err == nil && len(out) > 0 {
			parts := strings.Fields(string(out))
			if len(parts) > 1 {
				wm = parts[1]
			}
		}
	}

	if wm == "" {
		cmd = exec.Command("sh", "-c", "echo $XDG_CURRENT_DESKTOP")
		out, err = cmd.Output()
		if err == nil && len(out) > 0 {
			wm = strings.TrimSpace(string(out))
		}
	}

	if wm == "" {
		return fmt.Sprintf("%sUnknown%s", Red, Reset)
	}

	displayServer := "X11"
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		displayServer = "Wayland"
	}

	return fmt.Sprintf("%s (%s)", wm, displayServer)
}

// Check mountpoint in "Ignored mounts"
func isIgnoredMount(mount string) bool {
	for _, ignored := range ignoredMounts {
		if strings.HasPrefix(mount, ignored) {
			return true
		}
	}
	return false
}

// Get RAM and RAM type
func getRam() string {
	user, _ := user.Current()

	ram, err := mem.VirtualMemory()

	if err != nil || ram.Total == 0 {
		return fmt.Sprintf("%sRAM%s: %sError%s", Yellow, Reset, Red, Reset)
	}

	usedGiB := float64(ram.Used) / (1024 * 1024 * 1024)
	totalGiB := float64(ram.Total) / (1024 * 1024 * 1024)
	percent := ram.UsedPercent

	if user.Username == "root" {
		cmd := exec.Command("sudo", "dmidecode", "-t", "memory")
		output, err := cmd.CombinedOutput()

		if err != nil {
			return fmt.Sprintf("%sRAM%s: %sError%s", Yellow, Reset, Red, Reset)
		}

		seen := map[string]bool{}
		types := []string{}
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Type:") &&
				!strings.Contains(line, "Unknown") &&
				!strings.Contains(line, "Other") &&
				!strings.Contains(line, "RAM") {

				if !seen[line] {
					seen[line] = true
					typePart := strings.TrimPrefix(line, "Type: ")
					types = append(types, typePart)
				}
			}
		}

		ramType := "Unknown"
		if len(types) > 0 {
			ramType = strings.Join(types, ", ")
		}

		return fmt.Sprintf("%sRAM%s: %s%.2f GiB%s / %s%.2f Gib%s (%.0f%%) - %s%s%s",
			Yellow, Reset, Green, usedGiB,
			Reset, Cyan, totalGiB, Reset,
			percent, Purple, ramType, Reset)
	}

	return fmt.Sprintf("%sRAM%s: %s%.2f GiB%s / %s%.2f Gib%s (%.0f%%)%s",
		Yellow, Reset, Green, usedGiB,
		Reset, Cyan, totalGiB, Reset,
		percent, Reset)
}

// Get swap file info
func getSwap() string {
	swap, err := mem.SwapMemory()

	if err != nil {
		return fmt.Sprintf("%sSwap%s: %sError%s", Yellow, Reset, Red, Reset)
	}

	if swap.Total == 0 {
		return fmt.Sprintf("%sSwap%s: %sDisabled%s", Yellow, Reset, Red, Reset)
	}

	usedGiB := float64(swap.Used) / (1024 * 1024 * 1024)
	totalGiB := float64(swap.Total) / (1024 * 1024 * 1024)
	percent := swap.UsedPercent

	return fmt.Sprintf("%sSwap%s: %s%.2f GiB%s / %s%.2f Gib%s (%.0f%%)",
		Yellow, Reset, Green, usedGiB,
		Reset, Cyan, totalGiB, Reset,
		percent)
}

// Get info about disks
func getDisks() string {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return fmt.Sprintf("%sUnknown%s", Red, Reset)
	}

	var result strings.Builder
	for _, part := range partitions {
		if isIgnoredMount(part.Mountpoint) {
			continue
		}
		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			continue
		}
		usedGB := float64(usage.Used) / 1e9
		totalGB := float64(usage.Total) / 1e9
		percentUsed := usage.UsedPercent

		result.WriteString(fmt.Sprintf(
			"%sDisk%s (%s): %s%.2f GiB%s / %s%.2f GiB%s (%.0f%%) - %s%s%s\n",
			Yellow, Reset, part.Mountpoint,
			Green, usedGB, Reset,
			Cyan, totalGB, Reset,
			percentUsed,
			Purple, part.Fstype, Reset))
	}

	return result.String()
}
