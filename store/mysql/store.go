package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type Config struct {
	MysqlUsername string
	MysqlPassword string
	MysqlHost     string
	MysqlPort     int
	MysqlDb       string
}

type Store struct {
	*sqlx.DB
	Config *Config
}

func NewStore(config *Config) (*Store, error) {
	dsn := config.MysqlUsername + ":" + config.MysqlPassword + "@" + "(" + config.MysqlHost + ":" + strconv.Itoa(config.MysqlPort) + ")/" + config.MysqlDb + "?parseTime=true"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {

		return nil, err
	}

	return &Store{
		DB:     db,
		Config: config,
	}, nil
}

func (s *Store) InitializeDB() error {
	stmt := `
CREATE TABLE IF NOT EXISTS hd_wallets
(
    id    int auto_increment
        primary key,
    x_pub varchar(255) not null
)
CHARSET = utf8mb4;
`
	_, err := s.DB.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `
CREATE TABLE IF NOT EXISTS addresses
(
    id            int auto_increment
        primary key,
    created_at    datetime    not null,
    address       char(42)    not null,
    account_id    int         not null,
    account_index int         null,
    path          varchar(50) not null,
    hd_wallet_id  int         null,
    constraint addresses_hd_wallets_id_fk
        foreign key (hd_wallet_id) references hd_wallets (id)
)
CHARSET = utf8mb4;
`
	_, err = s.DB.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `
CREATE TABLE IF NOT EXISTS txs
(
    id           int auto_increment
        primary key,
    tx_time      datetime                         null,
    reflect_time datetime                         null,
    from_address char(42)                         null,
    to_address   char(42)                         not null,
    asset        varchar(10)                      not null,
    amount       double(15, 8) default 0.00000000 not null,
    block_no     bigint        default 0          not null,
    tx_hash      varchar(255)                     not null,
    is_reflected tinyint(1)    default 0          not null,
    tx_status    varchar(255)                     not null
)
    CHARSET = utf8mb4;
`
	_, err = s.DB.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `
CREATE TABLE IF NOT EXISTS merchants
(
    id       int auto_increment
        primary key,
    email    varchar(255) not null,
    password varchar(255) not null
)
    CHARSET = utf8mb4;
`
	_, err = s.DB.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}
