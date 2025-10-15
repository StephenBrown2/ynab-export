package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type guiState int

const (
	guiStateToken guiState = iota
	guiStateBudgetSelect
	guiStateExporting
	guiStateDone
)

type guiApp struct {
	app            fyne.App
	window         fyne.Window
	selectedBudget budget
	token          string
	exportPath     string
	budgets        []budget
	jsonData       []byte
	summary        budgetSummary
	state          guiState
}

func runGUI() {
	a := app.NewWithID("com.stephenbrown2.ynab-export")
	a.Settings().SetTheme(theme.DefaultTheme())

	w := a.NewWindow("YNAB Export Tool")
	w.Resize(fyne.NewSize(600, 500))
	w.CenterOnScreen()

	gui := &guiApp{
		app:    a,
		window: w,
		state:  guiStateToken,
	}

	// Check for token in environment variable
	if token := os.Getenv("YNAB_API_TOKEN"); token != "" {
		gui.token = token
		// Validate token
		if err := validateToken(token); err == nil {
			// Token is valid, fetch budgets
			go gui.fetchBudgetsAsync()
		} else {
			// Token invalid, show token entry
			gui.showTokenEntry()
		}
	} else {
		gui.showTokenEntry()
	}

	w.ShowAndRun()
}

func (g *guiApp) showTokenEntry() {
	g.state = guiStateToken

	title := widget.NewLabelWithStyle("YNAB Budget Exporter", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	intro := widget.NewLabel("This tool will help you export your YNAB budget for import into Actual Budget.")
	intro.Wrapping = fyne.TextWrapWord

	instructions := widget.NewLabel(`To get your API token:
1. Sign in to the YNAB web app
2. Go to Account Settings → Developer Settings
3. Under 'Personal Access Tokens', click 'New'
4. Enter your password and click 'Generate'
5. Copy the FULL token (not the obfuscated version)`)
	instructions.Wrapping = fyne.TextWrapWord

	tokenEntry := widget.NewPasswordEntry()
	tokenEntry.SetPlaceHolder("Enter your YNAB API token...")
	tokenEntry.OnChanged = func(_ string) {
		// Token validation feedback could go here
	}

	var progressDialog dialog.Dialog

	continueBtn := widget.NewButton("Continue", func() {
		token := strings.TrimSpace(tokenEntry.Text)
		if token == "" {
			dialog.ShowError(errors.New("please enter your API token"), g.window)
			return
		}

		if len(token) != ynabTokenLength {
			dialog.ShowError(fmt.Errorf("token should be %d characters long", ynabTokenLength), g.window)
			return
		}

		g.token = token

		// Show progress dialog
		progressBar := widget.NewProgressBarInfinite()
		progressContent := container.NewVBox(
			widget.NewLabel("Validating token..."),
			progressBar,
		)
		progressDialog = dialog.NewCustomWithoutButtons("Please wait", progressContent, g.window)
		progressDialog.Show()

		// Validate token and fetch budgets in background
		go func() {
			if err := validateToken(token); err != nil {
				progressDialog.Hide()
				dialog.ShowError(fmt.Errorf("token validation failed: %w", err), g.window)
				return
			}

			g.fetchBudgetsAsync()
			progressDialog.Hide()
		}()
	})
	continueBtn.Importance = widget.HighImportance

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(title),
		layout.NewSpacer(),
		intro,
		widget.NewSeparator(),
		instructions,
		layout.NewSpacer(),
		tokenEntry,
		layout.NewSpacer(),
		container.NewCenter(continueBtn),
		layout.NewSpacer(),
	)

	g.window.SetContent(content)
}

func (g *guiApp) fetchBudgetsAsync() {
	// Show loading state
	g.app.SendNotification(&fyne.Notification{
		Title:   "Fetching budgets",
		Content: "Retrieving your budgets from YNAB...",
	})

	// Fetch budgets in background
	go func() {
		budgets, err := fetchBudgetsSync(g.token)

		// Update UI in main thread
		fyne.Do(func() {
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to fetch budgets: %w", err), g.window)
				return
			}

			g.budgets = budgets
			g.showBudgetSelection()
		})
	}()
}

