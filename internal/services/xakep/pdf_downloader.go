package xakep

func NewPDFDownloader(dstDir string) *PDFDownloader {
	p := PDFDownloader{dstDir: dstDir}

	return &p
}

type PDFDownloader struct {
	dstDir string
}
