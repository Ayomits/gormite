package local_schema

import (
	"fmt"
	"github.com/KoNekoD/gormite/pkg/utils"
	"strings"
)

func applyMetadataMutatorsForNewColumn(
	columnTagsData *columnData,
	bag *tableBag,
) {
	if columnTagsData.IsForeignKey {
		bag.table.AddForeignKeyConstraint(
			getName(bag.store, columnTagsData.TypeName),
			[]string{columnTagsData.ColumnName},
			[]string{"id"},
			nil,
			nil,
		)
	}

	if columnTagsData.IsPrimaryKey {
		bag.primaryKeys = append(bag.primaryKeys, columnTagsData.ColumnName)
	}

	if columnTagsData.IsUnique {
		uniqNames := strings.Split(*columnTagsData.UniqueName, ",")
		for _, uniqNameItem := range uniqNames {
			uniqNameItem = strings.TrimSpace(uniqNameItem)
			if _, hasUniqMapKey := bag.uniqColumnsMap[uniqNameItem]; !hasUniqMapKey {
				bag.uniqColumnsMap[uniqNameItem] = make([]string, 0)
			}
			bag.uniqColumnsMap[uniqNameItem] = append(
				bag.uniqColumnsMap[uniqNameItem],
				columnTagsData.ColumnName,
			)
		}
	}

	if columnTagsData.IsUniqueCondition {
		conditions := strings.Split(*columnTagsData.UniqueCondition, ";")
		for _, condition := range conditions {
			conditionParts := strings.Split(condition, ":")
			if len(conditionParts) != 2 {
				panic(fmt.Sprintf("invalid uniq condition %s", condition))
			}

			bag.uniqConditionsMap[conditionParts[0]] = conditionParts[1]
		}
	}

	if columnTagsData.IsIndex {
		indexNames := strings.Split(*columnTagsData.IndexName, ",")
		for _, indexNameItem := range indexNames {
			indexNameItem = strings.TrimSpace(indexNameItem)
			if _, hasIndexMapKey := bag.indexColumnsMap[indexNameItem]; !hasIndexMapKey {
				bag.indexColumnsMap[indexNameItem] = make([]string, 0)
			}
			bag.indexColumnsMap[indexNameItem] = append(
				bag.indexColumnsMap[indexNameItem],
				columnTagsData.ColumnName,
			)
		}
	}

	if columnTagsData.IsIndexCondition {
		conditions := strings.Split(*columnTagsData.IndexCondition, ";")
		for _, condition := range conditions {
			conditionParts := strings.Split(condition, ":")
			if len(conditionParts) != 2 {
				panic(fmt.Sprintf("invalid index condition %s", condition))
			}

			bag.indexConditionsMap[conditionParts[0]] = conditionParts[1]
		}
	}
}

func applyMetadataMutatorsAfterColumnsIntrospection(bag *tableBag) {
	for indexName, columns := range bag.indexColumnsMap {
		options := make(map[string]interface{})

		if v, ok := bag.indexConditionsMap[indexName]; ok {
			options["where"] = v
		}

		bag.table.AddIndex(columns, &indexName, make([]string, 0), options)
	}

	for uniqPseudoName, columns := range bag.uniqColumnsMap {
		uniqIdxName := fmt.Sprintf(
			"idx__%s__%s__uniq",
			bag.table.GetName(),
			strings.Join(columns, "_"),
		)

		options := make(map[string]interface{})

		if v, ok := bag.uniqConditionsMap[uniqPseudoName]; ok {
			options["where"] = v
		}

		bag.table.AddUniqueIndex(columns, &uniqIdxName, options)
	}

	bag.table.SetPrimaryKey(
		bag.primaryKeys,
		utils.AsPtr(fmt.Sprintf("%s_pkey", bag.table.GetName())),
	)
}
