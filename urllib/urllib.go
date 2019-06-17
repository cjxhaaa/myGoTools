package urllib

import (
	"net/url"
	"strings"
)

//Urljoin is similiar to python's urljoin
//todo: Fragment not add
func Urljoin(base string, purl string) (string,error) {

	if base == "" {
		return purl,nil
	}

	if purl == "" {
		return base,nil
	}

	var buf strings.Builder
	bu,err := url.Parse(base)
	if err != nil {
		return base, err
	}

	pu,err := url.Parse(purl)
	if err != nil {
		return purl, err
	}

	if bu.Scheme != "" {
		buf.WriteString(bu.Scheme)
		buf.WriteString(":")
	} else if pu.Scheme != "" {
		buf.WriteString(pu.Scheme)
		buf.WriteString(":")
	}

	if bu.Scheme != "" || pu.Scheme != "" {
		buf.WriteString("//")
	}

	if h := bu.Host;h != "" {
		buf.WriteString(h)
	} else if h := pu.Host; h != "" {
		buf.WriteString(h)
	}

	path := pu.EscapedPath()
	if path != "" && path[0] != '/' && bu.Host != "" {
		buf.WriteByte('/')
	}
	if buf.Len() == 0 {
		if i := strings.IndexByte(path, ':'); i > -1 && strings.IndexByte(path[:i],'/') == -1 {
			buf.WriteString("./")
		}
	}
	buf.WriteString(path)

	if pu.ForceQuery || pu.RawQuery != "" {
		buf.WriteByte('?')
		buf.WriteString(pu.RawQuery)
	}
	return buf.String(), nil
}
