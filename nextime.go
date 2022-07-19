package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func MakeTable(b io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(b)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	return table
}

type CritPath struct {
	From string
	Path []*PathItem
	To   string
}

func (p *CritPath) String() string {
	b := &bytes.Buffer{}
	table := MakeTable(b)
	table.SetHeader([]string{"total", "routing", "logic", "net"})

	total := 0.0
	for i := 0; i < len(p.Path)-1; i++ {
		if p.Path[i].Type == "routing" {
			table.Append([]string{
				fmt.Sprintf("%0.2f", total),
				fmt.Sprintf("%0.2f", p.Path[i].Delay),
				fmt.Sprintf("%0.2f", p.Path[i+1].Delay),
				fmt.Sprintf("%s", p.Path[i].Net),
			})
		} else if i == 0 {
			table.Append([]string{
				fmt.Sprintf("%0.2f", total),
				fmt.Sprintf("%0.2f", p.Path[i].Delay),
				fmt.Sprintf("%0.2f", 0.0),
				fmt.Sprintf("%s", p.Path[i].From),
			})
		}
		total += p.Path[i].Delay
	}
	table.Append([]string{
		fmt.Sprintf("%0.2f", total),
		fmt.Sprintf("%0.2f", p.Path[len(p.Path)-1].Delay),
		fmt.Sprintf("%0.2f", 0.0),
		fmt.Sprintf("%s", p.Path[len(p.Path)-1].From),
	})
	table.Render()

	b.WriteString(fmt.Sprintf("Critical path: %s -> %s\n", p.Path[0].From, p.Path[len(p.Path)-1].To))
	b.WriteString(fmt.Sprintf("Max frequency: %0.2f MHz (%0.2f ns)\n", p.Fmax(), p.Period()))

	return b.String()
}

// Returns max frequency in MHz
func (p *CritPath) Fmax() float64 {
	return 1 / p.Period() * 1000
}

func (p *CritPath) Period() float64 {
	total := 0.0
	for _, item := range p.Path {
		total += item.Delay
	}
	return total
}

type PathItem struct {
	Budget float64
	Delay  float64
	From   *Cell
	To     *Cell
	Type   string
	Net    string
}

func (p *PathItem) String() string {
	return fmt.Sprintf("%0.2f %s -> %s", p.Delay, p.From, p.To)
}

type Cell struct {
	Cell string
	Port string
}

func (c *Cell) String() string {
	name := fmt.Sprintf("%s[%s]", c.Cell, c.Port)
	if len(name) > 30 {
		name = name[:30]
	}
	return name
}

type Freq struct {
	Achieved   float64
	Constraint float64
}

type Resource struct {
	Available int
	Used      int
}

func (r *Resource) Info() []string {
	return []string{
		fmt.Sprintf("%d", r.Available),
		fmt.Sprintf("%d", r.Used),
		fmt.Sprintf("%0.2f", float64(r.Used)/float64(r.Available)),
	}
}

type Utilization map[string]Resource

func (u Utilization) String() string {
	b := &bytes.Buffer{}

	table := MakeTable(b)
	table.SetHeader([]string{"cell", "total", "used", "utilization"})

	keys := make([]string, 0, len(u))
	for name := range u {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	for _, name := range keys {
		r := u[name]
		if r.Used != 0 {
			table.Append(append([]string{name}, r.Info()...))
		}
	}
	table.Render()

	return b.String()
}

type Info struct {
	CriticalPaths      []CritPath `json:"critical_paths"`
	Fmax               map[string]Freq
	Utilization        Utilization
	DetailedNetTimings []Driver `json:"detailed_net_timings"`
}

type Driver struct {
	Driver    string
	Endpoints []Endpoint
	Event     string
	Net       string
	Port      string
}

type Endpoint struct {
	Budget float64
	Cell   string
	Delay  float64
	Event  string
	Port   string
}

var util = flag.Bool("util", false, "show utilization")
var clk = flag.String("clk", "", "show summary for clock")
var crit = flag.Bool("crit", false, "also show critical path for selected clock")

func main() {
	flag.Parse()
	args := flag.Args()

	file := "report.json"
	if len(args) > 0 {
		file = args[0]
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	var info Info
	err = json.Unmarshal(data, &info)
	if err != nil {
		log.Fatal(err)
	}

	if *util {
		fmt.Print(info.Utilization.String())
	}

	if *clk == "" {
		return
	}

	for net, freq := range info.Fmax {
		if *clk == "all" || strings.Contains(net, *clk) {
			if net != "" {
				fmax := make([]float64, 0, len(info.CriticalPaths))
				cps := make([]CritPath, 0, len(info.CriticalPaths))
				for _, cp := range info.CriticalPaths {
					if strings.Contains(cp.From, net) {
						fmax = append(fmax, cp.Fmax())
						cps = append(cps, cp)
					}
				}

				i := argmin(fmax)
				if *crit {
					fmt.Print(cps[i].String())
				}
				if cps[i].Fmax() < freq.Constraint {
					fmt.Printf("%s is %.2f MHz (failed at %0.2f MHz)\n", net, cps[i].Fmax(), freq.Constraint)
				} else {
					fmt.Printf("%s is %.2f MHz (succeeded at %0.2f MHz)\n", net, cps[i].Fmax(), freq.Constraint)
				}
			}
		}
	}

}

func argmin(arr []float64) int {
	idx := 0
	for i, f := range arr {
		if f < arr[idx] {
			idx = i
		}
	}
	return idx
}
