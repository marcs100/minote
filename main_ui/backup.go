package main_ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/marcs100/minote/minotedb"
)

func BackupNotes(dbFile string, parentWindow fyne.Window) {
	d := dialog.NewFolderOpen(func(backupDir fyne.ListableURI, err error) {
		if err == nil {
			if backupDir != nil {
				bytesWrtten, dbBakFile, err := minotedb.BackupDatabase(backupDir.Path(), dbFile)
				if err == nil {
					log.Printf("Database backed up - %d bytes written\n", bytesWrtten)
					dialog.ShowInformation("Backup successful", dbBakFile, parentWindow)
				} else {
					dialog.ShowError(err, parentWindow)
				}
			}
		} else {
			dialog.ShowError(err, parentWindow)
		}

	}, parentWindow)
	d.SetConfirmText("Backup")
	d.Show()
}
