package bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"bot-financeiro/internal/database"
	"bot-financeiro/internal/expenses"
	"bot-financeiro/internal/charts"
)
var editingExpenses = make(map[int64]int64)
var buscandoGastos = make(map[int64]struct {
	Inicio string
	Fim    string
})
func IniciarBot() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Erro: TELEGRAM_BOT_TOKEN não encontrado no .env")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Erro ao iniciar bot:", err)
	}

	bot.Debug = true
	log.Printf("✅ Bot %s iniciado!", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	db := database.ConectarDB()
	defer db.Close()

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			msg := update.Message.Text
			if expenseID, exists := editingExpenses[chatID]; exists {
				args := strings.Split(msg, " ")
				if len(args) < 2 {
					enviarMensagem(bot, chatID, "Formato inválido! Use: `NovaCategoria DD-MM-YYYY`")
					continue
				}
	
				newCategory := args[0]
				newDateStr := args[1]
	
				newDate, err := time.Parse("02-01-2006", newDateStr)
				if err != nil {
					enviarMensagem(bot, chatID, "Formato de data inválido! Use `DD-MM-YYYY`.")
					continue
				}
	
				// 🔹 Atualizar a despesa no banco
				err = expenses.EditarDespesa(db, chatID, int64(expenseID), newCategory, newDate)
				if err != nil {
					enviarMensagem(bot, chatID, "Erro ao editar despesa.")
					continue
				}
	
				// 🔹 Remover do mapa temporário
				delete(editingExpenses, chatID)
	
				enviarMensagem(bot, chatID, fmt.Sprintf("✅ Despesa atualizada: %s - %s", newCategory, newDate.Format("02-01-2006")))
				continue // ✅ Agora não cai no comando desconhecido
			}
			switch {
			case strings.HasPrefix(msg, "/help"):
				resposta := "📌 *Comandos disponíveis:*\n\n" +
					"💰 *Gerenciar Despesas:*\n" +
					"➖ `/add VALOR CATEGORIA [DATA]` → Adicionar uma despesa\n" +
					"     Exemplo: `/add 50 Transporte hoje` ou `/add 100 Lazer 15-03-2025`\n" +
					"➖ `/gastos` → Listar todas as suas despesas\n" +
					"➖ `/gastos_categoria` → Ver o total gasto por categoria\n" +
					"➖ `/gastos_dia` → Listar despesas de hoje\n" +
					"➖ `/gastos_semana` → Listar despesas da última semana\n" +
					"➖ `/gastos_mes` → Listar despesas do mês atual\n\n" +
			
					"🗑 *Remover Despesas:*\n" +
					"➖ `/remover` → Exibir despesas e remover interativamente\n" +
					"➖ `/remover tudo` → Remover todas as suas despesas\n\n" +
			
					"🔔 *Orçamentos e Alertas:*\n" +
					"➖ `/definir_orcamento VALOR` → Definir um limite de gastos mensais\n" +
					"     Exemplo: `/definir_orcamento 1000`\n\n" +
			
					"📊 *Gráficos de Despesas:*\n" +
					"➖ `/grafico pizza` → Gráfico de gastos por categoria\n" +
					"➖ `/grafico linha` → Evolução dos gastos ao longo do tempo\n" +
					"➖ `/grafico barras` → Total de gastos por mês\n\n" +
			
					"⚠️ *Atenção!*\n" +
					"O comando `/editar_gastos` *não está funcionando no momento*. Estamos trabalhando para corrigir isso!\n\n" +
			
					"ℹ️ *Outros Comandos:*\n" +
					"➖ `/help` → Mostrar esta lista de comandos\n\n" +
			
					"✅ *Dica:* Antes de remover uma despesa, use `/gastos` para visualizar seus gastos."
			
				msg := tgbotapi.NewMessage(chatID, resposta)
				msg.ParseMode = "Markdown"
				bot.Send(msg)			
			

			case strings.HasPrefix(msg, "/start"):
				resposta := "Olá! Eu sou seu bot financeiro. Você pode usar:\n" +
					"/add valor categoria - Adicionar despesa\n" +
					"/gastos - Listar despesas\n" +
					"/remover - Remover uma despesa interativamente\n" +
					"Para mais detalhes /help."
				enviarMensagem(bot, chatID, resposta)

			case strings.HasPrefix(msg, "/gastos_categoria"):
				despesas, err := expenses.SomarDespesasPorCategoria(db, chatID)
				if err != nil || len(despesas) == 0 {
					enviarMensagem(bot, chatID, "Nenhuma despesa registrada ainda.")
					continue
				}
			
				var resposta strings.Builder
				resposta.WriteString("📊 *Total de Gastos por Categoria:*\n\n")
				for categoria, total := range despesas {
					resposta.WriteString(fmt.Sprintf("- *%s:* R$%.2f\n", categoria, total))
				}
			
				msg := tgbotapi.NewMessage(chatID, resposta.String())
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			
			case strings.HasPrefix(msg, "/gastos_dia"):
				lista, err := expenses.ListarDespesasDoDia(db, chatID)
				if err != nil || len(lista) == 0 {
					enviarMensagem(bot, chatID, "Nenhuma despesa registrada hoje.")
					continue
				}
			
				var resposta strings.Builder
				resposta.WriteString("📅 *Gastos de Hoje:*\n\n")
				for _, d := range lista {
					resposta.WriteString(fmt.Sprintf("- *%s*: R$%.2f\n", d.Category, d.Amount))
				}
			
				msg := tgbotapi.NewMessage(chatID, resposta.String())
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			
			case strings.HasPrefix(msg, "/gastos_semana"):
				lista, err := expenses.ListarDespesasDaSemana(db, chatID)
				if err != nil || len(lista) == 0 {
					enviarMensagem(bot, chatID, "Nenhuma despesa registrada nos últimos 7 dias.")
					continue
				}
			
				var resposta strings.Builder
				resposta.WriteString("📆 *Gastos da Última Semana:*\n\n")
				for _, d := range lista {
					resposta.WriteString(fmt.Sprintf("- *%s*: R$%.2f (%s)\n", d.Category, d.Amount, d.Date.Format("02-01-2006")))
				}
			
				msg := tgbotapi.NewMessage(chatID, resposta.String())
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			
			case strings.HasPrefix(msg, "/definir_orcamento"):
				args := strings.Split(msg, " ")
				if len(args) < 2 {
					enviarMensagem(bot, chatID, "Formato inválido! Use: `/definir_orcamento VALOR`\nExemplo: `/definir_orcamento 1000`")
					continue
				}
			
				limite, err := strconv.ParseFloat(args[1], 64)
				if err != nil {
					enviarMensagem(bot, chatID, "Erro ao converter o valor. Use um número válido.")
					continue
				}
			
				err = expenses.DefinirOrcamento(db, chatID, limite) 
				if err != nil {
					enviarMensagem(bot, chatID, "Erro ao definir orçamento.")
					continue
				}
			
				enviarMensagem(bot, chatID, fmt.Sprintf("✅ Orçamento mensal definido para R$%.2f!", limite))

			case strings.HasPrefix(msg, "/gastos_mes"):
				lista, err := expenses.ListarDespesasDoMes(db, chatID)
				if err != nil || len(lista) == 0 {
					enviarMensagem(bot, chatID, "Nenhuma despesa registrada neste mês.")
					continue
				}
			
				var resposta strings.Builder
				resposta.WriteString("📅 *Gastos do Mês Atual:*\n\n")
				for _, d := range lista {
					resposta.WriteString(fmt.Sprintf("- *%s*: R$%.2f (%s)\n", d.Category, d.Amount, d.Date.Format("02-01-2006")))
				}
			
				msg := tgbotapi.NewMessage(chatID, resposta.String())
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			
			
				
			case strings.HasPrefix(msg, "/add"):
				args := strings.Split(msg, " ")
				if len(args) < 3 {
					enviarMensagem(bot, chatID, "Formato inválido! Use: `/add VALOR CATEGORIA [DATA]`\nExemplo:\n`/add 50 Transporte hoje`\n`/add 30 Alimentação 10-03-2025`")
					continue
				}
			
				valor, err := strconv.ParseFloat(args[1], 64)
				if err != nil {
					enviarMensagem(bot, chatID, "Erro ao converter o valor. Use um número válido.")
					continue
				}
			
				categoria := args[2]
				var data time.Time
			
				// Se uma data for fornecida, interpretamos ela
				if len(args) > 3 {
					dataStr := args[3]
					switch strings.ToLower(dataStr) {
					case "hoje":
						data = time.Now()
					case "ontem":
						data = time.Now().AddDate(0, 0, -1)
					default:
						data, err = time.Parse("02-01-2006", dataStr)
						if err != nil {
							enviarMensagem(bot, chatID, "Formato de data inválido! Use `hoje`, `ontem` ou `DD-MM-YYYY`.")
							continue
						}
					}
				} else {
					data = time.Now() // Se nenhum valor for passado, usa a data de hoje
				}
			
				// 🔹 Adicionar a despesa no banco de dados
				err = expenses.AdicionarDespesa(db, chatID, valor, categoria, data)
				if err != nil {
					enviarMensagem(bot, chatID, "Erro ao adicionar despesa.")
					continue
				}
			
				// 🔹 Verificar orçamento após adicionar a despesa
				limite, err := expenses.ObterOrcamento(db, chatID) // Obtém o limite de orçamento definido pelo usuário
				if err != nil {
					enviarMensagem(bot, chatID, "Erro ao verificar orçamento.")
					continue
				}
			
				if limite > 0 {
					totalGastos, err := expenses.TotalGastosDoMes(db, chatID) // Obtém o total de gastos do mês
					if err == nil && totalGastos > limite {
						enviarMensagem(bot, chatID, fmt.Sprintf("🚨 *Atenção!* Você ultrapassou seu orçamento mensal de R$%.2f.\nGasto atual: R$%.2f", limite, totalGastos))
					}
				}
			
				enviarMensagem(bot, chatID, fmt.Sprintf("✅ Despesa de R$%.2f em %s registrada para %s!", valor, categoria, data.Format("02-01-2006")))
						

			case strings.HasPrefix(msg, "/gastos"):
				lista, err := expenses.ListarDespesas(db, chatID)
				if err != nil || len(lista) == 0 {
					enviarMensagem(bot, chatID, "Nenhuma despesa cadastrada ainda para você.")
					continue
				}
			
				var resposta strings.Builder
				resposta.WriteString("📊 *Suas despesas:*\n\n")
				for _, d := range lista {
					resposta.WriteString(fmt.Sprintf("- *%s*: R$%.2f | *%s*\n", d.Date.Format("02-01-2006"), d.Amount, d.Category))
				}
			
				msg := tgbotapi.NewMessage(chatID, resposta.String())
				msg.ParseMode = "Markdown"
				bot.Send(msg)			

			case strings.HasPrefix(msg, "/remover"):
				args := strings.Split(msg, " ")
			
				if len(args) == 1 {
					lista, err := expenses.ListarDespesas(db, chatID)
					if err != nil || len(lista) == 0 {
						enviarMensagem(bot, chatID, "Você não tem despesas cadastradas.")
						continue
					}
			
					keyboard := tgbotapi.NewInlineKeyboardMarkup()
					for _, d := range lista {
						// 🔹 Agora exibimos também a data no botão
						buttonText := fmt.Sprintf("🗑 R$%.2f | %s (%s)", d.Amount, d.Category, d.Date.Format("02-01-2006"))
						callbackData := fmt.Sprintf("remover_%d", d.ID) // 🔹 Modificado para remover pelo ID, garantindo precisão
			
						row := tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
						)
						keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
					}
			
					msg := tgbotapi.NewMessage(chatID, "Selecione a despesa para remover:")
					msg.ReplyMarkup = keyboard
					bot.Send(msg)
					continue
				}
			
				if args[1] == "tudo" {
					err := expenses.RemoverTodasDespesas(db, chatID)
					if err != nil {
						enviarMensagem(bot, chatID, fmt.Sprintf("Erro: %s", err))
					} else {
						enviarMensagem(bot, chatID, "✅ Todas as suas despesas foram removidas!")
					}
					continue
				}
			
				category := strings.Join(args[1:], " ")
				err := expenses.RemoverDespesasPorCategoria(db, chatID, category)
				if err != nil {
					enviarMensagem(bot, chatID, fmt.Sprintf("Erro: %s", err))
				} else {
					enviarMensagem(bot, chatID, fmt.Sprintf("✅ Todas as despesas da categoria '%s' foram removidas!", category))
				}

			case strings.HasPrefix(msg, "/editar_gastos"):
				lista, err := expenses.ListarDespesas(db, chatID)
				if err != nil || len(lista) == 0 {
					enviarMensagem(bot, chatID, "Você não tem despesas cadastradas para editar.")
					continue
				}

				keyboard := tgbotapi.NewInlineKeyboardMarkup()
				for _, d := range lista {
					buttonText := fmt.Sprintf("✏️ %s - R$%.2f (%s)", d.Category, d.Amount, d.Date.Format("02-01-2006"))
					callbackData := fmt.Sprintf("editar_%d", d.ID)

					row := tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
					)
					keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
				}

				msg := tgbotapi.NewMessage(chatID, "Selecione a despesa que deseja editar:")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			
			case strings.HasPrefix(msg, "/grafico"):
				args := strings.Fields(msg) // Usa Fields para evitar múltiplos espaços
			
				// Se o usuário não especificar um tipo de gráfico, exibe as opções e encerra
				if len(args) < 2 {
					enviarMensagem(bot, chatID, "📊 *Escolha um tipo de gráfico:*\n"+
						"- `/grafico pizza` *(Gastos por categoria)*\n"+
						"- `/grafico linha` *(Evolução dos gastos)*\n"+
						"- `/grafico barras` *(Total por mês)*")
					return
				}
			
				var filePath string
				var err error
			
				// Verifica qual gráfico foi solicitado
				switch args[1] {
				case "pizza":
					filePath, err = charts.GerarGraficoPizza(chatID)
				case "linha":
					filePath, err = charts.GerarGraficoLinha(chatID)
				case "barras":
					filePath, err = charts.GerarGraficoBarras(chatID)
				default:
					enviarMensagem(bot, chatID, "❌ *Tipo de gráfico inválido!*\n"+
						"Use `/grafico pizza`, `/grafico linha` ou `/grafico barras`.")
					return
				}
			
				// Se houve erro ao gerar o gráfico, exibir mensagem de erro
				if err != nil {
					enviarMensagem(bot, chatID, fmt.Sprintf("❌ Erro ao gerar gráfico: %s", err.Error()))
					return
				}
			
				// Enviar o gráfico gerado
				photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(filePath))
				bot.Send(photo)
			
			default:
				enviarMensagem(bot, chatID, "❌ Comando não reconhecido. Use /help para ver a lista de comandos disponíveis.")
			}
		
		}

		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			chatID := update.CallbackQuery.Message.Chat.ID

			// 🔹 Se for um botão de remover despesa por categoria
			if strings.HasPrefix(callbackData, "remover_") {
				category := strings.TrimPrefix(callbackData, "remover_")

				err := expenses.RemoverDespesasPorCategoria(db, chatID, category)
				if err != nil {
					enviarMensagem(bot, chatID, fmt.Sprintf("Erro ao remover despesas da categoria '%s': %s", category, err))
				} else {
					enviarMensagem(bot, chatID, fmt.Sprintf("✅ Todas as despesas da categoria '%s' foram removidas!", category))
				}
			}

			// 🔹 Se for um botão de editar despesa
			if strings.HasPrefix(callbackData, "editar_") {
				expenseIDStr := strings.TrimPrefix(callbackData, "editar_")
				expenseID, err := strconv.ParseInt(expenseIDStr, 10, 64)
				if err != nil {
					enviarMensagem(bot, chatID, "Erro ao processar a despesa selecionada.")
					return
				}

				// Armazena o ID da despesa no mapa temporário para edição
				editingExpenses[chatID] = expenseID
				enviarMensagem(bot, chatID, "✏️ Envie a nova categoria e a nova data no formato:\n`NovaCategoria DD-MM-YYYY`")
			}
		}
	}
}

func enviarMensagem(bot *tgbotapi.BotAPI, chatID int64, texto string) {
	msg := tgbotapi.NewMessage(chatID, texto)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}