package main

import (
	"fmt"
	"math/rand"
)

type Song struct {
	Title   string
	Artist  string
	Raiting int
	Length  int
}

type Node struct {
	Next *Node
	Prev *Node
	Song *Song
}

type Playlist struct {
	Head *Node
	Tail *Node

	Length int
}

type PlayStrategy interface {
	Next(*Node, *Playlist) *Node
	Name() string
}

type LinearStrategy struct{}

type ShuffleStrategy struct {
	List map[*Node]bool
}

type RepeatOneStrategy struct{}

type ByRaitingStrategy struct{}

func (l *LinearStrategy) Next(current *Node, playlist *Playlist) *Node {
	if current == nil || current.Next == nil {
		return playlist.Head
	}

	return current.Next
}

func (l *LinearStrategy) Name() string {
	return "Linear"
}

func NewShuffleStrategy() *ShuffleStrategy {
	return &ShuffleStrategy{List: make(map[*Node]bool)}
}

func (l *ShuffleStrategy) Next(current *Node, playlist *Playlist) *Node {
	if playlist.Length == 0 {
		return nil
	}

	if len(l.List) >= playlist.Length {
		l.List = make(map[*Node]bool)
	}

	unlistened := []*Node{}

	for current := playlist.Head; current != playlist.Tail; current = current.Next {
		if !l.List[current] {
			unlistened = append(unlistened, current)
		}
	}

	if len(unlistened) == 0 {
		return nil
	}

	picked_song := unlistened[rand.Intn(len(unlistened))]

	l.List[picked_song] = true
	return picked_song
}

func (l *ShuffleStrategy) Name() string {
	return "Shuffle"
}

func (l *RepeatOneStrategy) Next(current *Node, playlist *Playlist) *Node {
	if current == nil && playlist.Head != nil {
		return playlist.Head
	}
	return current
}

func (l *RepeatOneStrategy) Name() string {
	return "Repeat"
}

//	func (l *ByRaitingStrategy) Next(current *Node, playlist *Playlist) *Node {
//		return nil
//	}
//
//	func (l *ByRaitingStrategy) Name() string {
//		return "ByRaiting"
//	}
func NewPlaylist() *Playlist {
	return &Playlist{Length: 0}
}

func (p *Playlist) AddSong(song *Song) {
	node := &Node{Next: nil, Prev: nil, Song: song}

	if p.Length == 0 {

		p.Head = node
		p.Tail = node

	} else {
		node.Prev = p.Tail
		p.Tail = node

		node.Prev.Next = node
	}
	p.Length += 1
}

func (p *Playlist) RemoveSong(title string) bool {
	current := p.Head
	for current != nil {
		if current.Song.Title == title {
			if current.Prev != nil {
				current.Prev.Next = current.Next
			} else {
				p.Head = current.Next
			}

			if current.Next != nil {
				current.Next.Prev = current.Prev
			} else {
				p.Tail = current.Prev
			}

			p.Length--
			return true
		}
		current = current.Next
	}
	return false
}

func (p *Playlist) Display() {
	fmt.Println("\n📋 Текущий плейлист:")
	current := p.Head
	index := 1
	for current != nil {
		fmt.Printf("  %d. %s - %s (★%d, %d сек)\n",
			index, current.Song.Artist, current.Song.Title,
			current.Song.Raiting, current.Song.Length)
		current = current.Next
		index++
	}
	fmt.Printf("Всего песен: %d\n", p.Length)
}

type PlayerState int

const (
	Stopped = iota
	Playing
	Pause
)

type Player struct {
	playerState PlayerState

	strategy *PlayStrategy

	playlist *Playlist

	volume int
}

func (p *Player) Play() {
	p.playerState = Playing
	p.strategy.Next(nil, p.playlist)
}

func (p *Player) setStrategy(strategy *PlayStrategy) {
	p.strategy = strategy
}

func (p *Player) Pause() {
	p.playerState = Pause
}

func (p *Player) Stop() {
	p.playerState = Stopped
}

func (p *Player) SetVolume(volume int) {
	if volume > 100 {
		volume = 100
	} else if volume < 0 {
		volume = 0
	}

	p.volume = volume
}

func (p *Player) ShowStatus() {
	fmt.Printf("Статус: %s\n", p.playerState)
	fmt.Printf("Стратегия: %s\n", p.strategy.Name())
	fmt.Printf("Громкость: %d\n", p.volume)
}

func (p *Player) Display() {
	p.playlist.Display()
}

func (p *Player) AddSong(song *Song) {
	p.playlist.AddSong(song)
}

func (p *Player) Next() {
	p.strategy.Next(p.playlist.Head.Next, p.playlist)
}

func main() {
	// Начинаем воспроизведение
	fmt.Println("\n1. Начинаем воспроизведение:")
	player.Play()
	player.Next()
	player.Next()

	// Меняем стратегию
	fmt.Println("\n2. Меняем на случайное воспроизведение:")
	player.SetStrategy(NewShuffleStrategy())
	player.Next()
	player.Next()

	// Тестируем паузу и громкость
	fmt.Println("\n3. Тест паузы и громкости:")
	player.Pause()
	player.SetVolume(75)
	player.Play()

	// Показываем историю
	player.ShowHistory()

	// Тестируем повтор одной песни
	fmt.Println("\n4. Тест повторения одной песни:")
	player.SetStrategy(RepeatOneStrategy{})
	player.Next()
	player.Next() // Должна играть та же песня

	// Управление плейлистом
	fmt.Println("\n5. Управление плейлистом:")
	player.AddSong(&Song{"Hotel California", "Eagles", 5, 391})
	player.RemoveSong("Yesterday")
	playlist.Display()

	// Финальный статус
	fmt.Println("\n--- Финальный статус ---")
	player.ShowStatus()
	player.ShowHistory()

	// Останавливаем
	player.Stop()
	player.ShowStatus()
}
