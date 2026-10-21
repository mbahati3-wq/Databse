package database
import (
	" _ github.com/lib/pq" // postgres driver
	"pg"
	"os"
	"net/http"

)

type DbConfig struct {
	DbName	 string
	DbUserName string
	DbPassword string
	DbWorkingDir string
	DbPort uint32
	DbSSLMode string
	DbTimeZone string
}
type Error struct {
	Sucess bool
	Message string
	ExitCode uint32
	Recomendation string
}
func CreateDB(databse_config DbConfig) (Error) {
	pq

}

func CreateTable(table_name string, table_schema string) (Error) {
	