# 1.2.1

## :beetle: Bug fixes

- fix message deletion by not storing original message as a reply 

# 1.2.0

## :warning: Breaking

- **Changes in elasticsearch index mappings** : tools, labels, answers and team do not use the name of the field as index.

## :rocket: Features

- **answers** :
    - allow answers without tools or labels
    - add option to disable feedback for answers
    - add checks when creating / editing an answer to validate that the tool and label exists
- **tools / labels** : add check when deleting a resource to avoid deleting when used
- **olivere/elastic** : usage of this new library everywhere

## :beetle: Bug fixes

- **answers** : fix some issues with this feature

# 1.1.6

## :beetle: Bug fixes

- Increase slack id length

# 1.1.5

## :beetle: Bug fixes

- **answers** : Check for empty tools or labels in message
