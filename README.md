# logfx
Another simple tool to archive files. With it, you should be able to compress files from a directory and store them to a 
new directory. Beyond of that, you can process a cron job on a set of schedule.

# Features
- Compress files from a provided directory;
- Store the archive to another directory
- Perform cron job intuitively

# Usage
## Clone repository
First of all, clone the repository:
```bash 
git clone https://github.com/sineto/logfx.git

## and then access the cloned project
cd logfx
```

## Golang binary
You can see the list of command just running the binary file. Let's see a overview of it:
```bash
Usage:
  logfx [flags]
  logfx [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  cron        cron creates a execution job for logfx.
  help        Help about any command

Flags:
      --from string   path directory to compress logs (required)
  -h, --help          help for logfx
      --to string     path directory to move the compressed logs (default ".")
      --type string   compression file type - use (targz|zip) (default "targz")

Use "logfx [command] --help" for more information about a command.
```

So, justo do it:
```bash
sudo ./logfx [cron] --from=/var/log/audit --to=/home/johndoe [--expr="* 7 * * *"]
```
### Build from source
If you gave Golang well configured on you machine, you can create a new binary just running this:
```bash
GCO_ENABLED=0 GOOS=linux go build -buildvcs=false -o logfx .
```

# TODO
- [ ] Implement progress bar

# Important things
- I do not pretend port this project to be compatible with MacOS or Windows.

# References
- https://roadmap.sh/projects/log-archive-tool