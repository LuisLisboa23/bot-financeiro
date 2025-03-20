package charts

import (
	"fmt"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
	"bot-financeiro/internal/database"
	"bot-financeiro/internal/expenses"
)

// 🔹 Gráfico de Pizza
func GerarGraficoPizza(chatID int64) (string, error) {
	db := database.ConectarDB()
	defer db.Close()

	// Buscar os gastos por categoria
	gastos, err := expenses.SomarDespesasPorCategoria(db, chatID)
	if err != nil || len(gastos) == 0 {
		return "", fmt.Errorf("Nenhum dado encontrado para gráfico de pizza")
	}

	var valores []chart.Value
	for categoria, total := range gastos {
		valores = append(valores, chart.Value{
			Value: total,
			Label: fmt.Sprintf("%s (R$%.2f)", categoria, total),
		})
	}

	graph := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: valores,
	}

	fileName := fmt.Sprintf("gastos_categoria_%d.png", chatID)
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = graph.Render(chart.PNG, f)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// 🔹 Gráfico de Linha - Gastos ao longo do tempo
func GerarGraficoLinha(chatID int64) (string, error) {
    db := database.ConectarDB()
    defer db.Close()

    gastos, err := expenses.ListarDespesasPorData(db, chatID)
    if err != nil || len(gastos) == 0 {
        return "", fmt.Errorf("Nenhum dado encontrado para gráfico de linha")
    }

    var datas []time.Time
    var valores []float64

    for _, gasto := range gastos {
        datas = append(datas, gasto.Date) // Corrigido para "Date"
        valores = append(valores, gasto.Amount)
    }

    graph := chart.Chart{
        Width:  800,
        Height: 400,
        Series: []chart.Series{
            chart.TimeSeries{
                Name:    "Gastos ao longo do tempo",
                XValues: datas,
                YValues: valores,
            },
        },
        XAxis: chart.XAxis{
            Name:      "Data",
            NameStyle: chart.StyleShow(),
            Style:     chart.StyleShow(),
        },
        YAxis: chart.YAxis{
            Name:      "Valor (R$)",
            NameStyle: chart.StyleShow(),
            Style:     chart.StyleShow(),
        },
    }

    fileName := fmt.Sprintf("evolucao_gastos_%d.png", chatID)
    f, _ := os.Create(fileName)
    defer f.Close()

    _ = graph.Render(chart.PNG, f)
    return fileName, nil
}

// 🔹 Gráfico de Barras - Gastos por mês
func GerarGraficoBarras(chatID int64) (string, error) {
	db := database.ConectarDB()
	defer db.Close()

	// Buscar gastos por mês
	gastos, err := expenses.ListarDespesasPorMes(db, chatID)
	if err != nil || len(gastos) == 0 {
		return "", fmt.Errorf("Nenhum dado encontrado para gráfico de barras")
	}

	var meses []string
	var valores []float64

	for _, gasto := range gastos {
		meses = append(meses, gasto.Mes)  // Agora os meses aparecem corretamente
		valores = append(valores, gasto.Valor)
	}

	// Criar gráfico de barras com os rótulos dos meses
	graph := chart.BarChart{
		Width:  800,
		Height: 500,
		BarWidth: 50, // Ajusta a largura das barras para melhor legibilidade
		XAxis: chart.Style{
			Show: true,
			TextRotationDegrees: 45, // Rotação para não sobrepor os meses
			FontSize: 12,            // Aumenta a fonte dos meses
		},
		YAxis: chart.YAxis{
			Name:      "Gastos em R$",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Bars: []chart.Value{},
	}

	// Adiciona os valores das barras com melhor espaçamento
	for i, mes := range meses {
		graph.Bars = append(graph.Bars, chart.Value{
			Value: valores[i],
			Label: mes, // Agora os meses aparecem no eixo X
		})
	}

	// Gerar o arquivo da imagem
	fileName := fmt.Sprintf("gastos_mes_%d.png", chatID)
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = graph.Render(chart.PNG, file)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

