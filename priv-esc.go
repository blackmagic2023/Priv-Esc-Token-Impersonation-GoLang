package main

import (
    "fmt"
    "golang.org/x/sys/windows"
    "os"
)

func ImpersonateUser() {
    var hToken windows.Token
    proc := windows.CurrentProcess()
    
    err := windows.OpenProcessToken(proc, windows.TOKEN_DUPLICATE|windows.TOKEN_QUERY, &hToken)
    if err != nil {
        fmt.Printf("OpenProcessToken failed. Error: %v\n", err)
        return
    }
    defer hToken.Close()

    var hImpersonatedToken windows.Token
    err = windows.DuplicateTokenEx(hToken, windows.MAXIMUM_ALLOWED, nil, windows.SecurityImpersonation, windows.TokenImpersonation, &hImpersonatedToken)
    if err != nil {
        fmt.Printf("DuplicateToken failed. Error: %v\n", err)
        return
    }
    defer hImpersonatedToken.Close()

    err = windows.ImpersonateLoggedOnUser(hImpersonatedToken)
    if err != nil {
        fmt.Printf("ImpersonateLoggedOnUser failed. Error: %v\n", err)
        return
    }
    defer windows.RevertToSelf()

    var username [256]uint16
    var size uint32 = 256
    err = windows.GetUserNameEx(windows.NameSamCompatible, &username[0], &size)
    if err != nil {
        fmt.Printf("GetUserNameEx failed. Error: %v\n", err)
    } else {
        fmt.Printf("Impersonated user: %s\n", windows.UTF16ToString(username[:]))
    }
}

func main() {
    ImpersonateUser()
}
