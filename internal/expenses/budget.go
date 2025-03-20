package expenses

import (
	"database/sql"
	"log"
)

// DefinirOrcamento define ou atualiza o limite de orçamento do usuário
func DefinirOrcamento(db *sql.DB, userID int64, limitAmount float64) error {
	query := `
	INSERT INTO budgets (user_id, limit_amount) 
	VALUES ($1, $2) 
	ON CONFLICT (user_id) DO UPDATE SET limit_amount = $2`
	
	_, err := db.Exec(query, userID, limitAmount)
	if err != nil {
		log.Println("Erro ao definir orçamento:", err)
		return err
	}
	log.Println("✅ Orçamento definido:", limitAmount)
	return nil
}

// ObterOrcamento retorna o limite de orçamento do usuário
func ObterOrcamento(db *sql.DB, userID int64) (float64, error) {
	var limitAmount float64
	query := `SELECT limit_amount FROM budgets WHERE user_id = $1`
	err := db.QueryRow(query, userID).Scan(&limitAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Retorna 0 se o usuário não tiver orçamento definido
		}
		log.Println("Erro ao buscar orçamento:", err)
		return 0, err
	}
	return limitAmount, nil
}

// TotalGastosDoMes retorna a soma dos gastos do mês atual para um usuário
func TotalGastosDoMes(db *sql.DB, userID int64) (float64, error) {
	var total float64
	query := `
	SELECT COALESCE(SUM(amount), 0) 
	FROM expenses 
	WHERE user_id = $1 
	AND EXTRACT(YEAR FROM date) = EXTRACT(YEAR FROM CURRENT_DATE) 
	AND EXTRACT(MONTH FROM date) = EXTRACT(MONTH FROM CURRENT_DATE)`

	err := db.QueryRow(query, userID).Scan(&total)
	if err != nil {
		log.Println("Erro ao calcular gastos do mês:", err)
		return 0, err
	}
	return total, nil
}
