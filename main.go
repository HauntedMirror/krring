package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/jessevdk/go-flags"
	probing "github.com/prometheus-community/pro-bing"
)

// nolint:gochecknoglobals
var krrg = []string{
	`                                                                `,
	`                                          AA                    `,
	`                                    AA  AABBAAAA                `,
	`          AAAAAAAA        AAAAAAAAAABBAABBBBBBBBAAAA            `,
	`        AACCCCDDDDAA  AAAACCCCCCEEFFGGBBBBBBBBBBBBEEAA          `,
	`      AACCCCCCDDDDDDAACCCCCCCCCCEEEEFFGGBBBBHHBBBBEEEEAA        `,
	`      AACCCCCCDDDDCCCCCCCCCCCCCCEEEEEEFFGGHHHHHHFFEEEEAA        `,
	`      AACCCCDDDDDDCCCCCCCCCCCCCCEEEEEEEEFFGGHHHHEEFFEEAA        `,
	`        AADDDDDDCCCCCCCCCCCCCCCCEEEEEEEEEEFFDDEEFFEEAA          `,
	`          AADDDDCCCCCCCCCCCCCCCCEEEEEEEEEEDDDDDDEEEEEEAA        `,
	`        AAIIIIIICCCCCCCCCCCCCCCCEEEEEEEEDDEEDDEEDDEEEEAA        `,
	`        AAIIIIAACCCCCCCCCCCCCCCCEEEEEEEEEEEEDDEEAAEEEEAA        `,
	`        AAIIIIAACCCCCCCCCCCCCCCCEEEEEEEEEEEEDDEEAAEEEEAA        `,
	`        AAIIIIAACCCCCCJJJJJJKKKKKKKKJJJJJJEEEEEEAAEEEEAA        `,
	`        AAIIIIAACCCCCCLLMMLLKKKKKKKKLLMMLLEEEEEEAAEEEEAA        `,
	`        AAIIIIAACCCCCCLLNNNNKKKKKKKKNNNNLLEEEEEEAAEEEEAA        `,
	`        AAIIIIAACCCCCCLLOOOOKKKKKKKKOOOOLLEEEEEEAAEEEEAA        `,
	`        AAIIIIAAAAAAPPQQRRKKKKKKKKKKKKRRQQPPEEEEAAEEEEAA        `,
	`      AAIIIIAA      AACCQQQQQQQQQQQQQQQQEEAA      AAEEEEAA      `,
	`    AAIIIIIIAA      AACCAAAAAASSSSAAAAAAEEAA      AAEEEEAA      `,
	`  AAIIIIIIAA        AACCAAKKKKKKKKKKKKAAEEAA        AAEEEEAA    `,
	`  AAIIIIAA            AAAATTTTSSSSTTTTAAAA          AAEEEEAA    `,
	`AAIIIIAA            AAAAUUSSUUTTTTUUSSUUAAAA          AAEEEEAA  `,
	`AAIIIIAA          AAUUUUAAAASSTTTTSSAAAAUUUUAA        AAEEEEAA  `,
	`AAIIIIAA        AAKKKKUUAAUUUULLLLUUUUAAUUKKKKAA      AAEEEEAA  `,
	`  AAIIIIAA        AAAAAAUUUUUUSSSSUUUUUUAAAAAA      AAEEEEAA    `,
	`    AAIIIIAA        AAUUUUUUTTTTTTTTUUVVUUAA      AAEEEEAA      `,
	`      AAAAIIAA        AAAASSKKAAAAKKSSAAAA      AAEEAAAA        `,
	`          AA          AASSAAKKAAAAKKAASSAA        AA            `,
	`                        AAAAUUAAAAUUAAAA                        `,
	`                          AAUUAAAAUUAA                          `,
	`                          AAAAAAAAAAAA                          `,
}

// nolint:gochecknoglobals
var (
	appName        = "krring"
	appUsage       = "[OPTIONS] HOST"
	appDescription = "`ping` command but with krrg"
	appVersion     = "???"
	appRevision    = "???"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrPing
)

