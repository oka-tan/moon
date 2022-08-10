package main

import (
	"context"
	"database/sql"
	"log"
	"moon/config"
	"moon/db"
	"moon/lnx"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	log.Println("Starting Moon")

	conf, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.PostgresConfig.ConnectionString)))
	pg := bun.NewDB(sqldb, pgdialect.New())

	lnxService := lnx.NewService(conf.LnxConfig.Host, conf.LnxConfig.Port)

	napTime, err := time.ParseDuration(conf.LnxConfig.NapTime)
	if err != nil {
		napTime = 20 * time.Minute
	}

	for _, board := range conf.Boards {
		indexTracker := db.IndexTracker{
			Board:        board.Name,
			LastModified: time.UnixMicro(0),
			PostNumber:   0,
		}

		if board.ForceRecreate {
			_, err := pg.NewInsert().
				Model(&indexTracker).
				On("CONFLICT (board) DO UPDATE SET last_modified = EXCLUDED.last_modified, post_number = EXCLUDED.post_number").
				Returning("NULL").
				Exec(context.Background())

			if err != nil {
				log.Fatalf("Error creating index tracker for board %s", board.Name)
			}
		} else {
			_, err := pg.NewInsert().
				Model(&indexTracker).
				On("CONFLICT DO NOTHING").
				Returning("NULL").
				Exec(context.Background())

			if err != nil {
				log.Fatalf("Error creating index tracker for board %s", board.Name)
			}
		}

		lnxService.CreateIndex(board.Name, conf.LnxConfig.Configuration, board.ForceRecreate)
	}

	for {
		for _, board := range conf.Boards {
			dbPosts := make([]db.Post, 0, conf.LnxConfig.BatchSize)

			log.Printf("Indexing board %s\n", board.Name)

			maxTime := time.Now().Add(-5 * time.Second)

			tx, err := pg.BeginTx(context.Background(), &sql.TxOptions{})

			if err != nil {
				panic(err)
			}

			indexTracker := db.IndexTracker{}

			err = tx.NewSelect().
				Model(&indexTracker).
				Where("board = ?", board.Name).
				Scan(context.Background())

			if err != nil {
				tx.Rollback()
				panic(err)
			}

			if err := lnxService.Rollback(board.Name); err != nil {
				tx.Rollback()
				panic(err)
			}

			previousScrape := indexTracker.LastModified

			for {
				dbPosts = dbPosts[0:0]

				err := tx.NewSelect().
					Model(&dbPosts).
					Where("board = ?", board.Name).
					Where("last_modified < ?", maxTime).
					Where("(last_modified, post_number) > (?, ?)", indexTracker.LastModified, indexTracker.PostNumber).
					Order("last_modified ASC", "post_number ASC").
					Limit(conf.LnxConfig.BatchSize).
					Scan(context.Background())

				if err != nil {
					tx.Rollback()
					panic(err)
				}

				if len(dbPosts) == 0 {
					break
				}

				lastPost := dbPosts[len(dbPosts)-1]
				indexTracker.LastModified = lastPost.LastModified
				indexTracker.PostNumber = lastPost.PostNumber

				if err := lnxService.Upsert(dbPosts, board.Name, previousScrape); err != nil {
					tx.Rollback()
					panic(err)
				}
			}

			if err := lnxService.Commit(board.Name); err != nil {
				tx.Rollback()
				panic(err)
			}

			_, err = tx.NewUpdate().
				Model(&indexTracker).
				WherePK().
				Returning("NULL").
				Exec(context.Background())

			if err != nil {
				panic(err)
			}

			if err := tx.Commit(); err != nil {
				panic(err)
			}
		}

		log.Println("Napping")
		time.Sleep(napTime)
	}
}
