package services

type VisionService interface {
	AnalyseImage(imagePath string, prompt string) (string, error)
	GenerateMessage(message string) (string, error)
}
