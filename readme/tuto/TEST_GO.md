# Les tests en Go : créer, organiser et exécuter

Go embarque tout le nécessaire pour tester : le package `testing` (dans la bibliothèque standard) et la commande `go test`. Aucune dépendance externe n'est requise pour commencer.

---

## 1. Créer un fichier de test

Un fichier de test se place **dans le même dossier** que le code qu'il teste.

```
monprojet/
├── go.mod
├── calc.go          # code de production
└── calc_test.go     # tests associés
```

Règles à retenir :

- le nom du fichier se termine par `_test.go` ;
- il est ignoré par `go build` : il n'est compilé que par `go test` ;
- il peut utiliser le même package que le code testé (`package calc`) ou un package externe (`package calc_test`).

| Déclaration | Usage |
|---|---|
| `package calc` | Test « boîte blanche » : accès aux éléments non exportés (minuscules). |
| `package calc_test` | Test « boîte noire » : seule l'API publique est accessible, comme un vrai utilisateur du package. |

---

## 2. Conventions de nommage

| Élément | Convention | Exemple |
|---|---|---|
| Fichier | `<nom>_test.go` | `calc_test.go` |
| Test | `func TestXxx(t *testing.T)` | `TestAdd` |
| Benchmark | `func BenchmarkXxx(b *testing.B)` | `BenchmarkAdd` |
| Exemple | `func ExampleXxx()` | `ExampleAdd` |

Points importants :

- le nom d'une fonction de test **commence par `Test`**, suivi d'une **majuscule** (`TestAdd` est valide, `Testadd` ne sera pas exécuté) ;
- elle prend **un seul paramètre** : `t *testing.T` ;
- le nom décrit en général la fonction testée : `TestAdd`, `TestDivide`, `TestUser_Validate` (pour une méthode `Validate` du type `User`).

---

## 3. Écrire un test simple

Code à tester (`calc.go`) :

```go
package calc

import "errors"

func Add(a, b int) int {
	return a + b
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division par zéro")
	}
	return a / b, nil
}
```

Fichier de test (`calc_test.go`) :

```go
package calc

import "testing"

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5

	if got != want {
		t.Errorf("Add(2, 3) = %d ; attendu %d", got, want)
	}
}

func TestDivideParZero(t *testing.T) {
	_, err := Divide(10, 0)
	if err == nil {
		t.Fatal("une erreur était attendue, mais err est nil")
	}
}
```

### Les méthodes de `*testing.T` à connaître

| Méthode | Effet |
|---|---|
| `t.Error` / `t.Errorf` | Marque le test comme échoué, **mais continue** l'exécution. |
| `t.Fatal` / `t.Fatalf` | Marque le test comme échoué et **l'arrête immédiatement**. |
| `t.Log` / `t.Logf` | Affiche un message (visible avec `-v` ou en cas d'échec). |
| `t.Skip` | Ignore le test (ex. : dépendance absente). |
| `t.Run` | Lance un sous-test. |
| `t.Helper` | Indique qu'une fonction est un utilitaire de test. |
| `t.Cleanup` | Enregistre une fonction de nettoyage appelée à la fin du test. |

> **Astuce :** utilisez `t.Fatal` quand la suite du test n'a plus de sens (ex. : une erreur inattendue avant d'utiliser un résultat), et `t.Error` quand vous voulez voir tous les échecs d'un coup.

---

## 4. Le style « table-driven » (tests par tableau)

C'est **la** convention idiomatique en Go : on décrit les cas dans une table, puis on les parcourt avec des sous-tests.

```go
func TestAddTable(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"positifs", 2, 3, 5},
		{"négatifs", -2, -3, -5},
		{"avec zéro", 0, 7, 7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Add(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("Add(%d, %d) = %d ; attendu %d", tc.a, tc.b, got, tc.want)
			}
		})
	}
}
```

Avantages : ajouter un cas = ajouter une ligne, chaque cas a un nom, et on peut le lancer seul (voir section 6).

---

## 5. Exécuter les tests

Depuis la racine du projet (là où se trouve `go.mod`) :

```bash
# Tests du package du dossier courant
go test

# Tests de tous les packages du projet (récursif)
go test ./...

# Tests d'un package précis
go test ./internal/calc
```

Sortie typique :

```
ok  	monprojet/calc	0.002s
```

En cas d'échec :

```
--- FAIL: TestAdd (0.00s)
    calc_test.go:9: Add(2, 3) = 6 ; attendu 5
FAIL
FAIL	monprojet/calc	0.002s
```

