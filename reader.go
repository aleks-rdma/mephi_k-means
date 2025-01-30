package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func read(path string) (*[][]float64, int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("ошибка при открытии файла: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	records, err := reader.ReadAll()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("ошибка при чтении: %w", err)
	}

	var data [][]float64
	var dim, points_number int
	points_number = 0
	for _, record := range records {
		var floatRow []float64
		points_number += 1
		dim = 0
		for _, value := range record {
			dim += 1
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("ошибка при преобразовании значения в: %w", err)
			}
			floatRow = append(floatRow, floatValue)
		}

		data = append(data, floatRow)
	}
	return &data, dim, points_number, nil
}

func write2(path string, data *[]point_cluster, dim int) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	len_data := len(*data)

	res := make([][]string, len_data)
	for i := 0; i < len_data; i++ {
		res[i] = make([]string, 2)
		for j := 0; j < dim; j++ {
			res[i][j] = strconv.Itoa(int((*data)[i].point[j]))
		}
		res[i] = append(res[i], strconv.Itoa(int((*data)[i].cluster)))
	}

	for _, row := range res {
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}
}
