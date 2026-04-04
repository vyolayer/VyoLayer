package logger

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	appLevelInfo  = "INFO "
	appLevelWarn  = "WARN "
	appLevelError = "ERROR"
	appLevelDebug = "DEBUG"
)

var (
	serviceLoggerStore sync.Map
	servicePalette     = []string{
		"\033[38;5;39m",  // blue
		"\033[38;5;45m",  // cyan
		"\033[38;5;42m",  // green
		"\033[38;5;214m", // orange
		"\033[38;5;207m", // pink
		"\033[38;5;141m", // purple
		"\033[38;5;220m", // yellow
	}
)

// AppLogger prints colorful logs with service and optional context tags.
type AppLogger struct {
	serviceName string
	contextName string
	writer      *log.Logger
}

// ForService returns a reusable logger instance for a service name.
// The same logger is reused for repeated requests of the same service.
func ForService(serviceName string) *AppLogger {
	name := normalizeLoggerName(serviceName, "app")
	if existing, ok := serviceLoggerStore.Load(name); ok {
		return existing.(*AppLogger)
	}

	created := &AppLogger{
		serviceName: name,
		writer:      log.New(os.Stdout, "", 0),
	}
	actual, _ := serviceLoggerStore.LoadOrStore(name, created)
	return actual.(*AppLogger)
}

// NewAppLogger creates/returns a service-scoped app logger.
func NewAppLogger(serviceName string) *AppLogger {
	return ForService(serviceName)
}

// WithContext returns a child logger that keeps the same service name and adds a context label.
func (l *AppLogger) WithContext(contextName string) *AppLogger {
	if l == nil {
		return ForService("app").WithContext(contextName)
	}

	return &AppLogger{
		serviceName: l.serviceName,
		contextName: normalizeLoggerName(contextName, ""),
		writer:      l.writer,
	}
}

func (l *AppLogger) Info(message string, fields any) {
	l.log(appLevelInfo, message, mergeFields(fields))
}

func (l *AppLogger) Warn(message string, fields any) {
	l.log(appLevelWarn, message, mergeFields(fields))
}

func (l *AppLogger) Error(message string, fields any) {
	l.log(appLevelError, message, mergeFields(fields))
}

func (l *AppLogger) ErrorWithErr(message string, err error) {
	var fields map[string]any
	fields["error"] = err.Error()
	l.log(appLevelError, message, fields)
}

func (l *AppLogger) Debug(message string, fields any) {
	l.log(appLevelDebug, message, mergeFields(fields))
}

func (l *AppLogger) log(level, message string, fields map[string]any) {
	if l == nil {
		l = ForService("app")
	}

	timestamp := time.Now().Format(time.DateTime)
	level = strings.ToUpper(level)

	serviceTag := fmt.Sprintf("[%s]", l.serviceName)
	if isColorEnabled() {
		serviceTag = fmt.Sprintf("%s[%s]%s", serviceColor(l.serviceName), l.serviceName, colorReset)
	}

	contextTag := ""
	if l.contextName != "" {
		contextTag = fmt.Sprintf(" [%s]", l.contextName)
		if isColorEnabled() {
			contextTag = fmt.Sprintf(" %s[%s]%s", colorCyan, l.contextName, colorReset)
		}
	}

	levelTag := fmt.Sprintf("[%s]", level)
	if isColorEnabled() {
		levelTag = fmt.Sprintf("%s[%s]%s", levelColor(level), level, colorReset)
	}

	fieldString := renderFields(fields)
	line := fmt.Sprintf("%s %s %s%s %s%s", timestamp, levelTag, serviceTag, contextTag, message, fieldString)
	l.writer.Println(line)
}

func normalizeLoggerName(name, fallback string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fallback
	}
	return trimmed
}

func mergeFields(fields any) map[string]any {
	result := make(map[string]any)

	if fields == nil {
		return result
	}

	val := reflect.ValueOf(fields)
	typ := reflect.TypeOf(fields)

	// Dereference pointer
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return result
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	switch val.Kind() {

	case reflect.Map:
		for _, key := range val.MapKeys() {
			if key.Kind() == reflect.String {
				result[key.String()] = val.MapIndex(key).Interface()
			}
		}

	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)

			if !field.IsExported() {
				continue
			}

			fieldName := extractFieldName(field)

			if fieldName == "-" {
				continue
			}

			v := val.Field(i)
			if isZero(v) {
				continue
			}

			result[fieldName] = v.Interface()
		}
	}

	return result
}

func extractFieldName(field reflect.StructField) string {
	// Priority: json → db → fallback to field name

	if tag := field.Tag.Get("json"); tag != "" {
		name := strings.Split(tag, ",")[0]
		if name != "" {
			return name
		}
	}

	if tag := field.Tag.Get("db"); tag != "" {
		name := strings.Split(tag, ",")[0]
		if name != "" {
			return name
		}
	}

	return field.Name
}

func renderFields(fields map[string]any) string {
	if len(fields) == 0 {
		return ""
	}

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))

	for _, k := range keys {
		v := fields[k]

		switch val := v.(type) {
		case string:
			parts = append(parts, fmt.Sprintf("%s='%s'", k, val))
		case nil:
			parts = append(parts, fmt.Sprintf("%s=NULL", k))
		default:
			parts = append(parts, fmt.Sprintf("%s=%v", k, val))
		}
	}

	output := strings.Join(parts, " ")

	if isColorEnabled() {
		return " | " + colorBlue + output + colorReset
	}
	return " | " + output
}

func serviceColor(serviceName string) string {
	sum := 0
	for _, c := range serviceName {
		sum += int(c)
	}
	return servicePalette[sum%len(servicePalette)]
}

func isZero(v reflect.Value) bool {
	return v.IsZero()
}
