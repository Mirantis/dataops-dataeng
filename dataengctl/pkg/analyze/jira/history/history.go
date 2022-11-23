package history

// import (
// 	"context"
// 	"sort"
// 	"time"

// 	"github.com/andygrunwald/go-jira"
// )

// // map: issue key |-> relevant status transitions
// type IssueHistoryDetail map[*Issue][]IssueStatus

// type IssueStatus struct {
// 	TimeStart        time.Time
// 	TimeEnd          time.Time
// 	Duration         time.Duration
// 	Status           string
// 	Owner            string
// 	SustainingStatus string
// }

// type IssueTeamHistory struct {
// 	Closed                 bool
// 	EscalatedToEngineering bool
// 	TimeOpen               time.Time
// 	TimeClose              time.Time
// 	TimeEscalated          time.Time
// }

// type IssueCloseStatus struct {
// 	*Issue
// 	Closed    bool
// 	TimeClose time.Time
// }

// type IssueHistory struct {
// 	it        *IssueIterator
// 	client    *jira.Client
// 	detail    IssueHistoryDetail
// 	populated bool
// }

// type IssueHistoryOptions struct {
// 	Query  string
// 	Client *jira.Client
// }

// type OpenCloseInfo struct {
// 	TimeOpen  time.Time
// 	TimeClose time.Time
// }

// func (oc *OpenCloseInfo) Duration() time.Duration {
// 	return oc.TimeClose.Sub(oc.TimeOpen)
// }

// func NewIssueHistory(ctx context.Context, opts IssueHistoryOptions) (*IssueHistory, error) {
// 	var err error
// 	h := &IssueHistory{
// 		client: opts.Client,
// 	}

// 	h.it, err = NewIssueIterator(ctx, IssueIteratorOptions{
// 		Client: opts.Client,
// 		Query:  opts.Query,
// 		Expand: "changelog",
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return h, h.populate(ctx)
// }

// func (h *IssueHistory) OpenTickets() *IssueHistory {
// 	return h.NewSubHistory(func(_ *Issue, statuses []IssueStatus) bool {
// 		for _, status := range statuses {
// 			if terminalStatus(status.Status) {
// 				return true
// 			}
// 		}
// 		return false
// 	})
// }

// func (h *IssueHistory) ClosedTickets() *IssueHistory {
// 	return h.NewSubHistory(func(_ *Issue, statuses []IssueStatus) bool {
// 		for _, status := range statuses {
// 			if terminalStatus(status.Status) {
// 				return false
// 			}
// 		}
// 		return true
// 	})
// }

// func (h *IssueHistory) NewSubHistory(f func(*Issue, []IssueStatus) bool) *IssueHistory {
// 	history := &IssueHistory{
// 		it:        nil,
// 		client:    nil,
// 		populated: true,
// 		detail:    make(map[*Issue][]IssueStatus),
// 	}

// 	for issue, statuses := range h.detail {
// 		if f(issue, statuses) {
// 			history.detail[issue] = statuses
// 		}
// 	}
// 	return history
// }

// func (h *IssueHistory) Count() int {
// 	return len(h.detail)
// }

// func (h *IssueHistory) Raw() IssueHistoryDetail {
// 	return h.detail
// }

// func (h *IssueHistory) ByCustomer() map[Customer][]*jira.Issue {
// 	result := make(map[Customer][]*jira.Issue)
// 	for issue, _ := range h.detail {
// 		for _, customer := range issue.Customers {
// 			if _, ok := result[customer]; !ok {
// 				result[customer] = make([]*jira.Issue, 0)
// 			}
// 			result[customer] = append(result[customer], &issue.Issue)
// 		}
// 	}
// 	return result
// }

// // TODO: make this map[Customer]IssueHistoryDetail?
// func (h *IssueHistory) CustomerHistory() map[Customer][]IssueCloseStatus {
// 	th := h.TeamHistory()
// 	result := make(map[Customer][]IssueCloseStatus)
// 	for issue, _ := range h.detail {
// 		for _, customer := range issue.Customers {
// 			if _, ok := result[customer]; !ok {
// 				result[customer] = make([]IssueCloseStatus, 0)
// 			}
// 			result[customer] = append(result[customer], IssueCloseStatus{
// 				Issue:     issue,
// 				Closed:    th[issue.Key].Closed,
// 				TimeClose: th[issue.Key].TimeClose,
// 			})
// 		}
// 	}
// 	return result
// }

