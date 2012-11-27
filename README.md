Usage: `java -jar bodybuggbypass.jar`

The program writes all retrieved data to a file named: `{timestamp of last data retrieval}.log`

Once all data has been read, the device's memory is cleared and internal time as well as data retrieval times are set to the current date/time.

To avoid the same fate as the FreeTheBugg project, I am not packaging any of the jars this program depends on which belong to BodyMedia, instead you'll need to download and place them in the directory *bodybuggbypass_lib* which must be in the same directory as the executable *jar bodybuggbypass.jar*:
  * [armband-applets-1.10.0-SNAPSHOT.jar](http://application.bodybugg.com/bodybugg/files/static/install/armband-applets-1.10.0-SNAPSHOT.jar)
  * [common-applets-1.10.0-SNAPSHOT.jar](http://application.bodybugg.com//bodybugg/files/static/install/common-applets-1.10.0-SNAPSHOT.jar)
  * [common-shared-1.10.0-SNAPSHOT.jar](http://application.bodybugg.com/bodybugg/files/static/install/common-shared-1.10.0-SNAPSHOT.jar)