# About

This is a command line client to interact with the GOG.com api.

The end goal is to automate the backup of game files.

# Usage

1. Create a file named **cookie** with the following format:

```
sessions_gog_com=<value taken from gog.com>
gog-al=<value taken from gog.com>
```

To get the values, login to gog.com in your browser, then, get the cookie values.

Here is a guide to look at your cookies in Chrome: https://developers.google.com/web/tools/chrome-devtools/storage/cookies

The above is needed, because GOG.com does not yet have an official api for third-party tools with user-generated api keys.

So, any tool wishing to get some kind of api token or cookie programatically without a lot of user-involvement will need to scrape information from the login page and circumvent the recaptcha. At best, this functionality would be flaky and subject to frequent malfunction, so I opted not to go that direction, at the risk of being less user-friendly.

2. Go download a binary for your platform and put it in the same directory as your cookie file: https://github.com/Magnitus-/gogcli/releases

3. See what commands currently supported.

For Linux, you can run the following on the command prompt:

```
./gogcli --help
```

# Netscape Cookies

The client also supports Netscape cookie files. By default, it will use the format defined above, so to use the Netscape format, you need to specify it like so:

```
./gogcli gog-api user-info -y netscape
```

# Supported Storage Solutions

The client supports both the filesystem and s3-compatible object stores (tested with Minio, but should be compatible with Ceph, Swift, Amazon S3, Digital Ocean Spaces and others).

If you use the local filesystem, you just need to provided a path to commands.

If you use an s3 store, you need to provide a path to a configuration file in json format which is as follows:

```
{
    "Endpoint": "<The S3 endpoint of your object store>",
    "Region": "<The S3 region your bucket should be in>",
    "Bucket": "<The bucket in which your manifest and game files should be stored>",
    "Tls": true|false,
    "AccessKey": "<Your access key>",
    "SecretKey": "<Your secret key>"
}
```

# Building The Binaries Yourself

If you prefer, you can build the binary locally:
- Get golang: https://golang.org/dl/
- Run: ```go build```

# Common Use Case Examples

Here, I assume that the gogcli binary is in the **PATH** so that I just type **gogcli** without having to use a relative path (ex: **./gogcli**).

Furthermore, for Windows, the executable would be **gogcli.exe**, but the commands are otherwise the same.

Note that I'm using the long form arguments in all the examples to make everything more legible, but if you'd like to type less, all the arguments have a short form version, following the **POSIX** convention.

## Creating an Initial Manifest

If I want to generate an initial manifest for Linux and Windows files, for the French and English language, I would type:

```
gogcli manifest generate --lang=english --lang=french --os=windows --os=linux
```

This will take a couple of minutes and will produce a **manifest.json** manifest file. 

If GOG.com had some files where not downloadable (they sometimes forget to delete links to files that are no longer available), it will be indicated in the **manifest-404-warnings.json** file.

You can get a summary of your manifest by typing:

```
gogcli manifest summary
```

If you have a very large manifest and want to look for **Master of Orion** and **Master of Magic** entries without having to open the file, you can type:

```
gogcli manifest search --title="Master of Orion" --title="Master of Magic"
```

## Uploading Games From Your Manifest

So now, you are ready to upload your games in your storage.

If you want to store your games in your filesystem, say in the **/home/eric/games** directory, you would type:

```
gogcli storage apply --path=/home/eric/games --storage=fs
```

If instead, you have an s3 store whose configuration information is in the file s3.json (see instructions above for how to configure s3 stores), you would type:

```
gogcli storage apply --path=s3.json --storage=s3
```

If you have 1000 games to upload and would like to only upload 50 games for now and do the rest later, you would type:

```
gogcli storage apply --path=s3.json --storage=s3 --maximum=50
```

And if of those 1000 games, you to download your smaller indie games first and **Faster Than Light** first of all (because it is such a great game), you would first type the following to get FTL's game id:

```
gogcli manifest search --title=FTL
```

And then type:

```
gogcli storage apply --path=s3.json --storage=s3 --maximum=50 --sort-criterion=size --preferred-ids=1207659102
```

