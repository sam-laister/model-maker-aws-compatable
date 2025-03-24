package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	models "github.com/Soup666/diss-api/model"
	repositories "github.com/Soup666/diss-api/repository"
	"gorm.io/gorm"
)

type TaskServiceImpl struct {
	taskRepo       repositories.TaskRepository
	appFileService AppFileService
}

func NewTaskService(taskRepo repositories.TaskRepository, appFileService AppFileService) *TaskServiceImpl {
	return &TaskServiceImpl{taskRepo: taskRepo, appFileService: appFileService}
}

func (s *TaskServiceImpl) CreateTask(task *models.Task) (*models.Task, error) {
	err := s.taskRepo.CreateTask(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) GetTask(taskID uint) (*models.Task, error) {
	task, err := s.taskRepo.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) GetTasks(userID uint) ([]models.Task, error) {

	tasks, err := s.taskRepo.GetTasksByUser(userID)

	if err != nil {
		return nil, err
	}
	return tasks, nil

}

func (s *TaskServiceImpl) UpdateTask(task *models.Task) (*models.Task, error) {

	err := s.taskRepo.SaveTask(task)

	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskServiceImpl) UpdateMeta(task *models.Task, key string, value interface{}) error {
	if task.Metadata == nil {
		task.Metadata = make(map[string]interface{})
	}
	task.Metadata[key] = value
	_, err := s.UpdateTask(task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) ArchiveTask(taskID uint) error {

	task, err := s.taskRepo.GetTaskByID(taskID)

	if err != nil {
		return err
	}

	err = s.taskRepo.ArchiveTask(task)

	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) SaveTask(task *models.Task) error {
	err := s.taskRepo.SaveTask(task)

	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) DeleteTask(taskID *models.Task) error {
	err := s.taskRepo.ArchiveTask(taskID)

	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) FailTask(task *models.Task) error {
	task.Status = models.FAILED
	_, err := s.UpdateTask(task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskServiceImpl) RunPhotogrammetryProcess(task *models.Task) error {
	startTime := time.Now()

	TASK_COUNT := 7
	CURRENT_TASK := 0

	inputPath := filepath.Join("uploads", fmt.Sprintf("%d", task.Id))
	outputPath := filepath.Join("objects", fmt.Sprintf("%d", task.Id))
	mvsPath := filepath.Join(outputPath, "mvs")

	task.Status = models.INPROGRESS
	if _, err := s.UpdateTask(task); err != nil {
		log.Printf("Failed to update task status to INPROGRESS: %v\n", err)
		return err
	}

	// Clear the build directory
	if err := os.RemoveAll(outputPath); err != nil {
		log.Printf("Failed to clear directory %s: %v", outputPath, err)
		s.FailTask(task)
		return err
	}

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		log.Printf("Failed to create directory %s: %v", outputPath, err)
		s.FailTask(task)
		return err
	}

	if err := os.MkdirAll(mvsPath, os.ModePerm); err != nil {
		log.Printf("Failed to create directory %s: %v", mvsPath, err)
		s.FailTask(task)
		return err
	}

	// 1
	log.Println("Updating meta for task:", task.Id, " - ", CURRENT_TASK, "/", TASK_COUNT)
	if err := s.UpdateMeta(task, "opensfm-process", CURRENT_TASK/100); err != nil {
		log.Printf("Failed to update meta: %d: %v\n", CURRENT_TASK, err)
		return err
	}
	CURRENT_TASK++

	log.Println("# 1 ./bin/SfM_SequentialPipeline.py", inputPath, outputPath, "--opensfm-processes", "8")
	cmd := exec.Command("./bin/SfM_SequentialPipeline.py", inputPath, outputPath, "--opensfm-processes", "8")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Println("SfM_SequentialPipeline failed:", err)
		s.FailTask(task)
		return err
	}

	// 2
	log.Println("Updating meta for task:", task.Id, " - ", CURRENT_TASK, "/", TASK_COUNT)
	if err := s.UpdateMeta(task, "opensfm-process", CURRENT_TASK/100); err != nil {
		log.Printf("Failed to update meta: %d: %v\n", CURRENT_TASK, err)
		return err
	}
	CURRENT_TASK++

	log.Println("# 2 openMVG_main_openMVG2openMVS", "-i", filepath.Join(outputPath, "reconstruction_sequential/sfm_data.bin"), "-o", filepath.Join(mvsPath, "scene.mvs"), inputPath, outputPath, "-d", mvsPath)
	cmd = exec.Command("openMVG_main_openMVG2openMVS", "-i", filepath.Join(outputPath, "reconstruction_sequential/sfm_data.bin"), "-o", filepath.Join(mvsPath, "scene.mvs"), inputPath, outputPath, "-d", mvsPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		log.Println("openMVG_main_openMVG2openMVS failed:", err)
		s.FailTask(task)
		return err
	}

	// 3
	log.Println("Updating meta for task:", task.Id, " - ", CURRENT_TASK, "/", TASK_COUNT)
	if err := s.UpdateMeta(task, "opensfm-process", CURRENT_TASK/100); err != nil {
		log.Printf("Failed to update meta: %d: %v\n", CURRENT_TASK, err)
		return err
	}
	CURRENT_TASK++

	log.Println("# 3 DensifyPointCloud", "scene.mvs", "-o", "scene_dense.mvs", "-w", mvsPath)
	cmd = exec.Command("DensifyPointCloud", "scene.mvs", "-o", "scene_dense.mvs", "-w", mvsPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		log.Println("DensifyPointCloud failed:", err)
		s.FailTask(task)
		return err
	}

	// 4
	log.Println("Updating meta for task:", task.Id, " - ", CURRENT_TASK, "/", TASK_COUNT)
	if err := s.UpdateMeta(task, "opensfm-process", CURRENT_TASK/100); err != nil {
		log.Printf("Failed to update meta: %d: %v\n", CURRENT_TASK, err)
		return err
	}
	CURRENT_TASK++

	log.Println("# 4 ReconstructMesh", "scene_dense.mvs", "-o", "scene_mesh.ply", "-w", mvsPath)
	cmd = exec.Command("ReconstructMesh", "scene_dense.mvs", "-o", "scene_mesh.ply", "-w", mvsPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		log.Println("ReconstructMesh failed:", err)
		s.FailTask(task)
		return err
	}

	// 5
	log.Println("Updating meta for task:", task.Id, " - ", CURRENT_TASK, "/", TASK_COUNT)
	if err := s.UpdateMeta(task, "opensfm-process", CURRENT_TASK/100); err != nil {
		log.Printf("Failed to update meta: %d: %v\n", CURRENT_TASK, err)
		return err
	}
	CURRENT_TASK++

	log.Println("# 5 RefineMesh", "scene.mvs", "-m", "scene_mesh.ply", "-o", "scene_dense_mesh_refine.mvs", "-w", mvsPath, "--scales", "1", "--max-face-area", "16")
	cmd = exec.Command("RefineMesh", "scene.mvs", "-m", "scene_mesh.ply", "-o", "scene_dense_mesh_refine.mvs", "-w", mvsPath, "--scales", "1", "--max-face-area", "16")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		log.Println("RefineMesh failed:", err)
		s.FailTask(task)
		return err
	}

	// 6
	log.Println("Updating meta for task:", task.Id, " - ", CURRENT_TASK, "/", TASK_COUNT)
	if err := s.UpdateMeta(task, "opensfm-process", CURRENT_TASK/100); err != nil {
		log.Printf("Failed to update meta: %d: %v\n", CURRENT_TASK, err)
		return err
	}
	CURRENT_TASK++

	log.Println("# 6 TextureMesh", "scene_dense.mvs", "-m", "scene_dense_mesh_refine.ply", "-o", "scene_dense_mesh_refine_texture.mvs", "-w", mvsPath, "--export-type", "obj")
	cmd = exec.Command("TextureMesh", "scene_dense.mvs", "-m", "scene_dense_mesh_refine.ply", "-o", "scene_dense_mesh_refine_texture.mvs", "-w", mvsPath, "--export-type", "obj")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		log.Println("TextureMesh failed:", err)
		s.FailTask(task)
		return err
	}

	// 7
	log.Println("Updating meta for task:", task.Id, " - ", CURRENT_TASK, "/", TASK_COUNT)
	if err := s.UpdateMeta(task, "opensfm-process", CURRENT_TASK/100); err != nil {
		log.Printf("Failed to update meta: %d: %v\n", CURRENT_TASK, err)
		return err
	}
	CURRENT_TASK++

	fileName := filepath.Join(mvsPath, "final_model")
	fmt.Println("blender", "-b", "-P", "./bin/convert_obj_to_glb.py", "--", filepath.Join(mvsPath, "scene_dense_mesh_refine_texture.obj"), fileName)
	cmd = exec.Command("blender", "-b", "-P", "./bin/convert_obj_to_glb.py", "--", filepath.Join(mvsPath, "scene_dense_mesh_refine_texture.obj"), fileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		log.Println("MeshConversion failed:", err)
		s.FailTask(task)
		return err
	}

	mesh, err := s.appFileService.Save(&models.AppFile{
		Url:      fileName + ".glb",
		Filename: "final_model.glb",
		TaskId:   task.Id,
		FileType: "mesh",
	})

	if err != nil {
		log.Printf("Failed to save mesh: %v\n", err)
		return err
	}

	task.Mesh = mesh
	task.Completed = true
	task.Status = models.SUCCESS

	if _, err := s.UpdateTask(task); err != nil {
		log.Printf("Failed to update task: %v\n", err)
		return err
	}

	log.Println("Task updated successfully.")
	log.Printf("Processing completed in %s\n", time.Since(startTime))
	return nil
}

func (s *TaskServiceImpl) GetTaskFiles(taskID uint, fileType string) ([]models.AppFile, error) {
	files, err := s.appFileService.GetTaskFiles(taskID, fileType)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *TaskServiceImpl) GetTaskFile(taskID uint, fileType string) (*models.AppFile, error) {
	file, err := s.appFileService.GetTaskFile(taskID, fileType)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *TaskServiceImpl) FullyLoadTask(task *models.Task) (*models.Task, error) {
	files, err := s.GetTaskFiles(task.Id, "upload")
	if err != nil {
		return nil, err
	}
	task.Images = files

	mesh, err := s.GetTaskFile(task.Id, "mesh")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			task.Mesh = nil
		} else {
			return nil, err
		}
	} else {
		task.Mesh = mesh
	}

	return task, nil
}
