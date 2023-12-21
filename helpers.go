package acid

import "fmt"

func CSS(path string) string {
	fullPath := fmt.Sprintf("static/css/%s.css", path)
	reverseLink := AssetsWithDigests.ReverseMap[fullPath]
	return fmt.Sprintf(`<link rel="stylesheet" href="/assets/%s">`, reverseLink)
}

func Javascript(path string) string {
	// FIX: this does not support nested directories
	fullPath := fmt.Sprintf("static/javascript/%s.js", path)
	reverseLink := AssetsWithDigests.ReverseMap[fullPath]
	return fmt.Sprintf(`<script type="module" src="/assets/%s"></script>`, reverseLink)
}
