Neli webservices
================

Webservices pour le prototype N7.

# Installation

## Base de données

Créez une base de données mysql.

### Goose 

Goose est utilisé pour gérer les migrations de la base de données.

#### Si vous avez installé Go, vous pouvez installer Goose avec la commande suivante: 

```bash
go get -u github.com/pressly/goose/cmd/goose
``` 

Vous pouvez ensuite faire une migration : 

```bash
goose --dir database/migrations mysql DATABASE_CONNECTION up
``` 

N'oubliez pas de remplacer **DATABASE_CONNECTION** (par exemple : "username:password@/database")

Pour annuler une migration: 

```bash
goose --dir database/migrations mysql DATABASE_CONNECTION down
``` 

#### Si vous n'avez pas Goose alors patientez


### Seeds

Un utilitaire est disponible pour générer des données de tests. Il doit être lancé depuis le repertoire **src** à l'aide de la commande: 

```bash
go run database/migrations/seeds/seeds.go
``` 

## Tests

### Pré-requis

Mailhog (https://github.com/mailhog/MailHog) est requis afin de récupérer le token de réinitialisation du mot de passe dans la partie authentification.

### Stucture
Les tests sont générés à l'aide de postman. Le dossier **tests** contient les environments ainsi que les collections. 

Les scenarios de tests doivent être lancés dans l'ordre afin de rester les affectations de variables (notamment pour la partie authentification). 

### Description des variables

* login      : nom d'utilisateur utilisé pour se connecter
* password      : mot de passe utilisé pour se connecter
* token_life    : durée de vie d'un access token
* ttl           : permet d'allouer un jeton de raffraichissement plus court 
* access_token  : token d'accès à l'application sans authentification
* refresh_token : token de rafraichissement
* reset_token   : token de réinitialisation du mot de passe
* old_token     : token expiré
* mailhog       : url d'accès à l'API Mailhog
* neli_base_url : url utilisé dans les emails pour définir l'adresse de réinitialisation du token

### Newman
Lancer en ligne de commande

# Lancement

Le projet doit être lancé depuis le repertoire **src** à l'aide de la commande: 

```bash
go run main.go
``` 

# Authentification

Le token doit être passé dans le header "Authorization" et débuté par "Bearer "