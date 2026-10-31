package main

import {
    "testing"
    "net/http/httptest"
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d; esperado %d", res.StatusCode, http.StatusOK)
	}
	if tipo := res.Header.Get("Content-Type"); tipo != "application/json" {
		t.Errorf("Content-Type = %q", tipo)
	}
}