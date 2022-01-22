# SimpleFlow

SimpleFlow is a a collection of generic functions and patterns that help building common concurrent workflows.


## Worker Pools

Worker pools provide a way to spin up a finite set of go routines to process items in a collection.
It supports processing items in `slices`, `maps` or `channels`. The functions `WorkerPoolFromSlice`, 
`WorkerPoolFromMap` and `WorkerPoolFromChan` all block until all workers finish processing.

## Fan-Out and Fan-In

`FanOut` and FanIn` provide means of fanning-in and out channel to other channels. 

## Round Robin

`RoundRobin` distributes values from a channel over other channels in a round-robin fashion

## Batching

`BatchMap` and `BatchSlice` provide ways to break `maps` and `slices` into smaller components of at most `N` size.
