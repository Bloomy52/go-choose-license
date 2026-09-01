package license

type QuestionID string

const (
	Q1 QuestionID = "Q1"
	Q2 QuestionID = "Q2"
	Q3 QuestionID = "Q3"
	Q4 QuestionID = "Q4"
	Q5 QuestionID = "Q5"
	Q6 QuestionID = "Q6"
	Q7 QuestionID = "Q7"
	Q8 QuestionID = "Q8"
	Q9 QuestionID = "Q9"
)

type Option struct {
	Text         string
	Desc         string
	NextQuestion QuestionID
	LicenseIDs   []string
}

type Question struct {
	ID       QuestionID
	Title    string
	Subtitle string
	Category string
	StepNum  int
	Total    int
	Options  []Option
}

type LanguageNorm struct {
	Language   string
	LicenseIDs []string
	Note       string
}

func GetQuestionsMap() map[QuestionID]Question {
	return map[QuestionID]Question{
		Q1: {
			ID:       Q1,
			StepNum:  1,
			Total:    5,
			Category: "License Preference",
			Title:    "Do you want users and companies to use your code in closed-source / proprietary works?",
			Subtitle: "This determines whether you prefer a Permissive or Copyleft license.",
			Options: []Option{
				{
					Text:         "Yes, allow proprietary / closed-source usage",
					Desc:         "Permissive branch — minimal restrictions on downstream users",
					NextQuestion: Q2,
				},
				{
					Text:         "No, require open-source sharing",
					Desc:         "Copyleft branch — downstream derivative works must remain open source",
					NextQuestion: Q5,
				},
			},
		},
		Q2: {
			ID:       Q2,
			StepNum:  2,
			Total:    5,
			Category: "Permissive Licenses",
			Title:    "Would you like to put this work in the public domain?",
			Subtitle: "Public domain software waives all copyright claims where legal.",
			Options: []Option{
				{
					Text:       "Yes, dedicate to Public Domain / Zero Restrictions",
					Desc:       "Waive copyright completely or offer zero-clause license",
					LicenseIDs: []string{"unlicense", "bsd0"},
				},
				{
					Text:         "No, retain copyright & require disclaimers",
					Desc:         "Proceed to compare permissive options with attribution",
					NextQuestion: Q3,
				},
			},
		},
		Q3: {
			ID:       Q3,
			StepNum:  3,
			Total:    5,
			Category: "Permissive Licenses",
			Title:    "Do you want an Explicit Patent Grant?",
			Subtitle: "An explicit patent grant protects users against patent lawsuits from contributors.",
			Options: []Option{
				{
					Text:       "Yes, include explicit patent protection",
					Desc:       "Recommended for business/corporate backed projects",
					LicenseIDs: []string{"apache-2.0"},
				},
				{
					Text:         "No explicit patent clause needed",
					Desc:         "Proceed to check non-endorsement requirements",
					NextQuestion: Q4,
				},
			},
		},
		Q4: {
			ID:       Q4,
			StepNum:  4,
			Total:    5,
			Category: "Permissive Licenses",
			Title:    "Do you want a Non-Endorsement Clause?",
			Subtitle: "Prevents contributors' names or logos from being used to endorse derivative products.",
			Options: []Option{
				{
					Text:       "Yes, include non-endorsement clause",
					Desc:       "Popular choice in academic & research settings",
					LicenseIDs: []string{"bsd-3-clause"},
				},
				{
					Text:       "No, standard attribution license is sufficient",
					Desc:       "Simple & widely recognized permissive open-source licenses",
					LicenseIDs: []string{"mit", "bsd-2-clause"},
				},
			},
		},
		Q5: {
			ID:       Q5,
			StepNum:  2,
			Total:    5,
			Category: "Copyleft Licenses",
			Title:    "Do you intend for users to run your code directly or include it as a library?",
			Subtitle: "Libraries often use Weak Copyleft while applications use Strong Copyleft.",
			Options: []Option{
				{
					Text:         "Include as a library / package dependency",
					Desc:         "Weak copyleft branch — allowing proprietary application linking",
					NextQuestion: Q6,
				},
				{
					Text:         "Run directly as a standalone application / CLI",
					Desc:         "Strong copyleft branch — whole-project open source requirements",
					NextQuestion: Q7,
				},
			},
		},
		Q6: {
			ID:       Q6,
			StepNum:  3,
			Total:    5,
			Category: "Weak Copyleft Licenses",
			Title:    "Does your library require backwards compatibility with GPLv2 Software?",
			Subtitle: "Some older open-source ecosystems strictly adhere to GPL v2.0.",
			Options: []Option{
				{
					Text:       "Yes, require GPLv2 compatibility",
					Desc:       "Legacy LGPL v2.1 for older GPLv2 ecosystems",
					LicenseIDs: []string{"lgpl-2.1"},
				},
				{
					Text:       "No, use modern LGPL v3.0",
					Desc:       "Updated weak copyleft compatible with GPLv3",
					LicenseIDs: []string{"lgpl-3.0"},
				},
			},
		},
		Q7: {
			ID:       Q7,
			StepNum:  3,
			Total:    5,
			Category: "Strong Copyleft Licenses",
			Title:    "Your code ends up in someone else's project. Do you want to require that:",
			Subtitle: "Choose between file-level copyleft and project-level copyleft.",
			Options: []Option{
				{
					Text:       "A) Only the modified files stay open (File-level copyleft)",
					Desc:       "Users can combine your files with proprietary code without opening everything",
					LicenseIDs: []string{"mpl-2.0"},
				},
				{
					Text:         "B) The entire project must stay open (Project-level copyleft)",
					Desc:         "Requires any combined project to release complete source code",
					NextQuestion: Q8,
				},
			},
		},
		Q8: {
			ID:       Q8,
			StepNum:  4,
			Total:    5,
			Category: "Strong Copyleft Licenses",
			Title:    "Is this project used over a network (e.g., a web app or SaaS service)?",
			Subtitle: "AGPL closes the SaaS / network distribution loophole of standard GPL.",
			Options: []Option{
				{
					Text:       "Yes, closing the SaaS loophole is required",
					Desc:       "Ensures network services running your code must share modified source",
					LicenseIDs: []string{"agpl-3.0"},
				},
				{
					Text:         "No, standard distribution rules apply",
					Desc:         "Proceed to choose between GPL versions",
					NextQuestion: Q9,
				},
			},
		},
		Q9: {
			ID:       Q9,
			StepNum:  5,
			Total:    5,
			Category: "Strong Copyleft Licenses",
			Title:    "Do you need backwards compatibility with GPLv2-only Software?",
			Subtitle: "GPLv2 and GPLv3 are not directly compatible without explicit permission.",
			Options: []Option{
				{
					Text:       "Yes, require GPLv2-only compatibility",
					Desc:       "Required for Linux Kernel modules & legacy GPLv2 projects",
					LicenseIDs: []string{"gpl-2.0"},
				},
				{
					Text:       "No, standard GPLv3 (Modern Strong Copyleft)",
					Desc:       "Includes patent grants & anti-tivoization protections",
					LicenseIDs: []string{"gpl-3.0"},
				},
			},
		},
	}
}

