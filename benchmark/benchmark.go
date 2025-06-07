package benchmark

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

type OperationType string

const (
	Insert OperationType = "INSERT"
	Query  OperationType = "QUERY"
)

type DatabaseType string

const (
	Postgres  DatabaseType = "PostgreSQL"
	MongoDB   DatabaseType = "MongoDB"
	Cassandra DatabaseType = "Cassandra"
)

type BenchmarkResult struct {
	Database   DatabaseType
	Operation  OperationType
	Entity     string
	Duration   time.Duration
	RecordSize int
}

type BenchmarkLogger struct {
	results []BenchmarkResult
	logFile *os.File
}

func NewBenchmarkLogger(logFilePath string) (*BenchmarkLogger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar arquivo de log: %v", err)
	}

	return &BenchmarkLogger{
		results: make([]BenchmarkResult, 0),
		logFile: file,
	}, nil
}

func (b *BenchmarkLogger) AddResult(result BenchmarkResult) {
	b.results = append(b.results, result)
}

func (b *BenchmarkLogger) MeasureOperation(db DatabaseType, op OperationType, entity string, recordSize int, operation func() error) {
	start := time.Now()
	err := operation()
	duration := time.Since(start)

	if err != nil {
		log.Printf("Erro durante operação %s em %s para entidade %s: %v\n", op, db, entity, err)
		return
	}

	b.AddResult(BenchmarkResult{
		Database:   db,
		Operation:  op,
		Entity:     entity,
		Duration:   duration,
		RecordSize: recordSize,
	})
}

func (b *BenchmarkLogger) GenerateReport() error {
	if len(b.results) == 0 {
		return fmt.Errorf("nenhum resultado para gerar relatório")
	}

	// Cabeçalho do relatório
	header := "\n=== Relatório de Performance ===\n\n"
	if _, err := b.logFile.WriteString(header); err != nil {
		return err
	}

	w := tabwriter.NewWriter(b.logFile, 0, 0, 3, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "Banco de Dados\tOperação\tEntidade\tTamanho\tDuração\tRegistros/Segundo\t")
	fmt.Fprintln(w, strings.Repeat("-", 80))

	for _, r := range b.results {
		recordsPerSecond := float64(r.RecordSize) / r.Duration.Seconds()
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%.2f\t\n",
			r.Database,
			r.Operation,
			r.Entity,
			r.RecordSize,
			r.Duration.Round(time.Millisecond),
			recordsPerSecond,
		)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("erro ao gerar tabela: %v", err)
	}

	return nil
}

func (b *BenchmarkLogger) Close() error {
	if err := b.GenerateReport(); err != nil {
		return err
	}
	return b.logFile.Close()
}
