package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	wishbubbletea "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/ratelimiter"
	"golang.org/x/time/rate"
)

const (
	host        = "0.0.0.0"
	port        = "2222"
	keyPath     = "/data/host_key"
	maxSessions = 20
)

func main() {
	srv, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(keyPath),

		// Public portfolio: allow interactive access.
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithPasswordAuth(func(ctx ssh.Context, password string) bool {
			return true
		}),

		wish.WithMiddleware(
			wishbubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			ratelimiter.Middleware(ratelimiter.NewRateLimiter(
				rate.Every(time.Minute/10), 5, 1000,
			)),
			blockExecMiddleware(),
			loggingMiddleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Shutdown error", "error", err)
	}
}

// blockExecMiddleware rejects any SSH session that tries to run a command.
func blockExecMiddleware() wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			cmd := sess.Command()
			if len(cmd) > 0 {
				log.Warn("Blocked exec attempt",
					"remote", sess.RemoteAddr(),
					"cmd", sanitizeLog(cmd[0]),
				)
				_ = sess.Exit(1)
				return
			}
			next(sess)
		}
	}
}

// loggingMiddleware logs each connection with remote address and timing.
func loggingMiddleware() wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			start := time.Now()
			log.Info("Connection",
				"remote", sess.RemoteAddr(),
				"user", sanitizeLog(sess.User()),
			)
			next(sess)
			log.Info("Disconnected",
				"remote", sess.RemoteAddr(),
				"duration", time.Since(start).Round(time.Millisecond),
			)
		}
	}
}

// sanitizeLog strips non-printable characters before logging.
func sanitizeLog(s string) string {
	if len(s) > 64 {
		s = s[:64]
	}
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r >= 32 && r < 127 {
			out = append(out, r)
		}
	}
	return string(out)
}

func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := sess.Pty()

	m := newModel()
	m.width = pty.Window.Width
	m.height = pty.Window.Height

	opts := append(
		wishbubbletea.MakeOptions(sess),
		tea.WithAltScreen(),
	)

	return m, opts
}