func GetLanguageNorms() []LanguageNorm {
	return []LanguageNorm{
		{Language: "Python", LicenseIDs: []string{"mit"}, Note: "MIT is standard across PyPI packages."},
		{Language: "JavaScript", LicenseIDs: []string{"mit"}, Note: "MIT dominates the NPM ecosystem."},
		{Language: "TypeScript", LicenseIDs: []string{"mit"}, Note: "MIT is standard for TypeScript repositories."},
		{Language: "Java", LicenseIDs: []string{"apache-2.0"}, Note: "Apache 2.0 is popular across Maven & Apache Foundation projects."},
		{Language: "C", LicenseIDs: []string{"bsd-2-clause"}, Note: "BSD 2-Clause / BSD 3-Clause is customary for systems C."},
		{Language: "C++", LicenseIDs: []string{"apache-2.0"}, Note: "Apache 2.0 & MIT are standard for modern C++."},
		{Language: "C#", LicenseIDs: []string{"mit"}, Note: "MIT License is standard for NuGet & .NET packages."},
		{Language: "Go", LicenseIDs: []string{"apache-2.0"}, Note: "Apache 2.0 & BSD 3-Clause are standard in Go modules."},
		{Language: "Rust", LicenseIDs: []string{"mit", "apache-2.0"}, Note: "Dual MIT / Apache 2.0 is standard in Cargo crates."},
		{Language: "PHP", LicenseIDs: []string{"mit"}, Note: "MIT License is standard in Composer packages."},
		{Language: "Ruby", LicenseIDs: []string{"mit"}, Note: "MIT License is standard in RubyGems."},
		{Language: "Swift", LicenseIDs: []string{"apache-2.0"}, Note: "Apache 2.0 is standard for Apple Swift packages."},
		{Language: "Kotlin", LicenseIDs: []string{"apache-2.0"}, Note: "Apache 2.0 is standard for Kotlin & Android libraries."},
		{Language: "Objective-C", LicenseIDs: []string{"mit"}, Note: "MIT License is standard for iOS / Cocoa projects."},
		{Language: "Dart", LicenseIDs: []string{"bsd-3-clause"}, Note: "BSD 3-Clause is standard for Flutter & Dart pub packages."},
		{Language: "Shell/Bash", LicenseIDs: []string{"mit"}, Note: "MIT License is standard for shell scripts & CLI utilities."},
		{Language: "PowerShell", LicenseIDs: []string{"mit"}, Note: "MIT License is standard for PowerShell Gallery modules."},
		{Language: "Perl", LicenseIDs: []string{"artistic-2.0"}, Note: "Artistic License 2.0 is traditional for Perl CPAN modules."},
		{Language: "Lua", LicenseIDs: []string{"mit"}, Note: "MIT License is standard for Lua scripts & modules."},
		{Language: "Tcl", LicenseIDs: []string{"bsd-3-clause"}, Note: "BSD 3-Clause is standard in Tcl/Tk projects."},
		{Language: "SQL", LicenseIDs: []string{"postgresql"}, Note: "PostgreSQL License & MIT are standard for open SQL tools."},
		{Language: "R", LicenseIDs: []string{"gpl-3.0"}, Note: "GPL v3.0+ is standard across CRAN packages."},
		{Language: "MATLAB", LicenseIDs: []string{"bsd-3-clause"}, Note: "BSD 3-Clause is standard for MATLAB scripts."},
		{Language: "Julia", LicenseIDs: []string{"mit"}, Note: "MIT License is standard in Julia General registry."},
		{Language: "SAS", LicenseIDs: []string{"mit"}, Note: "MIT License is standard for open SAS scripts."},
		{Language: "Scala", LicenseIDs: []string{"apache-2.0"}, Note: "Apache 2.0 is standard for Scala & JVM ecosystems."},
		{Language: "Groovy", LicenseIDs: []string{"apache-2.0"}, Note: "Apache 2.0 is standard for Groovy & Gradle plugins."},
		{Language: "Clojure", LicenseIDs: []string{"epl-2.0"}, Note: "Eclipse Public License 2.0 is standard for Clojure libraries."},
		{Language: "Haskell", LicenseIDs: []string{"bsd-3-clause"}, Note: "BSD 3-Clause is standard on Hackage."},
		{Language: "F#", LicenseIDs: []string{"mit"}, Note: "MIT License is standard for F# libraries."},
		{Language: "Elixir", LicenseIDs: []string{"apache-2.0"}, Note: "Apache 2.0 is standard in Hex.pm packages."},
		{Language: "Erlang", LicenseIDs: []string{"apache-2.0"}, Note: "Apache 2.0 is standard for Erlang OTP packages."},
		{Language: "OCaml", LicenseIDs: []string{"mit"}, Note: "MIT License is standard in Opam packages."},
		{Language: "Elm", LicenseIDs: []string{"bsd-3-clause"}, Note: "BSD 3-Clause is standard in Elm package registry."},
		{Language: "Assembly", LicenseIDs: []string{"bsd-2-clause"}, Note: "BSD 2-Clause is standard for low-level assembly routines."},
		{Language: "Zig", LicenseIDs: []string{"mit"}, Note: "MIT License is standard in Zig packages."},
		{Language: "Ada", LicenseIDs: []string{"gpl-3.0"}, Note: "GPL v3.0+ with Runtime Exception is standard for Ada."},
		{Language: "Fortran", LicenseIDs: []string{"bsd-3-clause"}, Note: "BSD 3-Clause is standard for scientific Fortran code."},
		{Language: "COBOL", LicenseIDs: []string{"apache-2.0"}, Note: "Apache 2.0 is common for enterprise COBOL code."},
		{Language: "Visual Basic.NET", LicenseIDs: []string{"mit"}, Note: "MIT License is standard for VB.NET projects."},
	}
}
