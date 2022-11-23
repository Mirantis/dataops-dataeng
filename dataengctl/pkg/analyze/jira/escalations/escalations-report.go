package escalations

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"
// )

// func escalationsReport() {
// 	history, err := NewEscalationsHistory()
// 	if err != nil {
// 		log.Fatalf("failed getting escalations history: %v", err)
// 	}

// 	fmt.Printf("Quarter,Carry In,Opened,Closed,Escalated,Closed(all time),Avg Hours to Close,Open(SUS),Closed(SUS),Open(ENG),Closed(ENG),Opened(Crit),Closed(Crit),Triage Hours(SUS),Turnaround Hours(SUS),Turnaround Hours(Crit),Escalate Hours(SUS),Triage Hours(ENG),Turnaround Hours(ENG),Hours Pending Sus,Hours Pending Support,Hours Pending Engineering,Hours Pending 3rd Party, Hours Pending PR Review,\n")
// 	for _, quarter := range quarters {
// 		info, err := history.InfoForQuarter(quarter)
// 		if err != nil {
// 			log.Fatalf("failed getting history for quarter: %v", err)
// 		}
// 		fmt.Println(info.CSVString())
// 	}
// }

// type EscalationsHistory struct {
// 	Overall         *IssueHistory
// 	ByOpenQuarter   map[string]*IssueHistory
// 	ByClosedQuarter map[string]*IssueHistory
// }

// func NewEscalationsHistory() (*EscalationsHistory, error) {
// 	eh := &EscalationsHistory{
// 		ByOpenQuarter:   make(map[string]*IssueHistory),
// 		ByClosedQuarter: make(map[string]*IssueHistory),
// 	}

// 	client, err := getClient()
// 	if err != nil {
// 		return nil, err
// 	}

// 	query := `project = 'FIELD' AND "Component" in ("Mirantis Secure Registry", "Mirantis Secure Registry (DTR)", "Mirantis Kubernetes Engine (UCP)", "Mirantis Kubernetes Engine", "Mirantis Container Runtime (Engine)", "Mirantis Container Runtime")`

// 	eh.Overall, err = NewIssueHistory(context.Background(), IssueHistoryOptions{
// 		Client: client,
// 		Query:  query,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// need data by quarter
// 	for _, quarter := range quarters {
// 		eh.ByOpenQuarter[quarter.name] = eh.Overall.NewSubHistory(func(_ *Issue, statuses []IssueStatus) bool {
// 			return statuses[0].TimeStart.After(quarter.start) && statuses[0].TimeStart.Before(quarter.end)
// 		})

// 		eh.ByClosedQuarter[quarter.name] = eh.Overall.NewSubHistory(func(_ *Issue, statuses []IssueStatus) bool {
// 			for _, status := range statuses {
// 				if terminalStatus(status.Status) {
// 					return status.TimeStart.After(quarter.start) && status.TimeStart.Before(quarter.end)
// 				}
// 			}
// 			return false
// 		})
// 	}

// 	return eh, nil
// }

// type EscalationsInfo struct {
// 	Quarter                    string
// 	CarryIn                    int
// 	Opened                     int
// 	Closed                     int
// 	Escalated                  int
// 	ClosedAllTime              int
// 	AvgHoursToClose            float64
// 	OpenSus                    int
// 	ClosedSus                  int
// 	OpenEng                    int
// 	ClosedEng                  int
// 	OpenCrit                   int
// 	ClosedCrit                 int
// 	TriageHoursSus             float64
// 	TurnaroundHoursSus         float64
// 	TurnaroundHoursCrit        float64
// 	EscalateHoursSus           float64
// 	TriageHoursEng             float64
// 	TurnaroundHoursEng         float64
// 	AvgHoursPendingSus         float64
// 	AvgHoursPendingSupport     float64
// 	AvgHoursPendingEngineering float64
// 	AvgHoursPending3rdParty    float64
// 	AvgHoursPendingPRReview    float64
// }

