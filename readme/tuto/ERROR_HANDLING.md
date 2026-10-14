	# Système de gestion et de déduplication des erreurs

Ce document décrit le système à utiliser pour gérer les erreurs de la TUI et éviter qu'une même erreur soit enregistrée plusieurs fois en base.

## Principe général

Une erreur doit être traitée en trois étapes :

1. détecter qu'une erreur existe ;
2. l'identifier avec un code stable ;
3. décider si cette occurrence doit être enregistrée ou ignorée.

La présence d'une erreur se teste toujours avec `err != nil`. Le code sert uniquement à reconnaître son type.

```go
if err != nil {
	// Une erreur existe.
}
```

## Erreur applicative codée

Les messages sont destinés à l'affichage et peuvent changer ou être traduits. Le code doit rester stable.

```go
type AppError struct {
	Code string
	Err  error
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}
```

Exemple :

```go
err := &AppError{
	Code: "filter.invalid_value",
	Err:  fmt.Errorf("filtre invalide : %s", filter),
}
```

Récupération du code :

```go
var appErr *AppError

if errors.As(err, &appErr) {
	switch appErr.Code {
	case "filter.invalid_value":
		// Erreur de filtre.
	case "stream.closed":
		// Stream fermé.
	}
}
```

## Déduplication des erreurs

Un simple booléen par erreur devient difficile à maintenir dès qu'il existe plusieurs catégories. Il vaut mieux conserver un ensemble de clés déjà enregistrées.

```go
type errorKey struct {
	Source string
	Code   string
	Reason string
}

type skip struct {
	LoggedErrors map[errorKey]struct{}
}
```

Initialisation et vérification :

```go
if skip.LoggedErrors == nil {
	skip.LoggedErrors = make(map[errorKey]struct{})
}

key := errorKey{
	Source: "mill.filter",
	Code:   appErr.Code,
	Reason: filter,
}

if _, exists := skip.LoggedErrors[key]; exists {
	return // Cette cause a déjà été enregistrée.
}

logs.Error("error.unexpected:action", map[string]any{
	"Action": "unexpected:action.filter:logs",
	"Error":  err.Error(),
})

skip.LoggedErrors[key] = struct{}{}
```

## Choix de la clé

La clé doit correspondre au niveau de déduplication souhaité :

- `Source + Code` : toutes les erreurs de même type sont regroupées ;
- `Source + Code + Reason` : deux causes différentes sont enregistrées séparément ;
- le texte complet de l'erreur ne doit pas servir de code, car il peut être localisé ou contenir des détails variables.

Exemples de codes :

```text
filter.invalid_value
filter.missing_flag
tail.invalid_value
stream.closed
stream.connection_failed
```

## Dans `View()`

`View()` peut être appelée plusieurs fois. Elle ne doit donc pas enregistrer directement une erreur sans déduplication.

```go
mill.Filter()

if mill.HasErr() {
	m.handleMillError("filter", mill.GetErr())
}

if m.err != nil {
	m.handleStreamError(m.err)
}
```

Les sources doivent rester distinctes :

```text
mill.filter
mill.tail
stream
```

Ainsi, une erreur de filtre ne masque pas une erreur de stream.

## À éviter

- appeler `logs.Error(...)` directement dans une méthode appelée à chaque rendu ;
- utiliser un booléen différent pour chaque erreur ;
- utiliser le texte traduit comme identifiant ;
- remettre `m.err` à `nil` uniquement pour empêcher un doublon ;
- ajouter un délai arbitraire si l'objectif est de dédupliquer une cause précise.

Un délai de cinq minutes est du *rate limiting*. Il convient seulement si l'objectif est de limiter la fréquence des événements, et non de dédupliquer les mêmes causes.

## Ajouter un nouveau type d'erreur

1. définir un code stable ;
2. créer l'`AppError` avec ce code ;
3. choisir une source (`mill.tail`, `repository`, etc.) ;
4. décider si la raison doit faire partie de la clé ;
5. passer par le handler dédupliqué avant d'appeler `logs.Error`.
