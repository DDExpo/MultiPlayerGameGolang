package main

import "sync"

type Vec2 struct {
	X float32
	Y float32
}

type SpatialHash struct {
	cellSize   float32
	invCell    float32
	grid       map[int64]map[string]Vec2
	playerCell map[string]int64
	mu         sync.RWMutex
}

func NewSpatialHash(cellSize float32) *SpatialHash {
	return &SpatialHash{
		cellSize:   cellSize,
		invCell:    1.0 / cellSize,
		grid:       make(map[int64]map[string]Vec2, 1024),
		playerCell: make(map[string]int64, 1024),
	}
}

func (s *SpatialHash) hash(x, y float32) int64 {
	cx := int32(x * s.invCell)
	cy := int32(y * s.invCell)
	return cellKey(cx, cy)
}

func (s *SpatialHash) Update(id string, x, y float32) {
	key := s.hash(x, y)

	if oldKey, ok := s.playerCell[id]; ok && oldKey != key {
		if cell := s.grid[oldKey]; cell != nil {
			delete(cell, id)
			if len(cell) == 0 {
				delete(s.grid, oldKey)
			}
		}
	}

	cell := s.grid[key]
	if cell == nil {
		cell = make(map[string]Vec2, 8)
		s.grid[key] = cell
	}

	cell[id] = Vec2{x, y}
	s.playerCell[id] = key
}

func (s *SpatialHash) Remove(id string) {
	key, ok := s.playerCell[id]
	if !ok {
		return
	}
	if cell := s.grid[key]; cell != nil {
		delete(cell, id)
		if len(cell) == 0 {
			delete(s.grid, key)
		}
	}

	delete(s.playerCell, id)
}

func (s *SpatialHash) GetCell(x, y float32) map[string]Vec2 {
	return s.grid[s.hash(x, y)]
}

func cellKey(cx, cy int32) int64 {
	return int64(cx)<<32 | int64(uint32(cy))
}
