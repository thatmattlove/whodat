package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"

	q "github.com/checktheroads/whodat/query"
	"github.com/gookit/color"
	"github.com/mkideal/cli"
)

type argT struct {
	Help        bool `cli:"!h,help" usage:"Show this Help Menu" json:"-"`
	GetPrefixes bool `cli:"p,prefixes" usage:"Get ASN's Advertised Prefixes" dft:"false"`
}

var supportsColor bool = true

func (argv *argT) AutoHelp() bool {
	return argv.Help
}

func detectColorSupport() bool {
	out, err := exec.Command("tput", "colors").CombinedOutput()
	if err != nil {
		return false
	}
	r := string(out)
	return len(r) > 0
}
func init() {
	if !detectColorSupport() {
		supportsColor = false
	}
}

func parseAsn(raw string) (p string, e error) {
	asnPattern := regexp.MustCompile("^[0-9]+$")
	match := asnPattern.FindStringSubmatch(raw)

	if len(match) < 1 {
		prefixPattern := regexp.MustCompile("^AS([0-9]+)$")
		match = prefixPattern.FindStringSubmatch(raw)

		if len(match) < 2 {
			return "", errors.New("Failed to parse ASN")
		}
		p = match[1]
	} else {
		p = match[0]
	}

	return p, nil
}

func white(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgWhite}.Render(a...)
}

func whiteBold(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgWhite, color.Bold}.Render(a...)
}

func blue(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgBlue}.Render(a...)
}

func blueBold(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgBlue, color.Bold}.Render(a...)
}

func red(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgRed}.Render(a...)
}

func redBold(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgRed, color.Bold}.Render(a...)
}

func green(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgGreen}.Render(a...)
}

func greenBold(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgGreen, color.Bold}.Render(a...)
}

func magenta(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgMagenta}.Render(a...)
}

func magentaBold(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgMagenta, color.Bold}.Render(a...)
}

func cyan(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgCyan}.Render(a...)
}

func cyanBold(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgCyan, color.Bold}.Render(a...)
}

func yellow(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgYellow}.Render(a...)
}

func yellowBold(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgYellow, color.Bold}.Render(a...)
}

func gray(a ...interface{}) string {
	if !supportsColor {
		return fmt.Sprint(a...)
	}
	return color.Style{color.FgGray}.Render(a...)
}

func handleASNPrefixes(ctx *cli.Context, asn string) {
	prefixes, err := q.GetASNPrefixes(asn)

	if err != nil {
		ctx.String(redBold(err.Error()))
		os.Exit(1)
	}

	ctx.String("    %s\n", whiteBold("Prefixes:"))
	ctx.String("      %s\n", whiteBold("IPv4:"))

	for _, p := range prefixes.IPv4 {
		ctx.String("        %s\n", yellow(p))
	}

	ctx.String("      %s\n", whiteBold("IPv6:"))
	for _, p := range prefixes.IPv6 {
		ctx.String("        %s\n", cyan(p))
	}
}

func handleASN(ctx *cli.Context, asn string) {
	asnDetail, err := q.GetASNDetail(asn)

	if err != nil {
		ctx.String(redBold(err.Error()))
		os.Exit(1)
	}

	ctx.String("\n")
	ctx.String("  %s (%s)\n", blueBold(asnDetail.ASN), yellowBold(asnDetail.Org))
	ctx.String("\n")

	if asnDetail.LookingGlass != "" {
		ctx.String("    %s\n", whiteBold("Looking Glass: ")+green(asnDetail.LookingGlass))
	}

	if asnDetail.Website != "" {
		ctx.String("    %s\n", whiteBold("Website: ")+magenta(asnDetail.Website))
	}
}

func handleIP(ctx *cli.Context, ip string) {
	ipDetail, err := q.GetIPDetail(ip)

	if err != nil {
		ctx.String(redBold(err.Error()))
		os.Exit(1)
	}

	ctx.String("\n")
	ctx.String("  %s (%s)\n\n", magentaBold(ipDetail.IP), green(ipDetail.PTR))
	ctx.String("    %s\n", yellowBold(ipDetail.Org))
	ctx.String("    %s%s\n", whiteBold("AS"), cyanBold(ipDetail.ASN))
	ctx.String("\n")
	ctx.String("    %s %s (%s)\n", gray("Prefix:"), blueBold(ipDetail.Prefix), gray(ipDetail.Name))
	ctx.String("    %s %s\n", gray("RIR:"), redBold(ipDetail.RIR))
}

func handlePrefix(ctx *cli.Context, pfx string) {
	prefixDetail, err := q.GetPrefixDetail(pfx)

	if err != nil {
		ctx.String(redBold(err.Error()))
		os.Exit(1)
	}

	ctx.String("\n")
	ctx.String("  %s (%s)\n\n", magentaBold(prefixDetail.Prefix), green(prefixDetail.Name))
	ctx.String("    %s\n\n", yellowBold(prefixDetail.Org))
	ctx.String("    %s\n", whiteBold("Origins:"))

	for _, o := range prefixDetail.Origins {
		ctx.String("      %s (%s)\n", cyanBold(o.ASN), blue(o.Org))
	}

	ctx.String("\n")
	ctx.String("    %s\n", gray("RIR: ")+redBold(prefixDetail.RIR))
}

func main() {

	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		args := os.Args[1:]

		_, nw, err := net.ParseCIDR(args[0])

		if err != nil {
			ip := net.ParseIP(args[0])

			if ip != nil {
				handleIP(ctx, ip.String())
			} else {
				asn, err := parseAsn(args[0])

				if err != nil {
					ctx.String("%s '%s' %s", red("Failed to parse"), redBold(args[0]), red("as IP, Prefix, or ASN"))
					os.Exit(1)
				}

				handleASN(ctx, asn)

				if argv.GetPrefixes {
					handleASNPrefixes(ctx, asn)
				}

			}
		} else {
			handlePrefix(ctx, nw.String())
		}
		return nil
	}, fmt.Sprintf("\n%s\n  %s", magentaBold("whodat"), gray("Quickly get IP, Prefix, and ASN Information at the command-line."))))
}
