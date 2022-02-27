# koushoku

Source code of site [redacted] for those who are willing to run their own instance.

### How it serve and index the archives

Archives and its files are served directly, without writing the files inside the archives into the disk (except for thumbnails). Archives inside the specified data directory will be indexed as long as it follows the following naming formats:

- [Artist] Title (Magazine) [Foo] [Bar] [Crap] {tags kebab-case optional}
- [Circle (Artist)] Title (Magazine) [Foo] [Bar] [Crap] {tags kebab-case optional}

Archives will be indexed concurrently, and usually takes several minutes (~1m10s for around ~8k archives). You can decrease the maximum concurrent numbers if your server is overloaded.

## Prerequisites

- Git
- Go 1.17+
- libvips 8.3+ (8.8+ recommended)
- C compatible compiler such as gcc 4.6+ or clang 3.0+
- PostgreSQL

## Setup

### Install the prerequisites

```sh
# Arch-based distributions
sudo pacman -Syu
sudo pacman -S git go libvips postgresql

# Debian-based distributions
sudo apt-get install -y software-properties-common
sudo add-apt-repository -y ppa:strukturag/libde265
sudo add-apt-repository -y ppa:strukturag/libheif
sudo add-apt-repository -y ppa:tonimelisma/ppa
sudo add-apt-repository -y ppa:longsleep/golang-backports

sudo apt-get update -y
sudo apt-get install -y build-essential git golang-go libvips-dev postgresql
```

### Initialize database cluster

**Only for Arch-based distributions** - Before PostgreSQL can function correctly, the database cluster must be initialized - [wiki.archlinux.org](https://wiki.archlinux.org/title/PostgreSQL#Installation).

```sh
echo initdb -D /var/lib/postgres/data | sudo su - postgres
```

### Enable and start PostgreSQL

```sh
sudo systemctl --now enable postgresql
```

### Create a new database and user/role

```sh
sudo -u postgres psql --command "CREATE USER koushoku LOGIN SUPERUSER PASSWORD 'koushoku';"
sudo -u postgres psql --command "CREATE DATABASE koushoku OWNER koushoku;"
```

### Build the back-end

```sh
git clone https://github.com/rs1703/koushoku
cd koushoku
make build
```

## License

**koushoku** is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).
