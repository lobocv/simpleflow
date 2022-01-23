[![GoReportCard](https://goreportcard.com/badge/github.com/lobocv/simpleflow)](https://goreportcard.com/report/github.com/lobocv/simpleflow)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-95%25-brightgreen.svg?longCache=true&style=flat)</a>



# SimpleFlow

SimpleFlow is a a collection of generic functions and patterns that help building common concurrent workflows.


## Worker Pools

Worker pools provide a way to spin up a finite set of go routines to process items in a collection.
It supports processing items in `slices`, `maps` or `channels`. The functions `WorkerPoolFromSlice`, 
`WorkerPoolFromMap` and `WorkerPoolFromChan` all block until all workers finish processing.

## Fan-Out and Fan-In

`FanOut` and `FanIn` provide means of fanning-in and out channel to other channels. 

## Round Robin

`RoundRobin` distributes values from a channel over other channels in a round-robin fashion

## Batching

`BatchMap`, `BatchSlice` and `BatchChan` provide ways to break `maps`, `slices` and `channels` into smaller
components of at most `N` size.
