# Simple script to clone, install a node.js repository and pull GET
# requests as local json files from that server.

# WARNING this script will terminate all running node.exe applications

import os
from json import dump
from multiprocessing import Process, Value
from sys import argv

from Naked.toolshed.shell import execute, muterun
from requests import exceptions, get

# python load_data.py http://github.com/vincit/summer-2018.git 8080 species sightings


def getJSON(paths, sharedVal, port=8081, base="localhost"):
    'Pulls the JSON requests from the server and saves them as json files'

    if not os.path.exists("rsc/"):
        os.mkdir("rsc")

    while(True):
        try:
            resp = [get("http://" + base + ":" + str(port) +  "/" + p) for p in paths]
            for i in zip(resp, paths):
                with open("rsc/{:s}.json".format(i[1]), "w") as file:
                    dump(i[0].json(), file)
            # Increase the shared value thus breaking the control process from the loop
            sharedVal.value += 1
            return
        # if the getJSON process is faster than the startServer, there will be nothing to connect to
        except exceptions.ConnectionError:
            pass


def installRepo(repo, name):
    'Clone and install node repository'

    install = ["cd {:s}".format(name)]
    # If repo doesn't exist
    if not os.path.exists(name):
        install.insert(0, "git clone {:s}".format(repo))
        install.insert(2, "npm install")

    # if node server has not been installed
    elif not os.path.exists("{:s}/node_modules".format(name)):
        install.insert(1, "npm install")

    if len(install) > 1:
        execute(" && ".join(install))


def startNode(name, port=8081, script="server"):
    'Starts Node.js server with given port'

    commands = ["cd {:s}".format(name), "SET PORT={:d}".format(port), "node {:s}.js".format(script)]
    muterun(" && ".join(commands))


def controllProcess(name, var, port, files):

    #TODO terminate this process
    nodeProcess = Process(target=startNode, args=(name,port), daemon=True)
    nodeProcess.start()

    getProcess = Process(target=getJSON, args=(files, var, port), daemon=True)
    getProcess.start()

    while(var.value == 0):
        pass
    getProcess.join()

    # Close the node server without promting the keyboard interrupt
    muterun("TASKKILL /F /IM node.exe /T")

    return


def main():

    try:
        parts = argv[1].split("/")
        name = parts[len(parts) - 1].split(".")[0]
    except IndexError:
        print("Provide repository for script")
        return

    try:
        port = int(argv[2])
    except (ValueError, IndexError):
        print("Provide valid port number for server")
        return

    try:
        _ = argv[3]
        files = [x for x in argv[3:] if not os.path.exists("rsc/{:s}.json".format(x))]
    except IndexError:
        print("Provide paths for the GET method")
        return


    if files:

        installRepo(argv[1], name)
        var = Value("i", 0)

        # Spawn the main process thread.
        controll = Process(target=controllProcess, args=(name, var, port, files))
        controll.run()
        while(var.value < 1):
            pass

        print("Data fetching completed.")
        return

    else:
        print("Data already exists")


if __name__ == "__main__":
    main()
