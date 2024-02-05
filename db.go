package main

import (
    "fmt"

    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

type todo struct {
    Id int `sql:"AUTO_INCREMENT" gorm:"primaryKey"`
    User string
    Task string
    Completed bool
}

func ConnectDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to database!")
    }

    return db
}

func addTodo(user, content string) error {
    newTodo := todo{User: user, Task: content, Completed: false}

    fmt.Println("Adding todo: ", newTodo)
    result := db.Create(&newTodo)
    if result.Error != nil {
        return result.Error
    }

    return nil
}

func getTodos(user string) ([]todo, error) {
    var todo []todo
    result := db.Where("user = ?", user).Find(&todo)
    if result.Error != nil {
        return nil, result.Error
    }

    return todo, nil
}

func getAllTodos() ([]todo, error) {
    var todo []todo
    
    result := db.Find(&todo)
    if result.Error != nil {
        return nil, result.Error
    }
    return todo, nil
}

func completeTodoById(id int) error {
    return db.Model(&todo{}).Where("Id = ?", id).Update("completed", true).Error
}

func deleteTodos(user string) error {
    return db.Where("user = ?", user).Delete(&todo{}).Error
}

func updateTodoById(id int, content string) error {
    return db.Model(&todo{}).Where("Id = ?", id).Update("content", content).Error
}
