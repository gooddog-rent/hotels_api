package entity

/* func TestValidBuild(t *testing.T) {

	l := &LocationsStore{}

	tree := suffixtree.NewGeneralizedSuffixTree()
	searchTree := treedb.NewTreeRepo(tree).BuildTree(l)

	if searchTree.SuffixTree == nil {
		t.Errorf("BuildSuffixTree must be valid tree! got: %v", searchTree.SuffixTree)
	}
}

func TestValidWithDataBuild(t *testing.T) {

	l := &LocationsStore{
		Locations: []Location{
			{
				Region: "Punta Cana",
				Hotels: []string{"Riu Bavaro", "Hotel Bavaro", "Resort Beach 5"},
			},
			{
				Region: "Macao",
				Hotels: []string{"Beach Summer 4", "Macao Beach Resort", "Star Hotel Resort"},
			},
		},
	}

	loc := &LocationsStore{
		Locations: l.Locations,
	}

	tree := suffixtree.NewGeneralizedSuffixTree()
	searchTree := treedb.NewTreeRepo(tree).BuildTree(loc)

	result := map[string][]string{
		"Punta Cana": {"Riu Bavaro", "Hotel Bavaro", "Resort Beach 5"},
		"Macao":      {"Beach Summer 4", "Macao Beach Resort", "Star Hotel Resort"},
	}

	// check valid regions in tree
	for _, location := range l.Locations {
		if _, ok := result[location.Region]; !ok {
			t.Errorf("BuildSuffixTree must have a valid regions! got: %v", location)
		}
	}

	// check valid hotels in tree
	for _, location := range l.Locations {

		if _, ok := result[location.Region]; ok {
			for index, h := range result[location.Region] {
				if h != location.Hotels[index] {
					t.Errorf("BuildSuffixTree must have a valid hotels! got: %v", h)
				}
			}
		}
	}

	if searchTree.SuffixTree == nil {
		t.Errorf("BuildSuffixTree must be valid tree! got: %v", searchTree.SuffixTree)
	}
} */
