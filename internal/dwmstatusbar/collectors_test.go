package dwmstatusbar_test

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mhersson/dwmstatusbar/internal/dwmstatusbar"
	"github.com/mhersson/dwmstatusbar/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/spf13/afero"
)

const xsetOutput = `
Keyboard Control:
  auto repeat:  on    key click percent:  0    LED mask:  00000000
  XKB indicators:
    00: Caps Lock:   off    01: Num Lock:    off    02: Scroll Lock: off
    03: Compose:     off    04: Kana:        off    05: Sleep:       off
    06: Suspend:     off    07: Mute:        off    08: Misc:        off
    09: Mail:        off    10: Charging:    off    11: Shift Lock:  off
    12: Group 2:     off    13: Mouse Keys:  off
  auto repeat delay:  250    repeat rate:  25
  auto repeating keys:  00ffffffdffffbbf
                        fadfffefffedffff
                        9fffffffffffffff
                        fff7ffffffffffff
  bell percent:  0    bell pitch:  400    bell duration:  100
Pointer Control:
  acceleration:  2/1    threshold:  4
Screen Saver:
  prefer blanking:  yes    allow exposures:  yes
  timeout:  600    cycle:  600
Colors:
  default colormap:  0x20    BlackPixel:  0x0    WhitePixel:  0xffffff
Font Path:
  built-ins
DPMS (Display Power Management Signaling):
  Standby: 600    Suspend: 600    Off: 600
  DPMS is Disabled
  Monitor is On
`

var TestBattery = Describe("Battery", func() {
	var fs afero.Fs

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		fs.MkdirAll("/sys/class/power_supply/BAT0", 0o755)
		afero.WriteFile(fs, "/sys/class/power_supply/BAT0/capacity", []byte("100"), 0o644)
		afero.WriteFile(fs, "/sys/class/power_supply/BAT0/status", []byte("Discharging"), 0o644)
		dwmstatusbar.Fsys = fs
	})

	Context("when the battery is present", func() {
		It("returns the battery level with a percent sign", func() {
			Expect(dwmstatusbar.Battery("")).To(Equal("100%"))
		})

		It("says Charging when charging", func() {
			afero.WriteFile(fs, "/sys/class/power_supply/BAT0/status", []byte("Charging"), 0o644)
			Expect(dwmstatusbar.Battery("")).To(Equal("Charging"))
		})
	})

	Context("when the battery is not present", func() {
		It("returns No Battery", func() {
			fs.RemoveAll("/sys/class/power_supply/BAT0")
			Expect(dwmstatusbar.Battery("")).To(Equal("No Battery"))
		})
	})
})

var TestPIA = Describe("PIA", func() {
	var (
		execMock *mocks.MockExecInterface
		cmdMock  *mocks.MockCmdInterface
	)

	BeforeEach(func() {
		execMock = new(mocks.MockExecInterface)
		cmdMock = new(mocks.MockCmdInterface)
		dwmstatusbar.ExecCmd = execMock
		dwmstatusbar.Log.SetOutput(GinkgoWriter)
	})

	Context("when the command succeeds", func() {
		It("returns the IP address", func() {
			execMock.On("NewCommand", "piactl", "get", "vpnip").Return(cmdMock)
			cmdMock.On("CombinedOutput").Return([]byte("192.111.222.333"), nil)
			Expect(dwmstatusbar.PIA("")).To(Equal("192.111.222.333"))
		})
	})

	Context("when the command fails", func() {
		It("returns waiting for data", func() {
			execMock.On("NewCommand", "piactl", "get", "vpnip").Return(cmdMock)
			cmdMock.On("CombinedOutput").Return([]byte(""), fmt.Errorf("oops"))
			Expect(dwmstatusbar.PIA("")).To(Equal("No Data"))
		})
	})
})

