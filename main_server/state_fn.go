package main

type stateFn func(w *World) stateFn
