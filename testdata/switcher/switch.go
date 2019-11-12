package switcher

func switcher(in string) string {
	switch in {
	case "one":
		return "ok"
	default: // nocover
		return "unexpected"
	}
}