var TestExternalIP = Describe("ExternalIP", func() {
	var server *ghttp.Server

	BeforeEach(func() {
		server = ghttp.NewServer()
		dwmstatusbar.ExternalIPURL = server.URL()
	})

	AfterEach(func() {
		server.Close()
	})

	Context("when the command succeeds", func() {
		It("should retrieve the external IP", func() {
			expectedIP := "192.168.1.1"
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/"),
					ghttp.RespondWith(http.StatusOK, expectedIP),
				),
			)

			ip := dwmstatusbar.ExternalIP("")
			Expect(ip).To(Equal(expectedIP))
		})
	})

	Context("when the command fails", func() {
		It("should return an empty string", func() {
			expectedIP := ""
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/"),
					ghttp.RespondWith(http.StatusInternalServerError, expectedIP),
				),
			)

			ip := dwmstatusbar.ExternalIP("")
			Expect(ip).To(Equal(expectedIP))
		})
	})
})

var TestKeyboardLayout = Describe("KeyboardLayout", func() {
	var (
		execMock *mocks.MockExecInterface
		cmdMock  *mocks.MockCmdInterface
		out      string
	)

	BeforeEach(func() {
		execMock = new(mocks.MockExecInterface)
		cmdMock = new(mocks.MockCmdInterface)
		dwmstatusbar.ExecCmd = execMock
		dwmstatusbar.Log.SetOutput(GinkgoWriter)
		out = xsetOutput
	})

	Context("when the command succeeds", func() {
		It("returns the layout", func() {
			execMock.On("NewCommand", "xset", "q").Return(cmdMock)
			cmdMock.On("CombinedOutput").Return([]byte(out), nil)
			Expect(dwmstatusbar.KeyboardLayout(out)).To(Equal("US"))
		})
	})

	Context("when the command succeeds with a different layout", func() {
		It("returns the layout", func() {
			xsetOut := strings.Replace(out, "LED mask:  00000000", "LED mask:  00001000", 1)
			execMock.On("NewCommand", "xset", "q").Return(cmdMock)
			cmdMock.On("CombinedOutput").Return([]byte(xsetOut), nil)
			Expect(dwmstatusbar.KeyboardLayout(xsetOut)).To(Equal("NO"))
		})
	})

	Context("when the command fails", func() {
		It("returns default layout", func() {
			execMock.On("NewCommand", "xset", "q").Return(cmdMock)
			cmdMock.On("CombinedOutput").Return(nil, fmt.Errorf("oops"))
			Expect(dwmstatusbar.KeyboardLayout("")).To(Equal("US"))
		})
	})
})

var TestDPMS = Describe("DPMS", func() {
	var (
		execMock *mocks.MockExecInterface
		cmdMock  *mocks.MockCmdInterface
		out      string
	)

	BeforeEach(func() {
		execMock = new(mocks.MockExecInterface)
		cmdMock = new(mocks.MockCmdInterface)
		dwmstatusbar.ExecCmd = execMock
		dwmstatusbar.Log.SetOutput(GinkgoWriter)
		out = xsetOutput
	})

	Context("when the command succeeds", func() {
		It("returns DPMS OFF if DPMS is Disabled", func() {
			execMock.On("NewCommand", "xset", "q").Return(cmdMock)
			cmdMock.On("CombinedOutput").Return([]byte(out), nil)
			Expect(dwmstatusbar.DPMS(out)).To(Equal("DPMS OFF"))
		})
	})

	Context("when the command succeeds", func() {
		It("returns DPMS ON if DPMS is Enabled", func() {
			out := strings.Replace(out, "DPMS is Disabled", "DPMS is Enabled", 1)

			execMock.On("NewCommand", "xset", "q").Return(cmdMock)
			cmdMock.On("CombinedOutput").Return([]byte(out), nil)
			Expect(dwmstatusbar.DPMS(out)).To(Equal("DPMS ON"))
		})
	})

	Context("when the command fails", func() {
		It("returns default DPMS ON", func() {
			execMock.On("NewCommand", "xset", "q").Return(cmdMock)
			cmdMock.On("CombinedOutput").Return(nil, fmt.Errorf("oops"))
			Expect(dwmstatusbar.DPMS("")).To(Equal("DPMS ON"))
		})
	})
})
