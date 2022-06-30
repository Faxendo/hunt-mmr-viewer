# Hunt MMR Viewer

## What is it ?

This little tool allows you to check your live MMR from hunt files. It will load data from a file Hunt: Showdown uses to populate the match result screen, so it can also shows you MMR from your last match hunters

## How does it work ?

Hunt: Showdown stores a lot of data used to populate in-game screen on an XML file called `attributes.xml` and located in `<GAME_FOLDER>/user/profiles/default/attributes.xml`. MMR value is stored in plain integer, and later parsed by the game to only show you their "star-notation".

HuntMMRViewer will simply load this XML file, parse it to find the attributes we need and show it to you. Simple as that.

## How can I run it ?

* Go to "Releases" page on Github
* Download the last Release you see
* Extract it on any folder on your computer
* Run "huntmmr.exe"
* The tool will ask for the game base folder. You can find it by going to your Steam Library > Right click on "Hunt: Showdown" > Manage > Browse local files 
* If everything's ok, the HuntMMR loads the XML file and parse it on a table form to you.
* Press any key to exit
* Et voilà !

## I don't trust you and your f***ing exe file

And that's a VERY GOOD THING ! No, really, don't trust anybody on the Internet, I swear.

If you want to build it from sources, you can just `git clone` this repository, and assuming you have Go installed on your computer, simply `go build` it and run the exe file that comes out.

## What's next ?

Some ideas that came from discussions with friends and other players :

* Waiting Hunt v1.9 to add timestamps and more datas on the table
* Develop a web-based *"MMR Tracker"*
* Add a "Sync" option to upload MMR data along your games and track it on the *MMR Tracker*

## Your tool sucks ! It's not working !

Even if the best answer I have is "It's working on my PC", feel free to go to the **Issues** page and ask for help !

## I've got a brilliant idea for this tool

Same as above : go to **Issues** page and share it with the community :)