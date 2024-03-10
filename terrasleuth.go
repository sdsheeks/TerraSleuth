package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

type TerraformPlan struct {
	TerraformVersion string `json:"terraform_version"`
	PlannedValues    struct {
		RootModule struct {
			Resources []struct {
				Address      string      `json:"address"`
				Type         string      `json:"type"`
				Name         string      `json:"name"`
				ProviderName string      `json:"provider_name"`
				Values       interface{} `json:"values"`
			} `json:"resources"`
		} `json:"root_module"`
	} `json:"planned_values"`
	ResourceChanges []struct {
		Address string `json:"address"`
		Change  struct {
			Actions []string    `json:"actions"`
			After   interface{} `json:"after"`
		} `json:"change"`
	} `json:"resource_changes"`
}

func main() {

	file, err := os.Open("plan.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close()

	var plan TerraformPlan
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&plan)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	defer w.Flush()

	fmt.Fprintln(w, "\nTerraform Version\t")
	fmt.Fprintf(w, "%s\t\n", plan.TerraformVersion)

	fmt.Fprintln(w, "\nResources\t")
	fmt.Fprintln(w, "Resource\tType\tName\tProvider Name\t")
	for _, resource := range plan.PlannedValues.RootModule.Resources {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n",
			resource.Address, resource.Type, resource.Name, resource.ProviderName)
	}

	fmt.Fprintln(w, "\nResource Changes\t")
	fmt.Fprintln(w, "Resource\tActions\t")
	for _, change := range plan.ResourceChanges {
		fmt.Fprintf(w, "%s\t%v\t\n", change.Address, change.Change.Actions)
	}
}
