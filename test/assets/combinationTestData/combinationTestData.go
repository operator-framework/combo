package combinationTestData

var CombinationInput map[string][]string = map[string][]string{
	"TEST1": {"foo", "bar"},
	"TEST2": {"zip", "zap"},
	"TEST3": {"bip", "bap"},
}

var CombinationOutput []map[string]string = []map[string]string{
	{
		"TEST1": "foo",
		"TEST2": "zip",
		"TEST3": "bip",
	},
	{
		"TEST1": "foo",
		"TEST2": "zap",
		"TEST3": "bip",
	},
	{
		"TEST1": "bar",
		"TEST2": "zip",
		"TEST3": "bip",
	},
	{
		"TEST1": "bar",
		"TEST2": "zap",
		"TEST3": "bip",
	},
	{
		"TEST1": "foo",
		"TEST2": "zip",
		"TEST3": "bap",
	},
	{
		"TEST1": "foo",
		"TEST2": "zap",
		"TEST3": "bap",
	},
	{
		"TEST1": "bar",
		"TEST2": "zip",
		"TEST3": "bap",
	},
	{
		"TEST1": "bar",
		"TEST2": "zap",
		"TEST3": "bap",
	},
}
