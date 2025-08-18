package repository

import (
	"database/sql"
	"file-manager/entity"
)

func FileSystemRowMapper(rows *sql.Rows) (entity.FileSystem, error) {
	fileSystem := entity.FileSystem{}
	if err := rows.Scan(
		&fileSystem.Id,
		&fileSystem.OwnerId,
		&fileSystem.ParentId,
		&fileSystem.Rights,
		&fileSystem.IsFile,
		&fileSystem.Name,
		&fileSystem.Path,
		&fileSystem.CreatedAt,
		&fileSystem.UpdatedAt,
		&fileSystem.Deleted); err != nil {
		return fileSystem, err
	}
	return fileSystem, nil
}