// func (e EscalationsInfo) CSVString() string {
// 	return fmt.Sprintf(
// 		"%s,%d,%d,%d,%d,%d,%f,%d,%d,%d,%d,%d,%d,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,",
// 		e.Quarter,
// 		e.CarryIn,
// 		e.Opened,
// 		e.Closed,
// 		e.Escalated,
// 		e.ClosedAllTime,
// 		e.AvgHoursToClose,
// 		e.OpenSus,
// 		e.ClosedSus,
// 		e.OpenEng,
// 		e.ClosedEng,
// 		e.OpenCrit,
// 		e.ClosedCrit,
// 		e.TriageHoursSus,
// 		e.TurnaroundHoursSus,
// 		e.TurnaroundHoursCrit,
// 		e.EscalateHoursSus,
// 		e.TriageHoursEng,
// 		e.TurnaroundHoursEng,
// 		e.AvgHoursPendingSus,
// 		e.AvgHoursPendingSupport,
// 		e.AvgHoursPendingEngineering,
// 		e.AvgHoursPending3rdParty,
// 		e.AvgHoursPendingPRReview,
// 	)
// }

// func (h *EscalationsHistory) InfoForQuarter(quarter Quarter) (EscalationsInfo, error) {
// 	result := EscalationsInfo{}
// 	history, ok := h.ByOpenQuarter[quarter.name]
// 	if !ok {
// 		return result, fmt.Errorf("unrecognized quater %v", quarter.name)
// 	}
// 	closedHistory, ok := h.ByClosedQuarter[quarter.name]
// 	if !ok {
// 		return result, fmt.Errorf("unrecognized quarter %v", quarter.name)
// 	}

// 	carriedIn := h.Overall.NewSubHistory(func(iss *Issue, statuses []IssueStatus) bool {
// 		if time.Time(iss.Fields.Created).After(quarter.start) {
// 			return false
// 		}

// 		for _, status := range statuses {
// 			if terminalStatus(status.Status) {
// 				if status.TimeStart.Before(quarter.start) {
// 					return false
// 				}
// 				return true
// 			}
// 		}
// 		return true
// 	})

// 	result.Quarter = quarter.name
// 	result.CarryIn = carriedIn.Count()
// 	result.Opened = history.Count()
// 	result.Closed = closedHistory.Count()
// 	result.TriageHoursSus = history.AverageTimeInStatus(StatusSustainingTriage).Hours()
// 	result.TriageHoursEng = history.AverageTimeInStatus(StatusDevTriage).Hours()

// 	closed, avgTimeToClose := history.AverageTimeToClose()
// 	result.ClosedAllTime = closed
// 	result.AvgHoursToClose = avgTimeToClose.Hours()

// 	criticalIssues := history.NewSubHistory(func(iss *Issue, _ []IssueStatus) bool {
// 		for _, label := range iss.Fields.Labels {
// 			if label == LabelCriticalEscalation {
// 				return true
// 			}
// 		}
// 		return false
// 	})
// 	result.OpenCrit = criticalIssues.Count()
// 	closed, avg := criticalIssues.AverageTimeToClose()
// 	result.ClosedCrit = closed
// 	result.TurnaroundHoursCrit = avg.Hours()

// 	ss := history.AverageTimeInSustainingStatuses()
// 	result.AvgHoursPendingSus = ss[SustainingStatusPendingSus].Hours()
// 	result.AvgHoursPendingSupport = ss[SustainingStatusPendingSupport].Hours()
// 	result.AvgHoursPendingEngineering = ss[SustainingStatusPendingEngineering].Hours()
// 	result.AvgHoursPending3rdParty = ss[SustainingStatusPending3rdParty].Hours()
// 	result.AvgHoursPendingPRReview = ss[SustainingStatusPendingPRReview].Hours()

// 	teamHistory := history.TeamHistory()
// 	for _, th := range teamHistory {
// 		if th.EscalatedToEngineering {
// 			result.Escalated += 1
// 			if th.Closed {
// 				result.ClosedEng += 1
// 			} else {
// 				result.OpenEng += 1
// 			}
// 		} else {
// 			if th.Closed {
// 				result.ClosedSus += 1
// 			} else {
// 				result.OpenSus += 1
// 			}
// 		}
// 	}

