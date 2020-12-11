package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

const (
	NEW_FILE_PERM = 0644
)

type FileInfo os.FileInfo
type FilesInfo []FileInfo

func (f FilesInfo) Len() int      { return len(f) }
func (f FilesInfo) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

// GetFileInfo returns a FileInfo describing the named file
func GetFileInfo(name string) (FileInfo, error) {
	fi, err := os.Stat(name)
	return fi, err
}

// Exists checks if the given filename exists
func Exists(filename string) (bool, error) {
	if _, err := GetFileInfo(filename); err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
	}
	return true, nil
}

// Basename returns the last element of path
func Basename(path string) string {
	return filepath.Base(path)
}

// Dirname returns all but the last element of path, typically the path's directory
func Dirname(path string) string {
	return filepath.Dir(path)
}

// Extname returns the file name extension used by path
func Extname(path string) string {
	return filepath.Ext(path)
}

// Size return the size of the given filename.
// Returns -1 if the file does not exist or if the file size cannot be determined.
func Size(filename string) (int64, error) {
	if fi, err := GetFileInfo(filename); err != nil {
		return -1, err
	} else {
		return fi.Size(), nil
	}
}

// ModTime return the Last Modified Time of the given filename.
func ModTime(filename string) (time.Time, error) {
	if fi, err := GetFileInfo(filename); err != nil {
		return time.Unix(0, 0), err
	} else {
		return fi.ModTime(), nil
	}
}

// ModTimeUnix return the Last Modified Unix Timestamp of the given filename.
// Returns -1 if the file does not exist or if the file modtime cannot be determined.
func ModTimeUnix(filename string) (int64, error) {
	if ft, err := ModTime(filename); err != nil {
		return -1, err
	} else {
		return ft.Unix(), nil
	}
}

// ModTimeUnixNano return the Last Modified Unix Time of nanoseconds of the given filename.
// Returns -1 if the file does not exist or if the file modtime cannot be determined.
func ModTimeUnixNano(filename string) (int64, error) {
	if ft, err := ModTime(filename); err != nil {
		return -1, err
	} else {
		return ft.UnixNano(), nil
	}
}

// Mode return the FileMode of the given filename.
// Returns 0 if the file does not exist or if the file mode cannot be determined.
func Mode(filename string) (os.FileMode, error) {
	if fi, err := GetFileInfo(filename); err != nil {
		return 0, err
	} else {
		return fi.Mode(), nil
	}
}

// Perm return the Unix permission bits of the given filename.
// Returns 0 if the file does not exist or if the file mode cannot be determined.
func Perm(filename string) (os.FileMode, error) {
	if fi, err := GetFileInfo(filename); err != nil {
		return 0, err
	} else {
		return fi.Mode().Perm(), nil
	}
}

// Read reads the file named by filename and returns the contents.
func Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// Write writes data to a file named by filename.
func Write(filename string, content []byte) error {
	return ioutil.WriteFile(filename, content, NEW_FILE_PERM)
}

// ReadString reads the file named by filename and returns the contents as string.
func ReadString(filename string) (string, error) {
	buf, err := Read(filename)
	if err != nil {
		return "", err
	} else {
		return string(buf), nil
	}
}

// WriteString writes the contents of the string to filename.
func WriteString(filename, content string) error {
	return Write(filename, []byte(content))
}

// AppendString appends the contents of the string to filename.
func AppendString(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, NEW_FILE_PERM)
	if err != nil {
		return err
	}
	data := []byte(content)
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// GetContents reads the file named by filename and returns the contents as string.
// GetContents is equivalent to ReadString.
func GetContents(filename string) (string, error) {
	return ReadString(filename)
}

// PutContents writes the contents of the string to filename.
// PutContents is equivalent to WriteString.
func PutContents(filename, content string) error {
	return WriteString(filename, content)
}

// AppendContents appends the contents of the string to filename.
// AppendContents is equivalent to AppendString.
func AppendContents(filename, content string) error {
	return AppendString(filename, content)
}

// TempFile creates a new temporary file in the default directory for temporary files (see os.TempDir), opens the file for reading and writing, and returns the resulting *os.File.
func TempFile() (*os.File, error) {
	return ioutil.TempFile("", "")
}

// TempName creates a new temporary file in the default directory for temporary files (see os.TempDir), opens the file for reading and writing, and returns the filename.
func TempName() (string, error) {

	f, err := TempFile()
	if err != nil {
		return "", err
	}
	return f.Name(), nil

}

