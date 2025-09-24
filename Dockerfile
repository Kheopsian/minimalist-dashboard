# --- PHASE 1 : BUILD ---
# On part d'une image officielle Go pour la compilation
FROM golang:1.22-alpine AS builder

# On définit le dossier de travail à l'intérieur de l'image
WORKDIR /app

# On copie les fichiers de gestion des dépendances (go.mod et go.sum)
COPY go.mod go.sum ./
# On télécharge les dépendances
RUN go mod download

# On copie tout le reste du code source
COPY . .

# On compile l'application.
# CGO_ENABLED=0 crée un binaire statique, sans dépendance externe.
# -o /dashboard-api spécifie le nom du fichier de sortie.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /dashboard-api

# --- PHASE 2 : FINALE ---
# On part d'une image Alpine Linux, qui est très légère et sécurisée.
FROM alpine:latest

# On met à jour les paquets et on installe le certificat SSL (important pour les requêtes HTTPS)
RUN apk --no-cache add ca-certificates

# On définit le dossier de travail
WORKDIR /app

# On copie le dossier frontend depuis la phase de build
COPY --from=builder /app/frontend ./frontend

# On copie UNIQUEMENT le binaire compilé depuis la phase de build
COPY --from=builder /dashboard-api .

# On expose le port 8080 sur lequel notre serveur Go écoute
EXPOSE 8080

# La commande qui sera lancée au démarrage du conteneur
CMD [ "./dashboard-api" ]