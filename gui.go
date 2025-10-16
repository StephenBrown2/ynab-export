package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type guiState int

const (
	guiStateToken guiState = iota
	guiStateBudgetSelect
	guiStateExporting
	guiStateDone
	guiStateError
)

type guiApp struct {
	window           *app.Window
	theme            *material.Theme
	exportBtn        widget.Clickable
	backBtn          widget.Clickable
	errorOkBtn       widget.Clickable
	exportAnotherBtn widget.Clickable
	closeBtn         widget.Clickable
	continueBtn      widget.Clickable
	budgetClickables []widget.Clickable
	selectedBudget   budget
	token            string
	errorMsg         string
	exportPath       string
	tokenEditor      widget.Editor
	jsonData         []byte
	budgets          []budget
	budgetList       widget.List
	summary          budgetSummary
	selectedIndex    int
	state            guiState
}

func runGUI() {
	go func() {
		w := &app.Window{}
		w.Option(app.Title("YNAB Export Tool"))
		w.Option(app.Size(unit.Dp(360), unit.Dp(520)))

		gui := &guiApp{
			window:        w,
			theme:         material.NewTheme(),
			state:         guiStateToken,
			selectedIndex: -1,
		}

		// Load fonts
		gui.theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

		// Configure token editor
		gui.tokenEditor.SingleLine = true
		gui.tokenEditor.Submit = true
		gui.tokenEditor.Mask = '•' // Password mask

		// Configure budget list
		gui.budgetList.Axis = layout.Vertical

		// Check for token in environment
		if token := os.Getenv("YNAB_API_TOKEN"); token != "" {
			gui.token = token
			gui.tokenEditor.SetText(token)
			gui.state = guiStateExporting
			gui.window.Invalidate()
			if err := validateToken(token); err == nil {
				go gui.fetchBudgets()
			} else {
				// Token invalid, show token entry with error
				gui.showError(fmt.Errorf("environment token is invalid: %w", err))
			}
		} else {
			// No token in environment, show token entry screen
			gui.state = guiStateToken
		}

		if err := gui.run(); err != nil {
			fmt.Fprintf(os.Stderr, "GUI error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func (g *guiApp) run() error {
	var ops op.Ops

	for {
		switch e := g.window.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Handle global key events
			for {
				ev, ok := gtx.Event(key.Filter{Name: key.NameEscape})
				if !ok {
					break
				}
				if ev, ok := ev.(key.Event); ok && ev.State == key.Press {
					if g.state == guiStateError {
						g.state = guiStateToken
						g.window.Invalidate()
					}
				}
			}

			// Draw current state
			g.layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}

func (g *guiApp) layout(gtx layout.Context) layout.Dimensions {
	// Fill background with white
	paint.Fill(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255})

	switch g.state {
	case guiStateToken:
		return g.layoutTokenEntry(gtx)
	case guiStateBudgetSelect:
		return g.layoutBudgetSelection(gtx)
	case guiStateExporting:
		return g.layoutExporting(gtx)
	case guiStateDone:
		return g.layoutComplete(gtx)
	case guiStateError:
		return g.layoutError(gtx)
	default:
		return layout.Dimensions{}
	}
}

func (g *guiApp) layoutTokenEntry(gtx layout.Context) layout.Dimensions {
	inset := layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20), Left: unit.Dp(20), Right: unit.Dp(20)}

	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceAround}.Layout(gtx,
			layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H4(g.theme, "YNAB Budget Exporter")
				title.Alignment = text.Middle
				title.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
				return layout.Center.Layout(gtx, title.Layout)
			}),
			layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				body := material.Body1(g.theme, "Export your YNAB budget for\nimport into Actual Budget.")
				return body.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return g.drawSeparator(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				instructions := `To get your API token:
1. Sign in to YNAB web app
2. Go to Account Settings
3. Then Developer Settings
4. Click 'New' under Personal
   Access Tokens
5. Enter password, then Generate
6. Copy the FULL token from the
   top of the page`
				body := material.Body2(g.theme, instructions)
				return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, body.Layout)
			}),
			layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(g.theme, &g.tokenEditor, "Enter your YNAB API token...")
				return editor.Layout(gtx)
			}),
			layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// Handle button click
				for g.continueBtn.Clicked(gtx) {
					g.handleTokenSubmit()
				}

				// Check for Enter key in editor
				for {
					_, ok := g.tokenEditor.Update(gtx)
					if !ok {
						break
					}
					g.handleTokenSubmit()
				}

				btn := material.Button(g.theme, &g.continueBtn, "Continue")
				return layout.Center.Layout(gtx, btn.Layout)
			}),
			layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{}.Layout(gtx)
			}),
		)
	})
}

