// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/table/entities"
)

type StorageTableEntityResource struct{}

func TestAccTableEntity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")
	r := StorageTableEntityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTableEntity_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")
	r := StorageTableEntityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccTableEntity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")
	r := StorageTableEntityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTableEntity_update_typed(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")
	r := StorageTableEntityResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated_typed(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated_typedInt64(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated_typedDouble(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated_typedString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated_typedBoolean(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r StorageTableEntityResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := entities.ParseResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	account, err := client.Storage.FindAccount(ctx, id.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Table %q: %+v", id.AccountName, id.TableName, err)
	}
	if account == nil {
		return nil, fmt.Errorf("storage Account %q was not found", id.AccountName)
	}

	entitiesClient, err := client.Storage.TableEntityClient(ctx, *account)
	if err != nil {
		return nil, fmt.Errorf("building Table Entity Client: %+v", err)
	}

	input := entities.GetEntityInput{
		PartitionKey:  id.PartitionKey,
		RowKey:        id.RowKey,
		MetaDataLevel: entities.NoMetaData,
	}
	resp, err := entitiesClient.Get(ctx, id.AccountName, id.TableName, input)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Entity (Partition Key %q / Row Key %q) (Table %q / Storage Account %q / Resource Group %q): %+v", id.PartitionKey, id.RowKey, id.TableName, id.AccountName, account.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r StorageTableEntityResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo = "Bar"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r StorageTableEntityResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "import" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo = "Bar"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r StorageTableEntityResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo  = "Bar"
    Test = "Updated"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r StorageTableEntityResource) updated_typed(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo              = 123
    "Foo@odata.type" = "Edm.Int32"
    Test             = "Updated"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r StorageTableEntityResource) updated_typedInt64(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo              = 123
    "Foo@odata.type" = "Edm.Int64"
    Test             = "Updated"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r StorageTableEntityResource) updated_typedDouble(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo              = 123.123
    "Foo@odata.type" = "Edm.Double"
    Test             = "Updated"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r StorageTableEntityResource) updated_typedString(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo              = "123.123"
    "Foo@odata.type" = "Edm.String"
    Test             = "Updated"
  }
  lifecycle {
    ignore_changes = [entity]
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r StorageTableEntityResource) updated_typedBoolean(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = azurerm_storage_account.test.name
  table_name           = azurerm_storage_table.test.name

  partition_key = "test_partition%d"
  row_key       = "test_row%d"
  entity = {
    Foo              = "true"
    "Foo@odata.type" = "Edm.Boolean"
    Test             = "Updated"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r StorageTableEntityResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  storage_account_name = azurerm_storage_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
