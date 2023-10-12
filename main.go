package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func main() {
	db, err := InitDataBase(259200000)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(len(db.Table))
	fmt.Println(db)

	/*db.Table["test2"] = Task{
		TagTeam:   []string{"testtag"},
		Timestamp: 100000000000,
	}
	err = db.Save()
	if err != nil {
		fmt.Println(err.Error())
	}*/
	fmt.Println(db)
}

// Структура задачи, которые лежат в хранилище
type Task struct {
	TagTeam   []string
	Timestamp int64
}

// Структура хранилища задач (хранятся по ключу IdReadble) для их подавления в репорте
type SimpleDataBase struct {
	Table          map[string]Task
	SupressionTime int64
}

// Функция для чтения данных из файла, чтобы не мучать сервак
func InitDataBase(supTime int64) (SimpleDataBase, error) {

	// читаем файл
	data, err := os.ReadFile("data.json")
	if err != nil {
		return SimpleDataBase{}, fmt.Errorf("read of file: %v", err)
	}
	// парсим json в структуру
	var db SimpleDataBase
	err = json.Unmarshal(data, &db)
	if err != nil {
		return SimpleDataBase{}, fmt.Errorf("json Unmarshal data from file: %v", err)
	}

	// Встраиваем время
	db.SupressionTime = supTime
	// Удаляем значения, находящиеся в хранилище слишком долго
	db.checkKeyValue()

	return db, nil
}

func (db *SimpleDataBase) Save() error {

	// создаем файл для вывода
	out, err := os.Create("data.json")
	if err != nil {
		return fmt.Errorf("write to file: %v", err)
	}
	// в конце программы, закрываем файл вывода
	defer out.Close()

	jsonData, err := json.Marshal(db)
	if err != nil {
		return fmt.Errorf("encoding data:%s", jsonData)
	}

	_, err = out.Write(jsonData)
	if err != nil {
		return fmt.Errorf("write to file: %v", err)
	}

	return nil
}

// Функция проверяет значения в хранилище и удаляет те, которые лежат там больше установленного SupressionTime
func (db *SimpleDataBase) checkKeyValue() {
	nowTime := time.Now().UnixMilli()
	fmt.Println(nowTime)
	for id, task := range db.Table {
		fmt.Println(id)
		if task.Timestamp+db.SupressionTime < nowTime {
			fmt.Println("deleting")
			delete(db.Table, id)
		}
	}
}
