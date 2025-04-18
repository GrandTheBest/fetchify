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
	"   v1.0.3-generic   ",
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
	vmStat, _ := mem.VirtualMemory()

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
		fmt.Sprintf("%sRAM:%s %.2f GB / %.2f GB", Yellow, Reset, float64(vmStat.Used)/1e9, float64(vmStat.Total)/1e9),
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

// Get swap file info
func getSwap() string {
	swap, err := mem.SwapMemory()

	if err != nil {
		return fmt.Sprintf("%sSwap%s: %sError%s", Yellow, Reset, Red, Reset)
	}

	if swap.Total == 0 {
		return fmt.Sprintf("%sSwap%s: %sDisabled%s", Yellow, Reset, Red, Reset)
	}

	// value_used := swap.Used / 10000000
	// str_used := strconv.FormatUint(value_used, 10)

	// if len(str_used) > 1 {
	// 	str_used = str_used[:1] + "." + str_used[1:]
	// } else {
	// 	str_used = str_used + ".0"
	// }

	// value_total := swap.Total / 10000000
	// str_total := strconv.FormatUint(value_total, 10)

	// if len(str_total) > 1 {
	// 	str_total = str_total[:1] + "." + str_total[1:]
	// } else {
	// 	str_total = str_total + ".0"
	// }

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