// 	for _, th := range teamHistory {
// 		if th.EscalatedToEngineering {
// 			escalateTime := th.TimeEscalated.Sub(th.TimeOpen).Hours()
// 			result.EscalateHoursSus += escalateTime / float64(result.Escalated)
// 		}

// 		if th.Closed {
// 			if th.EscalatedToEngineering {
// 				engTurnaround := th.TimeClose.Sub(th.TimeEscalated).Hours()
// 				result.TurnaroundHoursEng += engTurnaround / float64(result.ClosedEng)
// 			} else {
// 				susTurnaround := th.TimeClose.Sub(th.TimeOpen).Hours()
// 				result.TurnaroundHoursSus += susTurnaround / float64(result.ClosedSus)
// 			}
// 		}
// 	}
// 	return result, nil
// }

// func escalationsReportOld() {
// 	client, err := getClient()
// 	if err != nil {
// 		log.Fatalf("error logging in: %v\n", err)
// 	}

// 	queryTpl := `project = 'FIELD' AND createdDate >= %s AND createdDate <= %s AND "Component" in ("Mirantis Secure Registry", "Mirantis Secure Registry (DTR)", "Mirantis Kubernetes Engine (UCP)", "Mirantis Kubernetes Engine", "Mirantis Container Runtime (Engine)", "Mirantis Container Runtime")`

// 	fmt.Printf("Quarter,Opened,Closed,Escalated,Closed(all time),Avg Hours to Close,Open(SUS),Closed(SUS),Open(ENG)Closed(ENG),Triage Hours(SUS),Turnaround Hours(SUS),Escalate Hours(SUS),Triage Hours(ENG),Turnaround Hours(ENG),\n")
// 	for _, quarter := range quarters {
// 		query := fmt.Sprintf(queryTpl, dateString(quarter.start), dateString(quarter.end))
// 		history, err := NewIssueHistory(context.Background(), IssueHistoryOptions{
// 			Client: client,
// 			Query:  query,
// 		})
// 		if err != nil {
// 			log.Fatalf("error getting issue history: %v\n", err)
// 		}

// 		result := EscalationsInfo{
// 			Quarter:        quarter.name,
// 			Opened:         history.Count(),
// 			TriageHoursSus: history.AverageTimeInStatus(StatusSustainingTriage).Hours(),
// 			TriageHoursEng: history.AverageTimeInStatus(StatusDevTriage).Hours(),
// 		}

// 		closed, avgTimeToClose := history.AverageTimeToClose()
// 		result.ClosedAllTime = closed
// 		result.AvgHoursToClose = avgTimeToClose.Hours()

// 		teamHistory := history.TeamHistory()
// 		for _, th := range teamHistory {
// 			if th.EscalatedToEngineering {
// 				result.Escalated += 1
// 				if th.Closed {
// 					result.ClosedEng += 1
// 				} else {
// 					result.OpenEng += 1
// 				}
// 			} else {
// 				if th.Closed {
// 					result.ClosedSus += 1
// 				} else {
// 					result.OpenSus += 1
// 				}
// 			}
// 		}

// 		for _, th := range teamHistory {
// 			if th.EscalatedToEngineering {
// 				escalateTime := th.TimeEscalated.Sub(th.TimeOpen).Hours()
// 				result.EscalateHoursSus += escalateTime / float64(result.Escalated)
// 			}

// 			if th.Closed {
// 				if th.EscalatedToEngineering {
// 					engTurnaround := th.TimeClose.Sub(th.TimeEscalated).Hours()
// 					result.TurnaroundHoursEng += engTurnaround / float64(result.ClosedEng)
// 				} else {
// 					susTurnaround := th.TimeClose.Sub(th.TimeOpen).Hours()
// 					result.TurnaroundHoursSus += susTurnaround / float64(result.ClosedSus)
// 				}
// 			}
// 		}

// 		// Partition by team (SUS or ENG)

// 		// TODO: compute result

// 		fmt.Println(result.CSVString())
// 	}
// }