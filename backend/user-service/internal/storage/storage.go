package storage

type Storage interface {
    CreateUser(uid uint) error
    
}