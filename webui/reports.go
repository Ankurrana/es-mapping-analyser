package webui

/*
	Create a graph of the usage pattern for the fields for easily visualisation
*/

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ankur-toko/es-mapping-analyser/reports"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type MyStringWrapper struct {
	str []byte
}

func (wrapper *MyStringWrapper) Write(k []byte) (n int, err error) {
	wrapper.str = k
	return len(k), nil
}

func (wrapper *MyStringWrapper) GetString() []byte {
	return wrapper.str
}

func QURMapToGraph(report map[string]reports.QMJSONReport) string {
	strWrapper := &MyStringWrapper{}
	page := components.NewPage()

	for index, r := range report {
		page.AddCharts(CreateGraph(index, r.QueriesCount, r.UsageMap))
	}
	page.Render(strWrapper)

	return string(strWrapper.GetString())
}

func CreateGraph(indexName string, QueriesAnalyzed int, m map[string]map[string]int) *charts.Bar {
	// range
	minValue := 1000000000
	maxValue := 0
	// list of items
	// orderise items
	// give them color
	fieldNames := []string{}
	usecaseMap := map[string]bool{}

	for field, usage := range m {
		fieldNames = append(fieldNames, field)
		for usecase, frequency := range usage {
			usecaseMap[usecase] = true

			if frequency > maxValue {
				maxValue = frequency
			}
			if frequency < minValue {
				minValue = frequency
			}
		}
	}

	sort.Strings(fieldNames)

	height := fmt.Sprintf("%vpx", len(fieldNames)*20)

	queriesCount := fmt.Sprintf("Queries Analyzed : %v", QueriesAnalyzed)

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    strings.ToUpper(indexName),
			Subtitle: queriesCount,
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1500px",
			Height: height,
		}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "80px"}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show: true,
			Left: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "Download",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "Show Data",
					// set the language
					// Chinese version: ["数据视图", "关闭", "刷新"]
					Lang: []string{"data view", "turn off", "refresh"},
				},
			}}),
	)

	k := bar.SetXAxis(fieldNames)

	for usecase, _ := range usecaseMap {

		items := make([]opts.BarData, 0)
		for i := 0; i < len(fieldNames); i++ {
			val := m[fieldNames[i]][usecase]
			items = append(items, opts.BarData{Value: val})
		}

		k.AddSeries(usecase, items)
	}
	k.
		SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
			Stack: "stackB",
		}))

	bar.XYReversal()
	return bar

}
