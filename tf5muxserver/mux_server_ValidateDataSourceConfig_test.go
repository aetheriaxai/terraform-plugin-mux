// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf5testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestMuxServerValidateDataSourceConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf5testserver.TestServer{
		DataSourceSchemas: map[string]*tfprotov5.Schema{
			"test_data_source_server1": {},
		},
	}
	testServer2 := &tf5testserver.TestServer{
		DataSourceSchemas: map[string]*tfprotov5.Schema{
			"test_data_source_server2": {},
		},
	}

	servers := []func() tfprotov5.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = muxServer.ProviderServer().ValidateDataSourceConfig(ctx, &tfprotov5.ValidateDataSourceConfigRequest{
		TypeName: "test_data_source_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.ValidateDataSourceConfigCalled["test_data_source_server1"] {
		t.Errorf("expected test_data_source_server1 ValidateDataSourceConfig to be called on server1")
	}

	if testServer2.ValidateDataSourceConfigCalled["test_data_source_server1"] {
		t.Errorf("unexpected test_data_source_server1 ValidateDataSourceConfig called on server2")
	}

	_, err = muxServer.ProviderServer().ValidateDataSourceConfig(ctx, &tfprotov5.ValidateDataSourceConfigRequest{
		TypeName: "test_data_source_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.ValidateDataSourceConfigCalled["test_data_source_server2"] {
		t.Errorf("unexpected test_data_source_server2 ValidateDataSourceConfig called on server1")
	}

	if !testServer2.ValidateDataSourceConfigCalled["test_data_source_server2"] {
		t.Errorf("expected test_data_source_server2 ValidateDataSourceConfig to be called on server2")
	}
}
