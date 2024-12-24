package main

import (
	"fmt"
	"go_parser_vacancies/utils"
	"sync"
)

func main() {
	searchText := "Golang"
	fmt.Println("Поиск вакансий...")

	vacancies, err := utils.FetchVacancies(searchText)
	if err != nil {
		fmt.Printf("Ошибка при получении списка вакансий: %v\n", err)
		return
	}

	fmt.Printf("Найдено %d вакансий:\n\n", len(vacancies))

	var wg sync.WaitGroup
	results := make(chan utils.VacancyResult, len(vacancies))

	// Ограничение количества параллельных горутин
	maxConcurrency := 20
	semaphore := make(chan struct{}, maxConcurrency)

	// Обрабатываем каждую вакансию в отдельной горутине
	for _, vacancy := range vacancies {
		wg.Add(1)

		// Используем семафор для ограничения параллельных горутин
		go func(vacancy utils.Vacancy) {
			defer wg.Done()

			// Получаем вакансию с помощью семафора
			semaphore <- struct{}{} // Захватываем семафор
			description := utils.FetchVacancyDetails(vacancy.ID)
			results <- utils.VacancyResult{
				ID:          vacancy.ID,
				Description: description,
			}
			<-semaphore // Освобождаем семафор
		}(vacancy)
	}

	// Ожидаем завершения всех горутин и закрываем канал
	go func() {
		wg.Wait()
		close(results)
	}()

	// Выводим результаты
	for result := range results {
		for _, vacancy := range vacancies {
			salary := utils.FormatSalary(vacancy.Salary.From, vacancy.Salary.To)
			if vacancy.ID == result.ID {
				fmt.Printf("Вакансия: %s\nКомпания: %s\nЗарплата: %v\nОпыт работы: %s\nЛокация: %s\nТип занятости: %s\nФормат работы: %s\nСсылка: %s\nОписание:\n     %s\n\n",
					vacancy.Name,
					vacancy.Employer.Name,
					salary,
					vacancy.Experience.Name,
					vacancy.Area.Name,
					vacancy.Employment.Name,
					vacancy.Schedule.Name,
					vacancy.URL,
					result.Description,
				)
			}
		}
	}
}