### Procédure spécifique à SoloPG

Dans ce projet, la configuration est chargée au démarrage des packages via
`CONFIG_FILE`. Le fichier `config/.config.yaml` contient lui-même des chemins
relatifs vers `locales/` et `data/`. Les binaires de test Go peuvent toutefois
être exécutés depuis le dossier de leur package et non depuis la racine du
projet. Il faut donc conserver la racine comme répertoire courant pendant
l'exécution de chaque binaire.

Depuis la racine du dépôt, la commande recommandée est :

```bash
PROJECT_ROOT="$(pwd)"
CONFIG_FILE="$PROJECT_ROOT/config/.config.yaml" \
GOCACHE="${GOCACHE:-/tmp/solopg-go-cache}" \
go test -exec "sh -c 'cd \"$PROJECT_ROOT\" && exec \"\$@\"' sh" ./...
```

Cette commande :

- exécute les tests de tous les packages avec `./...` ;
- fournit explicitement le fichier de configuration du projet ;
- garantit que les chemins `locales/` et `data/` sont résolus depuis la racine ;
- utilise un cache Go local si le cache système n'est pas accessible.

Pour un package précis, conservez les mêmes variables et remplacez `./...` :

```bash
PROJECT_ROOT="$(pwd)"
CONFIG_FILE="$PROJECT_ROOT/config/.config.yaml" \
GOCACHE="${GOCACHE:-/tmp/solopg-go-cache}" \
go test -exec "sh -c 'cd \"$PROJECT_ROOT\" && exec \"\$@\"' sh" \
./app/shared/components/models
```

Pour éviter toute ambiguïté, la version compacte équivalente est :

```bash
PROJECT_ROOT="$(pwd)"; CONFIG_FILE="$PROJECT_ROOT/config/.config.yaml" GOCACHE=/tmp/solopg-go-cache go test -exec "sh -c 'cd \"$PROJECT_ROOT\" && exec \"\$@\"' sh" ./...
```

En fonctionnement normal, après positionnement correct de l'environnement et
depuis la racine du dépôt, la commande générique suivante reste la forme la
plus simple :

```bash
go test ./...
```

---

## 6. Exécuter un test spécifique

L'option `-run` accepte une **expression régulière** appliquée au nom du test.

```bash
# Un test précis (correspondance partielle !)
go test -run TestAdd

# Correspondance exacte : ancres ^ et $
go test -run '^TestAdd$'

# Plusieurs tests
go test -run 'TestAdd|TestDivide'

# Un sous-test : séparer les niveaux par /
go test -run 'TestAddTable/négatifs'

# Dans tous les packages
go test ./... -run '^TestAdd$'
```

> **Attention :** `-run TestAdd` lance aussi `TestAddTable`, car c'est une regex. Utilisez `^...$` pour une correspondance exacte.

---

## 7. Options utiles de `go test`

| Option | Rôle |
|---|---|
| `-v` | Mode verbeux : affiche chaque test et ses `t.Log`. |
| `-run <regex>` | Filtre les tests par nom. |
| `-count=1` | Force l'exécution sans utiliser le **cache** des résultats. |
| `-race` | Active le détecteur de **data races** (conditions de concurrence). |
| `-cover` | Affiche le pourcentage de couverture. |
| `-coverprofile=c.out` | Écrit un profil de couverture dans un fichier. |
| `-short` | Active le mode court (voir `testing.Short()`). |
| `-timeout 30s` | Échoue si les tests dépassent la durée donnée (10 min par défaut). |
| `-failfast` | S'arrête au premier test en échec. |
| `-parallel N` | Nombre max de tests parallèles. |
| `-bench .` | Lance les benchmarks. |
| `-shuffle=on` | Exécute les tests dans un ordre aléatoire (détecte les dépendances cachées). |

Exemples courants :

```bash
# Le « classique » : tout, verbeux, sans cache
go test ./... -v -count=1

# Avec détection de data races
go test ./... -race

# Couverture dans le terminal
go test ./... -cover

# Rapport de couverture HTML
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Voir la couverture par fonction
go tool cover -func=coverage.out

# Benchmarks (sans lancer les tests)
go test -bench . -run '^$'
```

> **Le cache :** Go met en cache les résultats de `go test` si ni le code ni les fichiers lus n'ont changé (la sortie affiche alors `(cached)`). Utilisez `-count=1` pour forcer une nouvelle exécution.

---

## 8. Bonnes pratiques essentielles

### 8.1 Écrire des messages d'échec explicites

