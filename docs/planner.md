MVP:
- pour chaque bucket: DONE
  - Nom DONE
  - Date de création DONE
  - Nombre de fichiers DONE
  - Taille totale des fichiers DONE
  - Date de mise-à-jour de l'objet le plus récent Done
  - Et le plus important de tous, **combien ça coûte...** Done
- Cli interface DONE
- Concurrence DONE
- threading throttle DONE
- api throttle DONE
- Output data with flags DONE
- Test Done
- Affichage 
  - Possibilité de sortir les résultats en octets, kilooctets, Mégaoctets, etc. DONE
  - Pouvoir grouper les buckets par [régions](https://docs.aws.amazon.com/fr_fr/AWSEC2/latest/UserGuide/using-regions-availability-zones.html) DONE
  
- Filtres
  - Par nom de bucket DONE
  - Par [type de stockage](https://docs.aws.amazon.com/AmazonS3/latest/dev/storage-class-intro.html) (Standard, IA, RR). Tu peux fournir des stats sur les objets dans le bucket (la quantité par type de stockage) et / ou ajouter un filtre sur le type de stockage (les informations sur le bucket réflètent alors seulement les objets qui ont le type sélectionné) DONE

- DOCMENTATION A FAIRE

PLUS:
- TUI Halfway done
- Filtrer les fichiers considérés dans le calcul à l’aide d’un préfixe, un glob et / ou une expression régulière (ex: s3://mybucket/Folder/SubFolder/log*).
- Filtrer ou organiser les résultats selon le [type d'encryption](https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingEncryption.html)
- Obtenir des informations supplémentaires sur les buckets (Life cycle, cross-region replication, etc.)
- Tenir compte des [versions précédentes](https://docs.aws.amazon.com/AmazonS3/latest/UG/enable-bucket-versioning.html) des fichiers (nombre + taille).
- Des statistiques pour afficher le pourcentage de l’espace total occupé par un bucket ou toute autre bonne idée que tu pourrais avoir sont également les bienvenues.
- Metrics based on call per second