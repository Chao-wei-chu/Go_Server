package main

import (
	//"log"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

var (
	GLOBAL_DB             *sql.DB
	GLOBAL_BasicData      *Safe_BasicData
	GLOBAL_BasicData_FAKE *Safe_BasicData
	GLOBAL_CharactData    *Safe_CharactTableData
	GLOBAL_MeasureData    *Safe_MeasureTableData
	GLOBAL_ExecutionData  *Safe_ExecutionTableData
	GLOBAL_SenarioRules   RuleSet //建立包含所有SenarioRule的Map
	GLOBAL_ExecutionRules ExeSet  //建立包含所有ExeRule的Map
)

type Safe_BasicData struct {
	data   []*HealthBasicTableData
	length int
	mux    sync.Mutex
	isNew  bool
}

type Safe_CharactTableData struct {
	data   []*HealthCharactTableData
	length int
	mux    sync.Mutex
	isNew  bool
}

type Safe_MeasureTableData struct {
	data   []*HealthMeasureTableData
	length int
	mux    sync.Mutex
	isNew  bool
}

type Safe_ExecutionTableData struct {
	data   []*ServiceExecutionTableData
	length int
	mux    sync.Mutex
	isNew  bool
}

type HealthBasicTableData struct {
	RuleID        int
	SenarioID     string
	Code          int
	Message       string
	Description   string
	Href          string
	StartDateTime time.Time
	EndDateTime   time.Time
	Fake          int
}

type HealthCharactTableData struct {
	RuleID    int
	SenarioID string
	Name      string
	Value     string
}

type HealthMeasureTableData struct {
	RuleID            int
	SenarioID         string
	MetricName        string
	MetricDescription string
	ValueType         string
	UnitOfMeasure     string
	Value             string
	CaptureDateTime   time.Time
	ResourceID        string
	ResourceType      string
	ResourceName      string
}

type ServiceExecutionTableData struct {
	RuleID    int
	SenarioID string
	Code      int
	Message   string
	LastTime  time.Time
}

type RuleEntry struct {
	Code              string
	Message           string
	State             string
	ResourceID        string
	ResourceName      string
	MetricDescription string
}

type ExeEntry struct {
	Code    string
	Message string
	Action  string
}

type TestEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TestJSON struct {
	TestList []TestEntry `json:"serviceTestCharcteristic"`
}

type RuleSet map[string][]RuleEntry
type ExeSet map[string][]ExeEntry

func VAR_INIT() {
	log.Println(`VAR_INIT()`)
	/*----------Init DB用變數-----------*/
	GLOBAL_BasicData = new(Safe_BasicData)
	GLOBAL_BasicData.length = 0
	GLOBAL_BasicData.isNew = true
	GLOBAL_BasicData_FAKE = new(Safe_BasicData)
	GLOBAL_BasicData_FAKE.length = 0
	GLOBAL_BasicData_FAKE.isNew = true
	GLOBAL_CharactData = new(Safe_CharactTableData)
	GLOBAL_CharactData.length = 0
	GLOBAL_CharactData.isNew = true
	GLOBAL_MeasureData = new(Safe_MeasureTableData)
	GLOBAL_MeasureData.length = 0
	GLOBAL_MeasureData.isNew = true
	GLOBAL_ExecutionData = new(Safe_ExecutionTableData)
	GLOBAL_ExecutionData.length = 0
	GLOBAL_ExecutionData.isNew = true
	/*----Init Senario Rules for ServiceHealthState------*/
	if err := json.Unmarshal([]byte(ENV.SENARIO_RULES), &GLOBAL_SenarioRules); err != nil {
		log.Fatal("[ERROR] ", err, ", Program Shutdown.") //注意，這裡沒有rule甚麼都無法判斷，因此失敗回直接中斷程式
	}
	log.Printf("[INFO] %+v", GLOBAL_SenarioRules["cmp_rrn"])
}

func (t *RuleSet) UnmarshalJSON(b []byte) error {
	// Create a local struct that mirrors the data being unmarshalled
	type tcEntry struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		State   string `json:"state"`
	}

	type tcMain struct {
		Senario_id        string    `json:"senario_id"`
		ResourceID        string    `json:"resource_id"`
		ResourceName      string    `json:"resource_name"`
		MetricDescription string    `json:"metricDescription"`
		TcEntry           []tcEntry `json:"rules"`
	}
	// unmarshal the data into the slice
	var entries []tcMain
	if err := json.Unmarshal(b, &entries); err != nil {
		return err
	}
	//注意: Map一定要做make宣告
	totalMain := make(RuleSet)
	// loop over the slice and create the map of entries
	for _, ent := range entries {
		for _, ent2 := range ent.TcEntry {
			totalMain[ent.Senario_id] = append(totalMain[ent.Senario_id], RuleEntry{Code: ent2.Code, Message: ent2.Message, State: ent2.State, ResourceID: ent.ResourceID, ResourceName: ent.ResourceName, MetricDescription: ent.MetricDescription})
		}
	}
	*t = totalMain
	fmt.Println("This is (t *RuleSet)UnmarshalJSON ")
	return nil
}

func (t *ExeSet) UnmarshalJSON(b []byte) error {
	// Create a local struct that mirrors the data being unmarshalled
	type tcEntry struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Action  string `json:"action"`
	}

	type tcMain struct {
		Senario_id string    `json:"senario_id"`
		TcEntry    []tcEntry `json:"rules"`
	}
	// unmarshal the data into the slice
	var entries []tcMain
	if err := json.Unmarshal(b, &entries); err != nil {
		return err
	}
	//注意: Map一定要做make宣告
	totalMain := make(ExeSet)
	// loop over the slice and create the map of entries
	for _, ent := range entries {
		for _, ent2 := range ent.TcEntry {
			totalMain[ent.Senario_id] = append(totalMain[ent.Senario_id], ExeEntry{Code: ent2.Code, Message: ent2.Message, Action: ent2.Action})
		}
	}
	*t = totalMain
	fmt.Println("This is (t *ExeSet)UnmarshalJSON ")
	return nil
}

func ClearMemShowTimeStats(startingTime time.Time, endingTime time.Time) {
	var duration time.Duration = endingTime.Sub(startingTime)
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("[INFO] Time spent: [%v]ms [%.3f]sec,  Alloc = %v KiB TotalAlloc = %v KiB Sys = %v KiB, GCnum = %v \n", duration.Milliseconds(), duration.Seconds(), m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC)
}

func ClearMemStats() {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("[INFO] Alloc = %v KiB TotalAlloc = %v KiB Sys = %v KiB, GCnum = %v \n", m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC)
}
