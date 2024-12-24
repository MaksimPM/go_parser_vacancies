package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// Функция для удаления HTML-разметки и добавления разделений на новые строки
func stripHTMLTags(input string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(input))
	var result strings.Builder

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return result.String()
		case html.TextToken:
			result.WriteString(tokenizer.Token().Data)
		}
	}
}

// Функция для получения списка вакансий
func FetchVacancies(searchText string) ([]Vacancy, error) {
	url := fmt.Sprintf("https://api.hh.ru/vacancies?text=%s", searchText)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка запроса: %s", resp.Status)
	}

	var vacancies VacanciesResponse
	if err := json.NewDecoder(resp.Body).Decode(&vacancies); err != nil {
		return nil, err
	}

	return vacancies.Items, nil
}

// Функция для получения подробной информации о вакансии
func FetchVacancyDetails(id string) string {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("https://api.hh.ru/vacancies/%s", id)

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Ошибка запроса деталей вакансии %s: %v\n", id, err)
		return "Описание недоступно"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Вакансия с ID %s не найдена\n", id)
		return "Описание недоступно"
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка при получении деталей вакансии %s: %s\n", id, resp.Status)
		return "Описание недоступно"
	}

	var details VacancyDetails
	err = json.NewDecoder(resp.Body).Decode(&details)
	if err != nil {
		fmt.Printf("Ошибка декодирования JSON деталей вакансии %s: %v\n", id, err)
		return "Описание недоступно"
	}
	return stripHTMLTags(details.Description)
}

// Функция для обработки вакансии и получения результатов
func ProcessVacancy(id string, results chan<- VacancyResult, wg *sync.WaitGroup) {
	defer wg.Done()

	description := FetchVacancyDetails(id)
	results <- VacancyResult{
		ID:          id,
		Description: description,
	}
}

// Форматирование зарплаты
func FormatSalary(from, to int) string {
	if from == 0 && to == 0 {
		return "не указано"
	}
	if to == 0 {
		return fmt.Sprintf("%d", from)
	}
	return fmt.Sprintf("%d-%d", from, to)
}