// // Returns OpenCloseInfo only for issues which are closed.
// func (h *IssueHistory) OpenClose() map[*Issue]OpenCloseInfo {
// 	result := make(map[*Issue]OpenCloseInfo)
// 	for issue, statuses := range h.detail {
// 		for _, status := range statuses {
// 			if terminalStatus(status.Status) {
// 				result[issue] = OpenCloseInfo{
// 					TimeOpen:  time.Time(issue.Fields.Created),
// 					TimeClose: status.TimeStart,
// 				}
// 				break
// 			}
// 		}
// 	}
// 	return result
// }

// func (h *IssueHistory) TimeToCloseWithFilter(filter func(*Issue) bool) map[*Issue]time.Duration {
// 	result := make(map[*Issue]time.Duration)
// 	openClose := h.OpenClose()
// 	for issue, oc := range openClose {
// 		if filter(issue) {
// 			result[issue] = oc.Duration()
// 		}
// 	}
// 	return result
// }

// func (h *IssueHistory) AverageTimeToCloseWithFilter(filter func(*Issue) bool) (int, time.Duration) {
// 	var result float64
// 	closed := 0
// 	ttc := h.TimeToCloseWithFilter(filter)
// 	n := len(ttc)
// 	for _, duration := range ttc {
// 		result += float64(duration) / float64(n)
// 		closed += 1
// 	}
// 	return closed, time.Duration(int64(result))
// }

// func (h *IssueHistory) TimeToClose() map[*Issue]time.Duration {
// 	return h.TimeToCloseWithFilter(func(_ *Issue) bool { return true })
// }

// // Returns the average time from open to close, but ONLY FOR CLOSED TICKETS.
// func (h *IssueHistory) AverageTimeToClose() (int, time.Duration) {
// 	return h.AverageTimeToCloseWithFilter(func(_ *Issue) bool { return true })
// }

// func (h *IssueHistory) TeamHistory() map[string]IssueTeamHistory {
// 	result := make(map[string]IssueTeamHistory)
// 	for issue, statuses := range h.detail {
// 		history := IssueTeamHistory{
// 			TimeOpen: time.Time(issue.Fields.Created),
// 		}
// 		for _, status := range statuses {
// 			if terminalStatus(status.Status) {
// 				history.Closed = true
// 				history.TimeClose = status.TimeStart
// 				break
// 			} else if !history.EscalatedToEngineering && statusOwner(status.Status) == OwnerEngineering {
// 				history.EscalatedToEngineering = true
// 				history.TimeEscalated = status.TimeStart
// 			}
// 		}
// 		result[issue.Key] = history
// 	}
// 	return result
// }

// func (h *IssueHistory) Statuses() []string {
// 	result := make([]string, 0)
// 	seen := make(map[string]struct{})
// 	for _, statuses := range h.detail {
// 		for _, status := range statuses {
// 			if _, ok := seen[status.Status]; !ok {
// 				result = append(result, status.Status)
// 				seen[status.Status] = struct{}{}
// 			}
// 		}
// 	}
// 	return result
// }

// // Returns map: key |-> duration in status for key
// func (h *IssueHistory) TimeInStatus(targetStatus string) map[string]time.Duration {
// 	result := make(map[string]time.Duration)
// 	for issue, statuses := range h.detail {
// 		for _, status := range statuses {
// 			if status.Status == targetStatus {
// 				result[issue.Key] = result[issue.Key] + status.Duration
// 			}
// 		}
// 	}
// 	return result
// }

// // Returns map: key |-> map: status |-> duration for key in status
// func (h *IssueHistory) KeyTimeInStatus() map[string]map[string]time.Duration {
// 	result := make(map[string]map[string]time.Duration)
// 	for issue, statuses := range h.detail {
// 		r := make(map[string]time.Duration)
// 		for _, status := range statuses {
// 			r[status.Status] += r[status.Status] + status.Duration
// 		}
// 		result[issue.Key] = r
// 	}
// 	return result
// }

// // Returns map: status |-> map: key |-> duration in status for key
// func (h *IssueHistory) TimeInStatuses() map[string]map[string]time.Duration {
// 	result := make(map[string]map[string]time.Duration)
// 	for _, status := range h.Statuses() {
// 		result[status] = h.TimeInStatus(status)
// 	}
// 	return result
// }

