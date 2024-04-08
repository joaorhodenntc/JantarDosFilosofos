package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numFilosofos = 5
	numGarfos = 5
)

type Garfo struct {
	sync.Mutex
}

type Filosofo struct {
	id int
	garfoEsquerda *Garfo
	garfoDireita *Garfo
	quantidadeComidas int
	tempoTotalPensando time.Duration
	tempoTotalComendo time.Duration
}

type Mesa struct {
	garfos []*Garfo
}

func (m *Mesa) pegarGarfos(i int) {
	m.garfos[i].Lock()
	m.garfos[(i+1)%numGarfos].Lock()
}

func (m *Mesa) liberarGarfos(i int) {
	m.garfos[i].Unlock()
	m.garfos[(i+1)%numGarfos].Unlock()
}

func (f *Filosofo) pensar() {
	duracao := time.Duration(2) * time.Second
	fmt.Printf("Filósofo %d pensando por %v\n", f.id, duracao)
	time.Sleep(duracao)
	f.tempoTotalPensando += duracao
}

func (f *Filosofo) comer() {
	duracao := time.Duration(3) * time.Second
	f.quantidadeComidas++
	fmt.Printf("Filósofo %d comendo por %v (comidas: %d)\n", f.id,
	duracao, f.quantidadeComidas)
	time.Sleep(duracao)
	f.tempoTotalComendo += duracao
}

func (f *Filosofo) jantar(m *Mesa, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		f.pensar()
		m.pegarGarfos(f.id - 1)
		f.comer()
		m.liberarGarfos(f.id - 1)
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(numFilosofos)
	// Cria os garfos
	garfos := make([]*Garfo, numGarfos)
	for i := 0; i < numGarfos; i++ {
		garfos[i] = new(Garfo)
	}
	// Cria a mesa
	mesa := &Mesa{garfos: garfos}
	// Cria os filósofos
	filosofos := make([]*Filosofo, numFilosofos)
	start := time.Now()
	for i := 0; i < numFilosofos; i++ {
	filosofos[i] = &Filosofo{
		id: i + 1,
		garfoEsquerda: garfos[i],
		garfoDireita: garfos[(i+1)%numGarfos],
	}
}
	// Inicia o jantar
	for i := 0; i < numFilosofos; i++ {
		go filosofos[i].jantar(mesa, &wg)
	}
	// Aguarda todos os filósofos terminarem de jantar
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("\nChandry/Misra Dinner took %s\n\n", elapsed)
}