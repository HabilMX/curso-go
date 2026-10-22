package consulta

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/habil/revisor/internal/servicio"
)

// TestRevisar_elMotivoNoLlevaLaURL: el error de net/http trae la URL completa
// («Get "http://...?token=..."»), y el motivo sale tal cual en el reporte. Se
// prueba contra un servidor que acaba de cerrarse, para que la conexión falle
// de verdad y no por una URL mal escrita.
func TestRevisar_elMotivoNoLlevaLaURL(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	url := strings.Replace(srv.URL, "http://", "http://yo:supersecreto@", 1) + "/estado?token=abc123#ancla"
	srv.Close()

	e := NuevoHTTP(time.Second).Revisar(context.Background(), servicio.Servicio{Nombre: "pagos", URL: url})

	if e.OK() {
		t.Fatal("un servidor cerrado no puede salir sano")
	}
	for _, secreto := range []string{"yo:", "supersecreto", "/estado", "token", "abc123", "ancla"} {
		if strings.Contains(e.Motivo(), secreto) {
			t.Errorf("el motivo filtra %q: %q", secreto, e.Motivo())
		}
	}
	if !strings.HasPrefix(e.Motivo(), "no responde: ") {
		t.Errorf("motivo = %q, quería que empezara con «no responde: »", e.Motivo())
	}
}
