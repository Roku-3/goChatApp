package trace
import(
    "bytes"
    "testing"
)
func TestNew(t *testing.T) {
    var buf bytes.Buffer
    tracer := New(&buf)
    if tracer == nil {
        t.Error("Newからの戻り値がnil")
    } else {
        tracer.Trace("hello, trace package")
        if buf.String() != "hello, trace package\n" {
            t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
        }
    }
}

func TestOff(t *testing.T) {
    var silentTracer Tracer = Off()
    silentTracer.Trace("データ")
}

