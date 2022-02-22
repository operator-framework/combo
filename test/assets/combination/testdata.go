package testdata

var OneParameterCombinationInput = map[string][]string{
	"TEST1": {"foo", "bar"},
}

var OneParameterCombinationOutput = []map[string]string{
	{
		"TEST1": "foo",
	},
	{
		"TEST1": "bar",
	},
}

var CombinationInput = map[string][]string{
	"TEST1": {"foo", "bar"},
	"TEST2": {"zip", "zap"},
	"TEST3": {"bip", "bap"},
}

var CombinationOutput = []map[string]string{
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

var LongCombinationInput = map[string][]string{
	"TEST1": {"foo", "bar"},
	"TEST2": {"zip"},
	"TEST3": {"bip", "bap"},
	"TEST4": {"pip"},
	"TEST5": {"mip"},
	"TEST6": {"zip", "zap", "zop"},
}

var LongCombinationOutput = []map[string]string{
	{
		"TEST1": "foo",
		"TEST2": "zip",
		"TEST3": "bip",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zip",
	},
	{
		"TEST1": "foo",
		"TEST2": "zip",
		"TEST3": "bip",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zap",
	},
	{
		"TEST1": "foo",
		"TEST2": "zip",
		"TEST3": "bip",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zop",
	},
	{
		"TEST1": "foo",
		"TEST2": "zip",
		"TEST3": "bap",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zip",
	},
	{
		"TEST1": "foo",
		"TEST2": "zip",
		"TEST3": "bap",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zap",
	},
	{
		"TEST1": "foo",
		"TEST2": "zip",
		"TEST3": "bap",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zop",
	},
	{
		"TEST1": "bar",
		"TEST2": "zip",
		"TEST3": "bip",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zip",
	},
	{
		"TEST1": "bar",
		"TEST2": "zip",
		"TEST3": "bip",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zap",
	},
	{
		"TEST1": "bar",
		"TEST2": "zip",
		"TEST3": "bip",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zop",
	},
	{
		"TEST1": "bar",
		"TEST2": "zip",
		"TEST3": "bap",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zip",
	},
	{
		"TEST1": "bar",
		"TEST2": "zip",
		"TEST3": "bap",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zap",
	},
	{
		"TEST1": "bar",
		"TEST2": "zip",
		"TEST3": "bap",
		"TEST4": "pip",
		"TEST5": "mip",
		"TEST6": "zop",
	},
}
