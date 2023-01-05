package client

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"gopkg.in/yaml.v2"
)

//===========================================================================
// Edit the profile from the command line app
//===========================================================================

// The command line editors to search the $PATH for if $EDITOR is not specified.
var editorSearch = [3]string{"vim", "emacs", "nano"}

// Edit the profiles YAML at the specified path
func EditProfiles() (err error) {
	var path string
	if path, err = ProfilesPath(); err != nil {
		return err
	}
	return Edit(path, ValidProfiles)
}

// Edit the file at the specified path using a command line editor.
func Edit(path string, validate validator) error {
	return EditWith(path, "", validate)
}

// Edit the file at the specified path using the specified command line editor.
func EditWith(path, editor string, validate validator) (err error) {
	// Find the editor to use
	if editor, err = findEditor(editor); err != nil {
		return err
	}

	// Create a temporary file and copy the original file to it
	var tmpf string
	if tmpf, err = mktmpf(); err != nil {
		return fmt.Errorf("could not create temporary file for editing: %v", err)
	}
	defer os.Remove(tmpf)

	if err = copy2(path, tmpf); err != nil {
		return fmt.Errorf("could not copy source contents into temporary file for editing: %v", err)
	}

	// Execute the editor on the temporary file
	cmd := exec.Command(editor, tmpf)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("could not exec %s: %v", editor, err)
	}

	// Validate the written file before editing the original
	if validate != nil {
		if err = validate(tmpf); err != nil {
			return fmt.Errorf("validation error: %s", err)
		}
	}

	// If the editor exited succesfully, copy temporary file back to original file
	if err = copy2(tmpf, path); err != nil {
		return fmt.Errorf("could not copy temporary file contents back to source after editing: %v", err)
	}
	return nil
}

// Finds the path to the specified editor name, or if none is specified, uses the
// $EDITOR environment variable or a search for the standard editors. Returns an error
// if an editor can not be found in the $PATH.
func findEditor(name string) (string, error) {
	if name == "" {
		name = os.Getenv("EDITOR")
	}

	// Determine if the specified editor can be executed
	if name != "" {
		// Expand environment variables and ~ for the home directory.
		name = expand(name)

		// If name is a full path and the file is executable, return it.
		if isExecutable(name) {
			return name, nil
		}

		// Check if the name exists in the Path, if so, return it.
		return inPath(name)
	}

	// Search for one of the editors in the $PATH
	for _, name := range editorSearch {
		if path, err := inPath(name); err == nil {
			return path, nil
		}
	}

	// Could not find an editor
	return "", errors.New("could not find an editor")
}

//===========================================================================
// Validators
//===========================================================================

// A validator is applied to an editor to ensure that the edited file is in the correct
// format or meets other formatting specifications before the original file is modified.
type validator func(path string) error

// ValidProfiles loads the edited YAML into a profiles struct to ensure it's valid.
func ValidProfiles(path string) (err error) {
	var p *Profiles
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("invalid profiles YAML: %s", err)
	}
	return nil
}

//===========================================================================
// Helper Functions
//===========================================================================

// Returns true if the file exists and it can be executed on Unix systems.
func isExecutable(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		if !stat.IsDir() {
			return stat.Mode()&0111 != 0
		}
	}
	return false
}

// Searches for the specified editor in the $PATH
func inPath(name string) (path string, err error) {
	var fname string
	if fname, err = exec.LookPath(name); err != nil {
		return "", fmt.Errorf("could not find %q in $PATH", name)
	}
	if path, err = filepath.Abs(fname); err != nil {
		return fname, nil
	}
	return path, nil
}

// Expand the path from environment variables and handle ~ for the home directory.
func expand(path string) string {
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", "$HOME", 1)
	}
	return os.ExpandEnv(path)
}

func mktmpf() (_ string, err error) {
	var f *os.File
	if f, err = os.CreateTemp("", "goedit-*"); err != nil {
		return "", err
	}
	f.Close()
	return f.Name(), nil
}

// Copy the contents from the src path to the dst path
func copy2(src, dst string) (err error) {
	// Check the source path to make sure it is editable.
	var stat os.FileInfo
	if stat, err = os.Stat(src); err != nil {
		return fmt.Errorf("could not stat source file: %v", err)
	}

	if !stat.Mode().IsRegular() {
		return fmt.Errorf("%q is not a regular file", src)
	}

	var (
		source *os.File
		target *os.File
	)

	if source, err = os.Open(src); err != nil {
		return fmt.Errorf("could not open %q: %v", src, err)
	}
	defer source.Close()

	if target, err = os.Create(dst); err != nil {
		return fmt.Errorf("could not create %q: %v", dst, err)
	}
	defer target.Close()

	if _, err = io.Copy(target, source); err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	// Attempt to change the mode of the target file to the original mode (ignore errors)
	target.Close()
	os.Chmod(dst, stat.Mode())

	// Attempt to cahnge the owners of the target file to the original owners (ignore errors)
	if info, ok := stat.Sys().(*syscall.Stat_t); ok {
		os.Chown(dst, int(info.Uid), int(info.Gid))
	}

	return nil
}