If you'd like to look at what's left to download, you can download all the actions that have yet to run by typing:

```
gogcli storage download actions --path=s3.json --storage=s3
```

The remaining actions will be listed in a file called **actions.json**.

Then you decide to continue with the next 50 games tomorrow. Lets say you still want to download smaller indie games first, but this time, you'd like to download all the **Blackwell** games first of all. 

First, you will type the following to find out the game id of the **Blackwell** games:

```
gogcli manifest search --title=Blackwell
```

Then, you will type the following to download the next 50 games, with the **Blackwell** games first

```
gogcli storage resume --path=s3.json --storage=s3 --maximum=50 --sort-criterion=size --preferred-ids=1207662883 --preferred-ids=1207662893 --preferred-ids=1207662903 --preferred-ids=1207662913 --preferred-ids=1207664393
```

And so forth...

If, after you have uloaded everything, you get a little paranoid and want to make really sure that the files are still in the state indicated by the manifest, you can type:

```
gogcli storage validate --path=s3.json --storage=s3
```

If you don't see any output, you are good.

If you don't trust the tool yet and want feedback that it actually did something, you can type:

```
gogcli storage validate --path=s3.json --storage=s3 --debug
```

Be warned, you may get more output than you bargained for.

## Copy Your Files to A Secondary Storage

So now, lets say that you opted for the s3 storage in the example above, but you'd also like to copy your games on your local drive. You can type:

```
gogcli storage copy --source-path=s3.json --source-storage=s3 --destination-path=/home/eric/games --destination-storage=fs
```

Again, if you are feeling paranoid, you can validate that the files copied properly at the destination, by typing:

```
gogcli storage validate --path=/home/eric/games --storage=fs
```

## Updating Your Storage with GOG.com Updates

So now, **GOG.com** released some updates and you would like very much to update your storages.

Lets say you want to update your filesystem storage first.

You have two choices:

### Option 1: Trust GOG.com to report Updates Properly

So first, you download your manifest from your storage by typing:

```
gogcli storage download manifest --path=/home/eric/games --storage=fs
```

Your manifest will be downloaded in **manifest.json**.

After that, you get a list of your updated and new games by typing:

```
gogcli update generate
```

The game ids of your new and updated games will be listed in **updates.json**.

After that, you want to update your manifest with your updates, by typing:

```
gogcli manifest update --update=updates.json
```

Now, you are ready to apply the modifed manifest. Before that, you may wish to run a plan to look at the actions that will run against your storage by typing:

```
gogcli storage plan --path=/home/eric/games --storage=fs
```

Alternatively, if you are ok matching file name and file size couting as the same file when there are no checksums (the gog api doesn't provide it for extras... it should be ok 99.99%+ of the time), you can type:

```
gogcli storage plan --empty-checksum --path=/home/eric/games --storage=fs
```

The actions will be in the **actions.json** file.

You can apply your modifed manifest by typing:

```
gogcli storage apply --path=/home/eric/games --storage=fs
```

Again, if you wish for extras without checksum to count as the same file if the file name and file size match, you type this instead:

```
gogcli storage apply --empty-checksum --path=/home/eric/games --storage=fs
```

Afterwards, if you want to copy the modifications from your filesystem to your s3 store, you can type:

```
gogcli storage copy --source-path=/home/eric/games --source-storage=fs --destination-path=s3.json --destination-storage=s3
```

If you are surprised that it runs really fast, don't worry. The **gogcli storage copy** command doesn't just mindlessly copy files. It actually does a diff between the manifests of both storages and copies only what it must.

### Option 2: Don't Trust GOG.com, Just Check Everything

In this case, you'll just generate a new manifest from scratch and apply it.

Generate the manifest by running:

```
gogcli manifest generate --lang=english --lang=french --os=windows --os=linux
```

Afterwards, you can look at the actions that will run against your storage by typing:

```
gogcli storage plan --empty-checksum --path=/home/eric/games --storage=fs
```

Here, you probably always want to use the **empty-checksum** flag, because all the extras in your generated manifest won't have a checksum (the gog api doesn't provide it) and it will be a lot of extras to upload (for your entire game collection) if you decide that files with an empty checksum don't count as the same file even when the file name and file size matches.

