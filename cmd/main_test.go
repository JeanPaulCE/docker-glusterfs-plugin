package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	// Este es un test básico para verificar que el paquete main se compila correctamente
	// No podemos probar la función main directamente ya que inicia un servidor
	t.Run("package compiles", func(t *testing.T) {
		// Si llegamos aquí, el paquete se compiló correctamente
	})
} 