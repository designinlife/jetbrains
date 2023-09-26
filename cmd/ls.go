package cmd

import (
	"fmt"
	"github.com/designinlife/jetbrains/common"
	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strings"
)

const (
	JetbrainsApiBaseUrl = "https://data.services.jetbrains.com/products/releases"
)

type apiDataItem struct {
	Date      string `json:"date"`
	Type      string `json:"type"`
	Downloads struct {
		Linux struct {
			Link         string `json:"link"`
			Size         int64  `json:"size"`
			ChecksumLink string `json:"checksumLink"`
		} `json:"linux"`
		ThirdPartyLibrariesJson struct {
			Link         string `json:"link"`
			Size         int64  `json:"size"`
			ChecksumLink string `json:"checksumLink"`
		} `json:"thirdPartyLibrariesJson"`
		Windows struct {
			Link         string `json:"link"`
			Size         int64  `json:"size"`
			ChecksumLink string `json:"checksumLink"`
		} `json:"windows"`
		WindowsZip struct {
			Link         string `json:"link"`
			Size         int64  `json:"size"`
			ChecksumLink string `json:"checksumLink"`
		} `json:"windowsZip"`
		Mac struct {
			Link         string `json:"link"`
			Size         int64  `json:"size"`
			ChecksumLink string `json:"checksumLink"`
		} `json:"mac"`
		MacM1 struct {
			Link         string `json:"link"`
			Size         int64  `json:"size"`
			ChecksumLink string `json:"checksumLink"`
		} `json:"macM1"`
	} `json:"downloads"`
	Patches struct {
		Win []struct {
			FromBuild    string `json:"fromBuild"`
			Link         string `json:"link"`
			Size         int64  `json:"size"`
			ChecksumLink string `json:"checksumLink"`
		} `json:"win"`
		Mac []struct {
			FromBuild    string `json:"fromBuild"`
			Link         string `json:"link"`
			Size         int64  `json:"size"`
			ChecksumLink string `json:"checksumLink"`
		} `json:"mac"`
		Unix []struct {
			FromBuild    string `json:"fromBuild"`
			Link         string `json:"link"`
			Size         int64  `json:"size"`
			ChecksumLink string `json:"checksumLink"`
		} `json:"unix"`
	} `json:"patches"`
	NotesLink              string `json:"notesLink"`
	LicenseRequired        bool   `json:"licenseRequired"`
	Version                string `json:"version"`
	MajorVersion           string `json:"majorVersion"`
	Build                  string `json:"build"`
	Whatsnew               string `json:"whatsnew"`
	UninstallFeedbackLinks struct {
		WindowsJBR8             string `json:"windowsJBR8"`
		WindowsZipJBR8          string `json:"windowsZipJBR8"`
		Linux                   string `json:"linux"`
		ThirdPartyLibrariesJson string `json:"thirdPartyLibrariesJson"`
		Windows                 string `json:"windows"`
		WindowsZip              string `json:"windowsZip"`
		LinuxJBR8               string `json:"linuxJBR8"`
		Mac                     string `json:"mac"`
		MacJBR8                 string `json:"macJBR8"`
		MacM1                   string `json:"macM1"`
	} `json:"uninstallFeedbackLinks"`
	PrintableReleaseType interface{} `json:"printableReleaseType"`
}

type ApiDataSet map[string][]apiDataItem

type tableDataSet [][]string

func (t tableDataSet) Len() int {
	return len(t)
}

func (t tableDataSet) Less(i, j int) bool {
	return t[i][3] > t[j][3]
}

func (t tableDataSet) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Get a list of the latest versions of Jetbrains software.",
	Long: `This command will read the latest version number of the software 
through the Jetbrains HTTP-JSON interface and print the download address of each platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		names := map[string]string{"AC": "AppCode", "CL": "CLion", "RSU": "ReSharper Ultimate", "DG": "DataGrip",
			"GO": "Goland", "IIU": "IntelliJ IDEA", "PS": "PhpStorm", "PCP": "PyCharm", "RD": "Rider",
			"RM": "RubyMine", "WS": "WebStorm", "FL": "Fleet", "RR": "RustRover", "DS": "DataSpell"}

		var codes []string

		for k := range names {
			codes = append(codes, k)
		}

		params := map[string]string{
			"code":   strings.Join(codes, ","),
			"latest": "true",
			"type":   "release",
		}

		apiResult := ApiDataSet{}

		client := resty.New()

		resp, err := client.R().SetQueryParams(params).SetResult(&apiResult).Get(JetbrainsApiBaseUrl)

		if err != nil {
			panic(err)
		}
		if resp.StatusCode() != 200 {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("HTTP status code is not 200. (StatusCode: %d)", resp.StatusCode()))
			os.Exit(1)
		}

		windowsLinks := make([]string, 0)
		linuxLinks := make([]string, 0)
		macLinks := make([]string, 0)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"NAME", "SIZE", "VERSION", "RELEASE DATE"})
		table.SetAutoWrapText(false)
		table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER})

		var tableData tableDataSet

		for key, values := range apiResult {
			for _, value := range values {
				if key == "AC" {
					tableData = append(tableData, []string{names[key], common.ByteCountSI(value.Downloads.Mac.Size), value.Version, value.Date})
				} else {
					tableData = append(tableData, []string{names[key], common.ByteCountSI(value.Downloads.Windows.Size), value.Version, value.Date})
				}

				if len(value.Downloads.Windows.Link) > 0 {
					windowsLinks = append(windowsLinks, strings.ReplaceAll(value.Downloads.Windows.Link, "download.jetbrains.com", "download-cf.jetbrains.com"))
				}
				if len(value.Downloads.Linux.Link) > 0 {
					linuxLinks = append(linuxLinks, strings.ReplaceAll(value.Downloads.Linux.Link, "download.jetbrains.com", "download-cf.jetbrains.com"))
				}
				if len(value.Downloads.Mac.Link) > 0 {
					macLinks = append(macLinks, strings.ReplaceAll(value.Downloads.Mac.Link, "download.jetbrains.com", "download-cf.jetbrains.com"))
				}
			}
		}

		sort.Sort(tableData)

		table.AppendBulk(tableData)
		table.Render()

		fmt.Println()
		fmt.Println("The download link for \033[1;32mWindows\033[0m, follows as:")
		fmt.Println("------------------------------------------")

		for _, v1 := range windowsLinks {
			fmt.Println(v1)
		}

		fmt.Println()
		fmt.Println("The download link for \033[1;32mLinux\033[0m, follows as:")
		fmt.Println("------------------------------------------")

		for _, v1 := range linuxLinks {
			fmt.Println(v1)
		}

		fmt.Println()
		fmt.Println("The download link for \033[1;32mMac\033[0m, follows as:")
		fmt.Println("------------------------------------------")

		for _, v1 := range macLinks {
			fmt.Println(v1)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.
	lsCmd.Flags().Bool("help", false, "Show help.")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
