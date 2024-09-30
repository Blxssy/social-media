package photo

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/Blxssy/social-media/photo-service/api/photo"
	"github.com/Blxssy/social-media/photo-service/intenal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/grpc"
)

// Photo - интерфейс, описывающий методы для работы с фотографиями.
// Методы:
// - UploadPhoto: загружает фотографию и возвращает её ID или ошибку.
// - GetPhoto: возвращает данные фотографии по её ID.
// - GetUserPhotos: возвращает список фотографий для конкретного пользователя.
type Photo interface {
	UploadPhoto(ctx context.Context, photoData []byte, filename string, uid uint64) (uint64, error)
	GetPhoto(context context.Context, photoID uint64) (*models.Photo, error)
	GetUserPhotos(context context.Context, userID uint64) ([]*models.Photo, error)
}

// ServerAPI - структура, реализующая gRPC-сервер для работы с фотографиями.
// Поля:
// - UnimplementedPhotoServiceServer: встраиваемая структура для совместимости с gRPC.
// - photo: интерфейс Photo, отвечающий за операции с фото.
// - s3svc: клиент для взаимодействия с S3 (или другим облачным хранилищем).
type ServerAPI struct {
	pb.UnimplementedPhotoServiceServer
	photo Photo
	s3svc *s3.Client
}

// Register - функция регистрации gRPC-сервера для работы с фото.
// Параметры:
// - grpcServer: указатель на gRPC сервер, который нужно зарегистрировать.
// - photo: интерфейс Photo, предоставляющий реализацию логики работы с фотографиями.
func Register(grpcServer *grpc.Server, photo Photo) {
	pb.RegisterPhotoServiceServer(grpcServer, &ServerAPI{photo: photo})
}

// UploadPhoto - обработчик RPC-запроса для загрузки фотографии.
// Параметры:
// - ctx: контекст для контроля запроса.
// - req: запрос, содержащий данные фотографии, такие как байты файла, имя файла и идентификатор пользователя.
// Возвращает ответ с идентификатором загруженной фотографии или ошибку.
func (s *ServerAPI) UploadPhoto(ctx context.Context, req *pb.UploadPhotoRequest) (*pb.UploadPhotoResponse, error) {
	reader := bytes.NewReader(req.PhotoData)

	fileKey := fmt.Sprintf("photos/%s", req.Filename)
	_, err := s.s3svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("your-bucket-name"),
		Key:    aws.String(fileKey),
		Body:   reader,
	})
	if err != nil {
		return nil, fmt.Errorf("could not upload photo: %v", err)
	}

	return &pb.UploadPhotoResponse{}, nil
}

// GetPhoto - обработчик RPC-запроса для получения информации о фотографии.
// Параметры:
// - ctx: контекст для контроля запроса.
// - req: запрос, содержащий ID фотографии.
// Возвращает ответ с данными фотографии или ошибку.
func (s *ServerAPI) GetPhoto(context context.Context, req *pb.GetPhotoRequest) (*pb.GetPhotoResponse, error) {
	return &pb.GetPhotoResponse{}, nil
}

// GetUserPhotos - обработчик RPC-запроса для получения списка фотографий пользователя.
// Параметры:
// - ctx: контекст для контроля запроса.
// - req: запрос, содержащий ID пользователя, для которого нужно получить список фотографий.
// Возвращает ответ с массивом фотографий или ошибку.
func (s *ServerAPI) GetUserPhotos(context context.Context, req *pb.GetUserPhotosRequest) (*pb.GetUserPhotosResponse, error) {
	return &pb.GetUserPhotosResponse{}, nil
}
