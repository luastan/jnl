package jnl

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

/*

.jnl dir management

*/

var (
	JNL_DIR_NAME = ".jnl"
	JNL_DIR      string
	INFO_DIR     = "info"
)

func init() {
	path, err := filepath.Abs(JNL_DIR_NAME)

	if err != nil {
		JNL_DIR = JNL_DIR_NAME
	} else {
		JNL_DIR = path
	}
	// err = CleanWhileTesting()
	// if err != nil {
	// 	ErrorLogger.Fatalln(err.Error())
	// }
	err = DoInitialize()
	if err != nil {
		ErrorLogger.Fatalln(err.Error())
	}
}

func CleanWhileTesting() error {
	InfoLogger.Printf("Deleting '%s'...", JNL_DIR)
	return os.RemoveAll(JNL_DIR)
}

func IsInitialized() (bool, error) {
	_, err := os.Stat(JNL_DIR)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func DoInitialize() error {
	alreadyDone, err := IsInitialized()
	if err != nil {
		return err
	}

	if alreadyDone {
		return nil
	}
	return os.Mkdir(JNL_DIR, 0755)
}

func SaveState(ex *Execution) error {
	fullExDir := filepath.Join(JNL_DIR, ex.Info.Dir)
	_, err := os.Stat(fullExDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(fullExDir, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	f, err := os.Create(filepath.Join(fullExDir, INFO_DIR))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	b, err := json.Marshal(ex.Info)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func validPaths(root string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && path == JNL_DIR {
			return nil
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return filepath.SkipDir
	})

	return dirs, err
}

func unmarshallCommandInfo(dir string) (*CommandInfo, error) {
	contents, err := os.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	unmarshalled := &CommandInfo{}
	err = json.Unmarshal(contents, unmarshalled)
	if err != nil {
		return nil, err
	}
	return unmarshalled, nil
}

func EveryExecution() ([]*CommandInfo, error) {
	dirs, err := validPaths(JNL_DIR)
	if err != nil {
		return nil, err
	}

	var executions []*CommandInfo

	for _, dir := range dirs {
		data, err := unmarshallCommandInfo(filepath.Join(dir, INFO_DIR))
		if err != nil {
			ErrorLogger.Println(err.Error())
		} else {
			executions = append(executions, data)
		}
	}
	sort.Slice(executions, func(i, j int) bool {
		return executions[i].Start.Before(executions[j].Start)
	})
	return executions, nil
}
