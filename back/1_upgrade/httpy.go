package main

import "strings"

type HTTPy struct {
	Method      string
	HTTPVersion string
	UserAgent   string
	Host        string
	Origin      string
	SecWebKey   string
	Upgrade     string
}

func NewHTTPy() *HTTPy {
	return &HTTPy{}
}

func (h *HTTPy) GetHTTPy(lines string) (httpy *HTTPy) {

	linesArr := strings.Split(lines, "\n")
	httpy = &HTTPy{}
	for _, line := range linesArr {
		split := strings.SplitN(line, ":", 2)
		httpy.Method = "GET"
		httpy.HTTPVersion = "HTTP/1.1"

		switch split[0] {
		case "Host":
			{
				httpy.Host = strings.Trim(split[1], "\r")
			}
		case "Origin":
			{
				httpy.Origin = strings.Trim(split[1], "\r")
			}
		case "Sec-WebSocket-Key":
			{
				httpy.SecWebKey = strings.Trim(split[1], "\r")
			}
		case "Upgrade":
			{
				httpy.Upgrade = strings.Trim(split[1], "\r")
			}
		case "User-Agent":
			{
				httpy.UserAgent = strings.Trim(split[1], "\r")
			}

		}
	}
	return
}
