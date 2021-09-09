package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunCmd(t *testing.T) {
	dir := filepath.Join(testDir, "testsvc")
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		// go-doudou svc run
		_, _, err = ExecuteCommandC(rootCmd, []string{"svc", "run"}...)
		if err != nil {
			t.Error(err)
			return
		}
	}()
	time.Sleep(2 * time.Second)
}
