package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/TobiasYin/go-lsp/logs"
	"github.com/TobiasYin/go-lsp/lsp"
	"github.com/TobiasYin/go-lsp/lsp/defines"
)

func main() {
	server := lsp.NewServer(&lsp.Options{
		Network: "tcp",
		CompletionProvider: &defines.CompletionOptions{
			TriggerCharacters: &[]string{"."},
		},
	})

	logs.Init(log.Default())

	//server.OnHover(func(ctx context.Context, req *defines.HoverParams) (result *defines.Hover, err error) {
	//	logs.Println("==========================here OnHover======================")
	//	file := readFile(string(req.TextDocument.Uri))
	//	hoveredLine := file[int(req.Position.Line)]
	//	lineParts := strings.Split(hoveredLine, " ")
	//	logs.Println(req)
	//	logs.Println("line parts: ", lineParts)
	//	if len(lineParts) >= 3 && lineParts[0] == "//" && lineParts[1] == ">>>" {
	//		return &defines.Hover{
	//			Contents: defines.MarkupContent{
	//				Kind:  defines.MarkupKindMarkdown,
	//				Value: "Doctest - functionName(testParameters...)",
	//			},
	//		}, nil
	//	} else if len(lineParts) >= 3 && lineParts[0] == "//" && lineParts[1] == "--" {
	//		return &defines.Hover{
	//			Contents: defines.MarkupContent{
	//				Kind:  defines.MarkupKindMarkdown,
	//				Value: "Doctest - expected result",
	//			},
	//		}, nil
	//	}
	//
	//	return nil, nil
	//})

	server.OnCodeLens(func(ctx context.Context, req *defines.CodeLensParams) (result *[]defines.CodeLens, err error) {
		logs.Println("==========================here codelens======================")
		logs.Println(req)
		file := readFile(string(req.TextDocument.Uri))
		var codeLens []defines.CodeLens

		for num, testLine := range file {
			testLineParts := strings.Split(testLine, " ")
			//expectedResLine := file[num+1]
			//expectedResParts := strings.Split(expectedResLine, "")
			//var expectedRes string
			//if len(expectedResParts) >= 3 && expectedResParts[0] == "//" && expectedResParts[1] == "--" {
			//	expectedRes = strings.Join(testLineParts[2:], "")
			//}

			if len(testLineParts) >= 3 && testLineParts[0] == "//" && testLineParts[1] == ">>>" {
				codeLens = append(codeLens, defines.CodeLens{
					Range: defines.Range{
						Start: defines.Position{
							Line:      uint(num),
							Character: 0,
						},
						End: defines.Position{
							Line:      uint(num),
							Character: 0,
						},
					},
					Command: &defines.Command{
						Title:   "doctest: " + strings.Join(testLineParts[2:], ""),
						Command: getTestCommand(strings.Join(testLineParts[2:], "")),
					},
					Data: nil,
				})
			}
		}

		return &codeLens, nil
	})

	server.OnExecuteCommand(func(ctx context.Context, req *defines.ExecuteCommandParams) (err error) {
		logs.Printf("\n\nexecuting command\n\n")

		return nil
	})

	server.OnCompletion(func(ctx context.Context, req *defines.CompletionParams) (result *[]defines.CompletionItem, err error) {
		logs.Println("==========================here OnCompletion======================")
		logs.Println(req)
		d := defines.CompletionItemKindText
		return &[]defines.CompletionItem{defines.CompletionItem{
			Label:      "fill",
			Kind:       &d,
			InsertText: strPtr("Func()"),
		}}, nil
	})

	server.OnDidChangeTextDocument(func(ctx context.Context, req *defines.DidChangeTextDocumentParams) (err error) {
		logs.Println("==========================here OnDidChangeTextDocument======================")
		logs.Println("Text changed")
		return nil
	})

	server.OnExecuteCommand(func(ctx context.Context, req *defines.ExecuteCommandParams) (err error) {
		logs.Println("==========================here OnExecuteCommand======================")
		logs.Println(req)
		return nil
	})

	for _, method := range server.GetMethods() {
		if method != nil {
			logs.Println(*method)
		}
	}

	server.Run()
}

func strPtr(s string) *string {
	return &s
}

func readFile(path string) map[int]string {
	path = path[7:]
	fmt.Printf("\n\npath: %s\n\n", path)
	fileMap := map[int]string{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	i := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileMap[i] = scanner.Text()
		i++
		//fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return fileMap
}

func getTestCommand(funcName string) string {
	return funcName
}