Anyways, the actions will be listed in the **actions.json** file.

If you are satisfied with the actions that will run, you can now apply your manifest by typing:

```
gogcli storage apply --empty-checksum --path=/home/eric/games --storage=fs
```

Afterwards, you can copy the changes from your filesystem to your s3 store by typing:

```
gogcli storage copy --source-path=/home/eric/games --source-storage=fs --destination-path=s3.json --destination-storage=s3
```

## Updating Your Storage with GOG.com When You Have Pending Actions

Ok, so you started storing your games from your manifest using the **gogcli storage apply** command, but it did not complete, either because there was an error or you used the **--maximum** flag.

You took a break and in the interim, gog updated some games and now you're stuck with a storage that has some half-uploaded manifest that is no longer valid.

What do you do?

First of all, you produce an updated **manifest.json** manifest using the **gogcli storage download manifest** + **gogcli update generate** + **gogcli manifest update** commands (option 1)... or you just use a **gogcli manifest generate** command (option 2).

After that, you just run the following:

```
gogcli storage update-actions --path=s3.json --storage=s3
```

Or if you really don't want to redownload extras that don't have a checksum and trust that the file was not updated if the file name and size matches:

```
gogcli storage update-actions --empty-checksum --path=s3.json --storage=s3 
```

What just happened? Your storage's manifest got updated and the remaining actions got adjusted to include additional actions from the change in your manifest.

From there, you can just run **gogcli storage resume** normally.

## Repair Broken Storage

Ok, so ran **storage validate** and it returned some errors, maybe you deleted some files per accident or maybe, you generated a storage with a previous version of gogcli, then migrated your manifest and you'd like to make sure your storage is still ok.

Here, you might not be able to run a **gogli manifest appy**, because this commands only applies a differential between your local manifest and the manifest in the storage. It assumes that the storage's game files reflect the storage's manifest and do not check them separately.

Instead, to reconcile a situation where the game files diverge from what the manifest indicates, assuming your local file **manifest.json** contains a proper manifest (obtained either directly from your storage if you can manage it or otherwise from GOG.com running a **gogli manifest generate** command though in the later case, you'll have to redownload most of your extras as they don't have a checksum), you would run (assuming you have an s3 store):

```
gogcli storage repair --path=s3.json --storage=s3 
```

After running the command, you'll possibly have a bunch of pending actions in your storage if adjustments were needed.

You will run the pending actions in your storage just as you would resume an unfinished **gogli storage apply** by running:

```
gogcli storage resume --path=s3.json --storage=s3 
```

## Dealing With Repeated Download Mismatch

Sometimes, during a download, you might have to deal with repeated errors like this:

```
addFileAction(gameId=1207659025, fileInfo={Kind=extra, Name=darkstone_avatars.zip, ...}, ...) -> Download file size of 445 does not match expected file size of 184891
addFileAction(gameId=1207659025, fileInfo={Kind=extra, Name=darkstone_artworks.zip, ...}, ...) -> Download file size of 445 does not match expected file size of 6430403
```

If it occurs only once, it might be a bad download. If it occurs several times, it might mean that the game has changed since your manifest has been generated.

The best way to deal with the problem is to first get your manifest (using an s3 storage as an example here):

```
gogcli storage download manifest --storage=s3 --path=s3.json
```

Update that game entry in your manifest:

```
gogcli manifest update --id=1207659025
```

And finally, update the actions list in your storage with your updated manifest:

```
gogcli storage update-actions --empty-checksum --path=s3.json --storage=s3
```

## Migration From gogcli 0.9.x

The manifest has to be changed in version 0.10.0, because of an unforeseen situation with languages which forced me to change the manifest format.

To migrate your storage's manifest (using an s3 storage as an example here), copy the manifest to migrate in your current path. Then run:

```
gogcli storage manifest migrate
```

```
gogcli storage repair --path=s3.json --storage=s3
```