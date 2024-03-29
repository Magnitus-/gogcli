# Achitecture Documentation

This documentation is an architectural overview of how the tool operates.

# Goal of the Tool

The tool is meant to maintain a store of games purchased on GOG.com which is essentially a mirror of your game collection on the website.

There is some filtering support so that only the part of your collection that you want (ex: you can filter out MacOS or Italian installers) is mirrored.

# Concepts the Tool Uses

Conceptually, the tool uses the following components:

- **GOG.com**: The GOG website, where you purchase and download games
- **Storages**: A persistent store that will mirror either your entire GOG collection or a subset of it. It also stores some metadata files (the manifest and actions files) for efficiency's sake.
- **Manifest**: This is a file that describes either an existing or desired state for a collection of game files in your storage. It lists games and game files.
- **Actions**: This is a file that represents every atomic action the client will have to do in order to bridge the difference between what is described in a manifest file and the actual game files stored in your storage. After every action in the actions file is executed, the manifest will perfectly represent the state of the game files in your storage.
- **Updates**: This is a file that lists all the games that GOG.com marks as new or updated at a given point in time.

# Command Flow

The diagram below shows how the various commands affect the above components:

![Command Flow](command_flow.png)

The **gogcli manifest generate** will generate a fresh **manifest** from **GOG.com**. **gogcli manifest update** will update specific games in the **manifest** (either specified by game **id** or via an **updates** file generated based on what **GOG.com** flaged as new or updated in your collection) based on **GOG.com**.

The **gogcli storage apply manifest** command will move a **manifest** file to your **storage**, generate an **actions** file in your storage based on the different between the **manifest** you just applied and the previous manifest. In event that there is already a pre-existing **actions** file in the storage, **gogcli storage apply manifest** will reconcile it with the changes from the updated manifest.

The **gogcli storage execute-actions** command will run the actions in the **actions** file, downloading new game files in the storage, deleted deprecated game files and update the **actions** file accordingly as actions get run.

And should the **game files** in your **storage** become out of sync with your **manifest** by accident (file corruption, file deletion, etc), **gogcli storage repair** will generate a new **actions** file in your storage, based on the observed different between the **manifest** and the **game files** in your **storage** (the **manifest** is considered authoritative).