type options struct {
	Count     int  `short:"c" long:"count" default:"32" description:"Stop after <count> replies"`
	Privilege bool `short:"P" long:"privilege" description:"Enable privileged mode"`
	Version   bool `short:"V" long:"version" description:"Show version"`
}

func main() {
	code, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"[ %v ] %s\n",
			color.New(color.FgRed, color.Bold).Sprint("ERROR"),
			err,
		)
	}

	os.Exit(int(code))
}

func run(cliArgs []string) (exitCode, error) {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = appName
	parser.Usage = appUsage
	parser.ShortDescription = appDescription
	parser.LongDescription = appDescription

	args, err := parser.ParseArgs(cliArgs)
	if err != nil {
		if flags.WroteHelp(err) {
			return exitCodeOK, nil
		}

		return exitCodeErrArgs, fmt.Errorf("parse error: %w", err)
	}

	if opts.Version {
		// nolint:forbidigo
		fmt.Printf("%s: v%s-rev%s\n", appName, appVersion, appRevision)

		return exitCodeOK, nil
	}

	if len(args) == 0 {
		// nolint:goerr113
		return exitCodeErrArgs, errors.New("must requires an argument")
	}

	if 1 < len(args) {
		// nolint:goerr113
		return exitCodeErrArgs, errors.New("too many arguments")
	}

	pinger, err := initPinger(args[0], opts)
	if err != nil {
		return exitCodeOK, fmt.Errorf("an error occurred while initializing pinger: %w", err)
	}

	if err := pinger.Run(); err != nil {
		return exitCodeErrPing, fmt.Errorf("an error occurred when running ping: %w", err)
	}

	return exitCodeOK, nil
}

func initPinger(host string, opts options) (*probing.Pinger, error) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		return nil, fmt.Errorf("failed to init pinger %w", err)
	}

	pinger.Count = opts.Count

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		pinger.Stop()
	}()

	color.New(color.FgWhite, color.Bold).Printf(
		"PING %s (%s) type `Ctrl-C` to abort\n",
		pinger.Addr(),
		pinger.IPAddr(),
	)

	pinger.OnRecv = pingerOnrecv
	pinger.OnFinish = pingerOnFinish

	if opts.Privilege || runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	return pinger, nil
}

func pingerOnrecv(pkt *probing.Packet) {
	if runtime.GOOS == "windows" {
		fmt.Printf(
			"%s seq=%s %sbytes from %s: time=%s\n",
			renderASCIIArt(pkt.Seq),
			color.New(color.FgYellow, color.Bold).Sprintf("%d", pkt.Seq),
			color.New(color.FgBlue, color.Bold).Sprintf("%d", pkt.Nbytes),
			color.New(color.FgWhite, color.Bold).Sprintf("%s", pkt.IPAddr),
			color.New(color.FgMagenta, color.Bold).Sprintf("%v", pkt.Rtt),
		)
	} else {
		fmt.Printf(
			"%s seq=%s %sbytes from %s: ttl=%s time=%s\n",
			renderASCIIArt(pkt.Seq),
			color.New(color.FgYellow, color.Bold).Sprintf("%d", pkt.Seq),
			color.New(color.FgBlue, color.Bold).Sprintf("%d", pkt.Nbytes),
			color.New(color.FgWhite, color.Bold).Sprintf("%s", pkt.IPAddr),
			color.New(color.FgCyan, color.Bold).Sprintf("%d", pkt.TTL),
			color.New(color.FgMagenta, color.Bold).Sprintf("%v", pkt.Rtt),
		)
	}
}

