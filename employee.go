package controllers

import (
	"CRUD_API/app/https/domain"
	"database/sql"
	"fmt"
	"net/http"

)

var db *sql.DB

func NewSystemController(e *echo.Echo, Db *sql.DB) {
	db = Db
	e.GET("/employee/:id", GetEmployeeById)
	e.POST("/employee", InsertEmployee)
	e.PUT("/employee/:id", UpdateEmployee)
	e.DELETE("/employee/:id", DeleteEmployee)
}

func GetEmployeeById(c echo.Context) error {
	requested_id := c.Param("id")
	fmt.Println(requested_id)

	var (
		name   string
		id     string
		salary string
		age    string
	)
	err := db.QueryRow("SELECT id,employee_name, employee_salary, employee_age  FROM employee WHERE id = ?", requested_id).Scan(&id, &name, &salary, &age)
	if err != nil {
		fmt.Println(err)
	}

	response := domain.Employee{Id: id, Name: name, Salary: salary, Age: age}
	return c.JSON(http.StatusOK, response)
}

func InsertEmployee(c echo.Context) error {

	emp := new(domain.Employee)
	if err := c.Bind(emp); err != nil {
		return err
	}
	sql := "INSERT INTO employee(employee_name, employee_salary, employee_age ) VALUES( ?, ?, ?)"
	stmt, err := db.Prepare(sql)

	if err != nil {
		fmt.Print(err.Error())
	}
	defer stmt.Close()
	result, err2 := stmt.Exec(emp.Name, emp.Salary, emp.Age)

	if err2 != nil {
		panic(err2)
	}
	fmt.Println(result.LastInsertId())
	return c.JSON(http.StatusCreated, emp.Name)
}

func UpdateEmployee(c echo.Context) error {

	requested_id := c.Param("id")
	emp := new(domain.Employee)

	if err := c.Bind(emp); err != nil {
		return err
	}
	sql := "UPDATE employee SET employee_name = ? ,employee_salary= ?, employee_age = ? WHERE id= ? "
	stmt, err := db.Prepare(sql)

	if err != nil {
		fmt.Print(err.Error())
	}
	defer stmt.Close()

	result, err2 := stmt.Exec(emp.Name, emp.Salary, emp.Age, requested_id)

	if err2 != nil {
		panic(err2)
	}
	fmt.Println(result.LastInsertId())
	return c.JSON(http.StatusCreated, emp.Name)
}

func DeleteEmployee(c echo.Context) error {
	requested_id := c.Param("id")
	sql := "Delete FROM employee Where id = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Println(err)
	}
	result, err2 := stmt.Exec(requested_id)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(result.RowsAffected())

	return c.JSON(http.StatusOK, "Deleted")
}
