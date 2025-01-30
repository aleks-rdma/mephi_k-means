package main

import (
	"math"
	"math/rand"
	"time"
)

func vector_sum(a []float64, b []float64, dim int) []float64 {
	result := make([]float64, dim)
	for i := 0; i < dim; i++ {
		result[i] = a[i] + b[i]
	}
	return result
}

func vector_devide_by_number(a []float64, b int, dim int) []float64 {
	result := make([]float64, dim)
	for i := 0; i < dim; i++ {
		result[i] = a[i] / float64(b)
	}
	return result
}

func start_clusters(data *[][]float64, points_number int, clusters_number int) *[][]float64 {
	rand.Seed(time.Now().UnixNano())
	cluster_centres := make([][]float64, clusters_number)
	for i := 0; i < clusters_number; i++ {
		random_cluster := rand.Intn(points_number)
		cluster_centres[i] = (*data)[random_cluster]
	}
	return &cluster_centres
}

func distance_e(a []float64, b []float64, dim int) float64 {
	var summ float64
	for i := 0; i < dim; i++ {
		summ += math.Pow(a[i]-b[i], 2)
	}
	return summ
}

func point_to_cluster(data *[][]float64, cluster_centres *[][]float64, points_number int, dim int, clusters_number int) *[]int {
	clusters := make([]int, points_number)
	var cur_cluster int
	var p, min_dist float64
	cur_cluster = 0
	for i := 0; i < points_number; i++ {
		min_dist = 99999
		for j := 0; j < clusters_number; j++ {
			p = distance_e((*data)[i], (*cluster_centres)[j], dim)
			if p < min_dist {
				min_dist = p
				cur_cluster = j
			}
		}
		clusters[i] = cur_cluster
	}
	return &clusters
}

func new_clusters_centres(data *[][]float64, curr_clusters *[]int, points_number int, clusters_number int, dim int) *[][]float64 {
	new_centres := make([][]float64, clusters_number)
	for i := 0; i < clusters_number; i++ {
		new_centres[i] = make([]float64, dim)
		for j := 0; j < dim; j++ {
			new_centres[i][j] = 0
		}
	}
	count_points_in_clusters := make([]int, clusters_number)
	for i := 0; i < clusters_number; i++ {
		count_points_in_clusters[i] = 0
	}
	var cluster int
	var value []float64
	for i := 0; i < points_number; i++ {
		cluster = (*curr_clusters)[i]
		value = (*data)[i]
		new_centres[cluster] = vector_sum(new_centres[cluster], value, dim)
		count_points_in_clusters[cluster] += 1
	}
	result := make([][]float64, clusters_number)
	for i := 0; i < clusters_number; i++ {
		result[i] = vector_devide_by_number(new_centres[i], count_points_in_clusters[i], dim)
	}
	return &result
}

func is_equal_clusters_centres(a *[][]float64, b *[][]float64, clusters_number int, max_diff float64, dim int) bool {
	var flag bool
	flag = true
	for i := 0; i < clusters_number; i++ {
		if distance_e((*a)[i], (*b)[i], dim) > max_diff {
			flag = false
		}
	}
	return flag
}

func clustering(data *[][]float64, points_number int, clusters_number int, dim int) *[]point_cluster {
	result := make([]point_cluster, points_number)
	ptr_prev_cluster_centres := start_clusters(data, points_number, clusters_number)
	var ptr_clusters *[]int
	for {
		ptr_clusters = point_to_cluster(data, ptr_prev_cluster_centres, points_number, dim, clusters_number)
		ptr_new_cluster_centres := new_clusters_centres(data, ptr_clusters, points_number, clusters_number, dim)
		if is_equal_clusters_centres(ptr_new_cluster_centres, ptr_prev_cluster_centres, clusters_number, 0.001, dim) {
			break
		}
		copy(*ptr_prev_cluster_centres, *ptr_new_cluster_centres)
	}
	for i := 0; i < points_number; i++ {
		result[i].point = (*data)[i]
		result[i].cluster = (*ptr_clusters)[i]
	}
	return &result
}
