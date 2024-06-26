Requis:

- CLI tools
- AWS SDK
- BIN 
    - Windows
    - Linux
    - OSX

Specs:

- Concurrence
- throttle based on AWS API max rate
- Informations sur avancement

Logic:

- FinOps --> add stats versus bucket size, or overall size, etc...
- Output Files
- Search BY (name, storage class)
- Group by bucket
- Filter with regex on (s3://mybucket/Folder/SubFolder/log*)
- Filter by encryption type
- information mode to get more info on buckets
- add version logic (get version 2.0 of a blob)


  - Nom --> LS
  - Date de création --> 
  - Nombre de fichiers --> api_op_ListObjects.go
  - Taille totale des fichiers --> api_op_ListObjects.go
  - Date de mise-à-jour de l'objet le plus récent
  - Et le plus important de tous, **combien ça coûte...**

Plus:

- Interfaces to add multi cloud ressources
- Interfaces in AWS to add endpoint aswell

SDK : 
import "github.com/aws/aws-sdk-go-v2"

TUI: BubbleTea : https://github.com/charmbracelet/bubbletea


Output: 
Flags:
- Output group-by (Ex. :region, )
- Output filter-by (Ex.: name, storage-class) DONE
- Output order-by-dec (Ex.: price, name, storage-class, size) DONE
- Output order-by-inc (Ex.: price, name, storage-class, size) DONE
- Verbose ?

- groupby
    - region
- orderbyINC
    - price
    - name
    - storage-class
    - size
- orderbyDEC
    - price
    - name
    - storage-class
    - size
- filterByStorageClass
    - STANDARD
    - ...
Input//
- filterByNames
    Input

RATELIMIT AWS:
5500 GET per second per prefix ()