package acid

import "fmt"

func CSS(path string) string {
	fullPath := fmt.Sprintf("static/css/%s.css", path)
	reverseLink := assetsWithDigests.ReverseMap[fullPath]

	return fmt.Sprintf("/assets/%s", reverseLink)
}

func Javascript(path string) string {
	// FIX: this does not support nested directories
	fullPath := fmt.Sprintf("static/javascript/%s.js", path)
	reverseLink := assetsWithDigests.ReverseMap[fullPath]

	return fmt.Sprintf(`<script type="module" src="/assets/%s"></script>`, reverseLink)
}

func Image(path string) string {
	// remove the leading slash from path if it exists
	if path[0] == '/' {
		path = path[1:]
	}

	fullPath := fmt.Sprintf("static/images/%s", path)
	reverseLink := assetsWithDigests.ReverseMap[fullPath]

	return fmt.Sprintf("/assets/%s", reverseLink)
}
