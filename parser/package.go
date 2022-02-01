package parser

func getPackageName(file parserGoFile) string {
	return file.AstFile.Name.Name
}

func getAllGoFilesFromAllPackages(packages []*parserPackage) (goFiles []parserGoFile) {
	for _, pkg := range packages {
		for _, file := range pkg.GoFiles {
			goFiles = append(goFiles, *file)
		}
	}
	return goFiles
}
