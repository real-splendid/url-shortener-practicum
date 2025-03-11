// Пакет main реализует мультичекер для статического анализа.
//
// Этот мультичекер объединяет несколько статических анализаторов:
// - Стандартные анализаторы из golang.org/x/tools/go/analysis/passes
// - Все анализаторы класса SA из staticcheck.io
// - Выбранные анализаторы из других классов staticcheck.io (ST1000)
// - Дополнительные публичные анализаторы (errcheck, bodyclose, nilerr)
//
// Использование:
//
//	go run cmd/staticlint/main.go ./...
//
// Это запустит все анализаторы на указанных пакетах и сообщит о найденных проблемах.
//
// Мультичекер включает следующие категории анализаторов:
//
//  1. Стандартные анализаторы Go: Это встроенные анализаторы из инструментария Go, которые проверяют
//     на наличие распространённых ошибок и проблем в коде Go. К ним относятся:
//     - asmdecl: Сообщает о несоответствиях между файлами ассемблера и объявлениями Go
//     - assign: Обнаруживает бесполезные присваивания
//     - atomic: Проверяет на наличие распространённых ошибок при использовании пакета sync/atomic
//     - bools: Обнаруживает распространённые ошибки, связанные с булевыми операторами
//     - buildtag: Проверяет теги сборки
//     - cgocall: Обнаруживает некоторые нарушения правил передачи указателей в cgo
//     - composite: Проверяет на наличие неименованных составных литералов
//     - copylock: Проверяет на наличие блокировок, ошибочно переданных по значению
//     - errorsas: Проверяет, что второй аргумент errors.As является указателем на тип, реализующий error
//     - httpresponse: Проверяет на наличие ошибок при использовании HTTP-ответов
//     - loopclosure: Проверяет на наличие ссылок на переменные цикла из вложенных функций
//     - lostcancel: Проверяет на отсутствие вызова функции отмены контекста
//     - printf: Проверяет вызовы, похожие на printf
//     - shadow: Проверяет на наличие затенённых переменных
//     - structtag: Проверяет, что теги полей структуры корректно сформированы
//     - tests: Проверяет на наличие распространённых ошибок в использовании тестов и примеров
//     - unmarshal: Проверяет на передачу не-указателя или не-интерфейса в unmarshal
//     - unreachable: Проверяет на наличие недостижимого кода
//     - unusedresult: Проверяет на наличие неиспользуемых результатов вызовов определённых функций
//
//  2. Анализаторы Staticcheck SA: Эти анализаторы сосредоточены на проблемах корректности и потенциальных ошибках.
//     Все анализаторы SA включены, охватывая такие проблемы, как:
//     - Неправильное использование стандартных библиотек
//     - Проблемы с конкурентностью
//     - Проблемы с тестированием
//     - И многие другие проверки корректности
//
//  3. Анализаторы Staticcheck ST: Эти анализаторы сосредоточены на вопросах стиля и согласованности.
//     Включён: ST1000 (проверяет наличие корректных комментариев к пакетам)
//
//  4. Errcheck: Гарантирует, что ошибки, возвращаемые из вызовов функций, проверяются,
//     что помогает предотвратить тихие сбои в коде.
//
//  5. Bodyclose: Проверяет, что тела HTTP-ответов корректно закрываются,
//     предотвращая утечки ресурсов в приложениях, которые делают HTTP-запросы.
//
//  6. Nilerr: Обнаруживает возврат nil вместо корректного значения ошибки, когда функция
//     сталкивается с условием ошибки, что может привести к тихим сбоям.
package main

import (
	"strings"

	"github.com/gostaticanalysis/nilerr"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	var analyzers []*analysis.Analyzer
	standardAnalyzers := []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	}
	analyzers = append(analyzers, standardAnalyzers...)

	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	for _, v := range stylecheck.Analyzers {
		if v.Analyzer.Name == "ST1000" {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	analyzers = append(analyzers, errcheck.Analyzer)

	analyzers = append(analyzers, bodyclose.Analyzer)

	analyzers = append(analyzers, nilerr.Analyzer)

	multichecker.Main(analyzers...)
}
