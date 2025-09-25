package services

import (
	"minimalist-dashboard/internal/models"
	"minimalist-dashboard/internal/utils"
)

// DockerService gère les informations Docker
type DockerService struct{}

// NewDockerService crée une nouvelle instance du service Docker
func NewDockerService() *DockerService {
	return &DockerService{}
}

// GetDockerInfo récupère les informations Docker
func (d *DockerService) GetDockerInfo() models.DockerInfo {
	containersOut, _ := utils.RunCommand("docker", "ps", "--format", "{{.ID}}")
	imagesOut, _ := utils.RunCommand("docker", "images", "--format", "{{.ID}}")
	volumesOut, _ := utils.RunCommand("docker", "volume", "ls", "--format", "{{.Name}}")

	return models.DockerInfo{
		Containers: utils.CountLines(containersOut),
		Images:     utils.CountLines(imagesOut),
		Volumes:    utils.CountLines(volumesOut),
	}
}