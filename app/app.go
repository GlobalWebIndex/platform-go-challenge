package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"x-gwi/app/auth"
	"x-gwi/app/instance"
	"x-gwi/app/logs"
	"x-gwi/app/pki"
	"x-gwi/app/server"
	"x-gwi/app/storage"
)

type App struct {
	config  *Config
	inst    *instance.Instance
	pki     *pki.PKI
	auth    *auth.Auth
	storage *storage.AppStorage
	server  *server.Server

	cancel context.CancelFunc
	// chErr  chan error
	// chMsg  chan string
}

func Main() {
	if err := runApp(context.Background(), newConfig()); err != nil {
		logs.Error().Err(err).Msg("app exit ERROR")
		os.Exit(1)
	}

	logs.Info().Msg("app exit OK")
	os.Exit(0)
}

func runApp(ctx context.Context, config *Config) error {
	app, ctx, err := initApp(ctx, config)
	if err != nil {
		return err
	}

	app.server.Serve(ctx)

	// app.testLoad(ctx, 3) //nolint:gomnd

	<-ctx.Done()

	return nil
}

func (app *App) testLoad(ctx context.Context, xTimes int) { //nolint:revive,unused
	/* "x-gwi/test/load"
	go func() {
		time.Sleep(3 * time.Second) //nolint:gomnd // time to start server

		st, err := load.NewLoadHTTP(ctx, app.server, load.LoadNONE, xTimes) // LoadUA LoadNONE
		if err != nil {
			logs.Error().Err(err).Send()
		}

		// logs.Debug().Msg("closing app after tests")
		// app.close(ctx)
		logs.Debug().Interface("load-http", st).Send()
	}() */
}

func (app *App) closeApp(ctx context.Context) {
	if app.cancel != nil {
		app.cancel()
	}

	if app.server != nil {
		app.server.Close(ctx)
	}
}

func initApp(ctx context.Context, config *Config) (*App, context.Context, error) {
	var err error
	// 1
	if !config.Valid() {
		err = fmt.Errorf("invalid config") //nolint:goerr113
		logs.Error().Err(err).Interface("config", config).Send()
		context.Background()

		return nil, nil, err
	}
	// 2
	app := &App{ //nolint:exhaustruct
		config: config,
		inst:   config.Inst,
	}
	if app.inst.Mode() == instance.ModeDev.String() ||
		app.inst.Mode() == instance.ModeTest.String() {
		logs.Debug().Interface("config", app.config).Send()
	}
	// 3
	ctx, app.cancel = context.WithCancel(ctx)
	// 4
	app.interrupt(ctx)
	// 5
	if app.pki, err = pki.NewPKI(); err != nil {
		return nil, nil, fmt.Errorf("pki.NewPKI: %w", err)
	}
	// 6
	if app.auth, err = auth.NewAuth(ctx, config.Auth, app.inst); err != nil {
		return nil, nil, fmt.Errorf("auth.NewAuth: %w", err)
	}
	// 7
	if app.storage, err = storage.NewAppStorage(ctx, app.config.Storage, app.inst, app.auth); err != nil {
		return nil, nil, fmt.Errorf("storage.NewStorage: %w", err)
	}
	// 8
	if app.server, err = server.NewServer(ctx, app.config.Server, app.inst, app.pki); err != nil {
		return nil, nil, fmt.Errorf("server.NewServer: %w", err)
	}
	// 9
	if err = app.registerServices(app.server.ServiceRegistrar()); err != nil {
		return nil, nil, fmt.Errorf("app.registerServices: %w", err)
	}
	// ok
	return app, ctx, nil
}

/*
func (app *App) Valid() bool {
	return app.config.Valid() &&
		app.auth.Valid() &&
		app.storage.Valid() &&
		app.server.Valid()
}
*/

func (app *App) interrupt(ctx context.Context) {
	// app.chErr = make(chan error, 2)  //nolint:gomnd
	// app.chMsg = make(chan string, 8) //nolint:gomnd
	go func() {
		chSig := make(chan os.Signal, 1)
		signal.Notify(chSig, os.Interrupt)

		for {
			select {
			case sig := <-chSig:
				if sig == os.Interrupt {
					logs.Info().Msgf("os.Interrupt: (%v)", sig)
					app.closeApp(ctx)

					return
				}
			case <-ctx.Done():
				logs.Info().Msg("app.context.Done()")

				return
				/*
					case err := <-app.chErr:
						logs.Error().Err(err).Msg("app.chErr")
						app.closeApp(ctx)

						return
					case msg := <-app.chMsg:
						logs.Info().Msg(msg)
				*/
			} //nolint:wsl
		}
	}()
}