Indiquez toujours **l'entrée**, le résultat **obtenu** et le résultat **attendu** :

```go
t.Errorf("Add(%d, %d) = %d ; attendu %d", a, b, got, want)
```

### 8.2 Privilégier les tests par tableau

Dès que vous testez plusieurs cas d'une même fonction, utilisez le style table-driven avec `t.Run` et des noms de cas parlants.

### 8.3 Marquer les fonctions utilitaires avec `t.Helper()`

L'erreur est alors rapportée à la ligne de l'appelant, pas à celle de l'utilitaire :

```go
func assertEqual(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("obtenu %d ; attendu %d", got, want)
	}
}
```

### 8.4 Nettoyer proprement

Utilisez `t.Cleanup` pour libérer des ressources, et `t.TempDir()` pour obtenir un dossier temporaire supprimé automatiquement :

```go
func TestEcritureFichier(t *testing.T) {
	dir := t.TempDir() // supprimé à la fin du test
	t.Cleanup(func() { /* fermer une connexion, etc. */ })

	// ... utiliser dir
}
```

### 8.5 Garder les tests indépendants et déterministes

- un test ne doit pas dépendre de l'ordre d'exécution ni du résultat d'un autre test ;
- évitez de dépendre de l'heure, du réseau ou de valeurs aléatoires non maîtrisées (injectez-les plutôt).

### 8.6 Tester le comportement, pas l'implémentation

Testez ce que la fonction **promet** (entrées → sorties, erreurs), pas ses détails internes. Vos tests resteront valides lors d'un refactoring.

### 8.7 Tester aussi les cas d'erreur et les cas limites

Valeurs nulles ou vides, bornes, entrées invalides : c'est là que se cachent la plupart des bugs.

```go
if err == nil {
	t.Fatal("une erreur était attendue")
}
```

### 8.8 Séparer les tests longs

Les tests lents (base de données, réseau) peuvent être ignorés en mode court :

```go
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("test d'intégration ignoré en mode -short")
	}
	// ...
}
```

```bash
go test -short ./...
```

### 8.9 Utiliser `testdata/` pour les fichiers de test

Un dossier nommé `testdata` est **ignoré par l'outil Go** (pas compilé comme package). Placez-y vos fichiers d'exemple (JSON, CSV, etc.) :

```
calc/
├── calc.go
├── calc_test.go
└── testdata/
    └── cas1.json
```

### 8.10 Lancer régulièrement `-race` et `go vet`

```bash
go vet ./...
go test ./... -race
```

`go test` exécute déjà une partie de `go vet` automatiquement, mais un passage explicite reste utile, surtout en intégration continue.

---

## 9. Aller plus loin

### Tests parallèles

```go
func TestParallele(t *testing.T) {
	t.Parallel() // ce test peut s'exécuter en même temps que d'autres tests parallèles
	// ...
}
```

Dans un test par tableau avec `t.Parallel()` dans le sous-test, la variable de boucle est correctement capturée à partir de Go 1.22. Avec une version antérieure, ajoutez `tc := tc` au début de la boucle.

### Préparation globale avec `TestMain`

À définir une seule fois par package, pour une initialisation ou un nettoyage commun :

```go
func TestMain(m *testing.M) {
	// préparation (ex. : lancer une base de test)
	code := m.Run()
	// nettoyage
	os.Exit(code)
}
```

### Benchmarks

```go
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(2, 3)
	}
}
```

```bash
go test -bench . -benchmem -run '^$'
```

---

## 10. Aide-mémoire

```bash
go test                              # package courant
go test ./...                        # tout le projet
# SoloPG : voir la procédure spécifique ci-dessus si la configuration échoue
go test -v                           # verbeux
go test -run '^TestAdd$'             # un test précis
go test -run 'TestAddTable/négatifs' # un sous-test
go test -count=1                     # sans cache
go test -race                        # détecteur de races
go test -cover                       # couverture
go test -short                       # mode court
go test -failfast                    # stop au premier échec
go test -bench . -run '^$'           # benchmarks uniquement
```

**Checklist rapide pour un nouveau test :**

1. Créer `xxx_test.go` à côté de `xxx.go`.
2. Déclarer `func TestXxx(t *testing.T)`.
3. Préparer les entrées, appeler la fonction, comparer `got` et `want`.
4. Signaler l'échec avec `t.Errorf` ou `t.Fatalf` et un message clair.
5. Depuis la racine de SoloPG, lancer la commande spécifique ci-dessus.
6. Si la suite passe, lancer ensuite `go test ./... -race` avec le même contexte.
