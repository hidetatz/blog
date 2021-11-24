package main

import "fmt"

const indexPageMD = `
# hidetatz.io

hidetatz.io is my personal blog. The author Hidetatsu (pronounced he-day-tatz) is a software engineer mainly focuses on system architecture, reliability, performance and observability based in Japan. I write code around infrastructure, database, transaction, concurrent programming and distributed systems. My code is available in [GitHub](https://github.com/hidetatz).

I [do fail](https://hidetatz.fail/).

---

## Articles

%s

Some articles are available in Japanese also.

%s

---

## Some other pages

* [/inputs](/inputs.html)
  - What I've read, listened, watched, etc.
* [/distsys](/distsys.html)
  - Distributed systems learning meterials (in Japanese)

---

If you want to send me any feedback about this website, you can submit it as GitHub issue [here](https://github.com/hidetatz/blog/issues/new).

---

[feed](/feed.xml)
`

func generateIndexPageHTML(articles []*article) string {
	enblogsList := ""
	jablogsList := ""
	for _, a := range articles {
		switch a.typ {
		case inputType:
			continue
		case blogType:
			switch a.lang {
			case en:
				enblogsList += fmt.Sprintf("%s	- [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, link(a))
			case ja:
				jablogsList += fmt.Sprintf("%s	- [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, link(a))
			}
		}
	}

	return generateHTMLPage("hidetatz.io", fmt.Sprintf(indexPageMD, enblogsList, jablogsList))
}

func link(a *article) string {
	formattedTime := a.timestamp.Format(timeformat)
	switch a.typ {
	case blogType:
		// if blog, the link should be external URL or internal link
		switch {
		case a.url == nil:
			return fmt.Sprintf("/articles/%s/%s", formattedTime, trimExtension(a.fileName))
		default:
			return a.url.String()
		}
	default:
		// else, it is input. return internal link
		return fmt.Sprintf("/articles/%s/%s", formattedTime, trimExtension(a.fileName))
	}
}

const inputPageMD = `

[<- home](/)

# /inputs

%s

`

func generateInputPageHTML(articles []*article) string {
	list := ""
	for _, a := range articles {
		switch a.typ {
		case blogType:
			continue
		case inputType:
			list += fmt.Sprintf("%s	- [%s](%s)  \n", a.timestamp.Format(timeformat), a.title, link(a))
		}
	}

	return generateHTMLPage("hidetatz.io | inputs", fmt.Sprintf(inputPageMD, list))
}
