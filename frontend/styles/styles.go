package styles

var (
	//dims       Dimensions
	breakpoint int64 = 600
	mq         *MediaQuery
)

func init() {
	mq = NewMediaQuery([]Breakpoint{600})
}

func NavBar() *CSS {
	common := NewCSS(
		"list-style", "none",
		"background-color", "#444",
	)
	mobile := NewCSS(
		"margin", "auto",
		"width", "100%",
		"overflow", "auto",
	)
	desktop := NewCSS(
		"text-align", "center",
		"padding", "0",
		"margin", "0",
	)

	return mq.Query(common, mobile, desktop)
}

func NavItem() *CSS {
	common := NewCSS(
		"font-family", "'Oswald', sans-serif",
	)
	mobile := NewCSS(
		"font-size", "1.4em",
		"line-height", "50px",
		"height", "50px",
	)
	desktop := NewCSS(
		"width", "120px",
		// "float", "left",
		"font-size", "1.2em",
		"line-height", "40px",
		"height", "40px",
		"border-bottom", "1px solid #888",
	)

	return mq.Query(common, mobile, desktop)
}

func NavAnchor(hovered bool) *CSS {
	var color string
	if hovered {
		color = "rgb(135, 133, 133)"
	} else {
		color = "rgb(89, 89, 89)"
	}
	common := NewCSS(
		"text-decoration", "none",
		"display", "block",
		"transition", ".2s background-color",
		"color", "rgb(222, 222, 216)",
		"background-color", color,
	)

	return common
}

func ProxyForm() *CSS {

	common := NewCSS(
		"display", "grid",
		"grid-template-rows", "25% 25% 25% 25%",
		"background-color", "#444",
		"border-radius", "5px",
		"padding", "15px",
	)

	mobile := NewCSS()
	desktop := NewCSS()
	c := mq.Query(common, mobile, desktop)
	return c
}

func BackendList() *CSS {
	common := NewCSS(
		"grid-gap", "10px",
		"display", "grid",
		"margin", "10px",
	)
	mobile := NewCSS(
		"width", "100%",
		"grid-template-columns", "repeat(1, 80%)",
	)
	desktop := NewCSS(
		"width", "800px",
		"grid-template-columns", "repeat(4, 200px)",
	)
	return mq.Query(common, mobile, desktop)
}
