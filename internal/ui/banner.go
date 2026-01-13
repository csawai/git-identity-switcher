package ui

// Banner returns the ASCII art banner for gitx
func Banner() string {
	banner := `
 ██████╗ ██╗████████╗██╗  ██╗
 ██╔══██╗██║╚══██╔══╝╚██╗██╔╝
 ██████╔╝██║   ██║    ╚███╔╝ 
 ██╔══██╗██║   ██║    ██╔██╗ 
 ██████╔╝██║   ██║   ██╔╝ ██╗
 ╚═════╝ ╚═╝   ╚═╝   ╚═╝  ╚═╝
`
	return TitleStyle.Render(banner)
}

// SmallBanner returns a compact banner
func SmallBanner() string {
	return TitleStyle.Render("🔐 gitx")
}

// Subtitle returns a styled subtitle
func Subtitle(text string) string {
	return MutedText.Render(text)
}

// SectionDivider returns a styled section divider
func SectionDivider() string {
	return MutedText.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// Celebration returns a celebration message
func Celebration(message string) string {
	return SuccessBox.Render("🎉 " + message + " 🎊")
}