func pingerOnFinish(stats *probing.Statistics) {
	color.New(color.FgWhite, color.Bold).Printf(
		"\n───────── %s ping statistics ─────────\n",
		stats.Addr,
	)
	fmt.Printf(
		"%s: %v transmitted => %v received (%v loss)\n",
		color.New(color.FgWhite, color.Bold).Sprintf("PACKET STATISTICS"),
		color.New(color.FgBlue, color.Bold).Sprintf("%d", stats.PacketsSent),
		color.New(color.FgGreen, color.Bold).Sprintf("%d", stats.PacketsRecv),
		color.New(color.FgRed, color.Bold).Sprintf("%v%%", stats.PacketLoss),
	)
	fmt.Printf(
		"%s: min=%v avg=%v max=%v stddev=%v\n",
		color.New(color.FgWhite, color.Bold).Sprintf("ROUND TRIP"),
		color.New(color.FgBlue, color.Bold).Sprintf("%v", stats.MinRtt),
		color.New(color.FgCyan, color.Bold).Sprintf("%v", stats.AvgRtt),
		color.New(color.FgGreen, color.Bold).Sprintf("%v", stats.MaxRtt),
		color.New(color.FgMagenta, color.Bold).Sprintf("%v", stats.StdDevRtt),
	)
}

func renderASCIIArt(idx int) string {
	if len(krrg) <= idx {
		idx %= len(krrg)
	}

	line := krrg[idx]

	line = colorize(line, 'A', color.NewRGBStyle(color.RGB(234, 44, 134), color.RGB(234, 44, 134)))
	line = colorize(line, 'B', color.NewRGBStyle(color.RGB(32, 47, 47), color.RGB(32, 47, 47)))
	line = colorize(line, 'C', color.NewRGBStyle(color.RGB(90, 79, 90), color.RGB(90, 79, 90)))
	line = colorize(line, 'D', color.NewRGBStyle(color.RGB(29, 20, 30), color.RGB(29, 20, 30)))
	line = colorize(line, 'E', color.NewRGBStyle(color.RGB(254, 177, 212), color.RGB(254, 177, 212)))
	line = colorize(line, 'F', color.NewRGBStyle(color.RGB(255, 234, 209), color.RGB(255, 234, 209)))
	line = colorize(line, 'G', color.NewRGBStyle(color.RGB(253, 213, 234), color.RGB(253, 213, 234)))
	line = colorize(line, 'H', color.NewRGBStyle(color.RGB(165, 16, 65), color.RGB(165, 16, 65)))
	line = colorize(line, 'I', color.NewRGBStyle(color.RGB(85, 70, 82), color.RGB(85, 70, 82)))
	line = colorize(line, 'J', color.NewRGBStyle(color.RGB(73, 37, 58), color.RGB(73, 37, 58)))
	line = colorize(line, 'K', color.NewRGBStyle(color.RGB(255, 239, 226), color.RGB(255, 239, 226)))
	line = colorize(line, 'L', color.NewRGBStyle(color.RGB(255, 255, 255), color.RGB(255, 255, 255)))
	line = colorize(line, 'M', color.NewRGBStyle(color.RGB(105, 14, 98), color.RGB(105, 14, 98)))
	line = colorize(line, 'N', color.NewRGBStyle(color.RGB(119, 67, 115), color.RGB(119, 67, 115)))
	line = colorize(line, 'O', color.NewRGBStyle(color.RGB(255, 224, 248), color.RGB(255, 224, 248)))
	line = colorize(line, 'P', color.NewRGBStyle(color.RGB(192, 24, 98), color.RGB(192, 24, 98)))
	line = colorize(line, 'Q', color.NewRGBStyle(color.RGB(215, 146, 156), color.RGB(215, 146, 156)))
	line = colorize(line, 'R', color.NewRGBStyle(color.RGB(255, 202, 207), color.RGB(255, 202, 207)))
	line = colorize(line, 'S', color.NewRGBStyle(color.RGB(136, 35, 73), color.RGB(136, 35, 73)))
	line = colorize(line, 'T', color.NewRGBStyle(color.RGB(253, 213, 234), color.RGB(253, 213, 234)))
	line = colorize(line, 'U', color.NewRGBStyle(color.RGB(67, 52, 60), color.RGB(67, 52, 60)))
	line = colorize(line, 'V', color.NewRGBStyle(color.RGB(77, 64, 55), color.RGB(77, 64, 55)))

	return line
}

func colorize(text string, target rune, color color.PrinterFace) string {
	return strings.ReplaceAll(
		text,
		string(target),
		color.Sprint("#"),
	)
}
