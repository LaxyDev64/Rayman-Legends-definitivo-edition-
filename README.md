# rdv — prototipo de protocolo Quazal RDV para Rayman Legends (Switch)

Prototipo standalone en Go. Es un punto de partida, no una implementación
confirmada byte a byte.

## Qué hay aquí

- `stream.go` — lector/escritor binario mínimo para los tipos comunes
  NEX/RDV (String, List<T>, primitivos, DateTime). Standalone para que
  este prototipo compile solo; en la integración real con tu fork
  `nextendo-nex` debería sustituirse por tu propio `stream.go`/`rmc.go`,
  que ya implementan esto (probablemente de forma más completa).

- `ols_storage.go` — protocolo específico de juego **OLS Storage** (ID
  200). Esta parte está bien fundamentada: la spec completa (16 métodos,
  todos los structs) viene documentada por Kinnay's NintendoClients wiki
  para la versión Wii U de Rayman Legends. Se asume que la lógica de
  juego es la misma en Switch (Ubisoft comparte esa capa entre
  plataformas), pero **falta confirmar** que el formato exacto de los
  mensajes no cambió entre versiones.

- `simple_auth.go` — protocolo **Simple Authentication** (ID 16). Esta
  parte SÍ es un placeholder débil: no hay spec pública tan detallada
  como la de OLS Storage. La forma de los mensajes es una suposición
  razonable, no una confirmación.

## Qué falta / supuestos a validar

1. **Transporte.** Se asume que Quazal RDV comparte el mismo PRUDP/RMC
   que ya implementa `nextendo-nex` (cierto en términos generales — NEX
   deriva de RDV — pero puede haber diferencias de versión/config no
   confirmadas).
2. **Simple Authentication.** Placeholder — necesita reverse engineering
   real (ver `acb-rdv`/`QuazalWV` en C#, o capturar tráfico).
3. **OLS Storage en Switch específicamente.** Documentado para Wii U;
   asumido igual en Switch. Puede haber cambios de campos si Ubisoft
   actualizó el protocolo para el port de 2017.
4. **Autenticación de más alto nivel.** No está claro si Switch sigue
   usando el viejo esquema HTTP (`onlineconfigservice.ubi.com` style,
   como en `acb-rdv`) o algo más moderno (Uplay Win, protocolo 49 de la
   lista de Kinnay). Por investigar.
5. **Intercepción de tráfico en una Switch real** (o emulador) para
   confirmar todo lo anterior.

## Próximo paso sugerido

Capturar tráfico real del juego (Switch física con red controlada, o
emulador) para confirmar 2, 3 y 4 antes de invertir más tiempo en
lógica de servidor que podría no coincidir con el protocolo real.