func (g *guiApp) handleTokenSubmit() {
	token := strings.TrimSpace(g.tokenEditor.Text())
	if token == "" {
		g.showError(errors.New("please enter your API token"))
		return
	}

	if len(token) != ynabTokenLength {
		g.showError(fmt.Errorf("token should be %d characters long", ynabTokenLength))
		return
	}

	g.token = token
	g.state = guiStateExporting
	g.window.Invalidate()

	go func() {
		if err := validateToken(token); err != nil {
			g.showError(fmt.Errorf("token validation failed: %w", err))
			return
		}
		g.fetchBudgets()
	}()
}

func (g *guiApp) fetchBudgets() {
	budgets, err := fetchBudgetsSync(g.token)
	if err != nil {
		g.showError(fmt.Errorf("failed to fetch budgets: %w", err))
		return
	}

	g.budgets = budgets
	g.budgetClickables = make([]widget.Clickable, len(budgets))
	g.state = guiStateBudgetSelect
	g.window.Invalidate()
}

//nolint:gocognit // Complex layout function
func (g *guiApp) layoutBudgetSelection(gtx layout.Context) layout.Dimensions {
	inset := layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20), Left: unit.Dp(20), Right: unit.Dp(20)}

	// Ensure clickables are initialized
	if len(g.budgetClickables) != len(g.budgets) {
		g.budgetClickables = make([]widget.Clickable, len(g.budgets))
	}

	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H5(g.theme, "Select a Budget to Export")
				title.Alignment = text.Middle
				title.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
				return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, title.Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return g.drawSeparator(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				// Show message if no budgets
				if len(g.budgets) == 0 {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(g.theme, "No budgets found")
						return label.Layout(gtx)
					})
				}

				// Render list - NOTE: Known Gio rendering issue causes text to grey after ~33 chars
				return material.List(g.theme, &g.budgetList).Layout(gtx, len(g.budgets), func(gtx layout.Context, index int) layout.Dimensions {
					selectedBudget := g.budgets[index]
					isSelected := index == g.selectedIndex

					// Handle click
					if g.budgetClickables[index].Clicked(gtx) {
						g.selectedIndex = index
					}

					// Use material.Clickable for interaction
					return material.Clickable(gtx, &g.budgetClickables[index], func(gtx layout.Context) layout.Dimensions {
						// Add padding
						return layout.Inset{
							Top:    unit.Dp(12),
							Bottom: unit.Dp(12),
							Left:   unit.Dp(16),
							Right:  unit.Dp(16),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							// Truncate budget name to 32 chars to avoid Gio rendering issue
							budgetName := selectedBudget.Name
							if len(budgetName) > 32 {
								budgetName = budgetName[:32] + "..."
							}

							// Format last modified date concisely
							lastModifiedText := "Last modified: unknown"
							if selectedBudget.LastModifiedOn != "" {
								if t, err := time.Parse(time.RFC3339, selectedBudget.LastModifiedOn); err == nil {
									lastModifiedText = "Last modified: " + t.Format(time.DateOnly)
								}
							}

							// Render name and last modified on separate lines
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									label := material.Body1(g.theme, budgetName)
									if isSelected {
										label.Font.Weight = font.Bold
									}
									return label.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lastModified := material.Caption(g.theme, lastModifiedText)
									return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, lastModified.Layout)
								}),
							)
						})
					})
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return g.drawSeparator(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEnd}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if g.backBtn.Clicked(gtx) {
								g.state = guiStateToken
								g.tokenEditor.SetText("")
							}
							return material.Button(g.theme, &g.backBtn, "Back").Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if g.exportBtn.Clicked(gtx) {
								if g.selectedIndex >= 0 {
									g.selectedBudget = g.budgets[g.selectedIndex]
									g.exportBudget()
								} else {
									g.showError(errors.New("please select a budget first"))
								}
							}
							return material.Button(g.theme, &g.exportBtn, "Export Selected Budget").Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}

func (g *guiApp) exportBudget() {
	g.state = guiStateExporting
	g.window.Invalidate()

	go func() {
		summary, path, jsonData, err := exportBudgetSync(g.token, g.selectedBudget.ID, g.selectedBudget.Name)
		if err != nil {
			g.showError(fmt.Errorf("export failed: %w", err))
			return
		}

		g.summary = summary
		g.exportPath = path
		g.jsonData = jsonData
		g.state = guiStateDone
		g.window.Invalidate()
	}()
}

func (g *guiApp) layoutExporting(gtx layout.Context) layout.Dimensions {
	inset := layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20), Left: unit.Dp(20), Right: unit.Dp(20)}

	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceAround}.Layout(gtx,
			layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				var message string
				if g.selectedBudget.Name != "" {
					message = "Exporting budget: " + g.selectedBudget.Name
				} else {
					message = "Validating token..."
				}
				label := material.H6(g.theme, message)
				label.Alignment = text.Middle
				return layout.Center.Layout(gtx, label.Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(g.theme, "Please wait...")
				label.Alignment = text.Middle
				return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, label.Layout)
				})
			}),
			layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{}.Layout(gtx)
			}),
		)
	})
}

