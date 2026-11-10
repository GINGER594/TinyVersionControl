package main

import (
    "os"
    "os/exec"
    "fmt"
)

var tvcDirRel string = `.tvc` //relative path
var tvcDir string = `./`+tvcDirRel+`/`
var commitCmd string = "comm"
var commitErrMsg string = "error committing version '%s'\n"
var commitSuccessMsg string = "committed version '%s'\n"
var swapCmd string = "swap"
var swapErrMsg string = "error swapping to version '%s'\n"
var swapSuccessMsg string = "swapped to version '%s'\n"
var deleteCmd string = "dele"
var deleteErrMsg string = "error deleting version '%s'\n"
var deleteSuccessMsg string = "deleted version '%s'\n"
var lsCmd string = "ls"
var helpMsg string = commitCmd+" <version>: commit/overwrite\n"+swapCmd+" <version>: swap\n"+deleteCmd+" <version>: delete\n"+lsCmd+": list"

//applies bash command as a functor to dir (other args optional)
func fmapBashCmd(bashCmd, dir string, args ...string) error {
    entities, err := os.ReadDir(dir) //entities will be empty if not exists
    for _, entity := range entities {
        if !(entity.Name() == tvcDirRel && entity.IsDir()) { //not applying bash command to tvc dir
            if err := exec.Command(bashCmd, append([]string{dir+entity.Name()}, args...)...).Run(); err != nil {
                return err
            }
        }
    }
    return err
}

func main() {
    osArgs := append(os.Args, "", "") //+"", "" to prevent unsafe access
    if osArgs[1] == commitCmd && osArgs[2] != "" { //commit ver
        exec.Command("mkdir", tvcDir+osArgs[2]+`/`, "-p").Run() //creating ver dir (exit code will always be 0 due to -p)
        if err := fmapBashCmd("rm", tvcDir+osArgs[2]+`/`, "-r"); err != nil { //removing entities from ver dir
            fmt.Printf(commitErrMsg, osArgs[2])
            return
        } else if err = fmapBashCmd("cp", `./`, tvcDir+osArgs[2]+`/`, "-r"); err != nil { //copying entities from current dir to ver dir
            fmt.Printf(commitErrMsg, osArgs[2])
            return
        }
        fmt.Printf(commitSuccessMsg, osArgs[2])
    } else if osArgs[1] == swapCmd && osArgs[2] != "" { //swap to ver
        if err := exec.Command("ls", tvcDir+osArgs[2]+`/`).Run(); err != nil { //checking existence of ver dir
            fmt.Printf(swapErrMsg, osArgs[2])
            return
        } else if err = fmapBashCmd("rm", `./`, "-r"); err != nil { //removing entities from current dir
            fmt.Printf(swapErrMsg, osArgs[2])
            return
        } else if err = fmapBashCmd("cp", tvcDir+osArgs[2]+`/`, `./`, "-r"); err != nil { //copying entities from ver to current dir
            fmt.Printf(swapErrMsg, osArgs[2])
            return
        }
        fmt.Printf(swapSuccessMsg, osArgs[2])
    } else if osArgs[1] == deleteCmd && osArgs[2] != "" { //delete ver
        if err := exec.Command("rm", tvcDir+osArgs[2]+`/`, "-r").Run(); err != nil { //removing ver dir
            fmt.Printf(deleteErrMsg, osArgs[2])
            return
        }
        fmt.Printf(deleteSuccessMsg, osArgs[2])
    } else if osArgs[1] == lsCmd { //ls ver
        vers, _ := exec.Command("ls", tvcDir).Output() //does nothing if not exists
        fmt.Printf("%s", vers)
    } else { //help msg
        fmt.Println(helpMsg)
    }
}
