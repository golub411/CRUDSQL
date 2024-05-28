# CRUD-SQL

package main

import (
	"fmt"
	"log"
	"CrudSql"
)

func main() {
	// Открываем соединение с базой данных
	db, err := CrudSql.OpenDatabase("./path/to/your/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем таблицу
	columns := []string{"id INTEGER PRIMARY KEY", "name TEXT"}
	if err := db.CreateTable("your_table", columns); err != nil {
		log.Fatal(err)
	}

	// Вставляем значение
	values := []interface{}{1, "John Doe"}
	if err := db.InsertValue("your_table", []string{"id", "name"}, values); err != nil {
		log.Fatal(err)
	}

	// Выбираем значения
	results, err := db.SelectValue("your_table", []string{"id", "name"})
	if err != nil {
		log.Fatal(err)
	}

	// Выводим результаты
	for _, result := range results {
		fmt.Println(result)
	}

	// Обновляем значение
	set := map[string]interface{}{"name": "Jane Doe"}
	where := map[string]interface{}{"id": 1}
	if err := db.UpdateValue("your_table", set, where); err != nil {
		log.Fatal(err)
	}

	// Удаляем значение
	if err := db.DeleteValue("your_table", where); err != nil {
		log.Fatal(err)
	}
