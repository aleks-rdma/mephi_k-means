package main

import (
	"sync"
)

type point_cluster struct {
	point   []float64
	cluster int
}

func worker(wg *sync.WaitGroup, input_chan <-chan [][]float64, output_chan chan<- point_cluster, dim int, clusters_number int) {
	var cur_cluster int
	var p, min_dist float64
	cur_cluster = 0
	point_cluster1 := point_cluster{
		point:   nil,
		cluster: 0,
	}
	for data := range input_chan {
		min_dist = 99999
		for i := 0; i < clusters_number; i++ {
			p = distance_e(data[0], data[i+1], dim)
			if p < min_dist {
				min_dist = p
				cur_cluster = i
			}
		}
		point_cluster1.cluster = cur_cluster
		point_cluster1.point = data[0]
		output_chan <- point_cluster1
	}
	wg.Done()
}

func cluster_centres_to_workers(data [][]float64, cluster_centres [][]float64, data_to_workers chan<- [][]float64, clusters_number int, points_number int) {
	for i := 0; i < points_number; i++ {
		tmp_data := make([][]float64, clusters_number+1)
		tmp_data[0] = data[i]
		for j := 0; j < clusters_number; j++ {
			tmp_data[j+1] = cluster_centres[j]
		}
		data_to_workers <- tmp_data
	}
}

func new_clusters_centres2(data_from_worker <-chan point_cluster, dim int, points_number int, clusters_number int) (*[][]float64, *[]point_cluster) {
	var data []point_cluster
	var value point_cluster
	var point []float64
	var cluster int
	points_sum := make([][]float64, clusters_number)
	for i := 0; i < clusters_number; i++ {
		points_sum[i] = make([]float64, dim)
		for j := 0; j < dim; j++ {
			points_sum[i][j] = 0
		}
	}
	points_count := make([]int, clusters_number)
	for i := 0; i < clusters_number; i++ {
		points_count[i] = 0
	}
	for i := 0; i < points_number; i++ {
		value = <-data_from_worker
		data = append(data, value)
	}
	//for value := range data_from_forker {
	//	data = append(data, value)
	//}
	for i := 0; i < points_number; i++ {
		point = data[i].point
		cluster = data[i].cluster
		points_sum[cluster] = vector_sum(points_sum[cluster], point, dim)
		points_count[cluster] += 1
	}
	result := make([][]float64, clusters_number)
	for i := 0; i < clusters_number; i++ {
		result[i] = vector_devide_by_number(points_sum[i], points_count[i], dim)
	}
	return &result, &data
}

func clustering2(data *[][]float64, points_number int, clusters_number int, dim int, worker_count int) *[]point_cluster {
	ptr_prev_cluster_centres := start_clusters(data, points_number, clusters_number)
	var ptr_new_cluster_centres *[][]float64
	var result *[]point_cluster
	wg := sync.WaitGroup{}
	wg.Add(worker_count)

	data_to_workers := make(chan [][]float64, worker_count)
	data_from_workers := make(chan point_cluster, worker_count)
	for i := 0; i < worker_count; i++ {
		go worker(&wg, data_to_workers, data_from_workers, dim, clusters_number)
	}
	for {
		go cluster_centres_to_workers(*data, *ptr_prev_cluster_centres, data_to_workers, clusters_number, points_number)
		ptr_new_cluster_centres, result = new_clusters_centres2(data_from_workers, dim, points_number, clusters_number)
		if is_equal_clusters_centres(ptr_new_cluster_centres, ptr_prev_cluster_centres, clusters_number, 0.001, dim) {
			break
		}
		copy(*ptr_prev_cluster_centres, *ptr_new_cluster_centres)
	}
	close(data_to_workers)
	wg.Wait()
	close(data_from_workers)
	return result
}