// // Average time in targetStatus, rounded (or is it truncated? )to the minute
// func (h *IssueHistory) AverageTimeInStatus(targetStatus string) time.Duration {
// 	var (
// 		totalTime time.Duration
// 		count     int64
// 	)
// 	for _, duration := range h.TimeInStatus(targetStatus) {
// 		totalTime += duration
// 		count += 1
// 	}
// 	return time.Duration(int64(totalTime.Minutes())/count) * time.Minute
// }

// // Returns map: status |-> average duration in status
// func (h *IssueHistory) AverageTimeInStatuses() map[string]time.Duration {
// 	result := make(map[string]time.Duration)
// 	for _, status := range h.Statuses() {
// 		result[status] = h.AverageTimeInStatus(status)
// 	}
// 	return result
// }

// // Returns map: key |-> duration in sustaining status for key
// func (h *IssueHistory) TimeInSustainingStatus(targetStatus string) map[string]time.Duration {
// 	result := make(map[string]time.Duration)
// 	for issue, statuses := range h.detail {
// 		for _, status := range statuses {
// 			if status.SustainingStatus == targetStatus {
// 				result[issue.Key] = result[issue.Key] + status.Duration
// 			}
// 		}
// 	}
// 	return result
// }

// // Average time in targetStatus, rounded (or is it truncated? )to the minute
// func (h *IssueHistory) AverageTimeInSustainingStatus(targetStatus string) time.Duration {
// 	var (
// 		totalTime time.Duration
// 		count     int64
// 	)
// 	for _, duration := range h.TimeInSustainingStatus(targetStatus) {
// 		totalTime += duration
// 		count += 1
// 	}

// 	if count == 0 {
// 		return 0 * time.Minute
// 	}

// 	return time.Duration(int64(totalTime.Minutes())/count) * time.Minute
// }

// func (h *IssueHistory) AverageTimeInSustainingStatuses() map[string]time.Duration {
// 	statuses := []string{
// 		SustainingStatusPendingSus,
// 		SustainingStatusPendingSupport,
// 		SustainingStatusPendingEngineering,
// 		SustainingStatusPending3rdParty,
// 		SustainingStatusPendingPRReview,
// 	}
// 	result := make(map[string]time.Duration)
// 	for _, status := range statuses {
// 		result[status] = h.AverageTimeInSustainingStatus(status)
// 	}
// 	return result
// }

// // TODO: how to make this generic (could use an interface, but what
// //       about the IssueStatus type?)
// func (h *IssueHistory) populate(ctx context.Context) error {
// 	if h.populated {
// 		return nil
// 	}
// 	h.detail = make(map[*Issue][]IssueStatus)

// 	for issue := h.it.Next(ctx); issue != nil; issue = h.it.Next(ctx) {
// 		sort.SliceStable(issue.Changelog.Histories, func(i, j int) bool {
// 			iTime, err := issue.Changelog.Histories[i].CreatedTime()
// 			if err != nil {
// 				panic(err)
// 			}
// 			jTime, err := issue.Changelog.Histories[j].CreatedTime()
// 			if err != nil {
// 				panic(err)
// 			}
// 			return iTime.Before(jTime)
// 		})

// 		states := []IssueStatus{
// 			{
// 				TimeStart:        time.Time(issue.Fields.Created),
// 				Status:           InitialTicketStatus,
// 				SustainingStatus: InitialSustainingStatus,
// 				Owner:            statusOwner(InitialTicketStatus),
// 			},
// 		}

// 		status := InitialTicketStatus
// 		sustainingStatus := InitialSustainingStatus
// 		for _, history := range issue.Changelog.Histories {
// 			cTime, err := history.CreatedTime()
// 			if err != nil {
// 				return err
// 			}

// 			for _, item := range history.Items {
// 				if item.Field == FieldNameStatus {
// 					status = item.ToString
// 				} else if item.Field == FieldNameSustainingStatus {
// 					sustainingStatus = item.ToString
// 				} else {
// 					continue
// 				}

// 				states = append(states, IssueStatus{
// 					TimeStart:        cTime,
// 					Status:           status,
// 					SustainingStatus: sustainingStatus,
// 					Owner:            statusOwner(status),
// 				})
// 			}
// 		}

// 		for i := range states {
// 			timeEnd := time.Now()
// 			if i < len(states)-1 {
// 				timeEnd = states[i+1].TimeStart
// 			}

// 			states[i].TimeEnd = timeEnd
// 			states[i].Duration = timeEnd.Sub(states[i].TimeStart)
// 		}

// 		h.detail[issue] = states
// 	}
// 	h.populated = true
// 	return nil
// }