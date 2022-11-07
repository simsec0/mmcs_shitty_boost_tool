package console

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

func Clear() {
	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout
	_ = c.Run()
	return
}

func SetTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, nil
	}
	defer syscall.FreeLibrary(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, nil
	}
	newTitle, _ := syscall.UTF16PtrFromString(title)
	r, _, err := syscall.SyscallN(proc, 1, uintptr(unsafe.Pointer(newTitle)), 0, 0)
	return int(r), err
}

func ShowBanner() {
	fmt.Println()

	for _, line := range BannerMulti {
		fmt.Println(strings.Repeat(" ", 42), line)
	}
	fmt.Println()
}

func PromptInput(username string) string {
	var input string

	fmt.Printf("%v\u001B[4m%v\u001B[0m\u001B[38;5;59m@\u001B[0m\u001B[4mboosts\u001B[0m:~$ ", strings.Repeat(" ", 15), username)
	fmt.Scanln(&input)

	return strings.ToLower(input)
}

func Print(content string) {
	fmt.Printf("%v%v", strings.Repeat(" ", 15), content)
}
