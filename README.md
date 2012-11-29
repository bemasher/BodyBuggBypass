# Usage

### Device Dump/Clear
Usage: `java -jar bodybuggbypass.jar`

The program writes all retrieved data to a file named: `{timestamp of last data retrieval}.log`

Once all data has been read, the device's memory is cleared and internal time as well as data retrieval times are set to the current date/time.

### Dump Parser
Building the parser requires golang. Steps to build are as follows:

	go get github.com/bemasher/errhandler
	go build parser.go

Usage: `parser [INFILE] [[OUTFILE]]`

Parameters:

  * `INFILE` is required.
  * `OUTFILE` is optional and defaults to "data.json"

The input file is parsed and written as a json encoding to `OUTFILE`. The basic structure of the file is a list of the following struct:

	Session {
		Channel string
		Epoch int64
		Payload []arbitrary
	}

  * Channel represents the channel name of the session.
  * Epoch is the unix timestamp of the time the first data point was recorded.
  * Payload is a homogeneous list of any of the following types:
    * `uint16`
    * `int64`
    * `{int64, int, int}`

# Caveats

To avoid the same fate as the FreeTheBugg project, I am not packaging any of the jars this program depends on which belong to BodyMedia. Instead you'll need to download and place them in the directory *bodybuggbypass_lib* which must be in the same directory as the executable *jar bodybuggbypass.jar*:

  * [armband-applets-1.10.0-SNAPSHOT.jar](http://application.bodybugg.com/bodybugg/files/static/install/armband-applets-1.10.0-SNAPSHOT.jar)
  * [common-applets-1.10.0-SNAPSHOT.jar](http://application.bodybugg.com/bodybugg/files/static/install/common-applets-1.10.0-SNAPSHOT.jar)
  * [common-shared-1.10.0-SNAPSHOT.jar](http://application.bodybugg.com/bodybugg/files/static/install/common-shared-1.10.0-SNAPSHOT.jar)

# Information
Sample dump and parsed data are found in `sample.log` and `sample.json`.

More information can be found at: [Reverse Engineering the BodyBugg](http://www.bemasher.net/archives/1130)