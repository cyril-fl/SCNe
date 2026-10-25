# SCNe

**SCNe** est une CLI Go conçue comme un moteur pour piloter une session de **solo RPG** directement dans le terminal.

Grâce à son système de fichiers et de templates, il est possible d'étendre les règles et l'univers du jeu. Par exemple, il est possible de créer un ensemble de règles *fantasy* ou *sci-fi* pour adapter l'expérience au genre souhaité.

> ⚠️ Il n'est actuellement pas possible d'utiliser plusieurs ensembles de règles simultanément. Cette évolution n'est pas prévue pour le moment.

- **Dépôt :** https://github.com/cyril-fl/SCNe

## Documentation
Une documentation plus détaillée sur le fonctionnement du projet ainsi que sur certains concepts est disponible sur le [wiki](https://github.com/cyril-fl/SCNe/wiki).

## Fonctionnalités

La commande `play` ouvre une interface terminal (TUI) qui couvre toute la session de jeu :

- **Onboarding** — choix du nom, de l'archétype, forge du personnage et choix du lieu de départ.
- **Campagnes** — sauvegarde et chargement de parties, avec archives et journal de campagne.
- **Plateau de jeu** — zone de narration et panneau latéral en accordéon (actions, codex, informations sur le personnage et le lieu).
- **Codex** — bestiaire, lieux, PNJ, quêtes et objets découverts au fil de la partie.
- **Outils de jeu** — lancers de dés, oracles, indices et voyage entre les lieux.
- **Combat** — combat au tour par tour, au style narratif, dont les règles sont définies dans [`data/systems/rules/combat.yaml`](./data/systems/rules/combat.yaml).
- **Data-driven** — règles et univers décrits en YAML :
  - [`data/systems`](./data/systems) : attributs (classes, races, raretés, emplacements, stats…) et règles (dés, oracles, combat).
  - [`data/template`](./data/template) : personnages (PNJ, monstres), lieux, objets, panoplies, indices.
- **Interface personnalisable** — couleurs configurables et traduction (français, anglais).

## Installation et utilisation

SCNe est développé avec la dernière version de Go et s'exécute avec les services définis dans le fichier [`compose.yaml`](./compose.yaml).

```bash
git clone https://github.com/cyril-fl/SCNe.git
cd SCNe
make setup   # installe les hooks git (lefthook)
make up      # démarre Colima et les services Docker
make build
./scne play
```

### Prérequis

- [Go](https://go.dev/) — pour exécuter le projet directement sur la machine
- [Plum](https://github.com/cyril-fl/plum) — bibliothèque utilitaire (CLI, configuration, YAML…) utilisée par SCNe
- [Docker](https://www.docker.com/) et [Docker Compose](https://docs.docker.com/compose/) — pour exécuter l'environnement dans des conteneurs
- [Colima](https://github.com/abiosoft/colima) — optionnel, utilisé par `make up` / `make down` sur macOS
- MongoDB — utilisé par l'application et fourni par le `compose.yaml`
- [Compooss](#remerciements) — interface web permettant de consulter et manipuler MongoDB

> ℹ️ Le `go.mod` peut contenir une directive `replace` pointant vers une copie locale de Plum. Adapter ou supprimer cette ligne selon votre environnement.

### Configuration

#### Variables d'environnement

L'application nécessite un fichier d'environnement (voir [`.env.example`](./.env.example)) :

```env
# MongoDB
MONGO_HOST=localhost
MONGO_PORT=27017

# Config
CONFIG_FILE=./config/.config.yaml
```

Le nom du fichier dépend du mode d'exécution :

- `.env` — utilisé lorsque l'application est exécutée directement sur la machine avec Go.
- `.env.docker` — utilisé lorsque l'application est exécutée dans le conteneur Docker.

#### Fichier de configuration

Le fichier [`config/.config.yaml`](./config/.config.yaml) regroupe :

- `theme` — couleurs de l'interface (`accent`, `text`, `muted`, `subtle`, `error`, `success`), en hexadécimal ou en code ANSI. Ce bloc surcharge le thème par défaut défini dans `app/shared/services/theme` ; une couleur vide garde la valeur par défaut.

  ```yaml
  theme:
    accent: "#C9A26B"
    text: "#E3DED6"
    muted: "#8C877F"
    subtle: "#4B4843"
    error: "#D0614C"
    success: "#8FB573"
  ```

- `i18n` — langue par défaut (`fr`) et langues disponibles (`fr`, `en`), dont les traductions se trouvent dans [`locales/`](./locales).
- `commands` — métadonnées des commandes de la CLI (description, exemples, arguments).
- `documents` — chemins des fichiers et dossiers de règles et de templates chargés par le jeu.

### Environnement Docker

Le [`compose.yaml`](./compose.yaml) définit trois services :

- `app` — conteneur Go (`scne-app`) qui charge le fichier `.env.docker`.
- `mongo` — MongoDB (`scne-mongo`) lancé en replica set `rs0`, initialisé automatiquement par le healthcheck.
- `compooss` — interface web MongoDB, accessible sur http://localhost:6969 une fois `mongo` opérationnel.

Une fois les services démarrés, ouvrir un shell dans le conteneur avec :

```bash
docker compose exec -it app bash
```

## Commandes

SCNe expose plusieurs commandes via sa CLI. Leurs métadonnées sont définies dans la section `commands` de [`config/.config.yaml`](./config/.config.yaml).

### `scne`

Commande racine.

```yaml
root:
  short: "SCNe CLI"
  long: "SCNe is a CLI to manage a solo RPG session in the terminal."
```

### `config`

Affiche la configuration située dans [`config/.config.yaml`](./config/.config.yaml), au format JSON.

```yaml
config:
  use: "config"
  short: "Display the configuration of SCNe"
```
> ⚠️ Une évolution est prévue afin de permettre la personnalisation de la configuration directement depuis le terminal.

### `logs`

Affiche les logs générés par l'application.

```yaml
logs:
  use: "logs"
  short: "Display application logs"
  long: "Display the application logs generated by SCNe."
  example: "scne logs --tail 100 --filter error"
  args:
    - name: "tail"
      shorthand: "t"
      value: ""
      usage: "Show only the last N log entries"
      type: "int"
    - name: "filter"
      shorthand: "f"
      value: ""
      usage: "Show only logs at the specified kind of level (e.g., error, warning, info)"
      type: "string"
```

### `play`

Lance une session de solo RPG (voir [Fonctionnalités](#fonctionnalités)). L'alias `run` reste disponible.

```yaml
play:
  use: "play"
  short: "Play a solo RPG session"
```
> ⚠️ La documentation détaillée de cette commande sera ajoutée dans un fichier `.md` dédié.

### `serve`

Se connecte directement à la base de données afin d'afficher les enregistrements de logs en direct.

```yaml
serve:
  use: "serve"
  short: "Serve a solo RPG session"
  example: "scne serve --filter error"
  args:
    - name: "filter"
      shorthand: "f"
      value: ""
      usage: "Show only logs at the specified kind of level (e.g., error, warning, info)"
      type: "string"
```

### `try`

Commande destinée exclusivement à la **DX (Developer Experience)**.

Elle permet de tester aléatoirement de nouvelles fonctionnalités ou de réaliser différentes opérations utiles au développement. Cette commande est uniquement disponible dans l'environnement de développement.

```yaml
try:
  use: "try"
  short: "Try a solo RPG session"
```

### `version`

Affiche les métadonnées de version actuelles de la CLI.

```yaml
version:
  use: "version"
  short: "Display the version of SCNe"
```

Le résultat dépend de l'environnement d'exécution :

```text
Version      v1.0.1
Commit       46a59a369090219b64003c4edbb689a3ec76ff5f
Release date 28/09/2026
```
- En développement, les métadonnées sont récupérées à partir de Git.
- En production, les métadonnées sont injectées lors du build grâce aux `ldflags` (voir [Versioning](#versioning)).

## Makefile

| Cible        | Description                                                      |
|--------------|------------------------------------------------------------------|
| `make setup` | Installe les hooks git via lefthook (`gofmt` en pre-commit)      |
| `make build` | Compile le binaire `scne` avec les métadonnées de version        |
| `make test`  | Lance les tests unitaires (`ARGS=./app/...` pour cibler un dossier) |
| `make up`    | Démarre Colima puis les services Docker                          |
| `make down`  | Arrête les services Docker puis Colima                           |

## Versioning

Les métadonnées de build sont centralisées dans `app/cmdversion/services/process/generatebuildinfo`.

Sans `ldflags`, la commande `version` utilise des valeurs par défaut et tente un fallback via `runtime/debug`.

Avec `ldflags`, les valeurs de version, de commit et de date sont injectées lors du build. La variable `ENV` (`dev` par défaut) détermine l'environnement : toute autre valeur que `dev` / `development` / `develop` désactive les fonctionnalités de développement (comme `try`).

### Commandes recommandées

```bash
make build            # ENV=dev par défaut
make build ENV=prod
./scne version
./scne --version
```

### Build manuel

```bash
go build -o scne -ldflags "\
  -X scne/app/cmdversion/services/process/generatebuildinfo.AppName=scne \
  -X scne/app/cmdversion/services/process/generatebuildinfo.Version=v0.0.3 \
  -X scne/app/cmdversion/services/process/generatebuildinfo.Commit=$(git rev-parse --short HEAD) \
  -X scne/app/cmdversion/services/process/generatebuildinfo.Date=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X scne/config.env=prod" .
```

## Remerciements

### Compooss
SCNe utilise [Compooss](https://github.com/abdullahmia/compooss), une interface web permettant de consulter et manipuler les données MongoDB.

- **Dépôt :** https://github.com/abdullahmia/compooss
- **Site web :** https://compooss.vercel.app/

### Bibliothèques
- [Cobra](https://github.com/spf13/cobra) — CLI
- [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles) et [Lip Gloss](https://github.com/charmbracelet/lipgloss) — interface terminal
- [go-i18n](https://github.com/nicksnyder/go-i18n) — traductions
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) — base de données