// Copy makes a copy of the file source to dest.
func Copy(source, dest string) (err error) {

	// checks source file is regular file
	sfi, err := GetFileInfo(source)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		errors.New("cannot copy non-regular files.")
		return
	}

	// checks dest file is regular file or the same
	dfi, err := GetFileInfo(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !dfi.Mode().IsRegular() {
			errors.New("cannot copy to non-regular files.")
			return
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	// hardlink source to dest
	err = os.Link(source, dest)
	if err == nil {
		return
	}

	// cannot hardlink , copy contents
	in, err := os.Open(source)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return
	}

	// syncing file
	err = out.Sync()

	// trying chmod destination file
	err = out.Chmod(sfi.Mode())
	return
}

// Rename renames (moves) a file
func Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

// Remove removes the named file or directory.
func Remove(name string) error {
	return os.Remove(name)
}

// RemoveAll removes path and any children it contains.
func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Unlink removes the named file or directory.
// Unlink is equivalent to Remove.
func Unlink(name string) error {
	return os.Remove(name)
}

// Rmdir removes path and any children it contains.
// Rmdir is equivalent to RemoveAll.
func Rmdir(path string) error {
	return os.RemoveAll(path)
}

// Mkdir creates a new directory with the specified name and permission bits.
func Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// MkdirAll creates a directory named path, along with any necessary parents.
func MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// MkdirP creates a directory named path, along with any necessary parents.
// MkdirP is equivalent to MkdirAll.
func MkdirP(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Chmod changes the mode of the named file to mode.
func Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

// Chown changes the numeric uid and gid of the named file.
func Chown(name string, uid, gid int) error {
	return os.Chown(name, uid, gid)
}

// Find returns the FilesInfo([]FileInfo) of all files matching pattern or nil if there is no matching file. The syntax of patterns is the same as in Match. The pattern may describe hierarchical names such as /usr/*/bin/ed (assuming the Separator is '/').
func Find(pattern string) (FilesInfo, error) {

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	files := make(FilesInfo, 0)
	for _, f := range matches {
		fi, err := GetFileInfo(f)
		if err == nil {
			files = append(files, fi)
		}
	}
	return files, err
}

type byName struct{ FilesInfo }

type bySize struct{ FilesInfo }

type byModTime struct{ FilesInfo }

func (s byName) Less(i, j int) bool { return s.FilesInfo[i].Name() < s.FilesInfo[j].Name() }
func (s bySize) Less(i, j int) bool { return s.FilesInfo[i].Size() < s.FilesInfo[j].Size() }
func (s byModTime) Less(i, j int) bool {
	return s.FilesInfo[i].ModTime().Before(s.FilesInfo[j].ModTime())
}

// SortByName sorts a slice of files by filename in increasing order.
func (fis FilesInfo) SortByName() {
	sort.Sort(byName{fis})
}

// SortBySize sorts a slice of files by filesize in increasing order.
func (fis FilesInfo) SortBySize() {
	sort.Sort(bySize{fis})
}

// SortByModTime sorts a slice of files by file modified time in increasing order.
func (fis FilesInfo) SortByModTime() {
	sort.Sort(byModTime{fis})
}

// SortByNameReverse sorts a slice of files by filename in decreasing order.
func (fis FilesInfo) SortByNameReverse() {
	sort.Sort(sort.Reverse(byName{fis}))
}

// SortBySizeReverse sorts a slice of files by filesize in decreasing order.
func (fis FilesInfo) SortBySizeReverse() {
	sort.Sort(sort.Reverse(bySize{fis}))
}

// SortByModTimeReverse sorts a slice of files by file modified time in decreasing order.
func (fis FilesInfo) SortByModTimeReverse() {
	sort.Sort(sort.Reverse(byModTime{fis}))
}

// Exec runs the command and returns its standard output as string
func Exec(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.Output()
	return string(out), err
}

// Md5 returns a MD5 hash of file
func Md5(filename string) (string, error) {

	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return "", err
	}

	hash := md5.New()
	io.Copy(hash, file)

	result := hash.Sum(nil)

	return hex.EncodeToString(result), nil
}

// Sha1 returns a SHA1 hash of file
func Sha1(filename string) (string, error) {

	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return "", err
	}

	hash := sha1.New()
	io.Copy(hash, file)

	result := hash.Sum(nil)

	return hex.EncodeToString(result), nil
}

func Touch(filename string) error {
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	return err
}
