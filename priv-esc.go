package main

import (
    "fmt"
    "golang.org/x/sys/windows"
    "os/exec"
    "syscall"
)

var (
    advapi32                    = syscall.NewLazyDLL("advapi32.dll")
    procImpersonateLoggedOnUser = advapi32.NewProc("ImpersonateLoggedOnUser")
)

func ImpersonateAndRunCalc() {
    var hToken windows.Token
    proc := windows.CurrentProcess()

    err := windows.OpenProcessToken(proc, windows.TOKEN_DUPLICATE|windows.TOKEN_QUERY, &hToken)
    if err != nil {
        fmt.Printf("OpenProcessToken failed. Error: %v\n", err)
        return
    }
    defer hToken.Close()

    var impersonatedToken windows.Token
    err = windows.DuplicateTokenEx(hToken, windows.MAXIMUM_ALLOWED, nil, windows.SecurityImpersonation, windows.TokenImpersonation, &impersonatedToken)
    if err != nil {
        fmt.Printf("DuplicateToken failed. Error: %v\n", err)
        return
    }
    defer impersonatedToken.Close()

    // Call ImpersonateLoggedOnUser via syscall
    r1, _, err := procImpersonateLoggedOnUser.Call(uintptr(impersonatedToken))
    if r1 == 0 {
        fmt.Printf("ImpersonateLoggedOnUser failed. Error: %v\n", err)
        return
    }
    defer windows.RevertToSelf()

    cmd := exec.Command("calc.exe")
    err = cmd.Run()
    if err != nil {
        fmt.Printf("Failed to start calc.exe: %v\n", err)
    }
}

func main() {
    ImpersonateAndRunCalc()
}
