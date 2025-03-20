package expenses

import (
	"database/sql"
	"log"
	"fmt"
	"time"
)

// Expense representa uma despesa
type Expense struct {
	ID       int
	Amount   float64
	Category string
	Date     time.Time `json:"date"`
}

type Despesa struct {
	Date  time.Time
	Amount float64
	Category string
}

// Estrutura para representar gastos mensais
type GastoMensal struct {
	Mes   string
	Valor float64
}

// AdicionarDespesa insere uma nova despesa no banco de dados
func AdicionarDespesa(db *sql.DB, userID int64, amount float64, category string, date time.Time) error {
	query := `INSERT INTO expenses (user_id, amount, category, date) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, userID, amount, category, date)
	if err != nil {
		log.Println("Erro ao adicionar despesa:", err)
		return err
	}
	log.Println("✅ Despesa adicionada:", amount, category, date)
	return nil
}

// ListarDespesas retorna todas as despesas cadastradas
func ListarDespesas(db *sql.DB, userID int64) ([]Expense, error) {
	query := `SELECT id, amount, category, date FROM expenses WHERE user_id = $1 ORDER BY date DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Erro ao buscar despesas:", err)
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Date); err != nil { // Agora capturamos a data
			log.Println("Erro ao ler linha:", err)
			continue
		}
		expenses = append(expenses, e)
	}

	return expenses, nil
}
// RemoverDespesa remove uma despesa pelo ID
func RemoverDespesa(db *sql.DB, userID int64, id int) error {
	query := `DELETE FROM expenses WHERE id = $1 AND user_id = $2`
	result, err := db.Exec(query, id, userID)
	if err != nil {
		log.Println("Erro ao remover despesa:", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Nenhuma despesa encontrada com o ID %d para este usuário", id)
	}

	log.Println("✅ Despesa removida com sucesso! ID:", id, "User:", userID)
	return nil
}

func RemoverDespesasPorCategoria(db *sql.DB, userID int64, category string) error {
	query := `DELETE FROM expenses WHERE user_id = $1 AND category = $2`
	result, err := db.Exec(query, userID, category)
	if err != nil {
		log.Println("Erro ao remover despesas por categoria:", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Nenhuma despesa encontrada na categoria '%s'", category)
	}

	log.Println("✅ Despesas removidas da categoria:", category)
	return nil
}

func RemoverTodasDespesas(db *sql.DB, userID int64) error {
	query := `DELETE FROM expenses WHERE user_id = $1`
	result, err := db.Exec(query, userID)
	if err != nil {
		log.Println("Erro ao remover todas as despesas:", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Você não tem despesas cadastradas.")
	}

	log.Println("✅ Todas as despesas foram removidas para o usuário:", userID)
	return nil
}

func SomarDespesasPorCategoria(db *sql.DB, userID int64) (map[string]float64, error) {
	query := `SELECT category, SUM(amount) FROM expenses WHERE user_id = $1 GROUP BY category`
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Erro ao buscar despesas por categoria:", err)
		return nil, err
	}
	defer rows.Close()

	despesasPorCategoria := make(map[string]float64)
	for rows.Next() {
		var categoria string
		var total float64
		if err := rows.Scan(&categoria, &total); err != nil {
			log.Println("Erro ao ler linha:", err)
			continue
		}
		despesasPorCategoria[categoria] = total
	}

	return despesasPorCategoria, nil
}

func ListarDespesasDoDia(db *sql.DB, userID int64) ([]Expense, error) {
	query := `SELECT id, amount, category, date FROM expenses WHERE user_id = $1 AND date = CURRENT_DATE ORDER BY date DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Erro ao buscar despesas do dia:", err)
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var e Expense
		if err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Date); err != nil {
			log.Println("Erro ao ler linha:", err)
			continue
		}
		expenses = append(expenses, e)
	}

	return expenses, nil
}

func ListarDespesasDaSemana(db *sql.DB, userID int64) ([]Expense, error) {
	// 🔹 Obtém o intervalo correto de 7 dias
	dataInicio := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	dataFim := time.Now().Format("2006-01-02") // Hoje

	// 🔹 Query corrigida para filtrar despesas corretamente
	query := `SELECT id, category, amount, date 
			  FROM expenses 
			  WHERE user_id = $1 
			  AND date BETWEEN $2::date AND $3::date
			  ORDER BY date DESC`
	rows, err := db.Query(query, userID, dataInicio, dataFim)
	if err != nil {
		log.Println("Erro ao listar despesas da semana:", err)
		return nil, err
	}
	defer rows.Close()

	var despesas []Expense
	for rows.Next() {
		var d Expense
		err := rows.Scan(&d.ID, &d.Category, &d.Amount, &d.Date)
		if err != nil {
			log.Println("Erro ao escanear despesas:", err)
			return nil, err
		}
		despesas = append(despesas, d)
	}

	return despesas, nil
}

func ListarDespesasDoMes(db *sql.DB, userID int64) ([]Expense, error) {
	// 🔹 Obtém o primeiro e o último dia do mês atual
	dataInicio := time.Now().Format("2006-01") + "-01" // Primeiro dia do mês
	dataFim := time.Now().Format("2006-01-02")         // Hoje

	// 🔹 Query SQL corrigida para pegar apenas despesas do mês atual
	query := `SELECT id, category, amount, date 
			  FROM expenses 
			  WHERE user_id = $1 
			  AND date BETWEEN $2::date AND $3::date
			  ORDER BY date DESC`
	rows, err := db.Query(query, userID, dataInicio, dataFim)
	if err != nil {
		log.Println("Erro ao listar despesas do mês:", err)
		return nil, err
	}
	defer rows.Close()

	var despesas []Expense
	for rows.Next() {
		var d Expense
		err := rows.Scan(&d.ID, &d.Category, &d.Amount, &d.Date)
		if err != nil {
			log.Println("Erro ao escanear despesas:", err)
			return nil, err
		}
		despesas = append(despesas, d)
	}

	return despesas, nil
}

func EditarDespesa(db *sql.DB, userID int64, expenseID int64, newCategory string, newDate time.Time) error {
	query := `UPDATE expenses SET category = $1, date = $2 WHERE id = $3 AND user_id = $4`
	_, err := db.Exec(query, newCategory, newDate, expenseID, userID)
	if err != nil {
		log.Println("Erro ao editar despesa:", err)
		return err
	}
	return nil
}

func BuscarDespesasPorPeriodo(db *sql.DB, chatID int64, dataInicial, dataFinal time.Time) ([]Expense, error) {
    rows, err := db.Query(`
        SELECT id, category, amount, date FROM expenses
        WHERE user_id = ? AND date BETWEEN ? AND ?
        ORDER BY date ASC
    `, chatID, dataInicial, dataFinal)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var despesas []Expense
    for rows.Next() {
        var d Expense
        err := rows.Scan(&d.ID, &d.Category, &d.Amount, &d.Date)
        if err != nil {
            return nil, err
        }
        despesas = append(despesas, d)
    }

    return despesas, nil
}
func ListarDespesasPorData(db *sql.DB, chatID int64) ([]Despesa, error) {
	rows, err := db.Query(`
		SELECT date, SUM(amount) as total
		FROM expenses
		WHERE user_id = $1
		GROUP BY date
		ORDER BY date ASC
	`, chatID)

	if err != nil {
		fmt.Println("❌ Erro ao buscar despesas por data:", err)
		return nil, err
	}
	defer rows.Close()

	var despesas []Despesa
	for rows.Next() {
		var despesa Despesa
		if err := rows.Scan(&despesa.Date, &despesa.Amount); err != nil {
			fmt.Println("❌ Erro ao escanear despesas por data:", err)
			return nil, err
		}
		despesas = append(despesas, despesa)
	}

	fmt.Println("📊 Dados de despesas por data:", despesas)
	return despesas, nil
}
func ListarDespesasPorMes(db *sql.DB, chatID int64) ([]GastoMensal, error) {
	rows, err := db.Query(`
		SELECT TO_CHAR(date, 'YYYY-MM') as mes, SUM(amount) as total
		FROM expenses
		WHERE user_id = $1
		GROUP BY mes
		ORDER BY mes ASC
	`, chatID)

	if err != nil {
		fmt.Println("❌ Erro ao buscar despesas por mês:", err)
		return nil, err
	}
	defer rows.Close()

	var gastosMensais []GastoMensal
	for rows.Next() {
		var gasto GastoMensal
		if err := rows.Scan(&gasto.Mes, &gasto.Valor); err != nil {
			fmt.Println("❌ Erro ao escanear despesas por mês:", err)
			return nil, err
		}
		gastosMensais = append(gastosMensais, gasto)
	}

	fmt.Println("📊 Dados de despesas por mês:", gastosMensais)
	return gastosMensais, nil
}
