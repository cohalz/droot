package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"

	"github.com/yuuki/droot/archive"
	"github.com/yuuki/droot/aws"
	"github.com/yuuki/droot/deploy"
	"github.com/yuuki/droot/log"
)

var CommandArgPull = "--dest DESTINATION_DIRECTORY --src S3_ENDPOINT [--mode MODE]"
var CommandPull = cli.Command{
	Name:   "pull",
	Usage:  "Pull an extracted docker image from s3",
	Action: fatalOnError(doPull),
	Flags: []cli.Flag{
		cli.StringFlag{Name: "dest, d", Usage: "Local filesystem path (ex. /var/containers/app)"},
		cli.StringFlag{Name: "src, s", Usage: "Amazon S3 endpoint (ex. s3://drootexample/app.tar.gz)"},
		cli.StringFlag{Name: "mode, m", Usage: "Mode of deployment. 'rsync' or 'symlink'. default is 'rsync'"},
		cli.BoolFlag{Name: "same-owner", Usage: "Try extracting files with the same ownership as exists in the archive (default for superuser)"},
	},
}

func doPull(c *cli.Context) error {
	if c.String("dest") == "" || c.String("src") == "" {
		cli.ShowCommandHelp(c, "pull")
		return errors.New("--src and --dest option required ")
	}

	destDir, err := filepath.Abs(c.String("dest"))
	if err != nil {
		return err
	}

	s3URL, err := url.Parse(c.String("src"))
	if err != nil {
		return err
	}
	if s3URL.Scheme != "s3" {
		return fmt.Errorf("Not s3 scheme %s", s3URL.String())
	}

	mode := c.String("mode")
	if mode == "" {
		mode = "rsync"
	}
	if mode != "rsync" && mode != "symlink" {
		return fmt.Errorf("Invalid mode %s. '--mode' must be 'rsync' or 'symlink'.", mode)
	}

	tmp, err := ioutil.TempFile(os.TempDir(), "droot_gzip")
	if err != nil {
		return fmt.Errorf("Failed to create temporary file: %s", err)
	}
	defer func(f *os.File) {
		f.Close()
		os.Remove(f.Name())
	}(tmp)

	log.Info("-->", "Downloading", s3URL, "to", tmp.Name())

	if _, err := aws.NewS3Client().Download(s3URL, tmp); err != nil {
		return fmt.Errorf("Failed to download file(%s) from s3: %s", s3URL.String(), err)
	}

	rawDir, err := ioutil.TempDir(os.TempDir(), "droot_raw")
	if err != nil {
		return fmt.Errorf("Failed to create temporary dir: %s", err)
	}
	defer os.RemoveAll(rawDir)
	if err := os.Chmod(rawDir, 0755); err != nil {
		return err
	}

	log.Info("-->", "Extracting archive", tmp.Name(), "to", rawDir)

	if err := archive.ExtractTarGz(tmp, rawDir, c.Bool("same-owner")); err != nil {
		return fmt.Errorf("Failed to extract archive: %s", err)
	}

	if mode == "rsync" {
		log.Info("-->", "Syncing", "from", rawDir, "to", destDir)

		if err := deploy.Rsync(rawDir, destDir); err != nil {
			return fmt.Errorf("Failed to rsync: %s", err)
		}
	} else if mode == "symlink" {
		if err := deploy.DeployWithSymlink(rawDir, destDir); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unreachable code. invalid mode %s", mode)
	}

	return nil
}
