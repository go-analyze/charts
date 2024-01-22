package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-analyze/charts/chartdraw"
)

// Note: Additional examples on how to add Stylesheets are in the custom_stylesheets example

func inlineSVGWithClasses(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(
		"<!DOCTYPE html><html><head>" +
			"<link rel=\"stylesheet\" type=\"text/css\" href=\"/main.css\">" +
			"</head>" +
			"<body>"))

	pie := chartdraw.PieChart{
		// Notes: * Setting ClassName will cause all other inline styles to be dropped!
		//        * The following type classes may be added additionally: stroke, fill, text
		Background: chartdraw.Style{ClassName: "background"},
		Canvas: chartdraw.Style{
			ClassName: "canvas",
		},
		Width:  512,
		Height: 512,
		Values: []chartdraw.Value{
			{Value: 5, Label: "Blue", Style: chartdraw.Style{ClassName: "blue"}},
			{Value: 5, Label: "Green", Style: chartdraw.Style{ClassName: "green"}},
			{Value: 4, Label: "Gray", Style: chartdraw.Style{ClassName: "gray"}},
		},
	}

	err := pie.Render(chartdraw.SVG, res)
	if err != nil {
		fmt.Printf("Error rendering pie chart: %v\n", err)
	}
	res.Write([]byte("</body>"))
}

func css(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/css")
	res.Write([]byte("svg .background { fill: white; }" +
		"svg .canvas { fill: white; }" +
		"svg .blue.fill.stroke { fill: blue; stroke: lightblue; }" +
		"svg .green.fill.stroke { fill: green; stroke: lightgreen; }" +
		"svg .gray.fill.stroke { fill: gray; stroke: lightgray; }" +
		"svg .blue.text { fill: white; }" +
		"svg .green.text { fill: white; }" +
		"svg .gray.text { fill: white; }"))
}

func main() {
	http.HandleFunc("/", inlineSVGWithClasses)
	http.HandleFunc("/main.css", css)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
