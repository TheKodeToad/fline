package gateway

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/TheKodeToad/fline/internal/config"
	"github.com/gorilla/websocket"
)

type sessionHandle struct {
	shutdown         chan<- struct{}
	shutdownFinished <-chan struct{}
}

// Gateway keeps track of gateway sessions and can be shut down with Shutdown.
type Gateway struct {
	sessionsLock sync.Mutex
	sessions     map[sessionHandle]bool
}

var upgrader websocket.Upgrader

func (g *Gateway) ServeHTTP(conf *config.Config, w http.ResponseWriter, r *http.Request) {
	sesh, err := startSession(conf, w, r)
	if err != nil {
		slog.Warn("failed to start session", slog.Any("err", err))
		return
	}

	shutdown := make(chan struct{}, 1)
	shutdownFinished := make(chan struct{}, 1)
	g.registerSession(sessionHandle{shutdown, shutdownFinished})

	sesh.run(shutdown)

	err = sesh.close()
	if err != nil {
		slog.Warn("failed to close session", slog.Any("err", err))
	}

	g.unregisterSession(sessionHandle{shutdown, shutdownFinished})
	shutdownFinished <- struct{}{}
}

func (g *Gateway) registerSession(sesh sessionHandle) {
	g.sessionsLock.Lock()
	defer g.sessionsLock.Unlock()

	if g.sessions == nil {
		g.sessions = map[sessionHandle]bool{}
	}
	g.sessions[sesh] = true
}

func (g *Gateway) unregisterSession(sesh sessionHandle) {
	g.sessionsLock.Lock()
	defer g.sessionsLock.Unlock()

	if g.sessions != nil {
		delete(g.sessions, sesh)
	}
}

func (g *Gateway) Shutdown() {
	// NOTE: unfortunately we can't defer the unlock as if we only unlock it after wg.Wait() it causes a deadlock -
	// we are waiting for the sessions to shut down which are waiting
	g.sessionsLock.Lock()

	if g.sessions == nil {
		g.sessionsLock.Unlock()
		return
	}

	var wg sync.WaitGroup

	wg.Add(len(g.sessions))
	for sesh := range g.sessions {
		go func() {
			sesh.shutdown <- struct{}{}
			<-sesh.shutdownFinished
			wg.Done()
		}()
	}

	g.sessions = nil
	g.sessionsLock.Unlock()

	wg.Wait()
}