func (g *guiApp) layoutComplete(gtx layout.Context) layout.Dimensions {
	inset := layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20), Left: unit.Dp(20), Right: unit.Dp(20)}

	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H5(g.theme, "✓ Export Complete!")
				title.Alignment = text.Middle
				title.Color = color.NRGBA{R: 0, G: 200, B: 0, A: 255}
				return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, title.Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// Truncate budget name if needed
				budgetName := g.selectedBudget.Name
				if len(budgetName) > 25 {
					budgetName = budgetName[:25] + "..."
				}
				return material.Body1(g.theme, "Budget: "+budgetName).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					// Show just the filename, not full path
					filename := filepath.Base(g.exportPath)
					return material.Body2(g.theme, "Saved to Downloads folder: "+filename).Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, g.drawSeparator)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				summaryText := g.formatBudgetSummary()
				return material.Body2(g.theme, summaryText).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, g.drawSeparator)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				instructions := `Next Steps:
1. Open Actual Budget
2. Close any open budget file
3. Select 'Import file'
4. Choose 'nYNAB'
5. Select the exported JSON
6. Review budget and follow cleanup steps at:
   actualbudget.org/docs/migration/nynab#cleanup`
				return material.Body2(g.theme, instructions).Layout(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEnd}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if g.exportAnotherBtn.Clicked(gtx) {
								g.state = guiStateBudgetSelect
								g.selectedIndex = -1
							}
							return material.Button(g.theme, &g.exportAnotherBtn, "Export Another Budget").Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if g.closeBtn.Clicked(gtx) {
								os.Exit(0)
							}
							return material.Button(g.theme, &g.closeBtn, "Close").Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}

func (g *guiApp) layoutError(gtx layout.Context) layout.Dimensions {
	inset := layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20), Left: unit.Dp(20), Right: unit.Dp(20)}

	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceAround}.Layout(gtx,
			layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H5(g.theme, "Error")
				title.Alignment = text.Middle
				title.Color = color.NRGBA{R: 200, G: 0, B: 0, A: 255}
				return layout.Center.Layout(gtx, title.Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(g.theme, g.errorMsg)
				label.Alignment = text.Middle
				return layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, label.Layout)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if g.errorOkBtn.Clicked(gtx) {
					g.state = guiStateToken
				}
				btn := material.Button(g.theme, &g.errorOkBtn, "OK")
				return layout.Center.Layout(gtx, btn.Layout)
			}),
			layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{}.Layout(gtx)
			}),
		)
	})
}

func (g *guiApp) drawSeparator(gtx layout.Context) layout.Dimensions {
	line := widget.Border{
		Color: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		Width: unit.Dp(1),
	}
	return line.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{
			Size: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Dp(unit.Dp(1))},
		}
	})
}

func (g *guiApp) showError(err error) {
	g.errorMsg = err.Error()
	g.state = guiStateError
	g.window.Invalidate()
}

func (g *guiApp) formatBudgetSummary() string {
	// Format with shorter labels to fit
	var result strings.Builder
	result.WriteString("Budget Summary:\n\n")
	result.WriteString(fmt.Sprintf("Currency: %s\n", g.summary.Currency))
	result.WriteString(fmt.Sprintf("Accounts: %d", g.summary.AccountCount))
	if g.summary.ClosedAccountCount > 0 {
		result.WriteString(fmt.Sprintf("\n  (+%d closed)", g.summary.ClosedAccountCount))
	}
	result.WriteString(fmt.Sprintf("\nCategories: %d", g.summary.CategoryCount))
	var catDetails []string
	if g.summary.HiddenCategoryCount > 0 {
		catDetails = append(catDetails, fmt.Sprintf("%d hidden", g.summary.HiddenCategoryCount))
	}
	if g.summary.DeletedCategoryCount > 0 {
		catDetails = append(catDetails, fmt.Sprintf("%d deleted", g.summary.DeletedCategoryCount))
	}
	if len(catDetails) > 0 {
		result.WriteString(fmt.Sprintf("\n  (+%s)", strings.Join(catDetails, ", ")))
	}
	result.WriteString(fmt.Sprintf("\nPayees: %d", g.summary.PayeeCount))
	result.WriteString(fmt.Sprintf("\nTransactions: %d", g.summary.TransactionCount))
	result.WriteString("\nFrom: " + formatMonthYear(g.summary.FirstMonth))
	result.WriteString("\nTo: " + formatMonthYear(g.summary.LastMonth))
	result.WriteString("\nSize: " + humanizeFileSize(g.summary.FileSize))
	return result.String()
}