func (g *guiApp) showBudgetSelection() {
	g.state = guiStateBudgetSelect

	title := widget.NewLabelWithStyle("Select a Budget to Export", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Create list of budgets
	budgetNames := make([]string, 0, len(g.budgets))
	for _, b := range g.budgets {
		displayName := b.Title()
		if b.LastModifiedOn != "" {
			if t, err := time.Parse(time.RFC3339, b.LastModifiedOn); err == nil {
				displayName = fmt.Sprintf("%s\nLast Modified: %s", b.Name, t.Format("2006-01-02"))
			}
		}
		budgetNames = append(budgetNames, displayName)
	}

	selectedIndex := -1

	list := widget.NewList(
		func() int {
			return len(budgetNames)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Budget Name"),
				widget.NewLabelWithStyle("ID", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box, ok := obj.(*fyne.Container)
			if !ok {
				return
			}
			nameLabel, ok := box.Objects[0].(*widget.Label)
			if !ok {
				return
			}
			idLabel, ok := box.Objects[1].(*widget.Label)
			if !ok {
				return
			}

			nameLabel.SetText(g.budgets[id].Title())
			idLabel.SetText(g.budgets[id].Description())
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		selectedIndex = id
	}

	exportBtn := widget.NewButton("Export Selected Budget", func() {
		if selectedIndex < 0 {
			dialog.ShowInformation("No selection", "Please select a budget first", g.window)
			return
		}

		g.selectedBudget = g.budgets[selectedIndex]
		g.exportBudgetAsync()
	})
	exportBtn.Importance = widget.HighImportance

	backBtn := widget.NewButton("Back", func() {
		g.showTokenEntry()
	})

	content := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()),
		container.NewVBox(
			widget.NewSeparator(),
			container.NewHBox(
				layout.NewSpacer(),
				backBtn,
				exportBtn,
			),
		),
		nil,
		nil,
		list,
	)

	g.window.SetContent(content)
}

func (g *guiApp) exportBudgetAsync() {
	g.state = guiStateExporting

	// Show progress dialog
	progressBar := widget.NewProgressBarInfinite()
	progressContent := container.NewVBox(
		widget.NewLabel("Exporting budget: "+g.selectedBudget.Name),
		widget.NewLabel("Please wait..."),
		progressBar,
	)
	progressDialog := dialog.NewCustomWithoutButtons("Exporting", progressContent, g.window)
	progressDialog.Show()

	// Export in background
	go func() {
		summary, path, jsonData, err := exportBudgetSync(g.token, g.selectedBudget.ID, g.selectedBudget.Name)

		// Update UI in main thread
		fyne.Do(func() {
			progressDialog.Hide()

			if err != nil {
				dialog.ShowError(fmt.Errorf("export failed: %w", err), g.window)
				return
			}

			g.summary = summary
			g.exportPath = path
			g.jsonData = jsonData
			g.showComplete()
		})
	}()
}

func (g *guiApp) showComplete() {
	g.state = guiStateDone

	title := widget.NewLabelWithStyle("✓ Export Complete!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	title.Importance = widget.SuccessImportance

	budgetInfo := widget.NewLabel("Budget: " + g.selectedBudget.Name)
	pathInfo := widget.NewLabel("Saved to: " + g.exportPath)
	pathInfo.Wrapping = fyne.TextWrapWord

	summaryTitle := widget.NewLabelWithStyle("Budget Summary:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Build summary text
	summaryText := fmt.Sprintf(`Currency:     %s
Accounts:     %d%s
Categories:   %d%s
Payees:       %d
Transactions: %d
Date Range:   %s to %s
File Size:    %s`,
		g.summary.Currency,
		g.summary.AccountCount,
		func() string {
			if g.summary.ClosedAccountCount > 0 {
				return fmt.Sprintf(" (plus %d closed)", g.summary.ClosedAccountCount)
			}
			return ""
		}(),
		g.summary.CategoryCount,
		func() string {
			var details []string
			if g.summary.HiddenCategoryCount > 0 {
				details = append(details, fmt.Sprintf("%d hidden", g.summary.HiddenCategoryCount))
			}
			if g.summary.DeletedCategoryCount > 0 {
				details = append(details, fmt.Sprintf("%d deleted", g.summary.DeletedCategoryCount))
			}
			if len(details) > 0 {
				return fmt.Sprintf(" (plus %s)", strings.Join(details, ", "))
			}
			return ""
		}(),
		g.summary.PayeeCount,
		g.summary.TransactionCount,
		formatMonthYear(g.summary.FirstMonth),
		formatMonthYear(g.summary.LastMonth),
		humanizeFileSize(g.summary.FileSize),
	)

	summary := widget.NewLabel(summaryText)

	instructionsTitle := widget.NewLabel("Next Steps:")
	instructionsTitle.TextStyle = fyne.TextStyle{Bold: true}

	instructions := widget.NewLabel(`1. Open Actual Budget
2. If a budget is already open, select the dropdown menu and 'Close File'
3. Select 'Import file'
4. Choose 'nYNAB'
5. Select the exported JSON file
6. Once imported, review your budget and follow cleanup steps at
   actualbudget.org/docs/migration/nynab#cleanup`)
	instructions.Wrapping = fyne.TextWrapWord

	// Display JSON preview (limit to first 50KB to avoid freezing)
	const maxPreviewSize = 50 * 1024 // 50KB
	jsonPreview := string(g.jsonData)
	var truncatedMsg string
	if len(g.jsonData) > maxPreviewSize {
		jsonPreview = string(g.jsonData[:maxPreviewSize])
		truncatedMsg = fmt.Sprintf("\n\n... (showing first 50KB of %s file)", humanizeFileSize(int64(len(g.jsonData))))
	}

	// Use RichText with monospace for better JSON display
	jsonText := widget.NewRichTextFromMarkdown("```json\n" + jsonPreview + truncatedMsg + "\n```")
	jsonText.Wrapping = fyne.TextWrapOff

	jsonScroll := container.NewScroll(jsonText)
	jsonScroll.SetMinSize(fyne.NewSize(550, 200))

	jsonAccordion := widget.NewAccordion(
		widget.NewAccordionItem("View Exported JSON (Preview)", jsonScroll),
	)

	closeBtn := widget.NewButton("Close", func() {
		g.app.Quit()
	})
	closeBtn.Importance = widget.HighImportance

	exportAnotherBtn := widget.NewButton("Export Another Budget", func() {
		g.showBudgetSelection()
	})

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(title),
		layout.NewSpacer(),
		budgetInfo,
		pathInfo,
		widget.NewSeparator(),
		summaryTitle,
		summary,
		widget.NewSeparator(),
		jsonAccordion,
		widget.NewSeparator(),
		instructionsTitle,
		instructions,
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			exportAnotherBtn,
			closeBtn,
		),
	)

	scrollContainer := container.NewScroll(content)
	g.window.SetContent(scrollContainer)
}
