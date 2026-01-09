Voici le PRD final, nettoyé et structuré pour ton équipe de développement.

---

# PRD: grepai - Infrastructure de Recherche Sémantique Temps Réel

| Propriété | Détails |
| --- | --- |
| **Produit** | `grepai` |
| **Type** | CLI Tool (Installation Globale) |
| **Cible** | Développeurs & Agents IA (Cursor, Claude Code) |
| **Périmètre** | Indexation & Retrieval Vectoriel (Pas de génération/Chat) |
| **Stack** | Go (1.22+), No CGO |

---

## 1. Vision du Produit

**grepai** est un outil d'infrastructure "Privacy-first" conçu pour moderniser la recherche de code. Contrairement à `grep` (recherche textuelle exacte), `grepai` indexe le sens du code (embeddings) pour permettre une recherche par intention.

Il est conçu pour tourner en arrière-plan (daemon), maintenant une "carte mentale" du projet toujours à jour en temps réel. Sa fonction première est de servir de contexte fiable aux développeurs et aux Agents IA qui travaillent sur la codebase.

---

## 2. Contraintes Techniques & Architecture

### 2.1 Stack Imposée

* **Langage :** Go (Golang) uniquement.
* **Compilation :** Binaires statiques cross-platform (MacOS, Linux, Windows).
* **Interdiction :** Usage de **CGO** strictement interdit (ce qui exclut SQLite). L'outil doit être portable sans dépendance système.
* **Modèles IA :** Usage exclusif des APIs d'**Embedding**. Aucune utilisation d'APIs de Chat/Completion.

### 2.2 Architecture Hexagonale

L'application doit suivre une architecture modulaire pour découpler la logique métier des services tiers.

* **Domaine :** Logique de découpage (Chunking), détection de fichiers, gestion du Watcher.
* **Interfaces (Ports) :**
* `VectorStore` : Abstraction pour le stockage et la recherche des vecteurs.
* `Embedder` : Abstraction pour la transformation Texte vers Vecteur.


* **Adaptateurs :**
* Stockage : Implémentation Fichier (GOB) et Implémentation Postgres (pgvector).
* IA : Implémentation Ollama (Local) et Implémentation OpenAI (SaaS).



---

## 3. Stratégie de Données & Stockage

### 3.1 Contexte Local

Bien que l'outil soit installé globalement dans le `$PATH`, il opère contextuellement.

* Les données d'indexation sont stockées dans un dossier caché `.grepai/` situé à la racine du projet surveillé.
* Ce dossier contient la configuration (`config.yaml`) et l'index par défaut.

### 3.2 Options de Stockage (Backend)

**Option A : Standalone (Par défaut - Priorité 1)**

* **Mécanisme :** Tout l'index réside en RAM pour une performance maximale.
* **Persistance :** Sérialisation native Go (format GOB) sur le disque dans `.grepai/index.gob` à l'arrêt ou périodiquement.
* **Structure :** Optimisée pour le CRUD rapide (mise à jour par fichier).

**Option B : Centralisé (Priorité 2)**

* **Mécanisme :** Connexion à une base PostgreSQL externe avec l'extension `pgvector`.
* **Cas d'usage :** Monorepos massifs ou index partagé entre développeurs.
* **Isolation :** Les vecteurs doivent être tagués avec un identifiant unique du projet pour permettre le multi-tenant dans la même table.

---

## 4. Fonctionnalités Principales

### 4.1 Commande `init`

Initialise le projet courant.

* Génère le fichier de configuration par défaut.
* Demande à l'utilisateur de choisir son Provider (Ollama vs OpenAI) et son Store (Local vs Postgres).
* Ajoute automatiquement `.grepai/` au `.gitignore` du projet s'il existe.

### 4.2 Commande `watch` (Cœur du Système)

Processus bloquant qui maintient la synchronisation.

* **Scan Initial :** Compare l'état des fichiers sur le disque avec l'index existant. Supprime les entrées obsolètes et indexe les nouveaux fichiers.
* **Monitoring Temps Réel :** Écoute les événements du système de fichiers (Création, Modification, Suppression, Renommage).
* **Debouncing :** Applique une temporisation obligatoire (ex: 500ms) pour regrouper les rafales d'événements (ex: `git pull`, `Save All`) et éviter de surcharger le provider d'embedding.
* **Gestion Atomique :** Lors de la modification d'un fichier, l'outil doit garantir que les anciens vecteurs sont supprimés avant ou pendant l'insertion des nouveaux pour éviter les doublons.

### 4.3 Commande `search`

Moteur de recherche pur.

* **Entrée :** Une requête en langage naturel.
* **Traitement :** Vectorisation de la requête -> Calcul de similarité (Cosinus) -> Tri.
* **Sortie :** Liste des N meilleurs résultats affichant : Chemin du fichier, Numéro de ligne, Score de pertinence, Extrait de code.

### 4.4 Commande `agent-setup`

Configure l'environnement pour les Agents IA.

* Détecte la présence de fichiers de configuration d'agents (`.cursorrules`, `CLAUDE.md`).
* Injecte (en append) une instruction système expliquant à l'agent comment utiliser `grepai` pour chercher du contexte au lieu de deviner.
* Garantit l'idempotence (n'ajoute pas le texte s'il est déjà présent).

---

## 5. Règles de Gestion & Performance

### 5.1 Ignorance et Filtrage

* L'outil doit impérativement respecter les règles du `.gitignore`.
* Il doit ignorer par défaut les dossiers systèmes et de dépendances (`.git`, `node_modules`, `vendor`, `bin`).
* Il doit détecter et ignorer les fichiers binaires.

### 5.2 Chunking (Découpage)

* Le code doit être découpé en segments (chunks) d'une taille configurable (défaut : 512 tokens) avec un chevauchement (overlap) pour ne pas couper le contexte (défaut : 50 tokens).

### 5.3 Gestion des Erreurs

* Si le provider d'Embedding est inaccessible (ex: Ollama éteint), le `watch` ne doit pas crasher mais logger une erreur et réessayer ou attendre.
* Les fichiers trop volumineux (> 1 Mo texte) doivent être ignorés avec un avertissement (warning) pour éviter les problèmes de mémoire ou de timeout API.

---

## 6. Infrastructure Docker (Livrable annexe)

Un fichier `compose.yaml` doit être fourni pour les utilisateurs souhaitant héberger l'infrastructure backend. Il doit contenir :

1. **PostgreSQL** configuré avec l'extension `vector` pré-installée.
2. **Volumes** pour la persistance des données.
3. (Optionnel) **Ollama** pour les utilisateurs Linux/Windows Server.

---

## 7. Roadmap de Développement

1. **Skeleton :** Structure du CLI, gestion de la config et injection des dépendances.
2. **Core Domain :** Implémentation du parsing de fichiers, du respect du `.gitignore` et du chunking.
3. **Adapters :** Implémentation du Store GOB (Memory) et de l'Embedder Ollama.
4. **Feature `search` :** Implémentation de la recherche vectorielle basique.
5. **Feature `watch` :** Implémentation du système `fsnotify`, du debouncing et de la mise à jour CRUD.
6. **Feature `agent-setup` :** Implémentation de la modification des fichiers de règles Agents.
7. **Extension :** Ajout du support Postgres et OpenAI.