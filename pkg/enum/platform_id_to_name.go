package enum

const (
	// Platform ID.
	IOSPlatformID        = 1
	AndroidPlatformID    = 2
	WindowsPlatformID    = 3
	OSXPlatformID        = 4
	WebPlatformID        = 5
	MiniWebPlatformID    = 6
	LinuxPlatformID      = 7
	AndroidPadPlatformID = 8
	IPadPlatformID       = 9
	AdminPlatformID      = 10

	// Platform string match to Platform ID.
	IOSPlatformStr        = "IOS"
	AndroidPlatformStr    = "Android"
	WindowsPlatformStr    = "Windows"
	OSXPlatformStr        = "OSX"
	WebPlatformStr        = "Web"
	MiniWebPlatformStr    = "MiniWeb"
	LinuxPlatformStr      = "Linux"
	AndroidPadPlatformStr = "APad"
	IPadPlatformStr       = "IPad"
	AdminPlatformStr      = "Admin"

	// terminal types.
	TerminalPC     = "PC"
	TerminalMobile = "Mobile"
	TerminalPad    = "Pad"
)

var PlatformID2Name = map[int]string{
	IOSPlatformID:        IOSPlatformStr,
	AndroidPlatformID:    AndroidPlatformStr,
	WindowsPlatformID:    WindowsPlatformStr,
	OSXPlatformID:        OSXPlatformStr,
	WebPlatformID:        WebPlatformStr,
	MiniWebPlatformID:    MiniWebPlatformStr,
	LinuxPlatformID:      LinuxPlatformStr,
	AndroidPadPlatformID: AndroidPadPlatformStr,
	IPadPlatformID:       IPadPlatformStr,
	AdminPlatformID:      AdminPlatformStr,
}

var PlatformName2ID = map[string]int{
	IOSPlatformStr:        IOSPlatformID,
	AndroidPlatformStr:    AndroidPlatformID,
	WindowsPlatformStr:    WindowsPlatformID,
	OSXPlatformStr:        OSXPlatformID,
	WebPlatformStr:        WebPlatformID,
	MiniWebPlatformStr:    MiniWebPlatformID,
	LinuxPlatformStr:      LinuxPlatformID,
	AndroidPadPlatformStr: AndroidPadPlatformID,
	IPadPlatformStr:       IPadPlatformID,
	AdminPlatformStr:      AdminPlatformID,
}

var PlatformName2class = map[string]string{
	IOSPlatformStr:        TerminalMobile,
	AndroidPlatformStr:    TerminalMobile,
	MiniWebPlatformStr:    MiniWebPlatformStr,
	WebPlatformStr:        WebPlatformStr,
	WindowsPlatformStr:    TerminalPC,
	OSXPlatformStr:        TerminalPC,
	LinuxPlatformStr:      TerminalPC,
	AndroidPadPlatformStr: TerminalPad,
	IPadPlatformStr:       TerminalPad,
	AdminPlatformStr:      AdminPlatformStr,
}

var PlatformID2class = map[int]string{
	IOSPlatformID:        TerminalMobile,
	AndroidPlatformID:    TerminalMobile,
	MiniWebPlatformID:    MiniWebPlatformStr,
	WebPlatformID:        WebPlatformStr,
	WindowsPlatformID:    TerminalPC,
	OSXPlatformID:        TerminalPC,
	LinuxPlatformID:      TerminalPC,
	AndroidPadPlatformID: TerminalPad,
	IPadPlatformID:       TerminalPad,
	AdminPlatformID:      AdminPlatformStr,
}

func PlatformIDToName(num int) string {
	return PlatformID2Name[num]
}

func PlatformNameToID(name string) int {
	return PlatformName2ID[name]
}

func PlatformNameToClass(name string) string {
	return PlatformName2class[name]
}

func PlatformIDToClass(num int) string {
	return PlatformID2class[num]
}
