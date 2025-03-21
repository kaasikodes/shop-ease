package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/kaasikodes/shop-ease/internal/store"
)


func UnSeed(s store.Storage, db *sql.DB) error {
	emails := []string{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	for _, u := range randomUsers {
		emails = append(emails, u.Email)
		
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func ()  {
		if err != nil {
			tx.Rollback()
		}
		
	}()
	err = s.Users.RemoveMultipleUsers(ctx, tx, emails)
	err = tx.Commit()
	if err != nil {
	   return err
	}


	return nil
}
func Seed(s store.Storage, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	// create default roles
	_, err := s.Roles.CreateDefaultRoles(ctx)
	if err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func ()  {
		if err != nil {
			tx.Rollback()
		}
		
	}()

	// create users - > with different roles and should be activated users
	for i, u := range randomUsers {
		role := store.DefaultRoles[i%len(store.DefaultRoles)]
		userRole := store.UserRole{
			ID: role.ID,
			Name: role.Name,
			IsActive: true,

		}
		if err = u.Password.Set(randomPassword); err != nil {
			
			return err
		}
		err =s.Users.Create(ctx, tx, &u, &userRole)
		if err != nil {
			return err
		}
		err = s.Users.Verify(ctx, tx, &u)
		if err != nil {
			return err
		}
		
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
	

}