package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/algon-320/KIDE/language"
	"github.com/algon-320/KIDE/online_judge"
	"github.com/algon-320/KIDE/util"
	"github.com/urfave/cli"
)

func cmdRun(c *cli.Context) error {
	lang := language.GetLanguage(c.String("language"))
	if err := run(lang); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func cmdTester(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError(util.PrefixError+"few args", 1)
	}

	lang := language.GetLanguage(c.String("language"))
	problemID := c.Args().First()
	if err := tester(lang, problemID, c.Int("case")); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func cmdDl(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError(util.PrefixError+"few args", 1)
	}

	url := c.Args().First()
	if err := downloadSampleCase(url); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func cmdSubmit(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.NewExitError(util.PrefixError+"few args", 1)
	}

	lang := language.GetLanguage(c.String("language"))

	filename, err := language.FindSourceCode(lang)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	p, err := online_judge.LoadProblem(c.Args().First())
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = submit(filename, lang, p)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func cmdView(c *cli.Context) error {
	if c.NArg() < 1 {
		// 引数が無い場合はすべて表示
		list := online_judge.GetAllProblemID()

		// make table
		title := []string{"id", "name", "oj name", "url"}
		data := [][]string{}
		for _, v := range list {
			p, err := online_judge.LoadProblem(v)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			data = append(data, []string{v, p.Name, p.Oj.Name(), p.URL})
		}

		util.PrintTable(title, data, true)
	} else {
		// 引数がある場合は引数で指定された問題を表示する
		problemID := c.Args().First()
		p, err := online_judge.LoadProblem(problemID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		p.Print()
	}
	return nil
}

func cmdProcesser(c *cli.Context) error {
	lang := language.GetLanguage(c.String("language"))

	filename, err := language.FindSourceCode(lang)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	sourceCode, err := ioutil.ReadFile(filename)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	sourceCodeStr := string(sourceCode)

	sourceCodeStr = processSource(sourceCodeStr)
	fmt.Print(sourceCodeStr)
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "KIDE"
	app.Usage = "Kyopro-Iikanjini-Dekiru-Environment"

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run the source-code here",
			Action:  cmdRun,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "language, l",
					Value: "c++",
					Usage: "designate language name",
				},
			},
		},
		{
			Name:    "tester",
			Aliases: []string{"t"},
			Usage:   "test samplecases",
			Action:  cmdTester,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "language, l",
					Value: "c++",
					Usage: "designate language name",
				},
				cli.IntFlag{
					Name:  "case, c",
					Value: -1,
					Usage: "designate samplecase (1-indexed value) testing",
				},
			},
		},
		{
			Name:    "dl",
			Aliases: []string{"d"},
			Usage:   "download samplecases",
			Action:  cmdDl,
		},
		{
			Name:    "submit",
			Aliases: []string{"s"},
			Usage:   "submit solution",
			Action:  cmdSubmit,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "language, l",
					Value: "c++",
					Usage: "designate language name",
				},
			},
		},
		{
			Name:    "view",
			Aliases: []string{"v"},
			Usage:   "view problems",
			Action:  cmdView,
		},
		{
			Name:    "processer",
			Aliases: []string{"p"},
			Usage:   "proccess source code and output",
			Action:  cmdProcesser,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "language, l",
					Value: "c++",
					Usage: "designate language name",
				},
			},
		},
	}

	app.Run(os.Args)
}
