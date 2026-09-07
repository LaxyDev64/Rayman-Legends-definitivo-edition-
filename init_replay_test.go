package main

import "testing"

// TestCapturasEmbebidas fija que las capturas estan DENTRO del binario y con el tamano
// exacto que el servidor de agosto anunciaba en su registro.
//
// Estos numeros no son decorativos: salen de las lineas que imprimia el binario que
// funcionaba ("0x73.60 -> 15704 captured bytes"). Si alguien vuelve a poner `*.bin` en el
// .gitignore, o mueve los ficheros, o los trunca, esta prueba falla en la compilacion en vez
// de fallar en la consola de un jugador tres semanas despues.
func TestCapturasEmbebidas(t *testing.T) {
	esperado := map[string]int{
		"captures/resp_0x6e_07.bin":           2590,
		"captures/resp_0x6e_09.bin":           17,
		"captures/resp_0x6e_10.bin":           10,
		"captures/resp_0x73_21.bin":           8,
		"captures/resp_0x73_60.bin":           15704,
		"captures/resp_0x73_84.bin":           29,
		"captures/resp_0x73_95.bin":           6,
		"captures/resp_4bytes_A_00000000.bin": 4,
		"captures/resp_4bytes_B_64000000.bin": 4,
	}
	for ruta, n := range esperado {
		b, err := capturas.ReadFile(ruta)
		if err != nil {
			t.Errorf("%s: no esta embebida (%v)", ruta, err)

			continue
		}
		if len(b) != n {
			t.Errorf("%s: %d octetos, se esperaban %d", ruta, len(b), n)
		}
	}
}

// TestCadaMetodoTieneSuCaptura recorre los nueve metodos con captura y comprueba que el
// mapeo resuelve a un fichero que existe y no esta vacio.
func TestCadaMetodoTieneSuCaptura(t *testing.T) {
	metodos := []struct {
		proto  uint16
		metodo uint32
	}{{0x6E, 7}, {0x6E, 9}, {0x6E, 10}, {0x6E, 11}, {0x73, 21}, {0x73, 60}, {0x73, 61}, {0x73, 84}, {0x73, 95}}

	for _, m := range metodos {
		ruta := ficheroDe(m.proto, m.metodo)
		if ruta == "" {
			t.Errorf("0x%02x.%d no tiene captura asignada", m.proto, m.metodo)

			continue
		}
		b, err := capturas.ReadFile(ruta)
		if err != nil || len(b) == 0 {
			t.Errorf("0x%02x.%d -> %s: ilegible o vacia", m.proto, m.metodo, ruta)
		}
	}
}

// TestLasDosDeCuatroNoSeSolapan protege la unica asignacion que es una suposicion.
//
// Utility 11 y DataStore 61 comparten longitud, asi que no se pudieron separar leyendo el
// binario. Da igual cual sea la correcta, pero NO pueden ser la misma: si un cambio las
// hiciera coincidir, una de las dos estaria contestando la respuesta de la otra y el sintoma
// volveria a ser mudo.
func TestLasDosDeCuatroNoSeSolapan(t *testing.T) {
	a, b := ficheroDe(0x6E, 11), ficheroDe(0x73, 61)
	if a == b {
		t.Fatalf("Utility 11 y DataStore 61 apuntan a la misma captura (%s)", a)
	}
	if a == "" || b == "" {
		t.Fatalf("una de las dos quedo sin asignar: %q / %q", a, b)
	}
}

// TestMetodosDeCuerpoVacio: la 0x73.82 y la 0x73.83 respondian CERO octetos, no una lista
// vacia. Cuatro octetos de mas donde no va ninguno es exactamente la clase de fallo que la
// consola no reporta y que se manifiesta tres pantallas mas adelante.
func TestMetodosDeCuerpoVacio(t *testing.T) {
	for _, m := range []uint32{82, 83} {
		if !cuerpoVacio[claveMetodo{0x73, m}] {
			t.Errorf("0x73.%d deberia responder cuerpo vacio", m)
		}
		if ficheroDe(0x73, m) != "" {
			t.Errorf("0x73.%d no debe tener fichero: su respuesta son cero octetos", m)
		}
	}
}
