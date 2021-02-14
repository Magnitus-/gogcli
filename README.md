# WIP

This is a work in progress. It is at this point usable for some use-cases, but not yet complete.

The api is not yet stable and may be subject to change.

Please, do not use yet for serious backups.

# About

This is a command line client to interact with the GOG.com api.

The end goal is to automate the backup of game files.

# Usage

1. Create a file named **cookie** with the following format:

```
sessions_gog_com=<value taken from gog.com>
gog-al=<value taken from gog.com>
```

To get the values, login to gog.com, then, open up the developer console and look at the cookie values of network requests going to gog.com.

For Chrome, you can do so as follows: 
  - left click
  - click on "Inspect"
  - click on the "Network" tab
  - Reload the main page
  - You should see a request to **www.gog.com**, click on it
  - Click on the **Headers** tab
  - Go down to the **Request Headers** section
  - All the cookie values are next to the **cookie:** caption in the format ```<key>=<value>```. They are separated by a **;** character.

The above is needed, because GOG.com does not yet have an official api for third-party tools with user-generated api keys.

So, any tool wishing to get some kind of api token or cookie programatically without a lot of user-involvement will need to scrape information from the login page and circumvent the recaptcha. At best, this functionality would be flaky and subject to frequent malfunction, so I opted not to go that direction, at the risk of being less user-friendly.

2. Get golang: https://golang.org/dl/

3. Build the binary by running:

```
go build
```

4. See what commands currently supported.

For Linux, you can run the following on the command prompt:

```
./gogcli --help
```

Not yet sure what the Windows/MacOS equivalent are, but you should have a runable binary that you can use.

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

# Note

There will be a pipeline later on to generate and publish the binaries once the tool is more stable, so usage step 2 & 3 will no longer be required.