# koushoku

Source code of site [redacted] for those who are willing to run their own instance. You need at least 1 GB of RAM and 1 TB of storage space.

**server** - Use this to deploy the server, it runs on port 42073 by default.

**util** - Utility tool, use this to download, indexand delete archives, etc.

Put batch releases into a file named **batches.txt** and singles into a file named **singles.txt**. Magnet links are separated by a new line.

For more information, use `util --help`. You can combine multiple arguments, for instance, when you execute it using `--purge --index` arguments, it will purge symlinks then re-index the archives.

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
