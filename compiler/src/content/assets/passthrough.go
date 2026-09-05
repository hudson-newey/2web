package assets

import "slices"

// Some assets should be passed through without modification.
// Some assets like images DO NOT go through this pipeline and instead have
// their own content model so we can do neat things like compression and file
// type conversion.
func ShouldUseAssetPassthrough(filePath string) bool {
	assetExtensions := []string{".pdf", ".wasm"}
	return slices.Contains(assetExtensions, filePath)
}
