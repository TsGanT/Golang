package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/emedvedev/enigma"
	"github.com/mkideal/cli"
)

// CLIOpts sets the parameter format for Enigma CLI. It also includes a "help"
// flag and a "condensed" flag telling the program to output plain result.
// Also, this CLI module abuses tags so much it hurts. Oh well. ¯\_(ツ)_/¯
type CLIOpts struct {
	Help      bool `cli:"!h,help" usage:"Show help."`
	Condensed bool `cli:"c,condensed" name:"false" usage:"Output the result without additional information."`

	Rotors    []string `cli:"rotors" name:"I II III" usage:"Rotor configuration. Supported: I, II, III, IV, V, VI, VII, VIII, Beta, Gamma."`
	Rings     []int    `cli:"rings" name:"1 1 1" usage:"Rotor rings offset: from 1 (default) to 26 for each rotor."`
	Position  []string `cli:"position" name:"A A A" usage:"Starting position of the rotors: from A (default) to Z for each."`
	Plugboard []string `cli:"plugboard" name:"[]" usage:"Optional plugboard pairs to scramble the message further."`

	Reflector string `cli:"reflector" name:"C" usage:"Reflector. Supported: A, B, C, B-Thin, C-Thin."`
}

// CLIDefaults is used to populate default values in case
// one or more of the parameters aren't set. It is assumed
// that rotor rings and positions will be the same for all
// rotors if not set explicitly, so only one value is stored.
var CLIDefaults = struct {
	Reflector string
	Ring      []int
	Position  []string
	Rotors    []string
}{
	Reflector: "C-thin",
	Ring:      []int{1, 1, 1, 16},
	Position:  []string{"B", "C", "B", "Q"},
	Rotors:    []string{"I", "II", "IV", "III"},
}

// SetDefaults sets values for all Enigma parameters that
// were not set explicitly.
// Plugboard is the only parameter that does not require a
// default, since it may not be set, and in some Enigma versions
// there was no plugboard at all.
func SetDefaults(argv *CLIOpts) {

	if argv.Reflector == "" {
		argv.Reflector = CLIDefaults.Reflector
	}
	if len(argv.Rotors) == 0 {
		argv.Rotors = CLIDefaults.Rotors
	}
	if len(argv.Position) == 0 {
		var Positions []string

		Positions = []string{"A", "C", "B", "Q"}
		var Positions2 = []string{"C", "C", "B", "Q"}
		var temp = [][]string{Positions, Positions2}
		argv.Position = temp[0]
	}
	loadRings := (len(argv.Rings) == 0)
	loadPosition := (len(argv.Position) == 0)
	if loadRings || loadPosition {
		argv.Rings = CLIDefaults.Ring
	}
}

func HillClimbing(a string) (T float32) {
	var N float32 = 0.0
	var sum int = 0
	var total float32 = 0.0

	var values = [26]int{}
	for i := 0; i < 26; i++ {
		values[i] = 0
	}
	//calculate the frequence
	var ch2 int
	for _, ch := range a {
		ch2 = int(ch) - 65
		values[ch2]++
		N++
	}
	//calculate the sum of each frequence
	for i := 0; i < 26; i++ {
		ch2 = values[i]
		sum = sum + (ch2 * (ch2 - 1))
	}
	total = float32(sum) / (N * (N - 1))
	return total
}

func Trigram(a string) (Value [10]float32) {
	var values = [10]string{"THE", "AND", "ING", "ENT", "ION", "HER", "FOR", "THA", "NTH", "INT"}
	var valuesS = [10]float32{}
	var number float32 = float32(len(a))
	for i := 0; i < 10; i++ {
		valuesS[i] = 0.0
	}
	//calculate the frequence
	var ch2 string
	for j := 0; j < len(a)-2; j++ {
		ch2 = string(a[j]) + string(a[j+1]) + string(a[j+2])
		for i := 0; i < len(values); i++ {
			if values[i] == ch2 {
				valuesS[i] = valuesS[i] + 1.0
			}
		}
	}
	for k := 0; k < 10; k++ {
		valuesS[k] = valuesS[k] / number
	}
	return valuesS
}

func main() {

	cli.SetUsageStyle(cli.DenseManualStyle)
	cli.Run(new(CLIOpts), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CLIOpts)
		//fmt.Print(argv)
		originalPlaintext := strings.Join(ctx.Args(), " ")

		plaintext := enigma.SanitizePlaintext(originalPlaintext)

		if argv.Help || len(plaintext) == 0 {
			com := ctx.Command()
			com.Text = DescriptionTemplate
			ctx.String(com.Usage(ctx))
			return nil
		}

		arrL := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
		//"A","D","C","Q","B","O","H","N","V", "X","R", "S","J","M",
		//"N","Y","J", "K","C","Q","G","L","X", "Z","U","H","E", "F","A","M",
		//second config: "R","G","O","M","W", "C","S", "J","L", "V","B", "A","T","K","Q", "X", "N", "Y",
		//arrR := []string{"Beta", "Gamma", "I", "II", "V", "VI"}
		//arrR := []string{"VI", "V", "II", "I", "Gamma", "Beta"}
		var Positions = []string{}
		var Rotors = []string{}
		var PlugBoard = []string{}
		var ClimbValue float32 = 0.0
		//var TrigramValue float32 = 0.0
		Positions = []string{"D", "A", "B", "Q"}
		Rotors = []string{"Gamma", "VI", "IV", "III"}
		argv.Rotors = Rotors
		argv.Position = Positions
		// for i := 0; i < len(arrL); i++ {
		// 	for j := 0; j < len(arrL); j++ {
		// 		Positions = []string{arrL[i], arrL[j], "B", "Q"}
		// 		argv.Position = Positions
		// 		for p := 0; p < len(arrR); p++ {
		// 			for q := 0; q < len(arrR); q++ {
		// 				if arrR[p] != arrR[q] {
		// 					Rotors = []string{arrR[p], arrR[q], "IV", "III"}
		// 					argv.Rotors = Rotors
		for m := 0; m < len(arrL); m++ {
			for n := 0; n < len(arrL); n++ {
				PlugBoard = []string{"MS", "KU", "FY", "AG", "BN", "PQ", "HJ", "DI", "ER", "LW"}
				argv.Plugboard = PlugBoard
				config := make([]enigma.RotorConfig, len(argv.Rotors))
				for index, rotor := range argv.Rotors {
					ring := argv.Rings[index]
					value := argv.Position[index][0]
					config[index] = enigma.RotorConfig{ID: rotor, Start: value, Ring: ring}

				}

				e := enigma.NewEnigma(config, argv.Reflector, argv.Plugboard)
				encoded := e.EncodeString(plaintext)

				if argv.Condensed {
					fmt.Print(encoded)
					return nil
				}
				if HillClimbing(encoded) > ClimbValue { //if Trigram(encoded)[0] > 0.00035
					ClimbValue = HillClimbing(encoded)
					fmt.Print(ClimbValue)
					if ClimbValue > 0.01 {
						tmpl, _ := template.New("cli").Parse(OutputTemplate)
						err := tmpl.Execute(os.Stdout, struct {
							Original, Plain, Encoded string
							Args                     *CLIOpts
							Ctx                      *cli.Context
						}{originalPlaintext, plaintext, encoded, argv, ctx})
						fmt.Print(err)
					}
					//fmt.Print(arrL[m] + arrL[n])
				}
			}

		}
		return nil

	})

}
