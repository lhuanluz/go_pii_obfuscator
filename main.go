package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

type PersonalData struct {
    CPF       string `json:"CPF"`
    Email     string `json:"Email"`
    Birthdate string `json:"Birthdate"`
}

func obfuscateData(data *PersonalData) {
    cpfClean := strings.NewReplacer(".", "", "-", "").Replace(data.CPF)
    if len(cpfClean) == 11 {
        data.CPF = cpfClean[:3] + strings.Repeat("*", 6) + cpfClean[9:]
    }

    parts := strings.Split(data.Email, "@")
    namePart := parts[0]
    if len(namePart) > 5 {
        namePart = namePart[:3] + strings.Repeat("*", len(namePart)-5) + namePart[len(namePart)-2:]
    }
    data.Email = namePart + "@" + parts[1]

    dateParts := strings.Split(data.Birthdate, "/")
    if len(dateParts) == 3 {
        data.Birthdate = dateParts[1] + "/" + dateParts[2]
    }
}

func main() {
    var inputFilePath, outputFilePath string
    flag.StringVar(&inputFilePath, "i", "", "Input JSON file path")
    flag.StringVar(&outputFilePath, "o", "", "Output JSON file path")
    flag.Parse()

    if inputFilePath == "" {
        fmt.Println("Please provide an input file path.")
        os.Exit(1)
    }

    if outputFilePath == "" {
        base := filepath.Base(inputFilePath)
        ext := filepath.Ext(inputFilePath)
        name := base[0 : len(base)-len(ext)]
        outputFilePath = "ob_" + name + ".json"
    }

    file, err := ioutil.ReadFile(inputFilePath)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    var data PersonalData
    err = json.Unmarshal(file, &data)
    if err != nil {
        fmt.Println("Error parsing JSON:", err)
        return
    }

    obfuscateData(&data)

    obfuscatedData, err := json.MarshalIndent(data, "", "    ")
    if err != nil {
        fmt.Println("Error marshalling JSON:", err)
        return
    }

    err = ioutil.WriteFile(outputFilePath, obfuscatedData, 0644)
    if err != nil {
        fmt.Println("Error writing file:", err)
    } else {
        fmt.Printf("Obfuscated data written to %s\n", outputFilePath)
    }
}

