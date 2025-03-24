package services

type VisionService interface {
	AnalyseImage(imagePath string, prompt string) (string, error)
}
