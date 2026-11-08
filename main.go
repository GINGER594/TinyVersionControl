package main

import (
    "os"
    "os/exec"
    "fmt"
)

//applies bash command as a functor to dir (other args optional)
func fmapBashCommand(bashCmd, dir string, args []string) {
    entities, _ := os.ReadDir(dir) //entities will be empty if err - no err-handling needed
    for _, entity := range entities {
        if entity.Name()[0] != byte('.') {
            cmd := exec.Command(bashCmd, append([]string{dir+entity.Name()}, args...)...)
            cmd.Run() //malformed commands will fail/do nothing - no err-handling needed
        }
    }
}

func main() {
    if tvcDir := `./.tvc/`; len(os.Args) == 3 {
        os.Mkdir(tvcDir, 0750) //creating backup folder if not exists
        //backup version (creates new if not exists)
        if os.Args[1] == "bu" {
            os.Mkdir(tvcDir+os.Args[2]+`/`, 0750)
            fmapBashCommand("rm", tvcDir+os.Args[2]+`/`, []string{"-r"}) //removing entities from specified version dir
            fmapBashCommand("cp", `./`, []string{tvcDir+os.Args[2]+`/`, "-r"}) //copying entities from current dir to specified version dir
        //revert version (does nothing if not exists)
        } else if _, verExistErr := os.ReadDir(tvcDir+os.Args[2]+`/`); os.Args[1] == "re" && verExistErr == nil {
            fmapBashCommand("rm", `./`, []string{"-r"}) //removing entities from current dir
            fmapBashCommand("cp", tvcDir+os.Args[2]+`/`, []string{`./`, "-r"}) //copying entities from specified version dir to current dir
        }
    //-h flag
    } else if len(os.Args) > 1 {
        if os.Args[1] == "-h" {
            fmt.Println("backup version: bu <version>\nrevert version: re <version>")
        }
    }
}
