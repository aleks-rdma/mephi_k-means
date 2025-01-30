package main

import (
	"fmt"
	"time"
)

func main() {
	var ptrdata *[][]float64
	var err error
	var dim, points_number, clusters_number int
	clusters_number = 2
	ptrdata, dim, points_number, err = read("4_dim_random.csv")
	if err == nil {
		startTime := time.Now()
		clusters2 := clustering2(ptrdata, points_number, clusters_number, dim, 8)
		elapsedTime := time.Since(startTime)
		//write2("output.csv", clusters2, dim)
		startTime2 := time.Now()
		clusters := clustering(ptrdata, points_number, clusters_number, dim)
		elapsedTime2 := time.Since(startTime2)
		//write2("output2.csv", clusters, dim)
		fmt.Println(elapsedTime, elapsedTime2)
		fmt.Println(clusters2)
		fmt.Println(clusters)
	}

}
